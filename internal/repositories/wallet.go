package repositories

import (
	"pgpockets/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WalletRepository interface {
	CreateWallet(userID uuid.UUID) error
	GetWalletBalance(userID uuid.UUID) (string, error)
	GetWalletByID(walletID uuid.UUID) (*models.Wallet, error)
	GetWalletByEmail(email string) (*models.Wallet, error)
	GetWalletByUserID(userID uuid.UUID) (*models.Wallet, error)
	GetBalancesForAllWallets(userID uuid.UUID) ([]map[string]string, error)
	UpdateWalletBalance(walletID uuid.UUID, newBalance string) error
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
		UserID:  userID,
		Balance: decimal.Zero,
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
	return wallet.Balance.String(), nil
}

func (r *walletRepository) GetWalletByID(walletID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("id = ?", walletID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetWalletByUserID(userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) GetWalletByEmail(email string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Joins("JOIN users ON users.id = wallets.user_id").
		Where("users.email = ?", email).
		First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) UpdateWalletBalance(walletID uuid.UUID, newBalance string) error {
	if err := r.db.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", newBalance).Error; err != nil {
		return err
	}
	return nil
}
func (r *walletRepository) GetBalancesForAllWallets(userID uuid.UUID) ([]map[string]string, error) {
	var wallets []models.Wallet
	if err := r.db.Select("balance, currency").
		Where("user_id = ?", userID).
		Find(&wallets).Error; err != nil {
		return nil, err
	}

	var balances []map[string]string
	for _, wallet := range wallets {
		balances = append(balances, map[string]string{
			"balance":  wallet.Balance.String(),
			"currency": wallet.Currency,
		})
	}
	return balances, nil
}
