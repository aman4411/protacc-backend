package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
	"github.com/aman4411/protacc-backend/utils"
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
	utils.LogInfo("Starting user creation process", "email", req.Email)

	// Check if user exists
	exists, err := s.repo.CheckUserExists(ctx, req.Email)
	if err != nil {
		utils.LogError("Error checking if user exists", "error", err.Error(), "email", req.Email)
		return nil, err
	}
	if exists {
		utils.LogError("User already exists", "email", req.Email)
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	utils.LogInfo("User existence check passed", "email", req.Email)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError("Error hashing password", "error", err.Error(), "email", req.Email)
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	utils.LogInfo("Password hashed successfully", "email", req.Email)

	// Begin transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		utils.LogError("Error starting transaction", "error", err.Error(), "email", req.Email)
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	utils.LogInfo("Database transaction started", "email", req.Email)

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
		utils.LogError("Error creating user in database", "error", err.Error(), "email", req.Email)
		return nil, err
	}

	utils.LogInfo("User created in database", "userId", user.ID, "email", user.Email)

	// Generate and store OTP
	utils.LogInfo("Generating OTP for email verification", "email", req.Email)
	otp := generateOTP()
	utils.LogInfo("OTP generated successfully", "email", req.Email, "otp", otp)

	expiresAt := now.Add(15 * time.Minute)
	utils.LogInfo("Creating email verification record", "email", req.Email, "userId", user.ID, "expiresAt", expiresAt)

	if err := s.repo.CreateEmailVerification(ctx, tx, user.ID, user.Email, otp, expiresAt); err != nil {
		utils.LogError("Error creating email verification", "error", err.Error(), "email", req.Email, "userId", user.ID)
		return nil, fmt.Errorf("failed to create email verification: %v", err)
	}

	utils.LogInfo("Email verification record created successfully", "email", req.Email, "userId", user.ID)

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		utils.LogError("Error committing transaction", "error", err.Error(), "email", req.Email)
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	utils.LogInfo("Database transaction committed successfully", "email", req.Email)

	// Send verification email
	if err := s.mail.SendVerificationEmail(user.Email, user.FirstName, otp); err != nil {
		utils.LogError("Error sending verification email", "error", err.Error(), "email", req.Email)
		return nil, fmt.Errorf("error sending verification email: %v", err)
	}

	utils.LogInfo("Verification email sent successfully", "email", req.Email)

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

func (s *UserService) AuthenticateUser(ctx context.Context, req *models.LoginRequest) (*models.UserResponse, error) {
	utils.LogInfo("Starting user authentication", "email", req.Email)

	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		utils.LogError("Failed to get user", "error", err.Error(), "email", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		utils.LogError("Invalid password", "email", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	utils.LogInfo("User authenticated successfully", "userId", user.ID, "email", user.Email)

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

func (s *UserService) GetUserByID(ctx context.Context, userID string) (*models.UserResponse, error) {
	utils.LogInfo("Fetching user by ID", "userId", userID)

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		utils.LogError("Failed to get user", "error", err.Error(), "userId", userID)
		return nil, err
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

func (s *UserService) GetUsers(ctx context.Context) ([]models.UserResponse, error) {
	utils.LogInfo("Fetching all users")

	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		utils.LogError("Failed to get users", "error", err.Error())
		return nil, err
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			Phone:           user.Phone,
			IsEmailVerified: user.IsEmailVerified,
			Role:            user.Role,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		})
	}

	return userResponses, nil
}
