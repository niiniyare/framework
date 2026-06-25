package server

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"

	iamauth "awo.so/framework/internal/iam/auth"
	"awo.so/framework/internal/iam/apikey"
	"awo.so/framework/internal/iam/domain"
	"awo.so/framework/internal/tenant"
)

const sessionCookie = "awo_session"

// AuthDeps holds IAM services needed by auth and API key handlers.
type AuthDeps struct {
	Auth      *iamauth.Service
	APIKeys   *apikey.Service      // may be nil
	Lifecycle *tenant.Lifecycle    // may be nil; required for /auth/register
	Log       *slog.Logger
}

// MountAuthRoutes registers /auth/* and /api-keys/* on the given Fiber app.
// Must be called before the per-entity routes so the prefixes take priority.
func MountAuthRoutes(app *fiber.App, d *AuthDeps) {
	g := app.Group("/auth")
	g.Post("/login", handleLogin(d))
	g.Post("/logout", handleLogout(d))
	g.Post("/mfa/complete", handleMFAComplete(d))
	g.Get("/me", handleMe(d))
	g.Post("/register", handleRegister(d)) // public — tenant middleware exempted in server.go

	if d.APIKeys != nil {
		MountAPIKeyRoutes(app, &APIKeyDeps{APIKeys: d.APIKeys, Log: d.Log})
	}
}

// POST /auth/login
func handleLogin(d *AuthDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		if t == nil {
			return errUnauthorized(c, "tenant context missing")
		}

		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return errBadRequest(c, "invalid request body")
		}
		body.Email = strings.TrimSpace(strings.ToLower(body.Email))
		if body.Email == "" || body.Password == "" {
			return errBadRequest(c, "email and password are required")
		}

		resolved, rawToken, err := d.Auth.Login(c.UserContext(), domain.LoginParams{
			Email:     body.Email,
			Password:  body.Password,
			IP:        c.IP(),
			UserAgent: c.Get("User-Agent"),
		})

		switch err {
		case nil:
			// Full session issued
			setSessionCookie(c, rawToken, int(resolved.ExpiresAt.Unix()))
			return ok(c, fiber.Map{
				"user_id":      resolved.UserID,
				"display_name": resolved.DisplayName,
				"expires_at":   resolved.ExpiresAt,
			})

		case iamauth.ErrMFARequired:
			// rawToken is the pending token; don't set cookie yet
			return ok(c, fiber.Map{
				"mfa_required":  true,
				"pending_token": rawToken,
			})

		case iamauth.ErrInvalidCredentials:
			return errUnauthorized(c, "invalid credentials")

		case iamauth.ErrAccountLocked:
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "account temporarily locked",
			})

		default:
			d.Log.Error("login error", slog.Any("err", err))
			return errInternal(c)
		}
	}
}

// POST /auth/mfa/complete
// TODO: implement TOTP verification in Phase 2.
func handleMFAComplete(_ *AuthDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "mfa not yet implemented"})
	}
}

// POST /auth/logout
func handleLogout(d *AuthDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawToken := extractRawToken(c)
		if rawToken == "" {
			return ok(c, fiber.Map{"logged_out": true})
		}
		if err := d.Auth.Logout(c.UserContext(), domain.LogoutParams{
			RawToken: rawToken,
		}); err != nil {
			d.Log.Warn("logout error", slog.Any("err", err))
		}
		clearSessionCookie(c)
		return ok(c, fiber.Map{"logged_out": true})
	}
}

// GET /auth/me — returns the current session's identity.
func handleMe(d *AuthDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawToken := extractRawToken(c)
		if rawToken == "" {
			return errUnauthorized(c, "no session")
		}
		resolved, err := d.Auth.ValidateSession(c.UserContext(), domain.ValidateSessionParams{
			RawToken: rawToken,
		})
		if err != nil {
			return errUnauthorized(c, "invalid session")
		}
		return ok(c, fiber.Map{
			"user_id":      resolved.UserID,
			"tenant_id":    resolved.TenantID,
			"display_name": resolved.DisplayName,
			"user_type":    resolved.UserType,
			"expires_at":   resolved.ExpiresAt,
			"flags":        resolved.Configuration.Flags,
		})
	}
}

// POST /auth/register — public endpoint; creates a new tenant + admin user + session.
// Body: {workspace_name, admin_name, admin_email, password, currency?, timezone?}
func handleRegister(d *AuthDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if d.Lifecycle == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": 503, "msg": "registration unavailable (DB not configured)",
			})
		}

		var body struct {
			WorkspaceName string `json:"workspace_name"`
			AdminName     string `json:"admin_name"`
			AdminEmail    string `json:"admin_email"`
			Password      string `json:"password"`
			Currency      string `json:"currency"`
			Timezone      string `json:"timezone"`
		}
		if err := c.BodyParser(&body); err != nil {
			return errBadRequest(c, "invalid request body")
		}
		body.AdminEmail = strings.TrimSpace(strings.ToLower(body.AdminEmail))
		body.WorkspaceName = strings.TrimSpace(body.WorkspaceName)

		switch {
		case body.WorkspaceName == "":
			return errBadRequest(c, "workspace_name is required")
		case body.AdminEmail == "":
			return errBadRequest(c, "admin_email is required")
		case len(body.Password) < 8:
			return errBadRequest(c, "password must be at least 8 characters")
		}
		if body.AdminName == "" {
			body.AdminName = body.AdminEmail
		}

		hash, err := iamauth.HashPassword(body.Password)
		if err != nil {
			return errInternal(c)
		}

		t, err := d.Lifecycle.Provision(c.UserContext(), tenant.ProvisionParams{
			Name:         body.WorkspaceName,
			Email:        body.AdminEmail,
			AdminName:    body.AdminName,
			AdminEmail:   body.AdminEmail,
			PasswordHash: hash,
			Plan:         "Basic",
			CurrencyCode: body.Currency,
			Timezone:     body.Timezone,
		})
		if err != nil {
			if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"status": 409, "msg": "workspace name or email already taken",
				})
			}
			d.Log.Error("register: provision failed", slog.Any("err", err))
			return errInternal(c)
		}

		// Auto-login: issue a session for the new admin.
		resolved, rawToken, err := d.Auth.Login(c.UserContext(), domain.LoginParams{
			Email:     body.AdminEmail,
			Password:  body.Password,
			IP:        c.IP(),
			UserAgent: c.Get("User-Agent"),
		})
		if err != nil {
			// Tenant created but session failed — not critical, user can log in manually.
			d.Log.Warn("register: auto-login failed", slog.Any("err", err))
			return ok(c, fiber.Map{
				"tenant_slug": t.Slug,
				"admin_email": body.AdminEmail,
				"auto_login":  false,
			})
		}

		setSessionCookie(c, rawToken, int(resolved.ExpiresAt.Unix()))
		return okCreated(c, fiber.Map{
			"tenant_slug":  t.Slug,
			"user_id":      resolved.UserID,
			"display_name": resolved.DisplayName,
			"expires_at":   resolved.ExpiresAt,
			"auto_login":   true,
		})
	}
}

// --- cookie helpers ---

func setSessionCookie(c *fiber.Ctx, token string, maxAge int) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookie,
		Value:    token,
		MaxAge:   maxAge,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})
}

func clearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookie,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})
}

func extractRawToken(c *fiber.Ctx) string {
	if v := c.Cookies(sessionCookie); v != "" {
		return v
	}
	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

func errUnauthorized(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
}
