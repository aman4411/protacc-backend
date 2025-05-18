package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking existing user: %v", err)
	}
	return count > 0, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, tx pgx.Tx, user *models.User) error {
	err := tx.QueryRow(ctx, `
		INSERT INTO users (first_name, last_name, email, phone, password_hash, is_email_verified, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, user.FirstName, user.LastName, user.Email, user.Phone, user.PasswordHash, user.IsEmailVerified, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

func (r *UserRepository) CreateEmailVerification(ctx context.Context, tx pgx.Tx, userID, email, otp string, expiresAt time.Time) error {
	utils.LogInfo("Attempting to create email verification record", "userId", userID, "email", email)

	now := time.Now()
	_, err := tx.Exec(ctx, `
		INSERT INTO email_verifications (user_id, email, otp, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, email, otp, expiresAt, now)

	if err != nil {
		utils.LogError("Failed to create email verification record", "error", err.Error(), "userId", userID, "email", email)
		return fmt.Errorf("error storing verification details: %v", err)
	}

	utils.LogInfo("Email verification record created in database", "userId", userID, "email", email, "expiresAt", expiresAt)
	return nil
}

func (r *UserRepository) GetEmailVerification(ctx context.Context, tx pgx.Tx, email, otp string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := tx.QueryRow(ctx, `
		SELECT id, user_id, email, otp, expires_at, created_at
		FROM email_verifications
		WHERE email = $1 AND otp = $2 AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`, email, otp).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.Email,
		&verification.OTP,
		&verification.ExpiresAt,
		&verification.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("invalid or expired OTP")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting verification details: %v", err)
	}
	return &verification, nil
}

func (r *UserRepository) UpdateEmailVerificationStatus(ctx context.Context, tx pgx.Tx, userID string) error {
	result, err := tx.Exec(ctx, "UPDATE users SET is_email_verified = true WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) DeleteEmailVerification(ctx context.Context, tx pgx.Tx, verificationID string) error {
	_, err := tx.Exec(ctx, "DELETE FROM email_verifications WHERE id = $1", verificationID)
	if err != nil {
		return fmt.Errorf("error deleting verification record: %v", err)
	}
	return nil
}

func (r *UserRepository) GetUserFirstName(ctx context.Context, userID string) (string, error) {
	var firstName string
	err := r.pool.QueryRow(ctx, "SELECT first_name FROM users WHERE id = $1", userID).Scan(&firstName)
	if err != nil {
		return "", fmt.Errorf("error getting user details: %v", err)
	}
	return firstName, nil
}

func (r *UserRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	utils.LogInfo("Fetching user by email", "email", email)

	var user models.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, first_name, last_name, email, phone, password_hash, is_email_verified, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.IsEmailVerified,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		utils.LogError("User not found", "email", email)
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		utils.LogError("Error fetching user", "error", err.Error(), "email", email)
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	utils.LogInfo("User fetched successfully", "userId", user.ID, "email", user.Email)
	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	utils.LogInfo("Fetching user by ID", "userId", userID)

	var user models.User
	err := r.pool.QueryRow(ctx, `
		SELECT id, first_name, last_name, email, phone, password_hash, is_email_verified, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.IsEmailVerified,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		utils.LogError("User not found", "userId", userID)
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		utils.LogError("Error fetching user", "error", err.Error(), "userId", userID)
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	utils.LogInfo("User fetched successfully", "userId", user.ID, "email", user.Email)
	return &user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]*models.User, error) {
	utils.LogInfo("Fetching all users")

	rows, err := r.pool.Query(ctx, `
		SELECT id, first_name, last_name, email, phone, password_hash, is_email_verified, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		utils.LogError("Error fetching users", "error", err.Error())
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.PasswordHash,
			&user.IsEmailVerified,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			utils.LogError("Error scanning user row", "error", err.Error())
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		utils.LogError("Error iterating user rows", "error", err.Error())
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	utils.LogInfo("Users fetched successfully", "count", len(users))
	return users, nil
}
