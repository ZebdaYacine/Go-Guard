package controller

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// Custom error handler
func CustomErrorHandler(c *fiber.Ctx, err error) error {

	// Default error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Handle specific errors
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Handle validation errors
	if _, ok := err.(validator.ValidationErrors); ok {
		code = fiber.StatusBadRequest
		message = "Validation failed"
	}

	// Log error (in production, use proper logging)
	fmt.Printf("Error: %v, Path: %s, Method: %s\n", err, c.Path(), c.Method())

	// Return JSON error response
	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})

}
