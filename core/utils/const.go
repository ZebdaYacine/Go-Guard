package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

const (
	RoleGuest    = "guest"
	RoleClient   = "client"
	RoleArtisant = "artisant"
	RoleUser     = "user"
	RoleAdmin    = "admin"
)

func GetValidRoles(role string) int {
	switch role {
	case RoleGuest:
		return 1
	case RoleClient:
		return 2
	case RoleArtisant:
		return 3
	case RoleAdmin:
		return 4
	case RoleUser:
		return 5
	default:
		return 0
	}
}

func HashPasswordSHA256(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CheckPasswordHashSHA256 compares a password with its SHA256 hash
func CheckPasswordHashSHA256(password, hash string) bool {
	hashedPassword := HashPasswordSHA256(password)
	return hashedPassword == hash
}
