package handlers

import (
	"encoding/json"
	"log"
	"pgpockets/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService services.AuthService
	logger      *zap.Logger
	validator   *validator.Validate
}

func NewAuthHandler(authService services.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		validator:   validator.New(),
	}
}

type UserRegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	DateOfBirth string `json:"dob" validate:"omitempty,datetime=01-02-2006"`
	PhoneNumber string `json:"phone_number" validate:"omitempty"`
	Address     string `json:"address" validate:"omitempty"`
	Gender      string `json:"gender" validate:"omitempty"`
}

func (h *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	var req UserRegisterRequest
	// Validate the request body
	if err := json.Unmarshal(c.Body(), &req); err != nil {
		h.logger.Warn("Failed to parse request body",
			zap.Error(err),
			zap.String("raw_body", string(c.Body())),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}
	log.Print(req)
	// Validate the request data
	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	user, err := h.authService.Register(req.Email, req.Password, req.FirstName, req.LastName, req.DateOfBirth, req.PhoneNumber, req.Address, req.Gender)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			h.logger.Warn("User already exists", zap.String("email", req.Email))
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "User with this email already exists",
			})
		}
		h.logger.Error("Failed to register user", zap.Error(err))
		if err == gorm.ErrInvalidTransaction {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}
	h.logger.Info("User registered successfully", zap.String("email", user.Email))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"email": user.Email,
		},
	})
}

func (h *AuthHandler) LoginUser(c *fiber.Ctx) error {
	// Get request body
	// Validate request body
	// Check if user exists
	// Verify password
	// Generate JWT token and refresh token
	// Store in sessions table and return to client
	return nil 
}
