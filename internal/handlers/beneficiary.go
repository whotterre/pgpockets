package handlers

import (
	"pgpockets/internal/models"
	"pgpockets/internal/services"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BeneficiaryHandler struct {
	beneficiaryService services.BeneficiaryService
	walletService      services.WalletService
	logger             *zap.Logger
	validator          *validator.Validate
}

func NewBeneficiaryHandler(beneficiaryService services.BeneficiaryService, logger *zap.Logger) *BeneficiaryHandler {
	return &BeneficiaryHandler{
		beneficiaryService: beneficiaryService,
		logger:             logger,
		validator:          validator.New(),
	}
}

func (h *BeneficiaryHandler) GetBeneficiaries(c *fiber.Ctx) error {
	userID := c.Params("userID")
	beneficiaries, err := h.beneficiaryService.GetBeneficiaries(userID)
	if err != nil {
		h.logger.Error("Failed to get beneficiaries", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	h.logger.Info("Successfully retrieved beneficiaries", zap.Int("count", len(beneficiaries)))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Beneficiaries retrieved successfully",
		"beneficiaries": beneficiaries,
	})
}

type CreateBeneficiaryRequest struct {
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	RecipientEmail string    `json:"recipient_id"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

func (h *BeneficiaryHandler) AddBeneficiary(c *fiber.Ctx) error {
	var req CreateBeneficiaryRequest
	userID := c.Locals("userID").(uuid.UUID)
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get recipient wallet by email
	wallet, err := h.walletService.GetWalletByEmail(req.RecipientEmail)
	if err != nil {
		h.logger.Error("Failed to get wallet by email", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get wallet",
		})
	}

	newBeneficiary := models.Beneficiary{
		UserID:      userID,
		Description: req.Description,
		WalletID:    wallet.ID,
		Wallet:      *wallet,
	}

	newBen, err := h.beneficiaryService.Create(&newBeneficiary)
	if err != nil {
		h.logger.Error("Failed to create beneficiary", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create beneficiary",
		})
	}
	h.logger.Info("Beneficiary added successfully", zap.String("beneficiaryID", newBen.ID.String()))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":        "Beneficiary created successfully",
		"beneficiary_id": newBen.ID.String(),
		"user_id":        newBen.UserID.String(),
		"description":    newBen.Description,
		"created_at":     newBen.CreatedAt.Format(time.RFC3339),
		"wallet_id":      newBen.WalletID.String(),
		"wallet": fiber.Map{
			"id":       newBen.Wallet.ID.String(),
			"balance":  newBen.Wallet.Balance,
			"currency": newBen.Wallet.Currency,
		},
	})

}
