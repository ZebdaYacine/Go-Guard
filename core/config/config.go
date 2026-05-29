package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBUser: os.Getenv("USER_DB"),
		DBPass: os.Getenv("PASS_DB"),
		DBHost: os.Getenv("HOST_DB"),
		DBPort: os.Getenv("PORT_DB"),
		DBName: os.Getenv("DB_NAME"),
	}, nil
}
