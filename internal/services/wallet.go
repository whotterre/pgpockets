package services

import (
	"pgpockets/internal/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WalletService interface {

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