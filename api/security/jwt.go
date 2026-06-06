package security

import (
	"errors"
	"go-gaurd/core/config"
	"go-gaurd/core/utils"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//
// =====================
// CONFIG
// =====================
//

type CustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func InitSecurity() *config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return cfg
}

func getAccessTokenSecret() []byte {
	cfg := InitSecurity()
	return []byte(cfg.ACCESS_TOKEN_SECRET)
}

func getRefreshTokenSecret() []byte {
	cfg := InitSecurity()
	return []byte(cfg.REFRESH_TOKEN_SECRET)
}

//
// =====================
// GENERATE TOKENS
// =====================
//

// Access Token
func GenerateAccessToken(userID string, role string, accessJTI string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID cannot be empty")
	}
	if role == "" {
		return "", errors.New("user role cannot be empty")
	}

	claims := &CustomClaims{
		UserID: userID,
		Role:   role,
		JTI:    accessJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(utils.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-gaurd",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getAccessTokenSecret())
}

// Refresh Token
func GenerateRefreshToken(userID string, role string, refreshJTI string) (string, error) {
	if userID == "" {
		return "", errors.New("user ID cannot be empty")
	}
	if role == "" {
		return "", errors.New("user role cannot be empty")
	}

	claims := &CustomClaims{
		UserID: userID,
		Role:   role,
		JTI:    refreshJTI,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(utils.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-gaurd",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getRefreshTokenSecret())
}

//
// =====================
// VALIDATION
// =====================
//

// Validate Access Token
func ValidateAccessToken(tokenString string) (string, string, string, *jwt.NumericDate, error) {
	if tokenString == "" {
		return "", "", "", nil, errors.New("token cannot be empty")
	}

	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getAccessTokenSecret(), nil
	})

	if err != nil {
		return "", "", "", nil, err
	}

	if !token.Valid {
		return "", "", "", nil, errors.New("invalid token")
	}

	return claims.UserID, claims.Role, claims.JTI, claims.ExpiresAt, nil
}

// Validate Refresh Token
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

//
// =====================
// REFRESH ACCESS TOKEN
// =====================
//

func RefreshAccessToken(refreshTokenString string) (string, error) {
	valid, claims, err := ValidateRefreshToken(refreshTokenString)
	if err != nil || !valid {
		return "", errors.New("invalid refresh token")
	}
	accessJTI := uuid.NewString()

	return GenerateAccessToken(claims.UserID, claims.Role, accessJTI)
}

//
// =====================
// GENERIC CHECK
// =====================
//

func CheckToken(tokenString string) (bool, string, string, string, error) {
	if tokenString == "" {
		return false, "", "", "", errors.New("token cannot be empty")
	}

	claims := &CustomClaims{}

	// Try access token first
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getAccessTokenSecret(), nil
	})

	if err == nil && token.Valid {
		return true, claims.UserID, claims.Role, claims.JTI, nil
	}

	// Try refresh token
	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getRefreshTokenSecret(), nil
	})

	if err != nil || !token.Valid {
		return false, "", "", "", errors.New("invalid token")
	}

	return true, claims.UserID, claims.Role, claims.JTI, nil
}

//
// =====================
// HELPERS
// =====================
//

func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is empty")
	}

	const bearerPrefix = "Bearer "

	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):], nil
	}

	return "", errors.New("invalid authorization header format")
}

func ExtractUserID(tokenString string) (string, error) {
	claims := &CustomClaims{}

	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return "", err
	}

	if claims.UserID == "" {
		return "", errors.New("user ID not found")
	}

	return claims.UserID, nil
}

func ExtractRole(tokenString string) (string, error) {
	claims := &CustomClaims{}

	_, _, err := new(jwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return "", err
	}

	if claims.Role == "" {
		return "", errors.New("user role not found")
	}

	return claims.Role, nil
}
