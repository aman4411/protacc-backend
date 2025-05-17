package models

import "time"

// EmailVerification represents the email verification data
type EmailVerification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Email     string    `json:"email"`
	OTP       string    `json:"otp"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// EmailVerificationRequest represents the verification request
type EmailVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required"`
}
