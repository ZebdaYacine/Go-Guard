package server

import (
	"go-gaurd/api/controller"
	"go-gaurd/api/routers"
	"log"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
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

	// app.Use(recover.New())

	routers.SetupRoutes(app)

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
