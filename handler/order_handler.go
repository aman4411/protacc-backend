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

// PreviewCoupon validates a coupon against the user's cart and returns the
// resulting discount breakdown (no order created).
func (h *OrderHandler) PreviewCoupon(c *fiber.Ctx) error {
	userID := getOrderUserID(c)
	var req struct {
		Code string `json:"code"`
	}
	if err := c.BodyParser(&req); err != nil || req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Coupon code is required"})
	}
	amounts, coupon, err := h.svc.PreviewCoupon(c.Context(), userID, req.Code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"amounts": amounts,
		"coupon":  fiber.Map{"code": coupon.Code, "description": coupon.Description},
	})
}

func (h *OrderHandler) CreateOrderFromCart(c *fiber.Ctx) error {
	userID := getOrderUserID(c)

	var req struct {
		CouponCode string `json:"coupon_code"`
	}
	_ = c.BodyParser(&req) // body is optional (no coupon)

	order, err := h.svc.CreateOrderFromCart(c.Context(), userID, req.CouponCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// ordersResponse runs the filtered query and writes the standard paginated JSON.
func (h *OrderHandler) ordersResponse(c *fiber.Ctx, userID *string) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "")
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, total, err := h.svc.GetOrdersWithFilters(c.Context(), userID, page, limit, status, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

// GetOrders returns the authenticated user's own orders (user route). It always
// scopes to the current user, regardless of role — so an admin viewing their
// own profile sees only their orders.
func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	currentUserID := getOrderUserID(c)
	return h.ordersResponse(c, &currentUserID)
}

// GetAdminOrders returns all orders (admin route, role-gated), optionally
// filtered to a single user via ?user_id.
func (h *OrderHandler) GetAdminOrders(c *fiber.Ctx) error {
	var userID *string
	if uid := c.Query("user_id", ""); uid != "" {
		userID = &uid
	}
	return h.ordersResponse(c, userID)
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
