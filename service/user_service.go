package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
	mail *MailService
}

func NewUserService(repo *repository.UserRepository, mail *MailService) *UserService {
	return &UserService{
		repo: repo,
		mail: mail,
	}
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ValidatePassword checks if password and confirm password match
func ValidatePassword(password, confirmPassword string) error {
	if password != confirmPassword {
		return fmt.Errorf("password and confirm password do not match")
	}
	return nil
}

func (s *UserService) CreateUser(ctx context.Context, req *models.UserRequest) (*models.UserResponse, error) {
	// Check if user exists
	exists, err := s.repo.CheckUserExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Create user
	now := time.Now()
	user := &models.User{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		Phone:           req.Phone,
		PasswordHash:    string(hashedPassword),
		IsEmailVerified: false,
		Role:            "user",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.CreateUser(ctx, tx, user); err != nil {
		return nil, err
	}

	// Generate and store OTP
	otp := generateOTP()
	expiresAt := now.Add(15 * time.Minute)

	if err := s.repo.CreateEmailVerification(ctx, tx, user.ID, user.Email, otp, expiresAt); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	// Send verification email
	if err := s.mail.SendVerificationEmail(user.Email, user.FirstName, otp); err != nil {
		return nil, fmt.Errorf("error sending verification email: %v", err)
	}

	return &models.UserResponse{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Phone:           user.Phone,
		IsEmailVerified: user.IsEmailVerified,
		Role:            user.Role,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, req *models.EmailVerificationRequest) error {
	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Get verification details
	verification, err := s.repo.GetEmailVerification(ctx, tx, req.Email, req.OTP)
	if err != nil {
		return err
	}

	// Update user verification status
	if err := s.repo.UpdateEmailVerificationStatus(ctx, tx, verification.UserID); err != nil {
		return err
	}

	// Delete verification record
	if err := s.repo.DeleteEmailVerification(ctx, tx, verification.ID); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	// Get user details for welcome email
	firstName, err := s.repo.GetUserFirstName(ctx, verification.UserID)
	if err != nil {
		return err
	}

	// Send welcome email
	if err := s.mail.SendWelcomeEmail(req.Email, firstName); err != nil {
		return fmt.Errorf("error sending welcome email: %v", err)
	}

	return nil
}

// generateOTP generates a 6-digit OTP
func generateOTP() string {
	const otpLength = 6
	const digits = "0123456789"
	b := make([]byte, otpLength)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000) // fallback
	}
	for i := 0; i < otpLength; i++ {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}
