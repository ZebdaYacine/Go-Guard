package utils

import "time"

const (
	OTP_MAIL_HTML = "otp_mailer"

	RoleGuest = "guest"
	// RoleClient   = "client"
	// RoleArtisant = "artisant"
	RoleUser  = "user"
	RoleAdmin = "admin"

	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour // 7 days
)
