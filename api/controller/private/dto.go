package private

import "mime/multipart"

type UpdateProfileRequest struct {
	UserName string                `form:"user_name" validate:"required,min=3,max=50"`
	Email    string                `form:"email" validate:"required,email"`
	Phone    string                `form:"phone" validate:"required"`
	Picture  *multipart.FileHeader `form:"picture" validate:"omitempty"`
}
type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email"`
	OTP             string `json:"otp" validate:"required,len=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	DeviceId     string `json:"device_id" validate:"required"`
}

type LogoutRequest struct {
	DeviceId string `json:"device_id" validate:"required"`
}
