package public

// Request structs
type RegisterRequest struct {
	User_name string `json:"user_name" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Phone     string `json:"phone" validate:"required"`
	Password  string `json:"password" validate:"required,min=6"`
	Role      string `json:"role" validate:"required"`
	Sex       string `json:"sex" validate:"required,oneof=male female other"`
	Picture   string `json:"picture"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
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
}

type SendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	Type  string `json:"type" validate:"required"` // verification, forgetpassword, login
}

type CheckOTPRequest struct {
	Email       string `json:"email" validate:"required,email"`
	OTP         string `json:"otp" validate:"required,len=6"`
	Type        string `json:"type" validate:"required"`
}

// Response structs
type RegisterResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	User    interface{} `json:"user"`
}

type LoginResponse struct {
	Message      string      `json:"message"`
	Success      bool        `json:"success"`
	User         interface{} `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefrechToken string      `json:"refresh_token"`
}
