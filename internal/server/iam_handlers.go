package server

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"awo.so/framework/internal/authz"
	iamauth "awo.so/framework/internal/iam/auth"
	"awo.so/framework/internal/middleware"
	"awo.so/framework/internal/tenant"
	"awo.so/framework/pkg/permissions"
)

// IAMDeps holds services needed by user/role management handlers.
type IAMDeps struct {
	Authz *authz.Service
	DB    *pgxpool.Pool // system DB — users table lives here
	Log   *slog.Logger
}

// MountIAMRoutes registers /iam/* routes for user and role management.
// Requires authDeps.Auth for user creation (password hashing).
func MountIAMRoutes(app *fiber.App, d *IAMDeps) {
	g := app.Group("/iam")

	// Users
	g.Get("/users", handleListUsers(d))
	g.Post("/users", handleCreateUser(d))
	g.Get("/users/:id", handleGetUser(d))
	g.Patch("/users/:id", handleUpdateUser(d))
	g.Delete("/users/:id", handleDeactivateUser(d))

	// Roles
	g.Get("/roles", handleListRoles(d))

	// Role assignments
	g.Get("/users/:id/roles", handleGetUserRoles(d))
	g.Post("/users/:id/roles", handleAssignRole(d))
	g.Delete("/users/:id/roles/:role", handleRevokeRole(d))

	// Policies (read-only view of Casbin rules for this tenant)
	g.Get("/policies", handleListPolicies(d))
}

// GET /iam/users
func handleListUsers(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}
		if !p.IsSuper && !hasIAMAccess(d, p, t) {
			return errForbidden(c, "insufficient privileges")
		}

		// Repo is the tenant DB — accessed via the request context pool.
		// We read from the per-tenant schema via the tenant's DB pool injected by middleware.
		rows, err := d.DB.Query(c.UserContext(), `
			SELECT id, email, full_name, user_type,
			       is_active, is_suspended, is_super,
			       last_login_at, created_at
			FROM users
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2`,
			pageSize(c), c.QueryInt("offset", 0))
		if err != nil {
			d.Log.Error("list users", slog.Any("err", err))
			return errInternal(c)
		}
		defer rows.Close()

		type userRow struct {
			ID          string     `json:"id"`
			Email       string     `json:"email"`
			FullName    string     `json:"full_name"`
			UserType    string     `json:"user_type"`
			IsActive    bool       `json:"is_active"`
			IsSuspended bool       `json:"is_suspended"`
			IsSuper     bool       `json:"is_super"`
			LastLoginAt *time.Time `json:"last_login_at,omitempty"`
			CreatedAt   time.Time  `json:"created_at"`
		}

		var users []userRow
		for rows.Next() {
			var u userRow
			if err := rows.Scan(
				&u.ID, &u.Email, &u.FullName, &u.UserType,
				&u.IsActive, &u.IsSuspended, &u.IsSuper,
				&u.LastLoginAt, &u.CreatedAt,
			); err != nil {
				return errInternal(c)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			return errInternal(c)
		}
		return ok(c, users)
	}
}

// POST /iam/users  body: {email, full_name, password, user_type}
func handleCreateUser(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}
		if !p.IsSuper && !hasIAMAccess(d, p, t) {
			return errForbidden(c, "insufficient privileges")
		}

		var body struct {
			Email    string `json:"email"`
			FullName string `json:"full_name"`
			Password string `json:"password"`
			UserType string `json:"user_type"`
		}
		if err := c.BodyParser(&body); err != nil {
			return errBadRequest(c, "invalid request body")
		}
		body.Email = strings.TrimSpace(strings.ToLower(body.Email))
		if body.Email == "" || body.Password == "" {
			return errBadRequest(c, "email and password are required")
		}
		if len(body.Password) < 8 {
			return errBadRequest(c, "password must be at least 8 characters")
		}
		if body.UserType == "" {
			body.UserType = "INTERNAL"
		}

		hash, err := iamauth.HashPassword(body.Password)
		if err != nil {
			return errInternal(c)
		}

		var userID uuid.UUID
		err = d.DB.QueryRow(c.UserContext(), `
			INSERT INTO users (email, full_name, password_hash, user_type, tenant_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id`,
			body.Email, body.FullName, hash, body.UserType, t.ID,
		).Scan(&userID)
		if err != nil {
			if strings.Contains(err.Error(), "unique") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exists"})
			}
			d.Log.Error("create user", slog.Any("err", err))
			return errInternal(c)
		}

		return okCreated(c, fiber.Map{
			"id":        userID,
			"email":     body.Email,
			"full_name": body.FullName,
			"user_type": body.UserType,
		})
	}
}

// GET /iam/users/:id
func handleGetUser(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}
		// Users may view their own profile; admins view any.
		userID := c.Params("id")
		if userID != p.UserID && !p.IsSuper && !hasIAMAccess(d, p, t) {
			return errForbidden(c, "insufficient privileges")
		}

		row := d.DB.QueryRow(c.UserContext(), `
			SELECT id, email, full_name, user_type,
			       is_active, is_suspended, is_super,
			       mfa_enabled, last_login_at, created_at, updated_at
			FROM users WHERE id = $1`, userID)

		var u struct {
			ID          string     `json:"id"`
			Email       string     `json:"email"`
			FullName    string     `json:"full_name"`
			UserType    string     `json:"user_type"`
			IsActive    bool       `json:"is_active"`
			IsSuspended bool       `json:"is_suspended"`
			IsSuper     bool       `json:"is_super"`
			MFAEnabled  bool       `json:"mfa_enabled"`
			LastLoginAt *time.Time `json:"last_login_at,omitempty"`
			CreatedAt   time.Time  `json:"created_at"`
			UpdatedAt   time.Time  `json:"updated_at"`
		}
		if err := row.Scan(
			&u.ID, &u.Email, &u.FullName, &u.UserType,
			&u.IsActive, &u.IsSuspended, &u.IsSuper, &u.MFAEnabled,
			&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return errNotFound(c, "user not found")
		}
		return ok(c, u)
	}
}

// PATCH /iam/users/:id  body: {full_name?, is_active?, is_suspended?}
func handleUpdateUser(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}
		if !p.IsSuper && !hasIAMAccess(d, p, t) {
			return errForbidden(c, "insufficient privileges")
		}

		var body struct {
			FullName    *string `json:"full_name"`
			IsActive    *bool   `json:"is_active"`
			IsSuspended *bool   `json:"is_suspended"`
		}
		if err := c.BodyParser(&body); err != nil {
			return errBadRequest(c, "invalid request body")
		}

		userID := c.Params("id")
		tag, err := d.DB.Exec(c.UserContext(), `
			UPDATE users
			SET full_name    = COALESCE($2, full_name),
			    is_active    = COALESCE($3, is_active),
			    is_suspended = COALESCE($4, is_suspended),
			    updated_at   = NOW()
			WHERE id = $1`,
			userID, body.FullName, body.IsActive, body.IsSuspended)
		if err != nil {
			d.Log.Error("update user", slog.Any("err", err))
			return errInternal(c)
		}
		if tag.RowsAffected() == 0 {
			return errNotFound(c, "user not found")
		}
		return ok(c, fiber.Map{"updated": true})
	}
}

// DELETE /iam/users/:id — soft deactivate only, never hard delete.
func handleDeactivateUser(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}
		if !p.IsSuper && !hasIAMAccess(d, p, t) {
			return errForbidden(c, "insufficient privileges")
		}
		if c.Params("id") == p.UserID {
			return errBadRequest(c, "cannot deactivate your own account")
		}

		tag, err := d.DB.Exec(c.UserContext(),
			`UPDATE users SET is_active = FALSE, updated_at = NOW() WHERE id = $1`,
			c.Params("id"))
		if err != nil {
			return errInternal(c)
		}
		if tag.RowsAffected() == 0 {
			return errNotFound(c, "user not found")
		}
		return ok(c, fiber.Map{"deactivated": true})
	}
}

// GET /iam/roles — lists roles seeded for this tenant in Casbin.
func handleListRoles(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}

		domain := authz.TenantDomain(t.ID.String())
		policies, err := d.Authz.GetPolicies(c.UserContext(), domain)
		if err != nil {
			return errInternal(c)
		}

		// Collect unique role names from p-rules.
		seen := make(map[string]struct{})
		roles := []string{}
		for _, pol := range policies {
			if _, ok := seen[pol.Subject]; !ok {
				seen[pol.Subject] = struct{}{}
				roles = append(roles, pol.Subject)
			}
		}
		return ok(c, roles)
	}
}

// GET /iam/users/:id/roles
func handleGetUserRoles(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}

		subject := authz.TenantSubject(c.Params("id"))
		domain := authz.TenantDomain(t.ID.String())

		roles, err := d.Authz.GetRoles(c.UserContext(), subject, domain)
		if err != nil {
			return errInternal(c)
		}
		return ok(c, roles)
	}
}

// POST /iam/users/:id/roles  body: {role, expires_at?}
func handleAssignRole(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}
		if !p.IsSuper && !hasIAMAccess(d, p, t) {
			return errForbidden(c, "insufficient privileges")
		}

		var body struct {
			Role      string     `json:"role"`
			ExpiresAt *time.Time `json:"expires_at"`
		}
		if err := c.BodyParser(&body); err != nil || body.Role == "" {
			return errBadRequest(c, "role is required")
		}

		subject := authz.TenantSubject(c.Params("id"))
		domain := authz.TenantDomain(t.ID.String())

		opts := []authz.AssignOpt{
			authz.WithAssignedBy(authz.TenantSubject(p.UserID)),
		}
		if body.ExpiresAt != nil {
			opts = append(opts, authz.WithExpiry(*body.ExpiresAt))
		}

		if err := d.Authz.AssignRole(c.UserContext(), subject, body.Role, domain, opts...); err != nil {
			d.Log.Error("assign role", slog.Any("err", err))
			return errInternal(c)
		}
		return okCreated(c, fiber.Map{"assigned": true, "role": body.Role})
	}
}

// DELETE /iam/users/:id/roles/:role
func handleRevokeRole(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}
		if !p.IsSuper && !hasIAMAccess(d, p, t) {
			return errForbidden(c, "insufficient privileges")
		}

		subject := authz.TenantSubject(c.Params("id"))
		domain := authz.TenantDomain(t.ID.String())

		if err := d.Authz.RevokeRole(c.UserContext(), subject, c.Params("role"), domain); err != nil {
			d.Log.Error("revoke role", slog.Any("err", err))
			return errInternal(c)
		}
		return ok(c, fiber.Map{"revoked": true})
	}
}

// GET /iam/policies — read-only view of Casbin p-rules for this tenant.
func handleListPolicies(d *IAMDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}
		if !p.IsSuper && !hasIAMAccess(d, p, t) {
			return errForbidden(c, "insufficient privileges")
		}

		domain := authz.TenantDomain(t.ID.String())
		policies, err := d.Authz.GetPolicies(c.UserContext(), domain)
		if err != nil {
			return errInternal(c)
		}
		return ok(c, policies)
	}
}

// hasIAMAccess returns true when the principal holds the admin or tenant_admin role
// in the tenant's Casbin domain — without a full Enforce round-trip.
func hasIAMAccess(d *IAMDeps, p *permissions.Principal, t *tenant.Tenant) bool {
	if d.Authz == nil {
		return false
	}
	subject := authz.TenantSubject(p.UserID)
	domain := authz.TenantDomain(t.ID.String())
	roles, err := d.Authz.GetRoles(context.Background(), subject, domain)
	if err != nil {
		return false
	}
	for _, r := range roles {
		if r == "admin" || r == "tenant_admin" {
			return true
		}
	}
	return false
}
