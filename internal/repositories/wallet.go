package repositories

import (
	"pgpockets/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	CreateWallet(userID uuid.UUID) error 
	GetWalletBalance(userID uuid.UUID) (string, error)
	GetWalletByUserID(userID uuid.UUID) (*models.Wallet, error)
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


func (r *walletRepository) GetWalletBalance(userID uuid.UUID) (string, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return "", err
	}
	return wallet.Balance, nil
}

func (r *walletRepository) GetWalletByUserID(userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}