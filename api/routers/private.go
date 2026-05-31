package routers

import (
	"fmt"
	"go-gaurd/api/controller/public"
	"go-gaurd/api/security"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(app *fiber.App, authController *public.AuthController) {

	user := app.Group("/api/user")
	user.Use(security.DetectClientIP(authController.RedisCache))
	app.Use(security.RateLimitPerUser(authController.RedisCache, 10, 1*time.Minute))
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
