package security

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Custom logger middleware with more details
func CustomLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path} | ${ip} | ${user-agent}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     log.Writer(), // Write to your log output
	})
}

// Route logger middleware to log specific route access
func RouteLogger(routeName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Log before request
		log.Printf("[ROUTE] Entering %s - %s %s", routeName, c.Method(), c.Path())

		// Process request
		err := c.Next()

		// Log after request
		duration := time.Since(start)
		log.Printf("[ROUTE] Exiting %s - Status: %d - Duration: %v", routeName, c.Response().StatusCode(), duration)

		return err
	}
}
