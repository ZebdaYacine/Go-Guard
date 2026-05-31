package security

import (
	"errors"
	"go-gaurd/core/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Custom claims structure
type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func InitSecurity() *config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return cfg
}

// getAccessTokenSecret returns the access token secret from config
func getAccessTokenSecret() []byte {
	cfg := InitSecurity()
	return []byte(cfg.ACCESS_TOKEN_SECRET)
}

// getRefreshTokenSecret returns the refresh token secret from config
func getRefreshTokenSecret() []byte {
	cfg := InitSecurity()
	return []byte(cfg.REFRESH_TOKEN_SECRET)
}

// TokenExpiry defines token expiration durations
var (
	accessTokenExpiry  = 15 * time.Minute
	refreshTokenExpiry = 7 * 24 * time.Hour // 7 days
)

// SetTokenExpiry allows customizing token expiration times
func SetTokenExpiry(accessExpiry, refreshExpiry time.Duration) {
	accessTokenExpiry = accessExpiry
	refreshTokenExpiry = refreshExpiry
}

// GenerateAccessToken creates a new JWT access token for a user
func GenerateAccessToken(userID string, role string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID cannot be empty")
	}
	if role == "" {
		return "", errors.New("user role cannot be empty")
	}

	claims := &CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app-name",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getAccessTokenSecret())
}

// GenerateRefreshToken creates a new JWT refresh token for a user
func GenerateRefreshToken(userID string, role string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID cannot be empty")
	}
	if role == "" {
		return "", errors.New("user role cannot be empty")
	}

	claims := &CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app-name",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getRefreshTokenSecret())
}

// CheckToken validates and parses a JWT token (works for both access and refresh tokens)
// Returns (isValid, userID, role, error)
func CheckToken(tokenString string) (bool, string, string, error) {
	if tokenString == "" {
		return false, "", "", errors.New("token cannot be empty")
	}

	// Try to parse with access token secret first
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getAccessTokenSecret(), nil
	})

	// If access token validation fails, try refresh token secret
	if err != nil {
		token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return getRefreshTokenSecret(), nil
		})

		if err != nil {
			return false, "", "", errors.New("invalid or expired token")
		}
	}

	if !token.Valid {
		return false, "", "", errors.New("invalid token")
	}

	// Validate required claims
	if claims.UserID == "" {
		return false, "", "", errors.New("user ID not found in token")
	}
	if claims.Role == "" {
		return false, "", "", errors.New("user role not found in token")
	}

	return true, claims.UserID, claims.Role, nil
}

// ValidateAccessToken specifically validates an access token
func ValidateAccessToken(tokenString string) (bool, string, string, error) {
	if tokenString == "" {
		return false, "", "", errors.New("token cannot be empty")
	}

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getAccessTokenSecret(), nil
	})

	if err != nil {
		return false, "", "", err
	}

	if !token.Valid {
		return false, "", "", errors.New("invalid token")
	}

	return true, claims.UserID, claims.Role, nil
}

// ValidateRefreshToken specifically validates a refresh token
func ValidateRefreshToken(tokenString string) (bool, *CustomClaims, error) {
	if tokenString == "" {
		return false, nil, errors.New("token cannot be empty")
	}

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getRefreshTokenSecret(), nil
	})

	if err != nil {
		return false, nil, err
	}

	if !token.Valid {
		return false, nil, errors.New("invalid token")
	}

	return true, claims, nil
}

// RefreshAccessToken generates a new access token using a valid refresh token
func RefreshAccessToken(refreshTokenString string) (string, error) {
	valid, claims, err := ValidateRefreshToken(refreshTokenString)
	if err != nil || !valid {
		return "", errors.New("invalid refresh token")
	}

	// Generate new access token
	return GenerateAccessToken(claims.UserID, claims.Role)
}

// ExtractUserID extracts user ID from a token without full validation
func ExtractUserID(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token cannot be empty")
	}

	claims := &CustomClaims{}

	// Try parsing without validation to extract claims
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		if claims.UserID != "" {
			return claims.UserID, nil
		}
	}

	return "", errors.New("user ID not found in token")
}

// ExtractRole extracts user role from a token without full validation
func ExtractRole(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("token cannot be empty")
	}

	claims := &CustomClaims{}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		if claims.Role != "" {
			return claims.Role, nil
		}
	}

	return "", errors.New("user role not found in token")
}

// Middleware helper to extract token from Authorization header
// Expects header in format: "Bearer <token>"
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	// Check if the header has the Bearer prefix
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):], nil
	}

	return "", errors.New("invalid authorization header format")
}
