package handler

import (
	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/service"
	"github.com/gofiber/fiber/v2"
)

type SettingsHandler struct {
	svc *service.SettingsService
}

func NewSettingsHandler(svc *service.SettingsService) *SettingsHandler {
	return &SettingsHandler{svc: svc}
}

// GetAllSettings retrieves all system settings (admin only)
func (h *SettingsHandler) GetAllSettings(c *fiber.Ctx) error {
	settings, err := h.svc.GetAllSettings(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(settings)
}

// GetPublicSettings retrieves only public settings
func (h *SettingsHandler) GetPublicSettings(c *fiber.Ctx) error {
	settings, err := h.svc.GetPublicSettings(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(settings)
}

// GetSettingsByCategory retrieves settings grouped by category
func (h *SettingsHandler) GetSettingsByCategory(c *fiber.Ctx) error {
	categories, err := h.svc.GetSettingsByCategory(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(categories)
}

// GetSetting retrieves a specific setting
func (h *SettingsHandler) GetSetting(c *fiber.Ctx) error {
	category := c.Params("category")
	key := c.Params("key")

	if category == "" || key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category and key are required",
		})
	}

	setting, err := h.svc.GetSetting(c.Context(), category, key)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Setting not found",
		})
	}

	return c.JSON(setting)
}

// UpdateSetting updates a specific setting
func (h *SettingsHandler) UpdateSetting(c *fiber.Ctx) error {
	category := c.Params("category")
	key := c.Params("key")

	if category == "" || key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category and key are required",
		})
	}

	var req struct {
		Value string `json:"value"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.svc.UpdateSetting(c.Context(), category, key, req.Value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Setting updated successfully",
	})
}

// UpdateMultipleSettings updates multiple settings in a transaction
func (h *SettingsHandler) UpdateMultipleSettings(c *fiber.Ctx) error {
	var req models.SettingsUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.Settings) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No settings provided for update",
		})
	}

	err := h.svc.UpdateMultipleSettings(c.Context(), req.Settings)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":       "Settings updated successfully",
		"updated_count": len(req.Settings),
	})
}

// CreateSetting creates a new setting
func (h *SettingsHandler) CreateSetting(c *fiber.Ctx) error {
	var setting models.SystemSetting

	if err := c.BodyParser(&setting); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if setting.Category == "" || setting.SettingKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category and setting_key are required",
		})
	}

	err := h.svc.CreateSetting(c.Context(), &setting)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(setting)
}

// DeleteSetting deletes a specific setting
func (h *SettingsHandler) DeleteSetting(c *fiber.Ctx) error {
	category := c.Params("category")
	key := c.Params("key")

	if category == "" || key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category and key are required",
		})
	}

	err := h.svc.DeleteSetting(c.Context(), category, key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Setting deleted successfully",
	})
}

// TestEmailSettings tests the email configuration
func (h *SettingsHandler) TestEmailSettings(c *fiber.Ctx) error {
	// This would test the email configuration by sending a test email
	// For now, we'll just return a success message
	return c.JSON(fiber.Map{
		"message": "Email test functionality will be implemented based on email service integration",
		"status":  "pending",
	})
}

// ResetToDefaults resets settings to their default values
func (h *SettingsHandler) ResetToDefaults(c *fiber.Ctx) error {
	category := c.Query("category")

	if category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category parameter is required",
		})
	}

	// This would reset settings to their default values
	// For now, we'll return a message indicating this functionality needs implementation
	return c.JSON(fiber.Map{
		"message":  "Reset to defaults functionality will be implemented",
		"category": category,
		"status":   "pending",
	})
}
