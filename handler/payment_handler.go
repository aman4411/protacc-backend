package handler

import (
	"strconv"

	"github.com/aman4411/protacc-backend/service"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	paymentSvc *service.PaymentService
	orderSvc   *service.OrderService
}

func NewPaymentHandler(paymentSvc *service.PaymentService, orderSvc *service.OrderService) *PaymentHandler {
	return &PaymentHandler{
		paymentSvc: paymentSvc,
		orderSvc:   orderSvc,
	}
}

// Helper function to get user ID from context
func getPaymentUserID(c *fiber.Ctx) string {
	return c.Locals("userId").(string)
}

// CreatePaymentOrder creates a Razorpay order for payment
func (h *PaymentHandler) CreatePaymentOrder(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	userID := getPaymentUserID(c)

	// Get order details
	order, err := h.orderSvc.GetOrderByID(c.Context(), orderID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	// Check if order is already paid
	if order.PaymentStatus {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order is already paid",
		})
	}

	// Create Razorpay order
	razorpayOrder, err := h.paymentSvc.CreateRazorpayOrder(c.Context(), order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create payment order: " + err.Error(),
		})
	}

	// Update order with Razorpay order ID
	razorpayOrderID, ok := razorpayOrder["id"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid response from payment gateway",
		})
	}

	err = h.orderSvc.UpdateOrderRazorpayOrderID(c.Context(), order.ID, razorpayOrderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Return response with order details needed for frontend
	return c.JSON(fiber.Map{
		"razorpay_order_id": razorpayOrderID,
		"amount":            int(order.BookingAmount * 100), // Amount in paise
		"currency":          "INR",
		"order_number":      order.OrderNumber,
		"order_id":          order.ID,
	})
}

// VerifyPayment verifies the payment signature and updates order status
func (h *PaymentHandler) VerifyPayment(c *fiber.Ctx) error {
	var req struct {
		RazorpayOrderID   string `json:"razorpay_order_id"`
		RazorpayPaymentID string `json:"razorpay_payment_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Verify payment signature
	err := h.paymentSvc.HandlePaymentSuccess(c.Context(), req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment verification failed: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Payment verified successfully",
	})
}

// HandleWebhook processes incoming webhooks from Razorpay
func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	// Get the raw body
	body := c.Body()

	// Get signature from header
	signature := c.Get("X-Razorpay-Signature")
	if signature == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing signature header",
		})
	}

	// Verify webhook signature
	if !h.paymentSvc.VerifyWebhookSignature(body, signature) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid webhook signature",
		})
	}

	// Parse webhook payload
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid webhook payload",
		})
	}

	// Process the webhook
	err := h.paymentSvc.ProcessWebhook(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process webhook: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

// GetPaymentStatus gets the current payment status of an order
func (h *PaymentHandler) GetPaymentStatus(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	userID := getPaymentUserID(c)

	// Get order details
	order, err := h.orderSvc.GetOrderByID(c.Context(), orderID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.JSON(fiber.Map{
		"order_id":            order.ID,
		"order_number":        order.OrderNumber,
		"payment_status":      order.PaymentStatus,
		"status":              order.Status,
		"total_amount":        order.TotalAmount,
		"booking_amount":      order.BookingAmount,
		"remaining_amount":    order.RemainingAmount,
		"razorpay_order_id":   order.RazorpayOrderID,
		"razorpay_payment_id": order.RazorpayPaymentID,
		"payment_method":      order.PaymentMethod,
	})
}
