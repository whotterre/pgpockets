package handlers

import (
	"encoding/json"
	"pgpockets/internal/models"
	"pgpockets/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CardHandler struct {
	service   services.CardService
	logger    *zap.Logger
	validator *validator.Validate
}

func NewCardHandler(service services.CardService, logger *zap.Logger) *CardHandler {
	return &CardHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

type CreateCardRequest struct {
	UserID         string `json:"user_id" validate:"required,uuid"`
	LastFourDigits string `json:"last_four_digits" validate:"required,len=4"`
	CardToken      string `json:"card_token" validate:"required"`
	CardType       string `json:"card_type" validate:"required,oneof=visa mastercard amex discover"`
	ExpiryMonth    string `json:"expiry_month" validate:"required,len=2"`
	ExpiryYear     string `json:"expiry_year" validate:"required,len=4"`
	CardBrand      string `json:"card_brand" validate:"omitempty"`
	BankName       string `json:"bank_name" validate:"omitempty"`
	IsActive       bool   `json:"is_active" validate:"omitempty"`
}

func (h *CardHandler) CreateCard(c *fiber.Ctx) error {
	var card CreateCardRequest
	if err := json.Unmarshal(c.Body(), &card); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := c.Locals("userID").(uuid.UUID)
	card.UserID = userID.String()

	if err := h.validator.Struct(&card); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
		})
	}


	newCard := &models.Card{
		UserID:         userID,
		LastFourDigits: card.LastFourDigits,
		CardToken:      card.CardToken,
		CardType:       card.CardType,
		ExpiryMonth:    card.ExpiryMonth,
		ExpiryYear:     card.ExpiryYear,
		CardBrand:      card.CardBrand,
		BankName:       card.BankName,
		IsActive:       card.IsActive,
	}

	if err := h.service.CreateCard(newCard); err != nil {
		h.logger.Error("Failed to create card", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create card",
		})
	}
	h.logger.Info("Card added successfully", zap.String("cardID", newCard.ID.String()))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "Card created successfully",
		"card_id":      newCard.ID.String(),
		"user_id":      newCard.UserID.String(),
		"card_token":   newCard.CardToken,
		"card_type":    newCard.CardType,
		"expiry_month": newCard.ExpiryMonth,
		"expiry_year":  newCard.ExpiryYear,
		"card_brand":   newCard.CardBrand,
		"bank_name":    newCard.BankName,
		"is_active":    newCard.IsActive,
	})
}

func (h *CardHandler) RetrieveAllCards(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	cards, err := h.service.RetrieveAllCards(userID.String())
	if err != nil {
		h.logger.Error("Failed to retrieve cards", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cards",
		})
	}

	if len(cards) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No cards found for this user",
		})
	}

	h.logger.Info("Cards retrieved successfully", zap.Int("count", len(cards)))
	return c.Status(fiber.StatusOK).JSON(cards)
}

func (h *CardHandler) GetCardByID(c *fiber.Ctx) error {
	cardID := c.Params("cardID")
	if cardID == "" {
		h.logger.Error("Card ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Card ID is required",
		})
	}

	cardUUID, err := uuid.Parse(cardID)
	if err != nil {
		h.logger.Error("Invalid Card ID format", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Card ID format",
		})
	}
	card, err := h.service.GetCardByID(cardUUID.String())
	if err != nil {
		h.logger.Error("Failed to retrieve card", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve card",
		})
	}
	if card == (models.Card{}) {
		h.logger.Info("Card not found", zap.String("cardID", cardID))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Card not found",
		})
	}
	h.logger.Info("Card retrieved successfully", zap.String("cardID", cardID))
	return c.Status(fiber.StatusOK).JSON(card)
}