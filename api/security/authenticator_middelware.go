package security

import (
	"fmt"
	"go-gaurd/database"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

// Admin authentication middleware
// AuthenticatorMiddleware is a Fiber middleware for JWT authentication
func AuthenticatorMiddleware(e *casbin.Enforcer, redisCache *database.RedisCache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		obj := string(c.Request().URI().Path())
		act := c.Method()

		log.Printf("Processing request: path=%s, method=%s", obj, act)

		// Check if guest (public) access is allowed based on policy
		// Note: "guest" role in your policy = public endpoints
		guestAllowed, err := e.Enforce("guest", obj, act)
		if err != nil {
			log.Printf("Error checking guest access: %v", err)
			// Continue to auth check instead of failing
		}

		if guestAllowed {
			log.Printf("✅ Guest access granted for %s %s", act, obj)
			return c.Next()
		}

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

		// Optional: Store in Redis if needed for session management
		ctx := c.Context()
		sessionKey := fmt.Sprintf("user_session:%s", userID)

		// Store session info in Redis with expiration
		err = redisCache.Cache.Set(ctx, sessionKey, map[string]interface{}{
			"userID": userID,
			"role":   role,
			"token":  token,
		}, 5*time.Minute).Err()

		return c.Next()
	}
}
