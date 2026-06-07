package private

import (
	"fmt"
	"go-gaurd/api/security"
	"go-gaurd/core/utils"
	"go-gaurd/database"
	"go-gaurd/feature/profile/usecase"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProfileController struct {
	ProfileUsecase usecase.ProfileUseCaseInterface
	validate       *validator.Validate
	RedisCache     *database.RedisCache
	MinioDB        *database.MinioClient
}

type ProfileControllerInterface interface {
	GetProfile(c *fiber.Ctx) error
	ActiveProfile(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	LogoutAllDevices(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	UpdateProfilePicture(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
}

func NewProfileController(profileUsecase usecase.ProfileUseCaseInterface, redisCache *database.RedisCache, minioDB *database.MinioClient) ProfileControllerInterface {
	log.Println("Initializing new ProfileController")
	return &ProfileController{
		ProfileUsecase: profileUsecase,
		validate:       validator.New(),
		RedisCache:     redisCache,
		MinioDB:        minioDB,
	}
}

// GetProfile implements [ProfileControllerInterface].
func (p *ProfileController) GetProfile(c *fiber.Ctx) error {
	userId, _, jti, _, err := p.GetClaimFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	access_blacklist := "blacklist:access:" + jti

	status, err := p.RedisCache.Cache.Get(c.Context(), access_blacklist).Result()
	fmt.Println("RefrechJti ==>", access_blacklist)

	if status == "revoked" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Access token has been revoked",
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
	userId, _, _, _, err := p.GetClaimFromToken(c)
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
func (p *ProfileController) LogoutAllDevices(c *fiber.Ctx) error {
	log.Println("========== LOGOUT ALL DEVICES ==========")

	var logoutRequest LogoutRequest
	p.ValidateBody(c, &logoutRequest)

	userID, _, _, _, err := p.GetClaimFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	ctx := c.Context()

	pattern := fmt.Sprintf("refresh-token:%s:*", userID)

	iter := p.RedisCache.Cache.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()

		parts := strings.Split(key, ":")
		if len(parts) < 3 {
			log.Printf("Invalid key format: %s", key)
			continue
		}

		if logoutRequest.DeviceId != parts[2] {
			_ = p.RedisCache.Cache.Set(
				ctx,
				"blacklist:access:"+userID,
				"revoked",
				utils.AccessTokenExpiry,
			)
			err := p.RedisCache.Cache.Del(ctx, key).Err()
			if err != nil {
				log.Printf("Failed to delete %s: %v", key, err)
			}
		}
	}

	if err := iter.Err(); err != nil {
		log.Println("SCAN error:", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All devices logged out successfully",
		"success": true,
		"userId":  userID,
	})
}

// Logout implements [ProfileControllerInterface].
func (p *ProfileController) Logout(c *fiber.Ctx) error {

	userId, _, jti, expiresAt, err := p.GetClaimFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	var logoutRequest LogoutRequest
	p.ValidateBody(c, &logoutRequest)

	ctx := c.Context()

	authHeader := c.Get("Authorization")

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	refreshKey := fmt.Sprintf("refresh-token:%s:%s", userId, logoutRequest.DeviceId)
	refreshTokenID, err := p.RedisCache.Cache.Get(ctx, refreshKey).Result()
	fmt.Println("RefreshJti ==>", refreshTokenID)

	p.RedisCache.Cache.Del(ctx, refreshKey)

	ttl := time.Until(expiresAt.Time)

	if ttl < 0 {
		ttl = 0
	}

	if accessToken != "" {
		_ = p.RedisCache.Cache.Set(
			ctx,
			"blacklist:access:"+jti,
			"revoked",
			ttl,
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
	err := ac.ValidateBody(c, &req)
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

	ac.ValidateBody(c, &req)
	r := claims

	fmt.Println(r)

	refreshKey := fmt.Sprintf("refresh-token:%s:%s", claims.UserID, req.DeviceId)
	refreshTokenID, err := ac.RedisCache.Cache.Get(c.Context(), refreshKey).Result()
	fmt.Println("RefreshTokenID ==> claims.JTI", refreshTokenID, claims.JTI)

	if refreshTokenID == r.ID {
		fmt.Println("RefreshJti ==>", refreshTokenID)
		accessJTI := uuid.NewString()
		// Step 3: Generate new access token
		log.Println("Step 3: Generating new access token")
		newAccessToken, err := security.GenerateAccessToken(claims.UserID, claims.Role, accessJTI)
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

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid refresh token",
		"success": false,
	})

}

// UpdateProfile implements [ProfileControllerInterface].
func (p *ProfileController) UpdateProfile(c *fiber.Ctx) error {

	var profileUpdateRequest UpdateProfileRequest
	p.ValidateBody(c, &profileUpdateRequest)

	query := usecase.Query{}
	ctx := c.Context()

	result := p.ProfileUsecase.UpdateProfile(ctx, query)
	if !result.Success {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Failed to update profile",
			"success": false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"success": true,
	})
}

// UpdateProfilePicture implements [ProfileControllerInterface].
func (p *ProfileController) UpdateProfilePicture(c *fiber.Ctx) error {
	ctx := c.Context()
	_, url := p.UploadFile(c)
	result := p.ProfileUsecase.UpdateProfilePicture(ctx, url, "")
	if !result.Success {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Failed to update profile",
			"success": false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Profile updated successfully %s", url),
		"success": true,
	})
}

// UpdatePassword implements [ProfileControllerInterface].
func (p *ProfileController) UpdatePassword(c *fiber.Ctx) error {
	panic("unimplemented")
}
