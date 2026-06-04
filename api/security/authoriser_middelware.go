package security

import (
	"go-gaurd/api"
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

		if api.IsPublicEndpoint(obj, act) {
			return c.Next()
		} else {

			authHeader := c.Get("Authorization")
			token, err := ExtractTokenFromHeader(authHeader)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid or missing token",
				})
			}
			role, err := ExtractRole(token)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid or missing token",
				})
			}

			allowed, err := e.Enforce(role, obj, act)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			if !allowed {
				return c.Status(http.StatusForbidden).JSON(fiber.Map{
					"error": "access denied",
				})
			}

			return c.Next()
		}
	}
}
