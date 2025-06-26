package handlers

import (
	"pgpockets/internal/services"

	"github.com/gofiber/fiber/v2"
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

	return c.Status(fiber.StatusOK).JSON(profile)
}
