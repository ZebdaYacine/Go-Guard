package test

import (
	"fmt"
	"go-gaurd/core/utils"
	"testing"
)

func TestGenerateOTP(t *testing.T) {

	otp := utils.GenerateOTP()
	if otp == "" {
		t.Error("Expected a valid OTP, got empty string")
	}

	fmt.Println(otp)

}
