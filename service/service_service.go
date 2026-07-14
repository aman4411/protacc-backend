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

// Cache key prefix + TTLs for the public, read-heavy service catalog. Every
// service/category mutation busts the whole "services:" family (see invalidate).
const (
	svcCachePrefix = "services:"
	svcCacheTTL    = 60 * time.Minute
	svcCatCacheTTL = 60 * time.Minute
)

type ServiceService struct {
	repo  *repository.ServiceRepository
	cache *cache.Cache
}

func NewServiceService(repo *repository.ServiceRepository, c *cache.Cache) *ServiceService {
	return &ServiceService{
		repo:  repo,
		cache: c,
	}
}

// invalidate drops all cached service/category reads after a write.
func (s *ServiceService) invalidate() {
	if s.cache != nil {
		s.cache.InvalidatePrefix(svcCachePrefix)
	}
}

func (s *ServiceService) GetAllCategories(ctx context.Context) ([]models.ServiceCategory, error) {
	return cache.Load(s.cache, svcCachePrefix+"categories", svcCatCacheTTL, func() ([]models.ServiceCategory, error) {
		return s.repo.GetServiceCategories(ctx)
	})
}

func (s *ServiceService) CreateServiceCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error) {
	res, err := s.repo.CreateServiceCategory(ctx, category)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *ServiceService) UpdateServiceCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error) {
	res, err := s.repo.UpdateServiceCategory(ctx, category)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *ServiceService) DeleteServiceCategory(ctx context.Context, categoryID int) error {
	err := s.repo.DeleteServiceCategory(ctx, categoryID)
	if err == nil {
		s.invalidate()
	}
	return err
}

func (s *ServiceService) GetServices(ctx context.Context, categoryID *int, categorySlug string) ([]models.Service, error) {
	id := "all"
	if categoryID != nil {
		id = fmt.Sprintf("%d", *categoryID)
	}
	key := fmt.Sprintf("%slist:%s:%s", svcCachePrefix, id, categorySlug)
	return cache.Load(s.cache, key, svcCacheTTL, func() ([]models.Service, error) {
		return s.repo.GetServices(ctx, categoryID, categorySlug)
	})
}

func (s *ServiceService) GetServicesByCategory(ctx context.Context, categorySlug string) ([]models.Service, error) {
	return s.GetServices(ctx, nil, categorySlug)
}

func (s *ServiceService) GetServiceBySlug(ctx context.Context, serviceSlug string) (*models.Service, error) {
	return cache.Load(s.cache, svcCachePrefix+"slug:"+serviceSlug, svcCacheTTL, func() (*models.Service, error) {
		return s.repo.GetServiceBySlug(ctx, serviceSlug)
	})
}

func (s *ServiceService) GetServiceByID(ctx context.Context, serviceID int) (*models.Service, error) {
	key := fmt.Sprintf("%sid:%d", svcCachePrefix, serviceID)
	return cache.Load(s.cache, key, svcCacheTTL, func() (*models.Service, error) {
		return s.repo.GetServiceByID(ctx, serviceID)
	})
}

func (s *ServiceService) SearchServices(ctx context.Context, query string) ([]models.Service, error) {
	if query == "" {
		return []models.Service{}, nil
	}
	return s.repo.SearchServices(ctx, query)
}

func (s *ServiceService) CreateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	normalizeService(service)
	res, err := s.repo.CreateService(ctx, service)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *ServiceService) UpdateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	normalizeService(service)
	res, err := s.repo.UpdateService(ctx, service)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

// normalizeService cleans up the FAQ list (drops blank entries, never nil) and
// defaults the price type so downstream storage/rendering is consistent.
func normalizeService(service *models.Service) {
	cleaned := make([]models.ServiceFAQ, 0, len(service.FAQs))
	for _, f := range service.FAQs {
		q := strings.TrimSpace(f.Question)
		a := strings.TrimSpace(f.Answer)
		if q == "" || a == "" {
			continue
		}
		cleaned = append(cleaned, models.ServiceFAQ{Question: q, Answer: a})
	}
	service.FAQs = cleaned

	if strings.TrimSpace(service.PriceType) == "" {
		service.PriceType = "fixed"
	}
}

func (s *ServiceService) DeleteService(ctx context.Context, serviceID int) error {
	err := s.repo.DeleteService(ctx, serviceID)
	if err == nil {
		s.invalidate()
	}
	return err
}

// Priority Management Methods
func (s *ServiceService) UpdateCategoryPriority(ctx context.Context, categoryID int, priority int) error {
	err := s.repo.UpdateCategoryPriority(ctx, categoryID, priority)
	if err == nil {
		s.invalidate()
	}
	return err
}

func (s *ServiceService) UpdateServicePriority(ctx context.Context, serviceID int, priority int) error {
	err := s.repo.UpdateServicePriority(ctx, serviceID, priority)
	if err == nil {
		s.invalidate()
	}
	return err
}

// Helper functions
func generateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("ORD%s", timestamp)
}
