package routers

import (
	"fmt"
	"go-gaurd/api/controller/public"
	"go-gaurd/api/security"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRouter(app *fiber.App, authController *public.AuthController) {

	admin := app.Group("/api/admin")
	admin.Use(security.DetectClientIP(authController.RedisCache))
	app.Use(security.RateLimitPerUser(authController.RedisCache, 10, 1*time.Minute))
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		fmt.Println("YOU ARE INSIDE DASHBOARD ENDPOINT")

		//TODO ADD USECASE ACCESS FOR GET DASHBORD

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ready",
			"message": "Dashboard endpoint",
		})
	})

}
