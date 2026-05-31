package private

import (
	"go-gaurd/database"
	"go-gaurd/feature/profile/usecase"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProfileController struct {
	ProfileUsecase *usecase.ProfileUseCase
	validate       *validator.Validate
	RedisCache     *database.RedisCache
}

type ProfileControllerInterface interface {
	GetAccount(c *fiber.Ctx) error
	UpdateAccount(c *fiber.Ctx) error
	ActiveAccount(c *fiber.Ctx) error
}

func NewProfileController(profileUsecase *usecase.ProfileUseCase, redisCache *database.RedisCache) *ProfileController {
	log.Println("Initializing new ProfileController")
	return &ProfileController{
		ProfileUsecase: profileUsecase,
		validate:       validator.New(),
		RedisCache:     redisCache,
	}
}

// ActiveAccount implements [ProfileControllerInterface].
func (p *ProfileController) ActiveAccount(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetAccount implements [ProfileControllerInterface].
func (p *ProfileController) GetAccount(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateAccount implements [ProfileControllerInterface].
func (p *ProfileController) UpdateAccount(c *fiber.Ctx) error {
	panic("unimplemented")
}
