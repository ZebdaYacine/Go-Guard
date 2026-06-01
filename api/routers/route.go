package routers

import (
	"go-gaurd/api/security"
	"go-gaurd/core/di"
	"log"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func SetupRoutes(app *fiber.App) {

	redisCache, err := di.InitializeRedis()
	if err != nil {
		log.Fatal("Failed to initialize Redis:", err)
	}

	authController, err := di.InitializeAuthApplication(redisCache)
	if err != nil {
		log.Printf("Error initializing AuthApplication: %s", err)
	}

	profileController, err := di.InitializeProfileApplication(redisCache)
	if err != nil {
		log.Printf("Error initializing ProfileApplication: %s", err)
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

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${method} ${path} | ${latency} | ${ip} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use(requestid.New())
	app.Use(security.SecurityHeaders())
	app.Use(security.AuthenticatorMiddleware(authController.RedisCache))
	app.Use(security.AuthoriserMiddleware(e))

	SetupPublicRoutes(app, authController)
	SetupPrivateRoutes(app, profileController)
	SetupAdminRouter(app, authController, profileController)
}
