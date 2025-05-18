package handler

import (
	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/service"
	"github.com/aman4411/protacc-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// Signup handles user registration
func (h *UserHandler) Signup(c *fiber.Ctx) error {
	utils.LogInfo("Starting signup process")

	var req models.UserRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogError("Failed to parse signup request body", "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	utils.LogInfo("Parsed signup request", "email", req.Email, "firstName", req.FirstName)

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		utils.LogError("Signup validation failed", "error", err.Error(), "email", req.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	utils.LogInfo("Signup validation passed", "email", req.Email)

	// Create user
	user, err := h.userService.CreateUser(c.Context(), &req)
	if err != nil {
		utils.LogError("Failed to create user", "error", err.Error(), "email", req.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	utils.LogInfo("User created successfully", "userId", user.ID, "email", user.Email)

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Signup successful. Please verify your email.",
		"user":    user,
	})
}

// VerifyEmail handles email verification
func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	var req models.EmailVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	if err := h.userService.VerifyEmail(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email verified successfully",
	})
}
