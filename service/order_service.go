package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

type OrderService struct {
	orderRepo     *repository.OrderRepository
	serviceRepo   *repository.ServiceRepository
	cartRepo      *repository.CartRepository
	couponService *CouponService
	mailService   *MailService
}

func NewOrderService(orderRepo *repository.OrderRepository, serviceRepo *repository.ServiceRepository, cartRepo *repository.CartRepository, couponService *CouponService, mailService *MailService) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		serviceRepo:   serviceRepo,
		cartRepo:      cartRepo,
		couponService: couponService,
		mailService:   mailService,
	}
}

// notifyOrderEmailAsync loads the full order + customer and emails them, off the
// request path — email failures are logged and never block the order flow.
// Shared by OrderService (completion) and PaymentService (order placed on payment).
func notifyOrderEmailAsync(orderRepo *repository.OrderRepository, mail *MailService, orderID int, completed bool) {
	if mail == nil {
		return
	}
	go func() {
		order, err := orderRepo.GetOrderForNotification(context.Background(), orderID)
		if err != nil {
			fmt.Printf("order email: failed to load order %d: %v\n", orderID, err)
			return
		}
		if completed {
			if err := mail.SendOrderCompletedEmail(order); err != nil {
				fmt.Printf("order email: completed send failed for order %d: %v\n", orderID, err)
			}
			return
		}

		// Order booked: notify the customer and the ProtAcc team (FROM_EMAIL).
		if err := mail.SendOrderPlacedEmail(order); err != nil {
			fmt.Printf("order email: placed send failed for order %d: %v\n", orderID, err)
		}
		if err := mail.SendOrderPlacedAdminEmail(order); err != nil {
			fmt.Printf("order email: admin notify failed for order %d: %v\n", orderID, err)
		}
	}()
}

// PreviewCoupon validates a coupon against the user's current cart and returns
// the resulting amount breakdown (without creating an order).
func (s *OrderService) PreviewCoupon(ctx context.Context, userID, code string) (models.DiscountedAmounts, *models.Coupon, error) {
	lines, err := s.orderRepo.GetCartLines(ctx, userID)
	if err != nil {
		return models.DiscountedAmounts{}, nil, err
	}
	if len(lines) == 0 {
		return models.DiscountedAmounts{}, nil, fmt.Errorf("your cart is empty")
	}
	var subtotal, bookingBase float64
	for _, l := range lines {
		subtotal += l.Price
		bookingBase += l.BookingAmount
	}
	coupon, discount, err := s.couponService.Validate(ctx, code, userID, lines)
	if err != nil {
		return models.DiscountedAmounts{}, nil, err
	}
	amounts := models.ComputeDiscountedAmounts(subtotal, bookingBase, discount, coupon.ApplicationMode)
	return amounts, coupon, nil
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
		Status:          models.OrderStatusPendingBookingPayment,
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

func (s *OrderService) CreateOrderFromCart(ctx context.Context, userID, couponCode string) (*models.Order, error) {
	discount := 0.0
	mode := models.CouponModeFinal
	normalizedCode := ""

	if couponCode != "" {
		lines, err := s.orderRepo.GetCartLines(ctx, userID)
		if err != nil {
			return nil, err
		}
		coupon, d, err := s.couponService.Validate(ctx, couponCode, userID, lines)
		if err != nil {
			return nil, err
		}
		discount = d
		mode = coupon.ApplicationMode
		normalizedCode = coupon.Code
	}

	return s.orderRepo.CreateOrderFromCart(ctx, userID, normalizedCode, discount, mode)
}

func (s *OrderService) GetOrders(ctx context.Context, userID *string) ([]models.Order, error) {
	return s.orderRepo.GetOrders(ctx, userID)
}

func (s *OrderService) GetOrdersWithFilters(ctx context.Context, userID *string, page, limit int, status, search string) ([]models.Order, int, error) {
	return s.orderRepo.GetOrdersWithFilters(ctx, userID, page, limit, status, search)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID int, status models.OrderStatus, notes string, updatedBy string) error {
	if err := s.orderRepo.UpdateOrderStatus(ctx, orderID, status, notes, updatedBy); err != nil {
		return err
	}

	// Email the customer the completed-order summary when the order is marked complete.
	if status == models.OrderStatusCompleted {
		notifyOrderEmailAsync(s.orderRepo, s.mailService, orderID, true)
	}

	return nil
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
