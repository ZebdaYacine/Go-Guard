package public

import (
	"fmt"
	"go-gaurd/api/security"
	"go-gaurd/core/utils"
	"go-gaurd/database"
	"go-gaurd/feature/auth/usecase"
	"log"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthUsecase usecase.AuthUseCaseInterface
	validate    *validator.Validate
	RedisCache  *database.RedisCache
}

type AuthControllerInterface interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	ForgetPassword(c *fiber.Ctx) error
	CheckOTP(c *fiber.Ctx) error
	RestPassword(c *fiber.Ctx) error
	//TODO DOING THIS WORKFLOW
	RefreshAccessToken(c *fiber.Ctx) error
}

func NewAuthController(authUsecase usecase.AuthUseCaseInterface, redisCache *database.RedisCache) AuthControllerInterface {
	log.Println("Initializing new AuthController")
	return &AuthController{
		AuthUsecase: authUsecase,
		validate:    validator.New(),
		RedisCache:  redisCache,
	}
}

func (ac *AuthController) validateBody(c *fiber.Ctx, req interface{}) error {
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

func (ac *AuthController) Register(c *fiber.Ctx) error {
	log.Println("========== REGISTER ENDPOINT STARTED ==========")

	var req RegisterRequest
	ctx := c.Context()

	// Step 1: Validate request body
	log.Println("Step 1: Validating request body")
	err := ac.validateBody(c, &req)
	if err != nil {
		log.Printf("Validation failed, returning error response: %v", err)
		return err
	}
	log.Printf("Request body validated successfully for email: %s", req.Email)

	// Step 2: Prepare query for usecase
	log.Println("Step 2: Preparing query for usecase")
	query := usecase.Query{
		User: usecase.User_Entity{
			User_name: req.User_name,
			Email:     req.Email,
			Phone:     req.Phone,
			Password:  req.Password,
			Role:      req.Role,
			Sex:       req.Sex,
			Picture:   req.Picture,
		},
	}
	log.Printf("Query prepared for user: %s (Role: %s)", req.User_name, req.Role)

	// Step 3: Create account via usecase
	log.Println("Step 3: Creating account via AuthUsecase")
	result := ac.AuthUsecase.CreateAccount(ctx, query)

	// Step 4: Handle failure case
	if !result.Success {
		log.Printf("Account creation failed: %s", result.Message)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": result.Message,
			"success": result.Success,
			"user":    result.User,
		})
	}

	log.Printf("Account created successfully for user: %s (Email: %s)", req.User_name, req.Email)

	response := RegisterResponse{
		Message: result.Message,
		Success: result.Success,
		User:    result.User,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	log.Println("========== LOGIN ENDPOINT STARTED ==========")

	var req LoginRequest
	ctx := c.Context()

	// Step 1: Validate request body
	log.Println("Step 1: Validating request body")
	err := ac.validateBody(c, &req)
	if err != nil {
		log.Printf("Validation failed, returning error response: %v", err)
		return err
	}
	log.Printf("Request body validated successfully for email: %s", req.Email)

	// Step 2: Prepare query for usecase
	log.Println("Step 2: Preparing query for usecase")
	query := usecase.Query{
		User: usecase.Login_Entity{
			Email:    req.Email,
			Password: req.Password,
		},
	}
	log.Printf("Query prepared for user: %s ", req.Email)

	// Step 3: Create account via usecase
	log.Println("Step 3: Logging in account via AuthUsecase")
	result := ac.AuthUsecase.Login(ctx, query)

	// Step 4: Handle failure case
	if !result.Success {
		log.Printf("Account login failed: %s", result.Message)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": result.Message,
			"success": result.Success,
			"user":    result.User,
		})
	}

	uuid := result.Id

	access_token, err := security.GenerateAccessToken(uuid, utils.RoleUser)
	if err != nil {
		log.Printf("ERROR: Token generation failed: %v", err)
		log.Fatalf("%s", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating access token",
			"success": false,
		})
	}
	log.Println("Access Token generated successfully")

	refresh_token, err := security.GenerateRefreshToken(uuid, utils.RoleUser)
	if err != nil {
		log.Printf("ERROR: Token generation failed: %v", err)
		log.Fatalf("%s", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error generating refresh token",
			"success": false,
		})
	}
	log.Println("Refresh Token generated successfully")

	log.Printf("Account logged in successfully for Email: %s", req.Email)
	response := LoginResponse{
		Message:      result.Message,
		Success:      result.Success,
		User:         result.User,
		AccessToken:  access_token,
		RefrechToken: refresh_token,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (ac *AuthController) ForgetPassword(c *fiber.Ctx) error {
	log.Println("========== FORGET PASSWORD ENDPOINT STARTED ==========")

	var req ForgetPasswordRequest
	ctx := c.Context()

	// Step 1: Validate request body
	log.Println("Step 1: Validating request body")
	err := ac.validateBody(c, &req)
	if err != nil {
		log.Printf("Validation failed, returning error response: %v", err)
		return err
	}
	log.Printf("Request body validated successfully for email: %s", req.Email)

	// Step 2: Check if user exists
	log.Println("Step 2: Checking if user exists")
	result := ac.AuthUsecase.CheckUserExists(ctx, req.Email)
	if !result.Success {
		log.Printf("User not found: %s", req.Email)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found with this email",
			"success": false,
		})
	}

	// Step 3: Generate OTP
	log.Println("Step 3: Generating OTP")
	OTP := utils.GenerateOTP()
	log.Printf("Generated OTP for %s: %s", req.Email, OTP)

	// Step 4: Store OTP in Redis cache with proper error handling
	log.Println("Step 4: Storing OTP in Redis cache")
	redisKey := fmt.Sprintf("forgetpassword:%s", req.Email)

	// Delete any existing OTP for this email first (optional but good practice)
	ac.RedisCache.Cache.Del(ctx, redisKey)

	// Store new OTP with 10 minutes expiration
	err = ac.RedisCache.Cache.Set(ctx, redisKey, OTP, 10*time.Minute).Err()
	if err != nil {
		log.Printf("Failed to store OTP in Redis: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to process request. Please try again.",
			"success": false,
		})
	}

	// Step 5: Send OTP via email
	log.Println("Step 5: Sending OTP via email")
	err = utils.SendOTP(req.Email, OTP)
	if err != nil {
		log.Printf("Failed to send OTP email: %v", err)
		// Rollback: Delete the OTP from Redis if email sending fails
		ac.RedisCache.Cache.Del(ctx, redisKey)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to send OTP. Please try again later.",
			"success": false,
		})
	}

	log.Printf("Forget password OTP sent successfully to: %s", req.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP sent successfully to your email",
		"success": true,
		"email":   req.Email,
	})
}

func (ac *AuthController) CheckOTP(c *fiber.Ctx) error {
	log.Println("========== CHECK OTP ENDPOINT STARTED ==========")

	var req CheckOTPRequest
	ctx := c.Context()

	// Step 1: Validate request body
	log.Println("Step 1: Validating request body")
	err := ac.validateBody(c, &req)
	if err != nil {
		log.Printf("Validation failed, returning error response: %v", err)
		return err
	}
	log.Printf("Request body validated successfully for email: %s, type: %s", req.Email, req.Type)

	// Step 2: Validate OTP type
	if req.Type != "verification" && req.Type != "forgetpassword" && req.Type != "login" {
		log.Printf("Invalid OTP type: %s", req.Type)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid OTP type. Valid types: verification, forgetpassword, login",
			"success": false,
		})
	}

	// Step 3: Get OTP from Redis cache
	log.Println("Step 3: Getting OTP from Redis cache")
	redisKey := fmt.Sprintf("%s:%s", req.Type, req.Email)
	log.Println(redisKey)
	storedOTP := ac.RedisCache.Cache.Get(ctx, redisKey)
	if storedOTP.Err() != nil {
		log.Printf("Failed to get OTP from Redis: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid or expired OTP. Please request a new one.",
			"success": false,
		})
	}

	// Step 4: Check if OTP matches
	if storedOTP.Val() != req.OTP {
		log.Printf("OTP mismatch for email: %s, type: %s", req.Email, req.Type)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid OTP",
			"success": false,
		})
	}

	// Step 5: For login type, generate tokens
	if req.Type == "login" {
		log.Println("Step 5: Generating tokens for login OTP verification")

		// Get user details from database
		query := usecase.Query{
			User: usecase.User_Entity{
				Email: req.Email,
			},
		}

		userID, role := 0, ""
		if userID == 0 {
			log.Printf("Failed to get user ID: %s", req.Email)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to retrieve user details",
				"success": false,
			})
		}

		accessToken, err := security.GenerateAccessToken(strconv.Itoa(userID), role)
		if err != nil {
			log.Printf("Failed to generate access token: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error generating access token",
				"success": false,
			})
		}

		refreshToken, err := security.GenerateRefreshToken(strconv.Itoa(userID), role)
		if err != nil {
			log.Printf("Failed to generate refresh token: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error generating refresh token",
				"success": false,
			})
		}

		// Delete OTP after successful verification
		ac.RedisCache.Cache.Del(ctx, redisKey)

		log.Printf("OTP verified successfully for email: %s, type: %s", req.Email, req.Type)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":       "OTP verified successfully",
			"success":       true,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user":          query.User,
		})
	}

	// Step 6: For non-login types, just verify OTP
	// Don't delete OTP for forgetpassword as it will be used in reset password
	if req.Type != "forgetpassword" {
		ac.RedisCache.Cache.Del(ctx, redisKey)
	}
	log.Printf("OTP verified successfully for email: %s, type: %s", req.Email, req.Type)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP verified successfully",
		"success": true,
	})
}

func (ac *AuthController) RestPassword(c *fiber.Ctx) error {
	log.Println("========== RESET PASSWORD ENDPOINT STARTED ==========")

	var req ResetPasswordRequest
	ctx := c.Context()

	// Step 1: Validate request body
	log.Println("Step 1: Validating request body")
	err := ac.validateBody(c, &req)
	if err != nil {
		log.Printf("Validation failed, returning error response: %v", err)
		return err
	}
	log.Printf("Request body validated successfully for email: %s", req.Email)

	// Step 2: Verify OTP from Redis cache
	log.Println("Step 2: Verifying OTP from Redis cache")
	redisKey := fmt.Sprintf("forgetpassword:%s", req.Email)
	storedOTP := ac.RedisCache.Cache.Get(ctx, redisKey)
	if storedOTP.Err() != nil {
		log.Printf("Failed to get OTP from Redis: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid or expired OTP. Please request a new one.",
			"success": false,
		})
	}

	// Step 3: Check if OTP matches
	if storedOTP.Val() != req.OTP {
		log.Printf("OTP mismatch for email: %s", req.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid OTP",
			"success": false,
		})
	}

	log.Println("OTP verified successfully")

	// Step 4: Validate new password and confirm password match
	if req.NewPassword != req.ConfirmPassword {
		log.Printf("Password mismatch for email: %s", req.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "New password and confirm password do not match",
			"success": false,
		})
	}

	// Step 5: Update password in database
	log.Println("Step 5: Updating password")
	query := usecase.Query{
		User: usecase.ResetPassword_Entity{
			Email:            req.Email,
			NewPassword:      req.NewPassword,
			ConfirmePassword: req.ConfirmPassword,
		},
	}

	result := ac.AuthUsecase.RestPassword(ctx, query)
	if !result.Success {
		log.Printf("Failed to update password: %s", result.Message)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": result.Message,
			"success": false,
		})
	}

	// Step 6: Delete OTP from Redis after successful password reset
	log.Println("Step 6: Deleting used OTP from Redis")
	ac.RedisCache.Cache.Del(ctx, redisKey)

	log.Printf("Password reset successfully for email: %s", req.Email)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully",
		"success": true,
	})
}

func (ac *AuthController) RefreshAccessToken(c *fiber.Ctx) error {
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
