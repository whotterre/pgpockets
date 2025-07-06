package handlers

import (
	"pgpockets/internal/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletHandler struct {
	walletService services.WalletService
	logger        *zap.Logger
	apiKey        string
}

func NewWalletHandler(
	walletService services.WalletService,
	logger *zap.Logger,
	apiKey string,
) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
		logger:        logger,
		apiKey:        apiKey,
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

func (w *WalletHandler) ChangeWalletCurrency(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	desiredCurrency := strings.ToUpper(c.Params("desiredCurrency")) 

	newBalance, err := w.walletService.ChangeWalletCurrency(userID, desiredCurrency, w.apiKey)
	if err != nil {
		w.logger.Error("Failed to convert currency", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve wallet balance",
			"details": err.Error(),
		})
	}
	w.logger.Info("Successfully converted currency")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Successfully converted currency",
		"balance":     newBalance,
		"newCurrency": desiredCurrency,
	})

}


func (h *WalletHandler) GetBalancesForAllWallets(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	
	balances, err := h.walletService.GetBalancesForAllWallets(userID)
	if err != nil {
		h.logger.Error("Something went wrong while getting balances for wallets")
		return err
	}

	h.logger.Info("Successfully retrieved wallet balances")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully retrieved wallet balances",
		"balances": balances,
	})

}