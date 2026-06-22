// Package server contains the Fiber HTTP server, middleware wiring, and route handlers.
package server

import "github.com/gofiber/fiber/v2"

// PageMeta carries pagination metadata in list responses.
type PageMeta struct {
	Total  int64  `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// FieldError represents a single field-level validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ok sends a 200 JSON response with the standard success envelope.
func ok(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

// okCreated sends a 201 JSON response.
func okCreated(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": data})
}

// okList sends a paginated list response.
func okList(c *fiber.Ctx, data any, meta PageMeta) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data, "meta": meta})
}

// errBadRequest sends a 400 with a top-level error message.
func errBadRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
}

// errUnprocessable sends a 422 with field-level validation errors.
func errUnprocessable(c *fiber.Ctx, fields []FieldError) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": fields})
}

// errNotFound sends a 404.
func errNotFound(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": msg})
}

// errForbidden sends a 403.
func errForbidden(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": msg})
}

// errInternal sends a 500. Never include internal error details.
func errInternal(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}
