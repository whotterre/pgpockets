package services

import (
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"

	"go.uber.org/zap"
)

type CardService interface {
	CreateCard(card *models.Card) error
	RetrieveAllCards(userID string) ([]*models.Card, error)
	GetCardByID(cardID string) (models.Card, error)
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

func (s *cardService) RetrieveAllCards(userID string) ([]*models.Card, error){
	cards, err := s.repository.RetrieveAllCards(userID)
	if err != nil {
		s.logger.Error("Failed to retrieve cards", zap.Error(err))
		return []*models.Card{}, err
	}
	s.logger.Info("Cards retrieved successfully", zap.Int("count", len(cards)))
	return cards, nil
}

func (s *cardService) GetCardByID(cardID string) (models.Card, error) {
	card, err := s.repository.GetCardByID(cardID)
	if err != nil {
		s.logger.Error("Failed to retrieve card by ID", zap.Error(err))
		return models.Card{}, err
	}
	s.logger.Info("Card retrieved successfully", zap.String("cardID", card.ID.String()))
	return card, nil
}