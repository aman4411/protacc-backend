package handler

import (
	"net/http"
	"strconv"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/service"
	"github.com/aman4411/protacc-backend/utils"
	"github.com/gofiber/fiber/v2"
)

type ContactHandler struct {
	contactService service.ContactService
}

func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

// CreateContact creates a new contact message (public endpoint)
func (h *ContactHandler) CreateContact(c *fiber.Ctx) error {
	var req models.CreateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get client IP and User-Agent
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	contact, err := h.contactService.CreateContact(&req, ipAddress, userAgent)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create contact message: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.InfoLogger.Printf("Contact message created successfully - ID: %d, Email: %s", contact.ID, contact.Email)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Contact message sent successfully",
		"contact": contact,
	})
}

// GetContacts retrieves contact messages with filters (admin only)
func (h *ContactHandler) GetContacts(c *fiber.Ctx) error {
	var filters models.ContactFilters

	// Parse query parameters
	filters.Status = c.Query("status")
	filters.ServiceInterest = c.Query("service_interest")
	filters.DateFrom = c.Query("date_from")
	filters.DateTo = c.Query("date_to")
	filters.Search = c.Query("search")

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filters.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	contacts, totalCount, err := h.contactService.GetContacts(filters)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get contact messages: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve contact messages",
		})
	}

	// Calculate pagination info
	totalPages := (totalCount + filters.Limit - 1) / filters.Limit
	if filters.Limit == 0 {
		totalPages = 1
	}

	return c.JSON(fiber.Map{
		"contacts": contacts,
		"pagination": fiber.Map{
			"current_page": filters.Page,
			"total_pages":  totalPages,
			"total_count":  totalCount,
			"limit":        filters.Limit,
			"has_next":     filters.Page < totalPages,
			"has_previous": filters.Page > 1,
		},
	})
}

// GetContactByID retrieves a specific contact message (admin only)
func (h *ContactHandler) GetContactByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid contact ID",
		})
	}

	contact, err := h.contactService.GetContactByID(id)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get contact message: %v", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Contact message not found",
		})
	}

	return c.JSON(fiber.Map{
		"contact": contact,
	})
}

// UpdateContactStatus updates the status of a contact message (admin only)
func (h *ContactHandler) UpdateContactStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid contact ID",
		})
	}

	var req models.UpdateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get the admin user info from the context (set by auth middleware)
	userClaims := c.Locals("user")
	if userClaims != nil {
		if claims, ok := userClaims.(map[string]interface{}); ok {
			if email, exists := claims["email"].(string); exists {
				req.RespondedBy = &email
			}
		}
	}

	err = h.contactService.UpdateContactStatus(id, req)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update contact status: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.InfoLogger.Printf("Contact status updated successfully - ID: %d, Status: %s", id, req.Status)

	return c.JSON(fiber.Map{
		"message": "Contact status updated successfully",
	})
}

// DeleteContact deletes a contact message (admin only)
func (h *ContactHandler) DeleteContact(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid contact ID",
		})
	}

	err = h.contactService.DeleteContact(id)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to delete contact message: %v", err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.InfoLogger.Printf("Contact message deleted successfully - ID: %d", id)

	return c.JSON(fiber.Map{
		"message": "Contact message deleted successfully",
	})
}

// GetContactStats retrieves contact message statistics (admin only)
func (h *ContactHandler) GetContactStats(c *fiber.Ctx) error {
	stats, err := h.contactService.GetContactStats()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get contact stats: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve contact statistics",
		})
	}

	return c.JSON(fiber.Map{
		"stats": stats,
	})
}
