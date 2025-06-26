package repositories

import (
	"pgpockets/internal/models"

	"gorm.io/gorm"
)

type CardRepository interface {
	CreateCard(card *models.Card) error
	RetrieveAllCards(userID string) ([]*models.Card, error)
	GetCardByID(cardID string) (models.Card, error) 
}


type cardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{
		db: db,
	}
}

func (r *cardRepository) CreateCard(card *models.Card) error {
	if err := r.db.Create(card).Error; err != nil {
		return err
	}
	return nil
}

// Retrieves all cards owned by a particular user
func (r *cardRepository) RetrieveAllCards(userID string) ([]*models.Card, error) {
	var cards []*models.Card
	if err := r.db.Where("user_id = ?", userID).Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *cardRepository) GetCardByID(cardID string) (models.Card, error) {
	var card models.Card
	if err := r.db.Where("id = ?", cardID).First(&card).Error; err != nil {
		return models.Card{}, err
	}
	return card, nil
}