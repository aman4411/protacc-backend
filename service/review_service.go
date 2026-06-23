package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type ReviewService struct {
	repo *repository.ReviewRepository
}

func NewReviewService(repo *repository.ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

// Eligibility describes whether a user may review a service and any existing review.
type Eligibility struct {
	CanReview       bool           `json:"can_review"`
	AlreadyReviewed bool           `json:"already_reviewed"`
	Existing        *models.Review `json:"existing,omitempty"`
}

// SubmitReview validates eligibility and creates/updates the user's review.
func (s *ReviewService) SubmitReview(ctx context.Context, userID string, serviceID, rating int, comment string) (*models.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	hasOrder, err := s.repo.HasPaidOrderForService(ctx, userID, serviceID)
	if err != nil {
		return nil, err
	}
	if !hasOrder {
		return nil, fmt.Errorf("you can review a service only after purchasing it")
	}

	rev := &models.Review{
		ServiceID: serviceID,
		UserID:    userID,
		Rating:    rating,
		Comment:   strings.TrimSpace(comment),
	}
	return s.repo.UpsertReview(ctx, rev)
}

// GetServiceReviews returns published reviews and the aggregate for a service.
func (s *ReviewService) GetServiceReviews(ctx context.Context, serviceID int) ([]models.Review, models.ReviewSummary, error) {
	return s.repo.GetServiceReviews(ctx, serviceID)
}

// GetTopReviews returns up to limit highly-rated reviews for the homepage.
func (s *ReviewService) GetTopReviews(ctx context.Context, limit int) ([]models.Review, error) {
	if limit < 1 || limit > 20 {
		limit = 6
	}
	return s.repo.GetTopReviews(ctx, limit)
}

// GetEligibility reports whether the user can review the service.
func (s *ReviewService) GetEligibility(ctx context.Context, userID string, serviceID int) (*Eligibility, error) {
	hasOrder, err := s.repo.HasPaidOrderForService(ctx, userID, serviceID)
	if err != nil {
		return nil, err
	}
	existing, err := s.repo.GetUserReviewForService(ctx, userID, serviceID)
	if err != nil {
		return nil, err
	}
	return &Eligibility{
		CanReview:       hasOrder,
		AlreadyReviewed: existing != nil,
		Existing:        existing,
	}, nil
}

// ListReviews returns all reviews for admin moderation.
func (s *ReviewService) ListReviews(ctx context.Context) ([]models.Review, error) {
	return s.repo.ListReviews(ctx)
}

// AdminCreateReview creates a review authored by an admin (no purchase check).
func (s *ReviewService) AdminCreateReview(ctx context.Context, serviceID, rating int, comment, reviewerName string) (*models.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}
	reviewerName = strings.TrimSpace(reviewerName)
	if reviewerName == "" {
		return nil, fmt.Errorf("reviewer name is required")
	}
	if serviceID == 0 {
		return nil, fmt.Errorf("service is required")
	}
	return s.repo.AdminCreateReview(ctx, serviceID, rating, strings.TrimSpace(comment), reviewerName)
}

// SetReviewStatus shows or hides a review (admin).
func (s *ReviewService) SetReviewStatus(ctx context.Context, id int, status models.ReviewStatus) error {
	if status != models.ReviewStatusPublished && status != models.ReviewStatusHidden {
		return fmt.Errorf("invalid status")
	}
	return s.repo.SetReviewStatus(ctx, id, status)
}

// DeleteReview removes a review (admin).
func (s *ReviewService) DeleteReview(ctx context.Context, id int) error {
	return s.repo.DeleteReview(ctx, id)
}
