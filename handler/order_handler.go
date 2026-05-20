package handler

import (
	"strconv"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	svc      *service.OrderService
	docSvc   *service.OrderDocumentService
	validate *validator.Validate
}

func NewOrderHandler(svc *service.OrderService, docSvc *service.OrderDocumentService) *OrderHandler {
	return &OrderHandler{svc: svc, docSvc: docSvc, validate: validator.New()}
}

func getOrderRole(c *fiber.Ctx) string {
	if role, ok := c.Locals("role").(string); ok {
		return role
	}
	return "user"
}

// Helper function to get user ID from context
func getOrderUserID(c *fiber.Ctx) string {
	return c.Locals("userId").(string)
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	serviceID, err := strconv.Atoi(c.Params("serviceId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid service ID",
		})
	}

	userID := getOrderUserID(c)

	order, err := h.svc.CreateOrder(c.Context(), userID, serviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *OrderHandler) CreateOrderFromCart(c *fiber.Ctx) error {
	userID := getOrderUserID(c)

	order, err := h.svc.CreateOrderFromCart(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	// Parse query parameters for admin filtering
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "")
	search := c.Query("search", "")
	userIDParam := c.Query("user_id", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var userID *string

	// Check if user is admin
	currentUserID := getOrderUserID(c)

	// If user is admin, they can see all orders or filter by user_id
	if userIDParam != "" {
		userID = &userIDParam
	} else {
		// Regular users can only see their own orders
		userID = &currentUserID
	}

	orders, total, err := h.svc.GetOrdersWithFilters(c.Context(), userID, page, limit, status, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"pagination": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  (total + limit - 1) / limit,
		},
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var req struct {
		Status models.OrderStatus `json:"status"`
		Notes  string             `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := getOrderUserID(c)

	if err := h.svc.UpdateOrderStatus(c.Context(), orderID, req.Status, req.Notes, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated",
	})
}

func (h *OrderHandler) GetOrderStatusHistory(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	userID := getOrderUserID(c)
	role := getOrderRole(c)
	if err := h.docSvc.EnsureOrderAccess(c.Context(), orderID, userID, role); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	history, err := h.svc.GetOrderStatusHistory(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(history)
}

func (h *OrderHandler) GetOrderByNumber(c *fiber.Ctx) error {
	orderNumber := c.Params("orderNumber")
	userID := getOrderUserID(c)

	order, err := h.svc.GetOrderByNumber(c.Context(), orderNumber, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.JSON(order)
}

func (h *OrderHandler) GetOrderDocuments(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	userID := getOrderUserID(c)
	role := getOrderRole(c)

	docs, err := h.docSvc.ListDocuments(c.Context(), orderID, userID, role)
	if err != nil {
		if err.Error() == "access denied" || err.Error() == "order not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(h.docSvc.EnrichDocuments(docs))
}

func (h *OrderHandler) AddUserOrderDocument(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var req models.AddOrderDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	userID := getOrderUserID(c)
	doc, err := h.docSvc.AddUserDocument(c.Context(), orderID, userID, &req)
	if err != nil {
		if err.Error() == "access denied" || err.Error() == "order not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	enriched := h.docSvc.EnrichDocuments([]models.OrderDocument{*doc})
	return c.Status(fiber.StatusCreated).JSON(enriched[0])
}

func (h *OrderHandler) AddAdminOrderDocument(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var req models.AddOrderDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	adminID := getOrderUserID(c)
	doc, err := h.docSvc.AddAdminDocument(c.Context(), orderID, adminID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	enriched := h.docSvc.EnrichDocuments([]models.OrderDocument{*doc})
	return c.Status(fiber.StatusCreated).JSON(enriched[0])
}
