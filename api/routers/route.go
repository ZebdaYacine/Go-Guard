package routers

import (
	"fmt"
	"go-gaurd/api/controller/public"
	"go-gaurd/api/security"
	"go-gaurd/core/di"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// Add logger middleware globally (logs all requests)
	// app.Use(security.CustomLogger()) // Add logger middleware globally

	// Health check endpoints (no logging or minimal logging)
	app.Get("/health", public.HealthCheck)
	app.Get("/ready", public.ReadinessCheck)

	// API guest group (fixed typo: "geust" -> "guest")
	guest := app.Group("/api/guest") // Fixed spelling
	guest.Get("/public/info", public.PublicInfo)

	// API auth group
	auth := app.Group("/api/auth")
	AuthController, err := di.InitializeAuthApplication()
	if err != nil {
		log.Printf("Error initializing AuthApplication: %s", err)
		// Don't continue if auth controller is required
		// You might want to panic or handle differently
	} else {
		auth.Post("/register", AuthController.Register)
	}

	// User group with authentication
	user := app.Group("/api/user")
	user.Use(security.AuthMiddleware(AuthController.RedisCache))
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
