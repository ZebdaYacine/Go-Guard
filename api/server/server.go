package server

import (
	"go-gaurd/api/controller"
	"go-gaurd/api/routers"
	"go-gaurd/api/security"
	"log"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// main.go

func InitFibreServer() {

	// Create Fiber app with custom config
	app := fiber.New(fiber.Config{
		ErrorHandler:            controller.CustomErrorHandler,
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            10 * time.Second,
		IdleTimeout:             120 * time.Second,
		Prefork:                 false,
		StrictRouting:           false,
		CaseSensitive:           true,
		Immutable:               false,
		UnescapePath:            false,
		ETag:                    false,
		BodyLimit:               4 * 1024 * 1024, // 4MB
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "192.168.1.0/24"},
	})

	// Global middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${method} ${path} | ${latency} | ${ip} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(security.SecurityHeaders())

	// Rate limiting middleware
	app.Use(limiter.New(limiter.Config{
		Max:               40,                      // Maximum number of requests
		Expiration:        1 * time.Minute,         // Per minute
		LimiterMiddleware: limiter.SlidingWindow{}, // Sliding window algorithm
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Rate limit by IP
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":    "Too many requests",
				"message":  "Please try again later",
				"retry_in": 40,
			})
		},
		SkipSuccessfulRequests: false,
	}))

	// Setup routes with groups
	routers.SetupRoutes(app)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped gracefully")
}
