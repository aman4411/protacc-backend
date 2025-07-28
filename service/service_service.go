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

func (s *ServiceService) GetServicesByCategory(ctx context.Context, categorySlug string) ([]models.Service, error) {
	// For now, get all services and filter by category if needed
	// This could be optimized with a proper repository method later
	return s.repo.GetServices(ctx, nil)
}

func (s *ServiceService) GetServiceBySlug(ctx context.Context, serviceSlug string) (*models.Service, error) {
	return s.repo.GetServiceBySlug(ctx, serviceSlug)
}

func (s *ServiceService) GetServiceByID(ctx context.Context, serviceID int) (*models.Service, error) {
	return s.repo.GetServiceByID(ctx, serviceID)
}

func (s *ServiceService) AddToCart(ctx context.Context, userID string, serviceID int) error {
	return s.repo.AddToCart(ctx, userID, serviceID)
}

func (s *ServiceService) GetCartItems(ctx context.Context, userID string) ([]models.CartItem, error) {
	return s.repo.GetCartItems(ctx, userID)
}

func (s *ServiceService) RemoveFromCart(ctx context.Context, userID string, serviceID int) error {
	return s.repo.RemoveFromCart(ctx, userID, serviceID)
}

// Order Methods
func (s *ServiceService) CreateOrderFromCart(ctx context.Context, userID string) (*models.Order, error) {
	return s.repo.CreateOrderFromCart(ctx, userID)
}

func (s *ServiceService) CreateOrder(ctx context.Context, userID string, serviceID int) (*models.Order, error) {
	// Get service details
	service, err := s.GetServiceByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	// Create order with single item
	orderItem := models.OrderItem{
		ServiceID:     serviceID,
		Quantity:      1,
		Price:         service.Price,
		BookingAmount: service.BookingAmount,
	}

	order := &models.Order{
		UserID:          userID,
		OrderNumber:     generateOrderNumber(),
		TotalAmount:     service.Price,
		BookingAmount:   service.BookingAmount,
		RemainingAmount: service.Price - service.BookingAmount,
		Status:          models.OrderStatusPendingPayment,
		PaymentStatus:   false,
		Items:           []models.OrderItem{orderItem},
	}

	err = s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Remove item from cart after successful order creation (if it exists)
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
