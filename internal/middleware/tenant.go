package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"awo.so/framework/internal/tenant"
)

// TenantResolver is a function that looks up a Tenant by its slug.
// It is implemented in internal/tenant/resolver.go but passed here as an interface
// to keep the middleware package free of database dependencies.
type TenantResolver interface {
	Resolve(slug string) (*tenant.Tenant, error)
}

// Tenant extracts the tenant slug from the request and attaches the Tenant to
// the Go context. Resolution order:
//  1. X-Awo-Tenant header (for API clients and mobile)
//  2. Subdomain of the Host header (e.g. "acme" from "acme.awo.app")
//
// Returns 404 Not Found if the tenant cannot be resolved (not 403 — do not
// confirm existence to unauthenticated callers).
func Tenant(resolver TenantResolver) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Get("X-Awo-Tenant")
		if slug == "" {
			slug = extractSubdomain(c.Hostname())
		}
		if slug == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tenant not found"})
		}

		t, err := resolver.Resolve(slug)
		if err != nil {
			// Intentionally 404 regardless of error type to avoid tenant enumeration.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tenant not found"})
		}

		if !t.IsActive() {
			if t.Status == tenant.StatusSuspended {
				return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
					"error":  "tenant suspended",
					"status": t.Status,
				})
			}
			return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "tenant decommissioned"})
		}

		goCtx := tenant.WithTenant(c.UserContext(), t)
		c.SetUserContext(goCtx)
		return c.Next()
	}
}

// extractSubdomain returns the leftmost label of a hostname.
// Returns "" if the hostname has fewer than 2 labels (e.g. "localhost").
func extractSubdomain(host string) string {
	// Strip port if present
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}
	parts := strings.SplitN(host, ".", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}
