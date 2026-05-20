package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
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

// GetUsersWithFilters returns filtered users with pagination
func (s *UserService) GetUsersWithFilters(ctx context.Context, page, limit int, search, role string, emailVerified *bool) ([]models.UserResponse, int, error) {
	users, total, err := s.repo.GetUsersWithFilters(ctx, page, limit, search, role, emailVerified)
	if err != nil {
		return nil, 0, err
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = models.UserResponse{
			ID:              user.ID,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Email:           user.Email,
			Phone:           user.Phone,
			Role:            user.Role,
			IsEmailVerified: user.IsEmailVerified,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		}
	}

	return userResponses, total, nil
}

// UpdateUserRole updates a user's role
func (s *UserService) UpdateUserRole(ctx context.Context, userID, role string) error {
	return s.repo.UpdateUserRole(ctx, userID, role)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	if err := s.repo.UpdateUserProfile(ctx, userID, req.FirstName, req.LastName, req.Phone); err != nil {
		return nil, err
	}
	return s.GetUserByID(ctx, userID)
}

func generateResetToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func getFrontendURL() string {
	if url := os.Getenv("FRONTEND_URL"); url != "" {
		return url
	}
	return "http://localhost:3000"
}

func (s *UserService) sendPasswordResetEmail(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Do not reveal whether the email exists
		utils.LogInfo("Password reset requested for unknown email", "email", email)
		return nil
	}

	token, err := generateResetToken()
	if err != nil {
		return fmt.Errorf("error generating reset token: %v", err)
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	if err := s.repo.InvalidatePasswordResetTokens(ctx, tx, user.ID); err != nil {
		return err
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	if err := s.repo.CreatePasswordResetToken(ctx, tx, user.ID, token, expiresAt); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", getFrontendURL(), token)
	if err := s.mail.SendPasswordResetEmail(user.Email, user.FirstName, resetLink); err != nil {
		return fmt.Errorf("error sending password reset email: %v", err)
	}

	return nil
}

// RequestPasswordReset sends a reset link to the given email (public endpoint)
func (s *UserService) RequestPasswordReset(ctx context.Context, email string) error {
	return s.sendPasswordResetEmail(ctx, email)
}

// RequestPasswordResetForUser sends a reset link to the authenticated user's email
func (s *UserService) RequestPasswordResetForUser(ctx context.Context, userID string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	return s.sendPasswordResetEmail(ctx, user.Email)
}

func (s *UserService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) error {
	if err := ValidatePassword(req.Password, req.ConfirmPassword); err != nil {
		return err
	}

	reset, err := s.repo.GetPasswordResetByToken(ctx, req.Token)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	if err := s.repo.UpdateUserPasswordTx(ctx, tx, reset.UserID, string(hashedPassword)); err != nil {
		return err
	}

	if err := s.repo.MarkPasswordResetUsed(ctx, tx, reset.ID); err != nil {
		return err
	}

	if err := s.repo.InvalidatePasswordResetTokens(ctx, tx, reset.UserID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *UserService) ValidateResetToken(ctx context.Context, token string) error {
	_, err := s.repo.GetPasswordResetByToken(ctx, token)
	return err
}

// DashboardStats represents the dashboard statistics response
type DashboardStats struct {
	TotalUsers    int                   `json:"total_users"`
	TotalOrders   int                   `json:"total_orders"`
	TotalRevenue  float64               `json:"total_revenue"`
	PendingOrders int                   `json:"pending_orders"`
	RecentUsers   []models.UserResponse `json:"recent_users"`
	RecentOrders  []models.RecentOrder  `json:"recent_orders"`
	UserGrowth    string                `json:"user_growth"`
	OrderGrowth   string                `json:"order_growth"`
	RevenueGrowth string                `json:"revenue_growth"`
}

// GetDashboardStats returns dashboard statistics for admin
func (s *UserService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	// Get total users count
	totalUsers, err := s.repo.GetTotalUsersCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users count: %v", err)
	}

	// Get total orders count and revenue
	totalOrders, totalRevenue, err := s.repo.GetOrdersStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders stats: %v", err)
	}

	// Get pending orders count
	pendingOrders, err := s.repo.GetPendingOrdersCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending orders count: %v", err)
	}

	// Get recent users (last 5)
	recentUsers, err := s.repo.GetRecentUsers(ctx, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent users: %v", err)
	}

	// Get recent orders (last 5)
	recentOrders, err := s.repo.GetRecentOrders(ctx, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent orders: %v", err)
	}

	// Get current month stats
	currentMonthUsers, err := s.repo.GetCurrentMonthUsersCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current month users: %v", err)
	}

	// Get last month stats
	lastMonthUsers, err := s.repo.GetLastMonthUsersCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get last month users: %v", err)
	}

	// Get current month orders and revenue
	currentMonthOrders, currentMonthRevenue, err := s.repo.GetCurrentMonthOrdersStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current month orders stats: %v", err)
	}

	// Get last month orders and revenue
	lastMonthOrders, lastMonthRevenue, err := s.repo.GetLastMonthOrdersStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get last month orders stats: %v", err)
	}

	// Calculate percentage changes
	userGrowth := calculatePercentageChange(lastMonthUsers, currentMonthUsers)
	orderGrowth := calculatePercentageChange(lastMonthOrders, currentMonthOrders)
	revenueGrowth := calculatePercentageChangeFloat(lastMonthRevenue, currentMonthRevenue)

	return &DashboardStats{
		TotalUsers:    totalUsers,
		TotalOrders:   totalOrders,
		TotalRevenue:  totalRevenue,
		PendingOrders: pendingOrders,
		RecentUsers:   recentUsers,
		RecentOrders:  recentOrders,
		UserGrowth:    userGrowth,
		OrderGrowth:   orderGrowth,
		RevenueGrowth: revenueGrowth,
	}, nil
}

// calculatePercentageChange calculates the percentage change between old and new values
func calculatePercentageChange(oldValue, newValue int) string {
	if oldValue == 0 {
		if newValue > 0 {
			return "+100%"
		}
		return "0%"
	}

	change := float64(newValue-oldValue) / float64(oldValue) * 100
	if change > 0 {
		return fmt.Sprintf("+%.1f%%", change)
	} else if change < 0 {
		return fmt.Sprintf("%.1f%%", change)
	}
	return "0%"
}

// calculatePercentageChangeFloat calculates the percentage change for float values
func calculatePercentageChangeFloat(oldValue, newValue float64) string {
	if oldValue == 0 {
		if newValue > 0 {
			return "+100%"
		}
		return "0%"
	}

	change := (newValue - oldValue) / oldValue * 100
	if change > 0 {
		return fmt.Sprintf("+%.1f%%", change)
	} else if change < 0 {
		return fmt.Sprintf("%.1f%%", change)
	}
	return "0%"
}
