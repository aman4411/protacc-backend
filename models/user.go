package models

import "time"

// User represents the user data in the application
type User struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	PasswordHash    string    `json:"passwordHash"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	Role            string    `json:"role"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// UserRequest represents the signup request data
type UserRequest struct {
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

// UserResponse represents the user data sent back to client
type UserResponse struct {
	ID              string    `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	Role            string    `json:"role"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
