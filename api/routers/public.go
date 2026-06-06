package routers

import (
	"go-gaurd/api/controller/public"
	"go-gaurd/api/security"
	"go-gaurd/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(
	app *fiber.App,
	authController public.AuthControllerInterface,
	redisCache *database.RedisCache,
) {

	detectIP := security.DetectClientIP(redisCache)

	// Health
	app.Get("/health", public.HealthCheck)
	app.Get("/ready", public.ReadinessCheck)

	// 20 requests/minute
	publicGroup := app.Group("/api/public")
	publicGroup.Use(detectIP)
	publicGroup.Use(security.RateLimitPerGuest(redisCache, 20, time.Minute))

	publicGroup.Get("/info", public.PublicInfo)

	// 10 requests/minute
	authGroup := app.Group("/api/auth")
	authGroup.Use(detectIP)
	authGroup.Use(security.RateLimitPerGuest(redisCache, 10, time.Minute))

	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

	// 3 requests/minute (sensitive operations)
	recoveryGroup := app.Group("/api/auth/recovery")
	recoveryGroup.Use(detectIP)
	recoveryGroup.Use(security.RateLimitPerGuest(redisCache, 3, time.Minute))

	recoveryGroup.Post("/forget-password", authController.ForgetPassword)
	recoveryGroup.Post("/check-otp", authController.CheckOTP)
	recoveryGroup.Put("/reset-password", authController.RestPassword)
}
