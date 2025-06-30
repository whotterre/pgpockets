package services

import (
	"errors"
	"fmt"
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TransactionService interface {
	TransferFunds(
		userID, senderWalletID, recieverWalletID uuid.UUID,
		amount decimal.Decimal,
		currency, description string,
	) (*models.Transaction, error)
	GetTransactionHistory(
		userID uuid.UUID,
		limit, offset int,
	) (*[]models.Transaction, int64, error)
}

type transactionService struct {
	txnRepo    repositories.TransactionRepository
	walletRepo repositories.WalletRepository
	appLogger  *zap.Logger
	db         *gorm.DB
}

func NewTransactionService(
	txnRepo repositories.TransactionRepository,
	logger *zap.Logger,
	walletRepo repositories.WalletRepository,
	db *gorm.DB,
) *transactionService {
	return &transactionService{
		txnRepo:    txnRepo,
		walletRepo: walletRepo,
		appLogger:  logger,
		db:         db,
	}
}

func (s *transactionService) TransferFunds(
	userID, senderWalletID, recieverWalletID uuid.UUID,
	amount decimal.Decimal,
	currency, description string,
) (*models.Transaction, error) {
	var txn *models.Transaction

	// Start a database transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		senderWallet, err := s.walletRepo.GetWalletByID(senderWalletID)
		if err != nil {
			return errors.New("sender wallet not found")
		}
		receiverWallet, err := s.walletRepo.GetWalletByID(recieverWalletID)
		if err != nil {
			return errors.New("receiver wallet not found")
		}

		// Check if the initiator is actually the owner of the wallet o!!!
		err = s.txnRepo.VerifyOwnership(userID, senderWalletID)
		if err != nil {
			return errors.New("user is not owner of wallet")
		}
		
		senderBalance := senderWallet.Balance
		if senderBalance.LessThan(amount) {
			return errors.New("insufficient funds")
		}

		// Create new transaction record
		txn = &models.Transaction{
			SenderWalletID:   &senderWalletID,
			ReceiverWalletID: &recieverWalletID,
			Amount:           amount.String(),
			Currency:         currency,
			TransactionType:  "transfer",
			Status:           "pending",
			Description:      s.generateDescription(description, senderWallet.UserID, receiverWallet.UserID),
			ReferenceID:      s.generateReferenceID(),
		}

		newTxn, err := s.txnRepo.CreateTransaction(txn)
		if err != nil {
			return errors.New("failed to create transaction")
		}
		txn = newTxn

		newSenderBalance := senderBalance.Sub(amount).String()
		if err := s.walletRepo.UpdateWalletBalance(senderWalletID, newSenderBalance); err != nil {
			s.appLogger.Error("Failed to update sender's balance", zap.String("because", err.Error()))
			return errors.New("failed to update sender's balance")
		}

		newRecieverBalance := receiverWallet.Balance.Add(amount)
		if err := s.walletRepo.UpdateWalletBalance(recieverWalletID, newRecieverBalance.String()); err != nil {
			s.appLogger.Error("Failed to update reciever's balance", zap.String("because", err.Error()))
			return errors.New("failed to update reciever's balance")
		}

		if err := s.txnRepo.UpdateTransactionStatus(newTxn.ID, "completed"); err != nil {
			s.appLogger.Error("Failed to update transaction status", zap.String("because", err.Error()))
			return errors.New("failed to update transaction status")
		}

		return nil
	})

	if err != nil {
		s.appLogger.Error("failed to transfer funds", zap.Error(err))
		return nil, err
	}
	s.appLogger.Info(
		"funds transferred successfully",
		zap.String("sender", senderWalletID.String()),
		zap.String("receiver", recieverWalletID.String()),
		zap.String("amount", amount.String()))
	return txn, nil
}

func (s *transactionService) GetTransactionHistory(
	userID uuid.UUID,
	limit, offset int,
	) (*[]models.Transaction, int64, error) {
	history, count, err := s.txnRepo.GetAllTransactionsByUserID(userID, limit, offset)
	if err != nil {
		s.appLogger.Error("failed to get transaction history", zap.Error(err))
		return nil, 0, err
	}
	return history, count, nil
}


func (s *transactionService) generateDescription(
	description string,
	senderID, recieverID uuid.UUID,
) string {
	if description != "" {
		return description
	}
	return fmt.Sprintf("Transfer from %s to %s", senderID.String()[:8], recieverID.String()[:8])
}

func (s *transactionService) generateReferenceID() string {
	return fmt.Sprintf("Txn_%d_%s", time.Now().Unix(), uuid.New())
}
