package services

import (
	"pgpockets/internal/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletService interface {
	CreateWallet(userID uuid.UUID) error
	GetWalletBalance(userID uuid.UUID) (string, error)
}

type walletService struct {
	walletRepo repositories.WalletRepository
	logger     *zap.Logger
}

func NewWalletService(walletRepo repositories.WalletRepository, logger *zap.Logger) *walletService {
	return &walletService{
		walletRepo: walletRepo,
		logger:     logger,
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

