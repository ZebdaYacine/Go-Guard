package test

import (
	"context"
	"fmt"
	"go-gaurd/core/di/gen"
	"go-gaurd/feature/auth/domain"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	"github.com/joho/godotenv"
)

func TestLogin(t *testing.T) {

	ctx := context.Background()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	fmt.Println(dir)
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	appDependencies, err := gen.InitializeAll()
	if err != nil {
		log.Fatal("Error initializing all dependencies:", err)
	}
	repo := domain.NewAuthRepository(appDependencies.Database)
	result := repo.Login(ctx, domain.Login_Entity{
		Email:    "zebdaadam1996@gmail.com",
		Password: "StrongP@ss123",
	})

	fmt.Println(result)

}

func TestCreateAccount(t *testing.T) {

	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	appDependencies, err := gen.InitializeAll()
	if err != nil {
		log.Fatal("Error initializing all dependencies:", err)
	}
	repo := domain.NewAuthRepository(appDependencies.Database)
	result := repo.CreateAccount(ctx, domain.User_Entity{
		User_name: "Zed",
		Email:     "zebdaadam1996@gmail.com",
		Phone:     "1234567890",
		Password:  "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8",
		Role:      "user",
		Sex:       "Male",
		Picture:   "",
	})

	fmt.Println(result)

}
