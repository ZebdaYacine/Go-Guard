package public

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

func ReadinessCheck(c *fiber.Ctx) error {
	// Check database connection, etc.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ready",
	})
}

func PublicInfo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"app_name":    "Fiber API Example",
		"version":     "1.0.0",
		"description": "A complete API example with Fiber",
	})
}
