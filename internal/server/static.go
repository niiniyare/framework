package server

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// MountStatic serves the web frontend from dir (built dist) and optionally
// a separate publicDir (Vite's public/ folder — used in dev when dist isn't built).
//
// Precedence: dir files first, then publicDir files, then SPA catch-all.
// Mount this AFTER all API routes so API paths are never swallowed.
func MountStatic(app *fiber.App, dir, publicDir string) {
	app.Static("/", dir, fiber.Static{Index: "index.html", Browse: false})

	if publicDir != "" {
		app.Static("/", publicDir, fiber.Static{Browse: false})
	}

	indexFile := filepath.Join(dir, "index.html")
	app.Get("*", func(c *fiber.Ctx) error {
		// Has file extension → real asset, not a client-side route.
		if filepath.Ext(c.Path()) != "" {
			return c.Status(fiber.StatusNotFound).SendString("not found")
		}
		if _, err := os.Stat(indexFile); err != nil {
			return c.Status(fiber.StatusNotFound).SendString(
				"web UI not built — run: cd web && npm run build")
		}
		return c.SendFile(indexFile)
	})
}
