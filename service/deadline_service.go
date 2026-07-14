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

// Cache prefix + TTL for the public upcoming-deadlines widget. Admin List() is
// never cached.
const (
	deadlineCachePrefix = "deadlines:"
	deadlineCacheTTL    = 60 * time.Minute
)

type DeadlineService struct {
	repo  *repository.DeadlineRepository
	cache *cache.Cache
}

func NewDeadlineService(repo *repository.DeadlineRepository, c *cache.Cache) *DeadlineService {
	return &DeadlineService{repo: repo, cache: c}
}

func (s *DeadlineService) invalidate() {
	if s.cache != nil {
		s.cache.InvalidatePrefix(deadlineCachePrefix)
	}
}

func (s *DeadlineService) List(ctx context.Context) ([]models.Deadline, error) {
	return s.repo.List(ctx)
}

func (s *DeadlineService) ListUpcoming(ctx context.Context, limit int) ([]models.Deadline, error) {
	if limit < 1 || limit > 50 {
		limit = 8
	}
	key := fmt.Sprintf("%supcoming:%d", deadlineCachePrefix, limit)
	return cache.Load(s.cache, key, deadlineCacheTTL, func() ([]models.Deadline, error) {
		return s.repo.ListUpcoming(ctx, limit)
	})
}

func (s *DeadlineService) validate(d *models.Deadline) error {
	d.Title = strings.TrimSpace(d.Title)
	if d.Title == "" {
		return fmt.Errorf("title is required")
	}
	if d.DueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}
	return nil
}

func (s *DeadlineService) Create(ctx context.Context, d *models.Deadline) (*models.Deadline, error) {
	if err := s.validate(d); err != nil {
		return nil, err
	}
	res, err := s.repo.Create(ctx, d)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *DeadlineService) Update(ctx context.Context, d *models.Deadline) (*models.Deadline, error) {
	if err := s.validate(d); err != nil {
		return nil, err
	}
	res, err := s.repo.Update(ctx, d)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *DeadlineService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		s.invalidate()
	}
	return err
}
