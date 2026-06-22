package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"

	"awo.so/framework/internal/core"
)

// Recovery catches panics, logs the stack trace, and returns a 500 JSON response.
// Stack traces are never sent to the client.
func Recovery(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				requestID := core.RequestIDFromCtx(c.UserContext())
				log.Error("panic recovered",
					slog.String("request_id", requestID),
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)
				err = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "internal server error",
				})
			}
		}()
		return c.Next()
	}
}
