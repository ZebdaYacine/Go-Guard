package routers

import (
	"go-gaurd/api/controller/private"
	"go-gaurd/api/security"
	"go-gaurd/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(app *fiber.App, profileController private.ProfileControllerInterface, redis *database.RedisCache) {
	//TODO TAKE 2 DAYS FOR THIS WORKFLOW
	user := app.Group("/api/user")
	user.Use(security.DetectClientIP(redis))
	profileLimiter := security.RateLimitPerUser(redis, 10, time.Minute)
	updateLimiter := security.RateLimitPerUser(redis, 5, time.Minute)
	authLimiter := security.RateLimitPerUser(redis, 1, time.Minute)

	user.Get("/profile", profileLimiter, profileController.GetProfile)
	user.Put("/update-profile", updateLimiter, profileController.UpdateProfile)
	user.Post("/refresh-token", authLimiter, profileController.RefreshAccessToken)
}
