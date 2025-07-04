package repositories

import (
	"pgpockets/internal/models"

	"gorm.io/gorm"
)

type BeneficiaryRepository interface {
	Create(beneficiary *models.Beneficiary) (*models.Beneficiary, error)
	GetBeneficiaries(userID string) ([]models.Beneficiary, error)
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
