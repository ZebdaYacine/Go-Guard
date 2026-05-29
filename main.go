package main

import (
	"context"
	"fmt"
	"log"

	"go-gaurd/core/di"
	"go-gaurd/feature/auth/usecase"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	AuthUseCase, err := di.InitializeAuthApplication()
	if err != nil {
		log.Fatal("Error initializing auth application")
	}

	result := AuthUseCase.CreateAccount(ctx, usecase.Query{User: usecase.User_Entity{
		User_name: "ZEBDA",
		Email:     "zebda@example.com",
		Phone:     "1234567890",
		Password:  "Sedcw@ewr456546",
		Role:      "admin",
		Sex:       "male",
		Picture:   "http://example.com/profile.jpg",
	}})

	fmt.Println(result.Message)

}
