package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	serviceRepo *repository.ServiceRepository
	cartRepo    *repository.CartRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, serviceRepo *repository.ServiceRepository, cartRepo *repository.CartRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		serviceRepo: serviceRepo,
		cartRepo:    cartRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, serviceID int) (*models.Order, error) {
	// Get service details
	service, err := s.serviceRepo.GetServiceByID(ctx, serviceID)
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

	// Generate order number
	timestamp := time.Now().Format("20060102150405")
	orderNumber := fmt.Sprintf("ORD%s", timestamp)

	order := &models.Order{
		UserID:          userID,
		OrderNumber:     orderNumber,
		TotalAmount:     service.Price,
		BookingAmount:   service.BookingAmount,
		RemainingAmount: service.Price - service.BookingAmount,
		Status:          models.OrderStatusPendingPayment,
		PaymentStatus:   false,
		Items:           []models.OrderItem{orderItem},
	}

	err = s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Remove item from cart after successful order creation (if it exists)
	_ = s.cartRepo.RemoveFromCart(ctx, userID, serviceID)

	return order, nil
}

func (s *OrderService) CreateOrderFromCart(ctx context.Context, userID string) (*models.Order, error) {
	return s.orderRepo.CreateOrderFromCart(ctx, userID)
}

func (s *OrderService) GetOrders(ctx context.Context, userID *string) ([]models.Order, error) {
	return s.orderRepo.GetOrders(ctx, userID)
}

func (s *OrderService) GetOrdersWithFilters(ctx context.Context, userID *string, page, limit int, status, search string) ([]models.Order, int, error) {
	return s.orderRepo.GetOrdersWithFilters(ctx, userID, page, limit, status, search)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int, status models.OrderStatus, notes string, updatedBy string) error {
	return s.orderRepo.UpdateOrderStatus(ctx, orderID, status, notes, updatedBy)
}

func (s *OrderService) GetOrderStatusHistory(ctx context.Context, orderID int) ([]models.OrderStatusHistory, error) {
	return s.orderRepo.GetOrderStatusHistory(ctx, orderID)
}

// Add these methods to the OrderService

// GetOrderByID retrieves an order by ID for a specific user
func (s *OrderService) GetOrderByID(ctx context.Context, orderID int, userID string) (*models.Order, error) {
	return s.orderRepo.GetOrderByID(ctx, orderID, userID)
}

// UpdateOrderRazorpayOrderID updates the Razorpay order ID for an order
func (s *OrderService) UpdateOrderRazorpayOrderID(ctx context.Context, orderID int, razorpayOrderID string) error {
	return s.orderRepo.UpdateOrderRazorpayOrderID(ctx, orderID, razorpayOrderID)
}

func (s *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string, userID string) (*models.Order, error) {
	return s.orderRepo.GetOrderByNumber(ctx, orderNumber, userID)
}
