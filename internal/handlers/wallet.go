package handlers

import (
	"pgpockets/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletHandler struct {
	walletService services.WalletService
	logger        *zap.Logger
}

func NewWalletHandler(walletService services.WalletService, logger *zap.Logger) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
		logger:        logger,
	}
}

func (w *WalletHandler) GetWalletBalance(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	balance, err := w.walletService.GetWalletBalance(userID)
	if err != nil {
		w.logger.Error("Failed to get wallet balance", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve wallet balance",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"balance": balance,
	})
}
