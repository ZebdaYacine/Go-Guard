package routers

import (
	"go-gaurd/api/security"
	"go-gaurd/core/di"
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	authController, err := di.InitializeAuthApplication()
	if err != nil {
		log.Printf("Error initializing AuthApplication: %s", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	e, err :=
		casbin.NewEnforcer(
			cwd+"/api/security/rbac/rbac_model.conf", cwd+"/api/security/rbac/policy.csv")
	if err != nil {
		log.Fatal("Failed to create casbin enforcer:", err)
	}

	app.Use(security.AuthenticatorMiddleware(authController.RedisCache))
	app.Use(security.AuthoriserMiddleware(e))
	SetupPublicRoutes(app, authController)
	SetupPrivateRoutes(app, authController)
	SetupAdminRouter(app, authController)
}
