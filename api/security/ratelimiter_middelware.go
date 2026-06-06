package security

import (
	"context"
	"go-gaurd/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// Guest-specific rate limiting
func RateLimitPerGuest(redisCache *database.RedisCache, max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,        // 20 requests
		Expiration: expiration, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return redisCache.Cache.Get(context.Background(), "IP").Val()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{

				"error": "Rate limit exceeded for this IP",
			})
		},
	})
}

// User-specific rate limiting
func RateLimitPerUser(redisCache *database.RedisCache, max int, expiration time.Duration) fiber.Handler {
	error_msg := "Rate limit exceeded for this user"
	return limiter.New(limiter.Config{
		Max:        max,        // 20 requests
		Expiration: expiration, // per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := redisCache.Cache.Get(context.Background(), "userID")
			if userID == nil {
				error_msg = "Rate limit exceeded for this IP"
				return redisCache.Cache.Get(context.Background(), "IP").Val()
			}
			return "user:" + userID.Val()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{

				"error": error_msg,
			})
		},
	})
}

func DetectClientIP(redisCache *database.RedisCache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		redisCache.Cache.Set(c.Context(), "IP", ip, 5*time.Minute)

		log.Printf(
			"IP=%s Method=%s Path=%s",
			ip,
			c.Method(),
			c.OriginalURL(),
		)

		return c.Next()
	}
}
