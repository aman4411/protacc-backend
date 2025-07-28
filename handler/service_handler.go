package handler

import (
	"strconv"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/service"
	"github.com/gofiber/fiber/v2"
)

type ServiceHandler struct {
	svc *service.ServiceService
}

func NewServiceHandler(svc *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{svc: svc}
}

// Helper function to get user ID from context
func getUserID(c *fiber.Ctx) string {
	return c.Locals("userId").(string)
}

// Service Category Handlers
func (h *ServiceHandler) CreateServiceCategory(c *fiber.Ctx) error {
	// TODO: Implement category creation if needed
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "Service category creation not implemented",
	})
}

func (h *ServiceHandler) GetServiceCategories(c *fiber.Ctx) error {
	categories, err := h.svc.GetAllCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(categories)
}

// Service Handlers
func (h *ServiceHandler) CreateService(c *fiber.Ctx) error {
	// TODO: Implement service creation if needed
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "Service creation not implemented",
	})
}

func (h *ServiceHandler) GetServices(c *fiber.Ctx) error {
	categorySlug := c.Query("category")

	services, err := h.svc.GetServicesByCategory(c.Context(), categorySlug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(services)
}

func (h *ServiceHandler) GetServiceByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid service ID",
		})
	}

	service, err := h.svc.GetServiceByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(service)
}

func (h *ServiceHandler) GetServiceBySlug(c *fiber.Ctx) error {
	categorySlug := c.Params("category")
	serviceSlug := c.Params("slug")
	if categorySlug == "" || serviceSlug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid service or category slug",
		})
	}

	service, err := h.svc.GetServiceBySlug(c.Context(), categorySlug, serviceSlug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(service)
}

// Cart Handlers
func (h *ServiceHandler) AddToCart(c *fiber.Ctx) error {
	serviceID, err := strconv.Atoi(c.Params("serviceId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid service ID",
		})
	}

	userID := getUserID(c)

	if err := h.svc.AddToCart(c.Context(), userID, serviceID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Service added to cart",
	})
}

func (h *ServiceHandler) GetCartItems(c *fiber.Ctx) error {
	userID := getUserID(c)

	items, err := h.svc.GetCartItems(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(items)
}

func (h *ServiceHandler) RemoveFromCart(c *fiber.Ctx) error {
	serviceID, err := strconv.Atoi(c.Params("serviceId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid service ID",
		})
	}

	userID := getUserID(c)

	if err := h.svc.RemoveFromCart(c.Context(), userID, serviceID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Service removed from cart",
	})
}

// Order Handlers
func (h *ServiceHandler) CreateOrder(c *fiber.Ctx) error {
	serviceID, err := strconv.Atoi(c.Params("serviceId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid service ID",
		})
	}

	userID := getUserID(c)

	order, err := h.svc.CreateOrder(c.Context(), userID, serviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *ServiceHandler) CreateOrderFromCart(c *fiber.Ctx) error {
	userID := getUserID(c)

	order, err := h.svc.CreateOrderFromCart(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func (h *ServiceHandler) GetOrders(c *fiber.Ctx) error {
	var userID *string

	// If user is admin, they can see all orders or filter by user_id
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID = &userIDStr
	} else {
		// Regular users can only see their own orders
		id := getUserID(c)
		userID = &id
	}

	orders, err := h.svc.GetOrders(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(orders)
}

func (h *ServiceHandler) UpdateOrderStatus(c *fiber.Ctx) error {
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

	userID := getUserID(c)

	if err := h.svc.UpdateOrderStatus(c.Context(), orderID, req.Status, req.Notes, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated",
	})
}

func (h *ServiceHandler) GetOrderStatusHistory(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("orderId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	history, err := h.svc.GetOrderStatusHistory(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(history)
}
