package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type DeadlineService struct {
	repo *repository.DeadlineRepository
}

func NewDeadlineService(repo *repository.DeadlineRepository) *DeadlineService {
	return &DeadlineService{repo: repo}
}

func (s *DeadlineService) List(ctx context.Context) ([]models.Deadline, error) {
	return s.repo.List(ctx)
}

func (s *DeadlineService) ListUpcoming(ctx context.Context, limit int) ([]models.Deadline, error) {
	if limit < 1 || limit > 50 {
		limit = 8
	}
	return s.repo.ListUpcoming(ctx, limit)
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
	return s.repo.Create(ctx, d)
}

func (s *DeadlineService) Update(ctx context.Context, d *models.Deadline) (*models.Deadline, error) {
	if err := s.validate(d); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, d)
}

func (s *DeadlineService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
