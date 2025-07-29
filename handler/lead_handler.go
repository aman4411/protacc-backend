package handler

import (
	"strconv"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type LeadHandler struct {
	leadService *service.LeadService
	validate    *validator.Validate
}

func NewLeadHandler(leadService *service.LeadService) *LeadHandler {
	return &LeadHandler{
		leadService: leadService,
		validate:    validator.New(),
	}
}

// CreateLead creates a new business lead (public endpoint)
func (h *LeadHandler) CreateLead(c *fiber.Ctx) error {
	var req models.CreateLeadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	lead, err := h.leadService.CreateLead(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Lead created successfully",
		"lead":    lead,
	})
}

// GetLeads returns paginated list of leads with filters (admin only)
func (h *LeadHandler) GetLeads(c *fiber.Ctx) error {
	filters := models.LeadFilters{
		Status:     c.Query("status"),
		Priority:   c.Query("priority"),
		AssignedTo: c.Query("assigned_to"),
		Search:     c.Query("search"),
		Page:       1,
		Limit:      10,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filters.Page = p
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filters.Limit = l
		}
	}

	leads, total, err := h.leadService.GetLeads(c.Context(), filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"leads": leads,
		"pagination": fiber.Map{
			"page":  filters.Page,
			"limit": filters.Limit,
			"total": total,
			"pages": (total + filters.Limit - 1) / filters.Limit,
		},
	})
}

// GetLeadByID returns a specific lead by ID (admin only)
func (h *LeadHandler) GetLeadByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid lead ID",
		})
	}

	lead, err := h.leadService.GetLeadByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lead not found",
		})
	}

	return c.JSON(fiber.Map{
		"lead": lead,
	})
}

// UpdateLead updates a lead (admin only)
func (h *LeadHandler) UpdateLead(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid lead ID",
		})
	}

	var req models.UpdateLeadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	err = h.leadService.UpdateLead(c.Context(), id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Lead updated successfully",
	})
}

// DeleteLead deletes a lead (admin only)
func (h *LeadHandler) DeleteLead(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid lead ID",
		})
	}

	err = h.leadService.DeleteLead(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Lead deleted successfully",
	})
}

// GetLeadStats returns lead statistics (admin only)
func (h *LeadHandler) GetLeadStats(c *fiber.Ctx) error {
	stats, err := h.leadService.GetLeadStats(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"stats": stats,
	})
}
