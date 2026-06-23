package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// displayName picks the stored reviewer name (admin reviews) or falls back to
// the purchasing user's "First L." form.
func displayName(stored, first, last sql.NullString) string {
	if stored.Valid && strings.TrimSpace(stored.String) != "" {
		return strings.TrimSpace(stored.String)
	}
	return reviewerName(first.String, last.String)
}

type ReviewRepository struct {
	db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// reviewerName renders "Rahul S." from first/last name, avoiding any PII beyond
// the first name + last initial.
func reviewerName(first, last string) string {
	first = strings.TrimSpace(first)
	last = strings.TrimSpace(last)
	if last == "" {
		return first
	}
	initial := []rune(last)[0]
	return first + " " + string(initial) + "."
}

// HasPaidOrderForService reports whether the user has any non-cancelled, paid
// order containing the service (booking amount paid or beyond).
func (r *ReviewRepository) HasPaidOrderForService(ctx context.Context, userID string, serviceID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM order_items oi
			JOIN orders o ON oi.order_id = o.id
			WHERE oi.service_id = $1
			  AND o.user_id = $2
			  AND o.status NOT IN ('pending_payment', 'pending_booking_payment', 'cancelled')
		)`
	var exists bool
	if err := r.db.QueryRow(ctx, query, serviceID, userID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// GetUserReviewForService returns the user's existing review for a service, or
// nil if none exists.
func (r *ReviewRepository) GetUserReviewForService(ctx context.Context, userID string, serviceID int) (*models.Review, error) {
	query := `
		SELECT id, service_id, user_id, rating, comment, status, created_at, updated_at
		FROM reviews
		WHERE user_id = $1 AND service_id = $2`
	var rev models.Review
	err := r.db.QueryRow(ctx, query, userID, serviceID).Scan(
		&rev.ID, &rev.ServiceID, &rev.UserID, &rev.Rating, &rev.Comment,
		&rev.Status, &rev.CreatedAt, &rev.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rev, nil
}

// UpsertReview creates or updates the user's review for a service (one per pair).
func (r *ReviewRepository) UpsertReview(ctx context.Context, rev *models.Review) (*models.Review, error) {
	query := `
		INSERT INTO reviews (service_id, user_id, rating, comment, status)
		VALUES ($1, $2, $3, $4, 'published')
		ON CONFLICT (service_id, user_id)
		DO UPDATE SET rating = EXCLUDED.rating, comment = EXCLUDED.comment,
			status = 'published', updated_at = NOW()
		RETURNING id, status, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, rev.ServiceID, rev.UserID, rev.Rating, rev.Comment).Scan(
		&rev.ID, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rev, nil
}

// GetServiceReviews returns published reviews for a service plus the aggregate.
func (r *ReviewRepository) GetServiceReviews(ctx context.Context, serviceID int) ([]models.Review, models.ReviewSummary, error) {
	var summary models.ReviewSummary
	aggQuery := `SELECT COALESCE(AVG(rating), 0), COUNT(*) FROM reviews WHERE service_id = $1 AND status = 'published'`
	if err := r.db.QueryRow(ctx, aggQuery, serviceID).Scan(&summary.Average, &summary.Count); err != nil {
		return nil, summary, err
	}

	query := `
		SELECT r.id, r.service_id, r.rating, r.comment, r.created_at, r.reviewer_name, u.first_name, u.last_name
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.service_id = $1 AND r.status = 'published'
		ORDER BY r.created_at DESC`
	rows, err := r.db.Query(ctx, query, serviceID)
	if err != nil {
		return nil, summary, err
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var rev models.Review
		var stored, first, last sql.NullString
		if err := rows.Scan(&rev.ID, &rev.ServiceID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &stored, &first, &last); err != nil {
			return nil, summary, err
		}
		rev.ReviewerName = displayName(stored, first, last)
		reviews = append(reviews, rev)
	}
	return reviews, summary, nil
}

// GetTopReviews returns a random selection of highly-rated published reviews
// (with non-empty comments) across all services, for the homepage.
func (r *ReviewRepository) GetTopReviews(ctx context.Context, limit int) ([]models.Review, error) {
	query := `
		SELECT r.id, r.service_id, r.rating, r.comment, r.created_at, r.reviewer_name, u.first_name, u.last_name, s.name, s.slug
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		JOIN services s ON r.service_id = s.id
		WHERE r.status = 'published' AND r.rating >= 4
		  AND r.comment IS NOT NULL AND LENGTH(TRIM(r.comment)) > 0
		ORDER BY RANDOM()
		LIMIT $1`
	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var rev models.Review
		var stored, first, last sql.NullString
		if err := rows.Scan(&rev.ID, &rev.ServiceID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &stored, &first, &last, &rev.ServiceName, &rev.ServiceSlug); err != nil {
			return nil, err
		}
		rev.ReviewerName = displayName(stored, first, last)
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

// ListReviews returns all reviews (admin moderation).
func (r *ReviewRepository) ListReviews(ctx context.Context) ([]models.Review, error) {
	query := `
		SELECT r.id, r.service_id, r.rating, r.comment, r.status, r.created_at, r.updated_at,
			r.reviewer_name, u.first_name, u.last_name, s.name
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		JOIN services s ON r.service_id = s.id
		ORDER BY r.created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var rev models.Review
		var stored, first, last sql.NullString
		if err := rows.Scan(&rev.ID, &rev.ServiceID, &rev.Rating, &rev.Comment, &rev.Status,
			&rev.CreatedAt, &rev.UpdatedAt, &stored, &first, &last, &rev.ServiceName); err != nil {
			return nil, err
		}
		rev.ReviewerName = displayName(stored, first, last)
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

// AdminCreateReview inserts a review authored by an admin (no purchasing user;
// a custom display name is stored).
func (r *ReviewRepository) AdminCreateReview(ctx context.Context, serviceID, rating int, comment, reviewerName string) (*models.Review, error) {
	query := `
		INSERT INTO reviews (service_id, user_id, rating, comment, reviewer_name, status)
		VALUES ($1, NULL, $2, $3, $4, 'published')
		RETURNING id, status, created_at, updated_at`
	rev := &models.Review{
		ServiceID:    serviceID,
		Rating:       rating,
		Comment:      comment,
		ReviewerName: reviewerName,
	}
	err := r.db.QueryRow(ctx, query, serviceID, rating, comment, reviewerName).Scan(
		&rev.ID, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return rev, nil
}

// SetReviewStatus updates a review's visibility (published / hidden).
func (r *ReviewRepository) SetReviewStatus(ctx context.Context, id int, status models.ReviewStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE reviews SET status = $1, updated_at = NOW() WHERE id = $2`, status, id)
	return err
}

// DeleteReview permanently removes a review (admin).
func (r *ReviewRepository) DeleteReview(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	return err
}
