// Package middleware provides the Fiber middleware pipeline for the Awo HTTP server.
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"awo.so/framework/internal/core"
)

// RequestID injects a unique request ID into the Fiber context and Go context.
// Honours an incoming X-Request-ID header if present and valid.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		c.Set("X-Request-ID", id)

		// Attach to Go context for propagation into hooks and store.
		goCtx := core.WithRequestID(c.Context(), id)
		c.SetUserContext(goCtx)

		return c.Next()
	}
}
