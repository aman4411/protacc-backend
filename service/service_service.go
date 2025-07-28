package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type ServiceService struct {
	repo *repository.ServiceRepository
}

func NewServiceService(repo *repository.ServiceRepository) *ServiceService {
	return &ServiceService{
		repo: repo,
	}
}

func (s *ServiceService) GetAllCategories(ctx context.Context) ([]models.ServiceCategory, error) {
	return s.repo.GetServiceCategories(ctx)
}

func (s *ServiceService) CreateServiceCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error) {
	return s.repo.CreateServiceCategory(ctx, category)
}

func (s *ServiceService) UpdateServiceCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error) {
	return s.repo.UpdateServiceCategory(ctx, category)
}

func (s *ServiceService) DeleteServiceCategory(ctx context.Context, categoryID int) error {
	return s.repo.DeleteServiceCategory(ctx, categoryID)
}

func (s *ServiceService) GetServices(ctx context.Context, categoryID *int, categorySlug string) ([]models.Service, error) {
	return s.repo.GetServices(ctx, categoryID, categorySlug)
}

func (s *ServiceService) GetServicesByCategory(ctx context.Context, categorySlug string) ([]models.Service, error) {
	return s.repo.GetServices(ctx, nil, categorySlug)
}

func (s *ServiceService) GetServiceBySlug(ctx context.Context, serviceSlug string) (*models.Service, error) {
	return s.repo.GetServiceBySlug(ctx, serviceSlug)
}

func (s *ServiceService) GetServiceByID(ctx context.Context, serviceID int) (*models.Service, error) {
	return s.repo.GetServiceByID(ctx, serviceID)
}

func (s *ServiceService) SearchServices(ctx context.Context, query string) ([]models.Service, error) {
	if query == "" {
		return []models.Service{}, nil
	}
	return s.repo.SearchServices(ctx, query)
}

func (s *ServiceService) CreateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	return s.repo.CreateService(ctx, service)
}

func (s *ServiceService) UpdateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	return s.repo.UpdateService(ctx, service)
}

func (s *ServiceService) DeleteService(ctx context.Context, serviceID int) error {
	return s.repo.DeleteService(ctx, serviceID)
}

// Helper functions
func generateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("ORD%s", timestamp)
}
