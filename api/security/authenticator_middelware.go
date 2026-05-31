package security

import (
	"fmt"
	"go-gaurd/api"
	"go-gaurd/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Admin authentication middleware
// AuthenticatorMiddleware is a Fiber middleware for JWT authentication
func AuthenticatorMiddleware(redisCache *database.RedisCache) fiber.Handler {
	return func(c *fiber.Ctx) error {

		obj := string(c.Request().URI().Path())
		act := c.Method()

		if api.IsPublicEndpoint(obj, act) {
			return c.Next()
		} else {

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
			rdb.Set(c.Context(), "userID", userID, 5*time.Minute)
			rdb.Set(c.Context(), "role", role, 5*time.Minute)
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
					keyType, _ := rdb.Type(c.Context(), key).Result()
					fmt.Printf("  - %s (type: %s)\n", key, keyType)
				}
			}

			return c.Next()
		}
	}
}
