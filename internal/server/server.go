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
	WebDir       string // path to built web assets (dist); empty = no static serving
	WebPublicDir string // path to Vite public/ dir; merged alongside WebDir (dev mode)
}

// DefaultConfig returns a sensible default server config.
func DefaultConfig() Config {
	return Config{
		Port:         8080,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		AllowOrigins: "*",
		WebDir:       "./web/dist",
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

	// Tenant + Auth middleware only applies to API surface.
	// /auth/register is public (no tenant exists yet during signup).
	tenantMW := middleware.Tenant(tenantResolver)
	authMW := middleware.Auth(sessions)
	exempt := func(c *fiber.Ctx) bool { return c.Path() == "/auth/register" }
	for _, prefix := range []string{"/api", "/iam", "/auth"} {
		p := prefix
		app.Use(p, func(c *fiber.Ctx) error {
			if exempt(c) {
				return c.Next()
			}
			return tenantMW(c)
		})
		app.Use(p, func(c *fiber.Ctx) error {
			if exempt(c) {
				return c.Next()
			}
			return authMW(c)
		})
	}

	RegisterRoutes(app, deps)

	if cfg.WebDir != "" {
		MountStatic(app, cfg.WebDir, cfg.WebPublicDir)
	}

	return app
}

// Addr returns the listen address string for the given config.
func Addr(cfg Config) string {
	return fmt.Sprintf(":%d", cfg.Port)
}
