package security

import (
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

func AuthoriserMiddleware(e *casbin.Enforcer) fiber.Handler {
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

		// Require authentication for non-guest endpoints
		authHeader := c.Get("Authorization")
		token, err := ExtractTokenFromHeader(authHeader)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid path or missing token",
			})
		}

		role, err := ExtractRole(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid path or missing token",
			})
		}

		// Check role-based access
		allowed, err := e.Enforce(role, obj, act)
		if err != nil {
			log.Printf("Enforcement error for role %s: %v", role, err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "authorization error",
			})
		}

		if !allowed {
			log.Printf("❌ Access denied for role=%s, path=%s, method=%s", role, obj, act)
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "access denied",
			})
		}

		log.Printf("✅ Access granted for role=%s, path=%s, method=%s", role, obj, act)
		return c.Next()
	}
}
