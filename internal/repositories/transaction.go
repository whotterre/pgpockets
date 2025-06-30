package repositories

import (
	"fmt"
	"pgpockets/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) (*models.Transaction, error)
	GetTransactionByID(id uuid.UUID) (*models.Transaction, error)
	GetAllTransactionsByUserID(
		userID uuid.UUID,
		limit int,  
		offset int,
	)(*[]models.Transaction, int64, error) 
	UpdateTransactionStatus(walletID uuid.UUID, newStatus string) error
	
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) (*models.Transaction, error) {
    if err := r.db.Create(transaction).Error; err != nil {
        return nil, err
    }
    return transaction, nil
}

// Gets a single transaction by it's id 
func (r *transactionRepository) GetTransactionByID(id uuid.UUID) (*models.Transaction, error){
	var transaction models.Transaction
	err := r.db.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err 
	}
	return &transaction, nil
}

// Updates a transaction's status
func (r *transactionRepository) UpdateTransactionStatus(walletID uuid.UUID, newStatus string) error {
	return r.db.Model(&models.Transaction{}).Where("id = ?", walletID).Update("status", newStatus).Error
}


// Gets all transactions made by a user
func (r *transactionRepository) GetAllTransactionsByUserID(
    userID uuid.UUID,
    limit int,  
    offset int,
) (*[]models.Transaction, int64, error) {
    var transactions []models.Transaction
    query := r.db.
        Model(&models.Transaction{}).
        Joins("JOIN wallets ON wallets.id = transactions.sender_wallet_id").
        Where("wallets.user_id = ?", userID).
        Order("transactions.made_at DESC")  
    
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
    }
    
    // 3. Apply pagination and fetch
    if err := query.Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
        return nil, 0, err
    }
    
    return &transactions, total, nil
}