package routers

import (
	"fmt"
	"go-gaurd/api/controller/private"
	"go-gaurd/api/security"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(app *fiber.App, profileController *private.ProfileController) {

	user := app.Group("/api/user")
	user.Use(security.DetectClientIP(profileController.RedisCache))
	app.Use(security.RateLimitPerUser(profileController.RedisCache, 10, 1*time.Minute))
	user.Get("/profile", func(c *fiber.Ctx) error {
		fmt.Println("YOU ARE INSIDE PROFILE ENDPOINT")

		// Get user info from context (set by AuthMiddleware)
		userID := c.Locals("userID")
		role := c.Locals("role")

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ready",
			"message": "Profile endpoint",
			"user_id": userID,
			"role":    role,
		})
	})

}
