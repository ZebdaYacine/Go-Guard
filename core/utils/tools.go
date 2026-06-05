package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/go-sql-driver/mysql"
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

func GetCodeGender(sex string) string {
	if sex == "male" {
		return "M"
	}
	return "F"
}

func HashPasswordSHA256(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func CheckPasswordHashSHA256(password, hash string) bool {
	hashedPassword := HashPasswordSHA256(password)
	return hashedPassword == hash
}

func HandleMysqlError(err error) string {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		errorMsg := mysqlErr.Message
		switch mysqlErr.Number {
		case 1062:
			switch {
			case strings.Contains(errorMsg, "email"):
				return "Email is already registered"
			case strings.Contains(errorMsg, "username"):
				return "Username is already taken"
			case strings.Contains(errorMsg, "phone"):
				return "Phone number is already registered"
			default:
				return "Duplicate entry detected"
			}
		case 1452:
			return "Referenced record does not exist (Foreign key failure)"
		case 1048:
			return "Required fields cannot be empty"
		case 1406:
			return "Provided data exceeds the maximum allowed length"
		case 1040:
			return "Database is overloaded: too many connections active"
		case 2006:
			return "Database connection was lost or timed out"
		case 1045:
			return "Database authentication failed (Access denied)"

		default:
			return "A database restriction error occurred"
		}
	}
	return "An unexpected error occurred"
}

func GenerateOTP() string {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%06d", n.Int64()+100000)
}
