package server

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"awo.so/framework/internal/middleware"
)

// Config holds server configuration.
type Config struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	AllowOrigins string // comma-separated list for CORS
}

// DefaultConfig returns a sensible default server config.
func DefaultConfig() Config {
	return Config{
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		AllowOrigins: "*",
	}
}

// New constructs and configures a Fiber app with all framework middleware applied.
// The caller should call RegisterRoutes(app, deps) after this, then app.Listen().
func New(cfg Config, deps *Deps, tenantResolver middleware.TenantResolver, sessions middleware.SessionStore, log *slog.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:           cfg.ReadTimeout,
		WriteTimeout:          cfg.WriteTimeout,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Middleware order matters — see §13.1 of docs.
	app.Use(middleware.Recovery(log))
	app.Use(middleware.RequestID())
	app.Use(middleware.Logger(log))
	app.Use(cors.New(cors.Config{AllowOrigins: cfg.AllowOrigins}))
	app.Use(middleware.Tenant(tenantResolver))
	app.Use(middleware.Auth(sessions))

	RegisterRoutes(app, deps)

	return app
}

// Addr returns the listen address string for the given config.
func Addr(cfg Config) string {
	return fmt.Sprintf(":%d", cfg.Port)
}
