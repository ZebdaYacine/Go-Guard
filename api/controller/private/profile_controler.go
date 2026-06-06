package private

import (
	"go-gaurd/database"
	"go-gaurd/feature/profile/usecase"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProfileController struct {
	ProfileUsecase usecase.ProfileUseCaseInterface
	validate       *validator.Validate
	RedisCache     *database.RedisCache
}

type ProfileControllerInterface interface {
	GetProfile(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	UpdateProfilePicture(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
	ActiveProfile(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

func NewProfileController(profileUsecase usecase.ProfileUseCaseInterface, redisCache *database.RedisCache) ProfileControllerInterface {
	log.Println("Initializing new ProfileController")
	return &ProfileController{
		ProfileUsecase: profileUsecase,
		validate:       validator.New(),
		RedisCache:     redisCache,
	}
}

// ActiveProfile implements [ProfileControllerInterface].
func (p *ProfileController) ActiveProfile(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetProfile implements [ProfileControllerInterface].
func (p *ProfileController) GetProfile(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Logout implements [ProfileControllerInterface].
func (p *ProfileController) Logout(c *fiber.Ctx) error {
	panic("unimplemented")
}

// RefreshAccessToken implements [ProfileControllerInterface].
func (p *ProfileController) RefreshAccessToken(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateProfile implements [ProfileControllerInterface].
func (p *ProfileController) UpdateProfile(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdateProfilePicture implements [ProfileControllerInterface].
func (p *ProfileController) UpdateProfilePicture(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UpdatePassword implements [ProfileControllerInterface].
func (p *ProfileController) UpdatePassword(c *fiber.Ctx) error {
	panic("unimplemented")
}
