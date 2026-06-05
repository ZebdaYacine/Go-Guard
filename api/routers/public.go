package routers

import (
	"go-gaurd/api/controller/public"
	"go-gaurd/api/security"
	"go-gaurd/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App, authController public.AuthControllerInterface, redisCache *database.RedisCache) {

	app.Use(security.DetectClientIP(redisCache))
	app.Use(security.RateLimitPerGuest(redisCache, 20, 1*time.Minute))
	app.Get("/health", public.HealthCheck)
	app.Get("/ready", public.ReadinessCheck)

	// API guest group (fixed typo: "geust" -> "guest")
	guest := app.Group("/api/guest") // Fixed spelling
	guest.Use(security.DetectClientIP(redisCache))
	guest.Use(security.RateLimitPerGuest(redisCache, 20, 1*time.Minute))
	guest.Get("/public/info", public.PublicInfo)

	// API auth group
	auth := app.Group("/api/auth")
	auth.Use(security.DetectClientIP(redisCache))
	auth.Use(security.RateLimitPerGuest(redisCache, 10, 1*time.Minute))
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/forget-password", authController.ForgetPassword)
	auth.Get("/check-otp", authController.CheckOTP)
	auth.Put("/rest-password", authController.RestPassword)

}
