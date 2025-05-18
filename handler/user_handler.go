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

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		utils.LogError("Failed to generate token", "error", err.Error(), "userId", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Set token in Authorization header
	c.Set("Authorization", "Bearer "+token)

	// Return success response with user data only
	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login handles user authentication
func (h *UserHandler) Login(c *fiber.Ctx) error {
	utils.LogInfo("Starting login process")

	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogError("Failed to parse login request body", "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		utils.LogError("Login validation failed", "error", err.Error(), "email", req.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Authenticate user
	user, err := h.userService.AuthenticateUser(c.Context(), &req)
	if err != nil {
		utils.LogError("Authentication failed", "error", err.Error(), "email", req.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		utils.LogError("Failed to generate token", "error", err.Error(), "userId", user.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	// Set token in Authorization header
	c.Set("Authorization", "Bearer "+token)

	// Return success response with user data only
	return c.Status(fiber.StatusOK).JSON(user)
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

// GetProfile returns the authenticated user's profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(user)
}

// GetUsers returns all users (admin only)
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(users)
}
