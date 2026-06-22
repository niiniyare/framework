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

// ok sends a 200 JSON response in AMIS envelope: {status:0, data:...}
func ok(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 0,
		"msg":    "",
		"data":   data,
	})
}

// okCreated sends a 201 JSON response in AMIS envelope.
func okCreated(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": 0,
		"msg":    "",
		"data":   data,
	})
}

// okList sends a paginated list response.
// AMIS crud expects {status:0, data:{items:[...], total:N}}.
func okList(c *fiber.Ctx, data any, meta PageMeta) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": 0,
		"msg":    "",
		"data": fiber.Map{
			"items": data,
			"total": meta.Total,
		},
	})
}

// errBadRequest sends a 400 with AMIS error envelope.
func errBadRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status": 400,
		"msg":    msg,
	})
}

// errUnprocessable sends a 422 with field-level validation errors.
func errUnprocessable(c *fiber.Ctx, fields []FieldError) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"status": 422,
		"msg":    "Validation failed",
		"errors": fields,
	})
}

// errNotFound sends a 404.
func errNotFound(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status": 404,
		"msg":    msg,
	})
}

// errForbidden sends a 403.
func errForbidden(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"status": 403,
		"msg":    msg,
	})
}

// errInternal sends a 500. Never include internal error details.
func errInternal(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status": 500,
		"msg":    "internal server error",
	})
}
