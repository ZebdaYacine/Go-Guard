package test

import (
	"fmt"
	"go-gaurd/core/utils"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestGenerateOTP(t *testing.T) {

	otp := utils.GenerateOTP()
	if otp == "" {
		t.Error("Expected a valid OTP, got empty string")
	}
	fmt.Println(otp)
}

func TestSendOTP(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	r := utils.SendOTP("test@example.com", utils.GenerateOTP())
	if r != nil {
		t.Error("Expected successful OTP sending, got error", r.Error())
	}
}
