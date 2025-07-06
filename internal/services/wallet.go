package services

import (
	"errors"
	"fmt"
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"
	"pgpockets/internal/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrNoCurrencyProvided = errors.New("no currency provided")
	ErrInvalidCurrency    = errors.New("invalid currency format")
)

type WalletService interface {
	CreateWallet(userID uuid.UUID) error
	GetWalletBalance(userID uuid.UUID) (string, error)
	GetWalletByEmail(email string) (*models.Wallet, error)
	ChangeWalletCurrency(userID uuid.UUID, currency, apiKey string) (string, error)
	GetBalancesForAllWallets(userIDStr string) ([]map[string]string, error)
}

type walletService struct {
	walletRepo repositories.WalletRepository
	logger     *zap.Logger
	db         *gorm.DB
}
func NewWalletService(walletRepo repositories.WalletRepository, logger *zap.Logger, db *gorm.DB) *walletService {
	return &walletService{
		walletRepo: walletRepo,
		logger:     logger,
		db:         db,
	}
}

func (s *walletService) CreateWallet(userID uuid.UUID) error {
	if err := s.walletRepo.CreateWallet(userID); err != nil {
		s.logger.Error("Failed to create wallet", zap.Error(err))
		return err
	}
	return nil
}

func (s *walletService) GetWalletBalance(userID uuid.UUID) (string, error) {
	balance, err := s.walletRepo.GetWalletBalance(userID)
	if err != nil {
		s.logger.Error("Failed to get wallet balance", zap.Error(err))
		return "", err
	}
	return balance, nil
}

func (s *walletService) GetWalletByEmail(email string) (*models.Wallet, error) {
	wallet, err := s.walletRepo.GetWalletByEmail(email)
	if err != nil {
		s.logger.Error("Failed to retrieve wallet by email", zap.Error(err))
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) ChangeWalletCurrency(userID uuid.UUID, currency string, apiKey string) (string, error) {
	// Validate the currency format
	// Check if the wallet exists for the user
	// Get the current exchange rate or just call the api to do it for us
	// Update the wallet's currency in the database
	// Return the response

	if currency == "" {
		s.logger.Error("Currency cannot be empty")
		return "", ErrNoCurrencyProvided
	}
	if !utils.IsValidCurrencyFormat(currency) {
		s.logger.Error("Invalid currency format", zap.String("currency", currency))
		return "", ErrInvalidCurrency
	}

	// Check if wallet exists for the user
	wallet, err := s.walletRepo.GetWalletByUserID(userID)
	if err != nil {
		s.logger.Error("Failed to retrieve wallet", zap.Error(err))
		return "", err
	}

	// Get current exchange rate
	rateA, rateB, err := utils.GetExchangeRatesForPair(wallet.Currency, currency, apiKey, s.logger)
	if err != nil {
		s.logger.Error("Failed to get exchange rates", zap.Error(err))
		return "", err
		
	}
	wallet.Currency = currency
	if err := s.walletRepo.UpdateWalletBalance(userID, wallet.Balance.String()); err != nil {
		s.logger.Error("Failed to update wallet balance ", zap.String("because", err.Error()))
		return "", err
	}

	decimalBalance := wallet.Balance
	// Make the conversion 
	/*
		The operation below is equal to 
		(RateB / RateA) * oldCurrencyValueAmount 
	*/
	newBalance := rateB.Div(rateA).Mul(decimalBalance)

	return newBalance.String(), nil
}

func (s *walletService) GetBalancesForAllWallets(userIDStr string) ([]map[string]string, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to convert user ID from string to uuid while getting balances for a wallets belonging to user with %s", userIDStr))
		return nil, err
	}

	balances, err := s.walletRepo.GetBalancesForAllWallets(userID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get balances for all wallets belonging to user with user %s", userIDStr))
		return nil, err 
	}
	s.logger.Info("Successfully retrieved all balances for all wallets belonging to the user")
	return balances, nil 
}