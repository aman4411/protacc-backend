package service

import (
	"context"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type AnalyticsService struct {
	analyticsRepo *repository.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo *repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
	}
}

// Service Methods
func (s *AnalyticsService) GetRevenueAnalytics(ctx context.Context, days int) (*models.RevenueAnalytics, error) {
	return s.analyticsRepo.GetRevenueAnalytics(ctx, days)
}

func (s *AnalyticsService) GetOrderAnalytics(ctx context.Context, days int) (*models.OrderAnalytics, error) {
	return s.analyticsRepo.GetOrderAnalytics(ctx, days)
}

func (s *AnalyticsService) GetUserAnalytics(ctx context.Context, days int) (*models.UserAnalytics, error) {
	return s.analyticsRepo.GetUserAnalytics(ctx, days)
}

func (s *AnalyticsService) GetServiceAnalytics(ctx context.Context, days int) (*models.ServiceAnalytics, error) {
	return s.analyticsRepo.GetServiceAnalytics(ctx, days)
}

func (s *AnalyticsService) GetOverallMetrics(ctx context.Context, days int) (*models.OverallMetrics, error) {
	return s.analyticsRepo.GetOverallMetrics(ctx, days)
}

func (s *AnalyticsService) GetRecentActivity(ctx context.Context, limit int) (*models.RecentActivity, error) {
	return s.analyticsRepo.GetRecentActivity(ctx, limit)
}
