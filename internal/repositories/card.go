package repositories

import (
	"pgpockets/internal/models"

	"gorm.io/gorm"
)

type CardRepository interface {
	CreateCard(card *models.Card) error
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
