package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
	razorpay "github.com/razorpay/razorpay-go"
)

type PaymentService struct {
	razorpayClient *razorpay.Client
	orderRepo      *repository.OrderRepository
	keySecret      string
	webhookSecret  string
}

func NewPaymentService(orderRepo *repository.OrderRepository) *PaymentService {
	keyID := os.Getenv("RAZORPAY_KEY_ID")
	keySecret := os.Getenv("RAZORPAY_KEY_SECRET")
	webhookSecret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")

	if keyID == "" || keySecret == "" {
		// Don't panic in development - just log a warning
		fmt.Printf("WARNING: Razorpay credentials not found. Payment integration will be disabled.\n")
		return &PaymentService{
			razorpayClient: nil,
			orderRepo:      orderRepo,
			keySecret:      "",
			webhookSecret:  "",
		}
	}

	client := razorpay.NewClient(keyID, keySecret)

	return &PaymentService{
		razorpayClient: client,
		orderRepo:      orderRepo,
		keySecret:      keySecret,
		webhookSecret:  webhookSecret,
	}
}

// CreateRazorpayOrder creates a Razorpay order for payment
func (s *PaymentService) CreateRazorpayOrder(ctx context.Context, order *models.Order) (map[string]interface{}, error) {
	if s.razorpayClient == nil {
		return nil, fmt.Errorf("Razorpay client not initialized. Please set RAZORPAY_KEY_ID and RAZORPAY_KEY_SECRET")
	}

	// Convert amount to paise (Razorpay expects amount in smallest currency unit)
	amountInPaise := int(order.BookingAmount * 100)

	data := map[string]interface{}{
		"amount":          amountInPaise,
		"currency":        "INR",
		"receipt":         order.OrderNumber,
		"payment_capture": 1, // Auto capture payment
		"notes": map[string]interface{}{
			"order_id":     order.ID,
			"order_number": order.OrderNumber,
			"user_id":      order.UserID,
		},
	}

	razorpayOrder, err := s.razorpayClient.Order.Create(data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Razorpay order: %v", err)
	}

	return razorpayOrder, nil
}

// VerifyPaymentSignature verifies the Razorpay payment signature
func (s *PaymentService) VerifyPaymentSignature(orderID, paymentID, signature string) bool {
	message := orderID + "|" + paymentID
	h := hmac.New(sha256.New, []byte(s.keySecret))
	h.Write([]byte(message))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// HandlePaymentSuccess handles successful payment verification
func (s *PaymentService) HandlePaymentSuccess(ctx context.Context, orderID string, paymentID string, signature string) error {
	// First verify the signature
	if !s.VerifyPaymentSignature(orderID, paymentID, signature) {
		return fmt.Errorf("invalid payment signature")
	}

	// Extract order ID from the notes or find by order number
	// Note: orderID here is Razorpay order ID, we need to find our internal order
	order, err := s.orderRepo.GetOrderByRazorpayOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %v", err)
	}

	// Update order status
	err = s.orderRepo.UpdateOrderPaymentStatus(ctx, order.ID, true, paymentID, models.OrderStatusPaymentReceived)
	if err != nil {
		return fmt.Errorf("failed to update order payment status: %v", err)
	}

	return nil
}

// VerifyWebhookSignature verifies the webhook signature from Razorpay
func (s *PaymentService) VerifyWebhookSignature(payload []byte, signature string) bool {
	if s.webhookSecret == "" {
		return false
	}

	h := hmac.New(sha256.New, []byte(s.webhookSecret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// ProcessWebhook processes incoming webhook from Razorpay
func (s *PaymentService) ProcessWebhook(ctx context.Context, payload map[string]interface{}) error {
	event, ok := payload["event"].(string)
	if !ok {
		return fmt.Errorf("invalid webhook payload: missing event")
	}

	switch event {
	case "payment.captured":
		return s.handlePaymentCaptured(ctx, payload)
	case "payment.failed":
		return s.handlePaymentFailed(ctx, payload)
	case "order.paid":
		return s.handleOrderPaid(ctx, payload)
	default:
		// Log unhandled event but don't return error
		fmt.Printf("Unhandled webhook event: %s\n", event)
		return nil
	}
}

func (s *PaymentService) handlePaymentCaptured(ctx context.Context, payload map[string]interface{}) error {
	paymentData, ok := payload["payload"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payment data in webhook")
	}

	payment, ok := paymentData["payment"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payment object in webhook")
	}

	orderEntity, ok := payment["order_id"].(string)
	if !ok {
		return fmt.Errorf("missing order_id in payment")
	}

	paymentID, ok := payment["id"].(string)
	if !ok {
		return fmt.Errorf("missing payment id")
	}

	// Update order status
	order, err := s.orderRepo.GetOrderByRazorpayOrderID(ctx, orderEntity)
	if err != nil {
		return err
	}

	return s.orderRepo.UpdateOrderPaymentStatus(ctx, order.ID, true, paymentID, models.OrderStatusPaymentReceived)
}

func (s *PaymentService) handlePaymentFailed(ctx context.Context, payload map[string]interface{}) error {
	// Similar implementation for failed payments
	// For now, just log the failure
	fmt.Printf("Payment failed webhook received: %v\n", payload)
	return nil
}

func (s *PaymentService) handleOrderPaid(ctx context.Context, payload map[string]interface{}) error {
	// Handle when order is fully paid
	fmt.Printf("Order paid webhook received: %v\n", payload)
	return nil
}
