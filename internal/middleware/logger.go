package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"awo.so/framework/internal/core"
	"awo.so/framework/internal/tenant"
)

// Logger returns a Fiber middleware that emits a structured slog record after
// each request. Fields: request_id, tenant_id, method, path, status, latency_ms.
func Logger(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		requestID := core.RequestIDFromCtx(c.UserContext())
		tenantID := ""
		if t := tenant.FromCtx(c.UserContext()); t != nil {
			tenantID = t.ID
		}

		attrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("tenant_id", tenantID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Int64("latency_ms", latency.Milliseconds()),
		}

		level := slog.LevelInfo
		if c.Response().StatusCode() >= 500 {
			level = slog.LevelError
		} else if c.Response().StatusCode() >= 400 {
			level = slog.LevelWarn
		}

		log.LogAttrs(c.Context(), level, "request", attrs...)
		return err
	}
}
