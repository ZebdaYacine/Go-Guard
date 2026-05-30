package security

import (
	"fmt"
	"go-gaurd/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Custom logger middleware with more details
func CustomLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path} | ${ip} | ${user-agent}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     log.Writer(), // Write to your log output
	})
}

// Route logger middleware to log specific route access
func RouteLogger(routeName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Log before request
		log.Printf("[ROUTE] Entering %s - %s %s", routeName, c.Method(), c.Path())

		// Process request
		err := c.Next()

		// Log after request
		duration := time.Since(start)
		log.Printf("[ROUTE] Exiting %s - Status: %d - Duration: %v", routeName, c.Response().StatusCode(), duration)

		return err
	}
}

// Security headers middleware
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set security headers
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		return c.Next()
	}
}

// Admin authentication middleware
// AuthMiddleware is a Fiber middleware for JWT authentication
func AuthMiddleware(redisCache *database.RedisCache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")

		// Extract token
		token, err := ExtractTokenFromHeader(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Missing or invalid authorization header",
			})
		}

		// Validate token
		valid, userID, role, err := ValidateAccessToken(token)
		if err != nil || !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid or expired token",
				"message": err.Error(),
			})
		}

		// Store user info in Fiber's context locals
		// TODO USE REDIS SERVER
		rdb := redisCache.Cache
		rdb.Set(c.Context(), "userID", userID, 0)
		rdb.Set(c.Context(), "role", role, 0)
		// List all keys
		keys, err := rdb.Keys(c.Context(), "*").Result()
		if err != nil {
			log.Fatal(err)
		}

		if len(keys) == 0 {
			fmt.Println("No keys found in Redis")
		} else {
			fmt.Printf("Found %d keys:\n", len(keys))
			for _, key := range keys {
				// Get type of each key
				keyType, _ := rdb.Type(c.Context(), key).Result()
				fmt.Printf("  - %s (type: %s)\n", key, keyType)
			}
		}
		// c.Locals("userID", userID)
		// c.Locals("role", role)

		// Continue to next handler
		return c.Next()
	}
}

// User-specific rate limiting
func RateLimitPerUser() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        20,              // 20 requests
		Expiration: 1 * time.Minute, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := c.Locals("user_id")
			if userID == nil {
				return c.IP()
			}
			return "user:" + userID.(string)
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded for this user",
			})
		},
	})
}

// Admin-specific rate limiting (stricter)
func RateLimitAdmin() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return "admin:" + c.IP()
		},
	})
}
