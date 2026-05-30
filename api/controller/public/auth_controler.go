package public

import (
	"go-gaurd/feature/auth/usecase"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthUsecase *usecase.AuthUseCase
	validate    *validator.Validate
}

func NewAuthController(authUsecase *usecase.AuthUseCase) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
		validate:    validator.New(),
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": result.Message,
		"user":    result.User,
		"success": result.Success,
	})
}
