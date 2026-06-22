package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	"awo.so/framework/internal/tenant"
	"awo.so/framework/pkg/permissions"
)

type principalContextKey struct{}

// SessionStore retrieves a Principal for a given tenant + session ID.
type SessionStore interface {
	GetSession(ctx context.Context, tenantID, sessionID string) (*permissions.Principal, error)
}

// Auth validates the session cookie or Bearer token and attaches the Principal
// to the Go context. Returns 401 if credentials are absent or invalid.
func Auth(sessions SessionStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := tenant.FromCtx(c.UserContext())
		if t == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		sessionID := extractSessionID(c)
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		principal, err := sessions.GetSession(c.UserContext(), t.ID, sessionID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		goCtx := context.WithValue(c.UserContext(), principalContextKey{}, principal)
		c.SetUserContext(goCtx)
		return c.Next()
	}
}

// PrincipalFromCtx extracts the authenticated Principal from the context.
// Returns nil if no principal is set (unauthenticated request).
func PrincipalFromCtx(ctx context.Context) *permissions.Principal {
	p, _ := ctx.Value(principalContextKey{}).(*permissions.Principal)
	return p
}

// extractSessionID looks for a session identifier in:
//  1. Cookie "awo_session"
//  2. Authorization: Bearer <token>
func extractSessionID(c *fiber.Ctx) string {
	if cookie := c.Cookies("awo_session"); cookie != "" {
		return cookie
	}
	auth := c.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}
