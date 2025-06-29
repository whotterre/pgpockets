package handlers

import (
	"pgpockets/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type TransactionHandler struct {
	transactionService services.TransactionService
	logger             *zap.Logger
	validator          *validator.Validate
}

func NewTransactionHandler(
	transactionService services.TransactionService,
	logger *zap.Logger,
) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		logger:             logger,
		validator:          validator.New(),
	}
}

type MakeTransferRequest struct {
	SenderWalletID   string `json:"sender_wallet_id" validate:"required"`
	ReceiverWalletID string `json:"rwi" validate:"required"`
	Amount           string `json:"amount" validate:"required,numeric"`
	Currency         string `json:"currency" validate:"required,len=3"`
	TransactionType  string `json:"transaction_type" validate:"required"`
	Description      string `json:"description" validate:"required,max=255"`
}

func (h *TransactionHandler) TransferFunds(c *fiber.Ctx) error {
	// Get relevant data from the request body
	var req MakeTransferRequest
	
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body for fund transfer", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate request body before making service calls
	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Type casting
	senderUUID, err := uuid.Parse(req.SenderWalletID)
	if err != nil {
		h.logger.Error("Something went wrong while trying to convert sender id to uuid")
	}
	recieverUUID, err := uuid.Parse(req.ReceiverWalletID)
	if err != nil {
		h.logger.Error("Something went wrong while trying to convert reciever id to uuid", zap.Error(err))
	}
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		h.logger.Error("Something went wrong while trying to convert amount to decimal", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid amount format",
		})
	}

	// Actually make the transfer via service call
	txnDetails, err := h.transactionService.TransferFunds(
		senderUUID,
		recieverUUID,
		amount,
		req.Currency,
		req.Description)
	if err != nil {
		h.logger.Error("Failed to transfer funds", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to transfer funds",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Funds transferred successfully",
		"transaction": txnDetails,
	})
}

