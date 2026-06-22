package server

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"awo.so/framework/internal/iam/apikey"
	"awo.so/framework/internal/iam/domain"
	"awo.so/framework/internal/middleware"
	"awo.so/framework/internal/tenant"
)

// APIKeyDeps holds services needed by API key handlers.
type APIKeyDeps struct {
	APIKeys *apikey.Service
	Log     *slog.Logger
}

// MountAPIKeyRoutes registers /api-keys/* on the given Fiber app.
func MountAPIKeyRoutes(app *fiber.App, d *APIKeyDeps) {
	g := app.Group("/api-keys")
	g.Get("/", handleListAPIKeys(d))
	g.Post("/", handleCreateAPIKey(d))
	g.Delete("/:id", handleRevokeAPIKey(d))
}

// GET /api-keys
func handleListAPIKeys(d *APIKeyDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		if t == nil {
			return errUnauthorized(c, "tenant context missing")
		}
		keys, err := d.APIKeys.List(c.UserContext(), t.ID)
		if err != nil {
			d.Log.Error("list api keys", slog.Any("err", err))
			return errInternal(c)
		}
		return ok(c, keys)
	}
}

// POST /api-keys  body: {name, scopes, expires_at}
func handleCreateAPIKey(d *APIKeyDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		p := middleware.PrincipalFromCtx(c.UserContext())
		if t == nil || p == nil {
			return errUnauthorized(c, "unauthorized")
		}

		var body struct {
			Name      string    `json:"name"`
			Scopes    []string  `json:"scopes"`
			ExpiresAt *time.Time `json:"expires_at"`
		}
		if err := c.BodyParser(&body); err != nil {
			return errBadRequest(c, "invalid request body")
		}
		if body.Name == "" {
			return errBadRequest(c, "name is required")
		}

		createdBy, err := uuid.Parse(p.UserID)
		if err != nil {
			return errBadRequest(c, "invalid principal")
		}

		key, rawToken, err := d.APIKeys.Create(c.UserContext(), domain.CreateAPIKeyParams{
			TenantID:  t.ID,
			Name:      body.Name,
			Scopes:    body.Scopes,
			ExpiresAt: body.ExpiresAt,
			CreatedBy: createdBy,
		})
		if err != nil {
			d.Log.Error("create api key", slog.Any("err", err))
			return errInternal(c)
		}

		// raw_token is returned ONCE — client must store it; never retrievable again
		return okCreated(c, fiber.Map{
			"id":         key.ID,
			"name":       key.Name,
			"scopes":     key.Scopes,
			"expires_at": key.ExpiresAt,
			"created_at": key.CreatedAt,
			"raw_token":  rawToken,
		})
	}
}

// DELETE /api-keys/:id
func handleRevokeAPIKey(d *APIKeyDeps) fiber.Handler {
	return func(c *fiber.Ctx) error {
		keyID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return errBadRequest(c, "invalid key id")
		}
		if err := d.APIKeys.Revoke(c.UserContext(), keyID); err != nil {
			d.Log.Error("revoke api key", slog.Any("err", err))
			return errInternal(c)
		}
		return ok(c, fiber.Map{"revoked": true})
	}
}
