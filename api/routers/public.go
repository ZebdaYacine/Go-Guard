package routers

import (
	"go-gaurd/api/controller/public"
	"go-gaurd/api/security"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App, authController *public.AuthController) {
	app.Use(security.RateLimitPerUser(authController.RedisCache, 20, 1*time.Minute))
	app.Get("/health", public.HealthCheck)
	app.Get("/ready", public.ReadinessCheck)

	// API guest group (fixed typo: "geust" -> "guest")
	guest := app.Group("/api/guest") // Fixed spelling
	guest.Use(security.RateLimitPerUser(authController.RedisCache, 20, 1*time.Minute))
	guest.Get("/public/info", public.PublicInfo)

	// API auth group
	auth := app.Group("/api/auth")
	auth.Use(security.RateLimitPerUser(authController.RedisCache, 10, 1*time.Minute))
	auth.Post("/register", authController.Register)

}
