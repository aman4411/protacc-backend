package handler

import (
	"strconv"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/service"
	"github.com/gofiber/fiber/v2"
)

type DeadlineHandler struct {
	svc *service.DeadlineService
}

func NewDeadlineHandler(svc *service.DeadlineService) *DeadlineHandler {
	return &DeadlineHandler{svc: svc}
}

// ListUpcoming returns active, upcoming deadlines (public).
func (h *DeadlineHandler) ListUpcoming(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "8"))
	items, err := h.svc.ListUpcoming(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"deadlines": items})
}

// List returns all deadlines (admin).
func (h *DeadlineHandler) List(c *fiber.Ctx) error {
	items, err := h.svc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"deadlines": items})
}

func (h *DeadlineHandler) Create(c *fiber.Ctx) error {
	var req models.Deadline
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	item, err := h.svc.Create(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

func (h *DeadlineHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid deadline ID"})
	}
	var req models.Deadline
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	req.ID = id
	item, err := h.svc.Update(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(item)
}

func (h *DeadlineHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid deadline ID"})
	}
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Deadline deleted"})
}
