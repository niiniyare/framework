package server

import "github.com/gofiber/fiber/v2"

// RegisterRoutes mounts entity CRUD routes on the Fiber app.
// All routes live under /api/v1/.
func RegisterRoutes(app *fiber.App, d *Deps) {
	api := app.Group("/api/v1")

	// Entity routes — :entity maps to EntityDefinition.Name (case-sensitive).
	api.Get("/:entity", HandleList(d))
	api.Post("/:entity", HandleCreate(d))
	api.Get("/:entity/:id", HandleGet(d))
	api.Put("/:entity/:id", HandleUpdate(d))
	api.Delete("/:entity/:id", HandleDelete(d))

	// Document lifecycle — only available for submittable entities.
	api.Post("/:entity/:id/submit", HandleSubmit(d))
	api.Post("/:entity/:id/cancel", HandleCancel(d))

	// Health endpoints — no tenant/auth required.
	app.Get("/health/live", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
