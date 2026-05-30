package public

import (
	"go-gaurd/api/security"
	"go-gaurd/core/utils"
	"go-gaurd/database"
	"go-gaurd/feature/auth/usecase"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthUsecase *usecase.AuthUseCase
	validate    *validator.Validate
	RedisCache  *database.RedisCache
}

func NewAuthController(authUsecase *usecase.AuthUseCase, redisCache *database.RedisCache) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
		validate:    validator.New(),
		RedisCache:  redisCache,
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	log.Println("Register endpoint called")

	var query usecase.Query

	ctx := c.Context()
	result := ac.AuthUsecase.CreateAccount(ctx, query)

	if !result.Success {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": result.Message,
			"success": result.Success,
			"user":    result.User,
		})
	}

	token, err := security.GenerateAccessToken("23434", utils.RoleClient)
	if err != nil {
		log.Fatalf(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating access token",
			"success": false,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": result.Message,
		"user":    result.User,
		"success": result.Success,
		"token":   token,
	})
}
