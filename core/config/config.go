package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser               string
	DBPass               string
	DBHost               string
	DBPort               string
	DBName               string
	REDIS_HOST           string
	REDIS_PORT           string
	ACCESS_TOKEN_SECRET  string
	REFRESH_TOKEN_SECRET string
	FROM                 string
	SMTP_PASS            string
	SMTP_USER            string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBUser:               os.Getenv("USER_DB"),
		DBPass:               os.Getenv("PASS_DB"),
		DBHost:               os.Getenv("HOST_DB"),
		DBPort:               os.Getenv("PORT_DB"),
		DBName:               os.Getenv("DB_NAME"),
		REDIS_HOST:           os.Getenv("REDIS_HOST"),
		REDIS_PORT:           os.Getenv("REDIS_PORT"),
		ACCESS_TOKEN_SECRET:  os.Getenv("ACCESS_TOKEN_SECRET"),
		REFRESH_TOKEN_SECRET: os.Getenv("REFRESH_TOKEN_SECRET"),
		FROM:                 os.Getenv("FROM"),
		SMTP_USER:            os.Getenv("SMTP_USER"),
		SMTP_PASS:            os.Getenv("SMTP_PASS"),
	}, nil
}
