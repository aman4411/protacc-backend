package handler

import (
	"strconv"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/service"
	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	svc *service.ReviewService
}

func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

// SubmitReview creates or updates the authenticated user's review for a service.
func (h *ReviewHandler) SubmitReview(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var req struct {
		ServiceID int    `json:"service_id"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.ServiceID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "service_id is required"})
	}

	review, err := h.svc.SubmitReview(c.Context(), userID, req.ServiceID, req.Rating, req.Comment)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(review)
}

// GetServiceReviews returns published reviews + aggregate for a service (public).
func (h *ReviewHandler) GetServiceReviews(c *fiber.Ctx) error {
	serviceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid service ID"})
	}
	reviews, summary, err := h.svc.GetServiceReviews(c.Context(), serviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"reviews": reviews, "summary": summary})
}

// GetTopReviews returns highly-rated reviews for the homepage (public).
func (h *ReviewHandler) GetTopReviews(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "6"))
	reviews, err := h.svc.GetTopReviews(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"reviews": reviews})
}

// GetEligibility reports whether the authenticated user can review a service.
func (h *ReviewHandler) GetEligibility(c *fiber.Ctx) error {
	userID, ok := c.Locals("userId").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	serviceID, err := strconv.Atoi(c.Query("service_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid service_id"})
	}
	elig, err := h.svc.GetEligibility(c.Context(), userID, serviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(elig)
}

// ListReviews returns all reviews for admin moderation.
func (h *ReviewHandler) ListReviews(c *fiber.Ctx) error {
	reviews, err := h.svc.ListReviews(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"reviews": reviews})
}

// AdminCreateReview lets an admin publish a review with a custom name (admin).
func (h *ReviewHandler) AdminCreateReview(c *fiber.Ctx) error {
	var req struct {
		ServiceID    int    `json:"service_id"`
		Rating       int    `json:"rating"`
		Comment      string `json:"comment"`
		ReviewerName string `json:"reviewer_name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	review, err := h.svc.AdminCreateReview(c.Context(), req.ServiceID, req.Rating, req.Comment, req.ReviewerName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(review)
}

// UpdateReviewStatus shows/hides a review (admin).
func (h *ReviewHandler) UpdateReviewStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid review ID"})
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if err := h.svc.SetReviewStatus(c.Context(), id, models.ReviewStatus(req.Status)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Review updated"})
}

// DeleteReview removes a review (admin).
func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid review ID"})
	}
	if err := h.svc.DeleteReview(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Review deleted"})
}
