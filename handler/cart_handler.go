package handler

import (
	"strconv"

	"github.com/aman4411/protacc-backend/service"
	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	svc *service.CartService
}

func NewCartHandler(svc *service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

// Helper function to get user ID from context
func getCartUserID(c *fiber.Ctx) string {
	return c.Locals("userId").(string)
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	serviceID, err := strconv.Atoi(c.Params("serviceId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid service ID",
		})
	}

	userID := getCartUserID(c)

	err = h.svc.AddToCart(c.Context(), userID, serviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item added to cart",
	})
}

func (h *CartHandler) GetCartItems(c *fiber.Ctx) error {
	userID := getCartUserID(c)

	items, err := h.svc.GetCartItems(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(items)
}

func (h *CartHandler) RemoveFromCart(c *fiber.Ctx) error {
	serviceID, err := strconv.Atoi(c.Params("serviceId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid service ID",
		})
	}

	userID := getCartUserID(c)

	err = h.svc.RemoveFromCart(c.Context(), userID, serviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item removed from cart",
	})
}
