package routers

import (
	"go-gaurd/api/controller/public"
	"go-gaurd/core/di"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// Health check endpoint
	app.Get("/health", public.HealthCheck)
	app.Get("/ready", public.ReadinessCheck)

	// API v1 group
	v1 := app.Group("/api/v1")
	v1.Get("/public/info", public.PublicInfo)

	// API auth group
	auth := app.Group("/api/auth")
	AuthController, err := di.InitializeAuthApplication()
	if err != nil {
		log.Println("ERROR %s", err)
	}
	auth.Post("/register", AuthController.Register)

}
