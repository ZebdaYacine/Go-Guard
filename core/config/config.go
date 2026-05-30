package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser            string
	DBPass            string
	DBHost            string
	DBPort            string
	DBName            string
	EZAUTH_ADDR       string
	EZAUTH_API_KEY    string
	EZAUTH_BASE_URL   string
	EZAUTH_DEBUG      string
	EZAUTH_JWT_SECRET string
	EZAUTH_TIMEOUT    string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBUser:            os.Getenv("USER_DB"),
		DBPass:            os.Getenv("PASS_DB"),
		DBHost:            os.Getenv("HOST_DB"),
		DBPort:            os.Getenv("PORT_DB"),
		DBName:            os.Getenv("DB_NAME"),
		EZAUTH_ADDR:       os.Getenv("EZAUTH_ADDR"),
		EZAUTH_API_KEY:    os.Getenv("EZAUTH_API_KEY"),
		EZAUTH_BASE_URL:   os.Getenv("EZAUTH_BASE_URL"),
		EZAUTH_DEBUG:      os.Getenv("EZAUTH_DEBUG"),
		EZAUTH_JWT_SECRET: os.Getenv("EZAUTH_JWT_SECRET"),
		EZAUTH_TIMEOUT:    os.Getenv("EZAUTH_TIMEOUT"),
	}, nil
}
