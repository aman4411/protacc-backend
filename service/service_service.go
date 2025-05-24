package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"

	"github.com/gosimple/slug"
)

type ServiceService struct {
	repo *repository.ServiceRepository
}

func NewServiceService(repo *repository.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

// Service Category Methods
func (s *ServiceService) CreateServiceCategory(ctx context.Context, category *models.ServiceCategory) error {
	category.Slug = slug.Make(category.Name)
	category.Status = models.ServiceStatusActive
	return s.repo.CreateServiceCategory(ctx, category)
}

func (s *ServiceService) GetServiceCategories(ctx context.Context) ([]models.ServiceCategory, error) {
	return s.repo.GetServiceCategories(ctx)
}

// Service Methods
func (s *ServiceService) CreateService(ctx context.Context, service *models.Service) error {
	service.Slug = slug.Make(service.Name)
	service.Status = models.ServiceStatusActive
	if service.BookingAmount == 0 {
		service.BookingAmount = 99.00 // Default booking amount
	}
	return s.repo.CreateService(ctx, service)
}

func (s *ServiceService) GetServices(ctx context.Context, categoryID *int) ([]models.Service, error) {
	return s.repo.GetServices(ctx, categoryID)
}

func (s *ServiceService) GetServiceByID(ctx context.Context, id int) (*models.Service, error) {
	return s.repo.GetServiceByID(ctx, id)
}

func (s *ServiceService) GetServiceBySlug(ctx context.Context, slug string) (*models.Service, error) {
	return s.repo.GetServiceBySlug(ctx, slug)
}

// Cart Methods
func (s *ServiceService) AddToCart(ctx context.Context, userID string, serviceID int) error {
	// Verify service exists and is active
	service, err := s.GetServiceByID(ctx, serviceID)
	if err != nil {
		return err
	}
	if service.Status != models.ServiceStatusActive {
		return fmt.Errorf("service is not available")
	}
	return s.repo.AddToCart(ctx, userID, serviceID)
}

func (s *ServiceService) GetCartItems(ctx context.Context, userID string) ([]models.CartItem, error) {
	return s.repo.GetCartItems(ctx, userID)
}

func (s *ServiceService) RemoveFromCart(ctx context.Context, userID string, serviceID int) error {
	return s.repo.RemoveFromCart(ctx, userID, serviceID)
}

// Order Methods
func (s *ServiceService) CreateOrder(ctx context.Context, userID string, serviceID int) (*models.Order, error) {
	// Get service details
	service, err := s.GetServiceByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	// Create order
	order := &models.Order{
		UserID:          userID,
		ServiceID:       serviceID,
		OrderNumber:     generateOrderNumber(),
		TotalAmount:     service.Price,
		BookingAmount:   service.BookingAmount,
		RemainingAmount: service.Price - service.BookingAmount,
		Status:          models.OrderStatusPendingPayment,
		PaymentStatus:   false,
	}

	err = s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Remove item from cart after successful order creation
	_ = s.RemoveFromCart(ctx, userID, serviceID)

	return order, nil
}

func (s *ServiceService) GetOrders(ctx context.Context, userID *string) ([]models.Order, error) {
	return s.repo.GetOrders(ctx, userID)
}

func (s *ServiceService) UpdateOrderStatus(ctx context.Context, orderID int, status models.OrderStatus, notes string, updatedBy string) error {
	return s.repo.UpdateOrderStatus(ctx, orderID, status, notes, updatedBy)
}

func (s *ServiceService) GetOrderStatusHistory(ctx context.Context, orderID int) ([]models.OrderStatusHistory, error) {
	return s.repo.GetOrderStatusHistory(ctx, orderID)
}

// Helper functions
func generateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("ORD%s", timestamp)
}
