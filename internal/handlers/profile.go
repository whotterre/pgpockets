package handlers

import (
	"encoding/json"
	"pgpockets/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	profileService services.ProfileService
	logger         *zap.Logger
}

func NewProfileHandler(profileService services.ProfileService, logger *zap.Logger) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         logger,
	}
}

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	profile, err := h.profileService.GetProfileByUserID(userID)
	if err != nil {
		h.logger.Error("Failed to get profile", zap.String("userID", userID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve profile",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"profile": profile,
	})
}

type ProfileUpdateRequest struct {
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	DateOfBirth *time.Time `json:"dob"`
	Address     string     `json:"address"`
	PhoneNumber string     `json:"phone_number"`
	Gender      string     `json:"gender"`
}

func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	var profileUpdate ProfileUpdateRequest
	if err := json.Unmarshal(c.Body(), &profileUpdate); err != nil {
		h.logger.Error("Failed to parse profile update request", zap.String("userID", userID.String()), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}
	profile, err := h.profileService.GetProfileByUserID(userID.String())
	if err != nil {
		h.logger.Error("Failed to get profile for update", zap.String("userID", userID.String()), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve profile for update",
			"details": err.Error(),
		})
	}
	profile.FirstName = profileUpdate.FirstName
	profile.LastName = profileUpdate.LastName
	profile.DateOfBirth = profileUpdate.DateOfBirth
	profile.Address = profileUpdate.Address
	profile.PhoneNumber = profileUpdate.PhoneNumber
	profile.Gender = profileUpdate.Gender
	if err := h.profileService.UpdateProfile(profile, userID); err != nil {
		h.logger.Error("Failed to update profile", zap.String("userID", userID.String()), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update profile",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}
