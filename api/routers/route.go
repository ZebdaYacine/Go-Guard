package routers

import (
	"go-gaurd/core/di"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	authController, err := di.InitializeAuthApplication()
	if err != nil {
		log.Printf("Error initializing AuthApplication: %s", err)
	}

	SetupPublicRoutes(app, authController)
	SetupPrivateRoutes(app, authController)
	SetupAdminRouter(app, authController)
}
