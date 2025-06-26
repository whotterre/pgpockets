package services

import (
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"

	"go.uber.org/zap"
)

type CardService interface {
	CreateCard(card *models.Card) error
}

type cardService struct {
	repository repositories.CardRepository
	logger     *zap.Logger
}

func NewCardService(repo repositories.CardRepository, appLogger *zap.Logger) *cardService {
	return &cardService{
		repository: repo,
		logger: appLogger,
	}
}


func (s *cardService) CreateCard(card *models.Card) error {
	if err := s.repository.CreateCard(card); err != nil {
		s.logger.Error("Failed to create card", zap.Error(err))
		return err
	}
	s.logger.Info("Card created successfully", zap.String("cardID", card.ID.String()))
	return nil
}
