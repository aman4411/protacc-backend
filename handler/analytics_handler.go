package handler

import (
	"strconv"

	"github.com/aman4411/protacc-backend/service"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	svc *service.AnalyticsService
}

func NewAnalyticsHandler(svc *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

// GetRevenueAnalytics returns revenue data for charts
func (h *AnalyticsHandler) GetRevenueAnalytics(c *fiber.Ctx) error {
	period := c.Query("period", "30") // Default to 30 days
	days, err := strconv.Atoi(period)
	if err != nil || days <= 0 {
		days = 30
	}

	analytics, err := h.svc.GetRevenueAnalytics(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(analytics)
}

// GetOrderAnalytics returns order statistics
func (h *AnalyticsHandler) GetOrderAnalytics(c *fiber.Ctx) error {
	period := c.Query("period", "30")
	days, err := strconv.Atoi(period)
	if err != nil || days <= 0 {
		days = 30
	}

	analytics, err := h.svc.GetOrderAnalytics(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(analytics)
}

// GetUserAnalytics returns user growth and registration data
func (h *AnalyticsHandler) GetUserAnalytics(c *fiber.Ctx) error {
	period := c.Query("period", "30")
	days, err := strconv.Atoi(period)
	if err != nil || days <= 0 {
		days = 30
	}

	analytics, err := h.svc.GetUserAnalytics(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(analytics)
}

// GetServiceAnalytics returns popular services and category performance
func (h *AnalyticsHandler) GetServiceAnalytics(c *fiber.Ctx) error {
	period := c.Query("period", "30")
	days, err := strconv.Atoi(period)
	if err != nil || days <= 0 {
		days = 30
	}

	analytics, err := h.svc.GetServiceAnalytics(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(analytics)
}

// GetOverallMetrics returns key performance indicators
func (h *AnalyticsHandler) GetOverallMetrics(c *fiber.Ctx) error {
	period := c.Query("period", "30")
	days, err := strconv.Atoi(period)
	if err != nil || days <= 0 {
		days = 30
	}

	metrics, err := h.svc.GetOverallMetrics(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(metrics)
}

// GetRecentActivity returns recent system activity
func (h *AnalyticsHandler) GetRecentActivity(c *fiber.Ctx) error {
	limit := c.Query("limit", "10")
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 10
	}

	activity, err := h.svc.GetRecentActivity(c.Context(), limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(activity)
}
