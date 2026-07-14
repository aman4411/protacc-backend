package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aman4411/protacc-backend/cache"
	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

// Cache prefix + TTL for public review reads (per-service + homepage top). Any
// review write busts the family. Eligibility / admin list stay uncached.
const (
	reviewCachePrefix = "reviews:"
	reviewCacheTTL    = 60 * time.Minute
)

type ReviewService struct {
	repo  *repository.ReviewRepository
	cache *cache.Cache
}

func NewReviewService(repo *repository.ReviewRepository, c *cache.Cache) *ReviewService {
	return &ReviewService{repo: repo, cache: c}
}

func (s *ReviewService) invalidate() {
	if s.cache != nil {
		s.cache.InvalidatePrefix(reviewCachePrefix)
	}
}

// serviceReviewsResult bundles the multi-return of GetServiceReviews for caching.
type serviceReviewsResult struct {
	Reviews []models.Review
	Summary models.ReviewSummary
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
	res, err := s.repo.UpsertReview(ctx, rev)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

// GetServiceReviews returns published reviews and the aggregate for a service.
func (s *ReviewService) GetServiceReviews(ctx context.Context, serviceID int) ([]models.Review, models.ReviewSummary, error) {
	key := fmt.Sprintf("%sservice:%d", reviewCachePrefix, serviceID)
	res, err := cache.Load(s.cache, key, reviewCacheTTL, func() (serviceReviewsResult, error) {
		reviews, summary, err := s.repo.GetServiceReviews(ctx, serviceID)
		return serviceReviewsResult{Reviews: reviews, Summary: summary}, err
	})
	if err != nil {
		return nil, models.ReviewSummary{}, err
	}
	return res.Reviews, res.Summary, nil
}

// GetTopReviews returns up to limit highly-rated reviews for the homepage.
func (s *ReviewService) GetTopReviews(ctx context.Context, limit int) ([]models.Review, error) {
	if limit < 1 || limit > 20 {
		limit = 6
	}
	key := fmt.Sprintf("%stop:%d", reviewCachePrefix, limit)
	return cache.Load(s.cache, key, reviewCacheTTL, func() ([]models.Review, error) {
		return s.repo.GetTopReviews(ctx, limit)
	})
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
	res, err := s.repo.AdminCreateReview(ctx, serviceID, rating, strings.TrimSpace(comment), reviewerName)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

// SetReviewStatus shows or hides a review (admin).
func (s *ReviewService) SetReviewStatus(ctx context.Context, id int, status models.ReviewStatus) error {
	if status != models.ReviewStatusPublished && status != models.ReviewStatusHidden {
		return fmt.Errorf("invalid status")
	}
	err := s.repo.SetReviewStatus(ctx, id, status)
	if err == nil {
		s.invalidate()
	}
	return err
}

// DeleteReview removes a review (admin).
func (s *ReviewService) DeleteReview(ctx context.Context, id int) error {
	err := s.repo.DeleteReview(ctx, id)
	if err == nil {
		s.invalidate()
	}
	return err
}
