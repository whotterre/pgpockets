package repositories

import (
	"pgpockets/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BeneficiaryRepository interface {
	Create(beneficiary *models.Beneficiary) (*models.Beneficiary, error)
	GetBeneficiaries(userID string) ([]models.Beneficiary, error)
	DeleteBeneficiary(beneID uuid.UUID, userID uuid.UUID) error 
}

type beneficiaryRepository struct {
	db *gorm.DB
}

func NewBeneficiaryRepository(db *gorm.DB) BeneficiaryRepository {
	return &beneficiaryRepository{db: db}
}

// Create implements BeneficiaryRepository.
func (r *beneficiaryRepository) Create(beneficiary *models.Beneficiary) (*models.Beneficiary, error) {
	if err := r.db.Create(beneficiary).Error; err != nil {
		return nil, err
	}
	return beneficiary, nil
}

func (r *beneficiaryRepository) GetBeneficiaries(userID string) ([]models.Beneficiary, error) {
	var beneficiaries []models.Beneficiary
	if err := r.db.Where("user_id = ?", userID).Find(&beneficiaries).Error; err != nil {
		return nil, err
	}
	return beneficiaries, nil
}

// Deletes a beneficiary by their beneficiary ID
func (r *beneficiaryRepository) DeleteBeneficiary(beneID uuid.UUID, userID uuid.UUID) error {
	if err := r.db.Where("id = ? AND user_id = ?", beneID, userID).Delete(&models.Beneficiary{}).Error; err != nil {
		return err
	}

	return nil
}
