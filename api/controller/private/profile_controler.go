package private

import (
	"go-gaurd/api/security"
	"go-gaurd/core/utils"
	"go-gaurd/database"
	"go-gaurd/feature/profile/usecase"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func (ac *ProfileController) validateBody(c *fiber.Ctx, req interface{}) error {
	log.Println("Validating request body")

	if err := c.BodyParser(req); err != nil {
		log.Printf("Body parsing failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"success": false,
		})
	}

	if err := ac.validate.Struct(req); err != nil {
		log.Printf("Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"success": false,
		})
	}

	log.Println("Request body validation successful")
	return nil
}

func (ac *ProfileController) getClaimFromToken(c *fiber.Ctx) (string, string, *jwt.NumericDate, error) {
	authHeader := c.Get("Authorization")
	token, err := security.ExtractTokenFromHeader(authHeader)
	if err != nil {
		log.Printf("Failed to extract token: %v", err)
		return "", "", nil, err
	}
	userId, role, expiresAt, err := security.ValidateAccessToken(token)
	if err != nil {
		log.Printf("Failed to validate access token: %v", err)
		return "", "", nil, err
	}
	return userId, role, expiresAt, nil
}

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

// GetProfile implements [ProfileControllerInterface].
func (p *ProfileController) GetProfile(c *fiber.Ctx) error {
	userId, _, _, err := p.getClaimFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	result := p.ProfileUsecase.GetProfile(c.Context(), userId)
	if !result.Success {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": result.Message,
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile retrieved successfully",
		"success": true,
		"user":    result.User,
	})

}

// ActiveProfile implements [ProfileControllerInterface].
func (p *ProfileController) ActiveProfile(c *fiber.Ctx) error {
	userId, _, _, err := p.getClaimFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	result := p.ProfileUsecase.ActiveProfile(c.Context(), userId)
	if !result.Success {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": result.Message,
			"success": false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile retrieved successfully",
		"success": true,
		"user":    result.User,
	})
}

// Logout implements [ProfileControllerInterface].
func (p *ProfileController) Logout(c *fiber.Ctx) error {
	userId, _, expiresAt, err := p.getClaimFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	accessToken := c.Get("Authorization")

	refreshToken := c.Cookies("refresh_token") // if you use cookies

	// 2. Calculate TTL from JWT exp
	exp := expiresAt
	ttl := time.Until(exp.Time)

	// 3. Store access token in blacklist (cache)
	if accessToken != "" {
		_ = p.RedisCache.Cache.Set(
			c.Context(),
			"blacklist:access:"+accessToken,
			"revoked",
			ttl,
		).Err()
	}

	// 4. Store refresh token in blacklist too (longer TTL if needed)
	if refreshToken != "" {
		_ = p.RedisCache.Cache.Set(
			c.Context(),
			"blacklist:refresh:"+refreshToken,
			"revoked",
			utils.RefreshTokenExpiry,
		).Err()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
		"success": true,
		"userId":  userId,
	})
}

// RefreshAccessToken implements [ProfileControllerInterface].
func (ac *ProfileController) RefreshAccessToken(c *fiber.Ctx) error {
	log.Println("========== REFRESH ACCESS TOKEN ENDPOINT STARTED ==========")

	var req RefreshTokenRequest

	// Step 1: Validate request body
	log.Println("Step 1: Validating request body")
	err := ac.validateBody(c, &req)
	if err != nil {
		log.Printf("Validation failed, returning error response: %v", err)
		return err
	}

	// Step 2: Validate refresh token
	log.Println("Step 2: Validating refresh token")
	status, claims, err := security.ValidateRefreshToken(req.RefreshToken)
	if err != nil || !status {
		log.Printf("Invalid refresh token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired refresh token",
			"success": false,
		})
	}

	// Step 3: Generate new access token
	log.Println("Step 3: Generating new access token")
	newAccessToken, err := security.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		log.Printf("Failed to generate new access token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating new access token",
			"success": false,
		})
	}

	log.Printf("New access token generated successfully for user: %s", claims.UserID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Access token refreshed successfully",
		"success":      true,
		"access_token": newAccessToken,
	})
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
