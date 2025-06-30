package handlers

import (
	"net/http"
	"pgpockets/internal/services"
	"strconv"

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
	userID := c.Locals("userID")
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
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.Error("Something went wrong while trying to convert user id to uuid", zap.Error(err))
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
		userUUID,
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

func (h *TransactionHandler) GetUserTransactionHistory(c *fiber.Ctx) error {
	userUUID := c.Locals("userID").(uuid.UUID)
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid limit value", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit value",
		})
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid offset value", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid offset value",
		})
	}

	history, count, err := h.transactionService.GetTransactionHistory(userUUID, int(limit), int(offset))
	if err != nil {
		h.logger.Error("failed to get transaction history because", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get transaction history",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Transaction history loaded successfully",
		"history": history,
		"count":   count,
	})

}

func (h *TransactionHandler) GetTransactionsInDateRange(c *fiber.Ctx) error {
	userUUID := c.Locals("userID").(uuid.UUID)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid limit value", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid limit value",
		})
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid offset value", zap.Error(err))
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid offset value",
		})
	}

	txnHistory, err := h.transactionService.GetTransactionsInDateRange(
		userUUID,
		startDate,
		endDate,
		int(limit),
		int(offset),
	)
	if err != nil {
		h.logger.Error("failed to get transactions in date range because", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get transactions in date range",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Transactions in date range loaded successfully",
		"history": txnHistory,
	})
}

func (h *TransactionHandler) GetTransactionByID(c *fiber.Ctx) error {
	txnIDStr := c.Params("txnID")
	if txnIDStr == "" {
		h.logger.Error("Transaction ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Transaction ID is required",
		})
	}

	txnID, err := uuid.Parse(txnIDStr)
	if err != nil {
		h.logger.Error("Invalid transaction ID format", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction ID format",
		})
	}

	txn, err := h.transactionService.RetrieveSingleTransaction(txnID)
	if err != nil {
		h.logger.Error("Failed to retrieve transaction by ID", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Transaction retrieved successfully",
		"transaction": txn,
	})
}