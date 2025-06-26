package repositories

import (
	"pgpockets/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	CreateWallet(userID uuid.UUID) error 
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{
		db: db,
	}
}


func (r *walletRepository) CreateWallet(userID uuid.UUID) error {
	wallet := &models.Wallet{
		UserID: userID,
		Balance: "0.0",
	}

	if err := r.db.Create(wallet).Error; err != nil {
		return err
	}
	return nil
}
