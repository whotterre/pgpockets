package services

import (
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"
)

type BeneficiaryService interface {
	Create(beneficiary *models.Beneficiary) (*models.Beneficiary, error)
	GetBeneficiaries(userID string) ([]models.Beneficiary, error)
}

type beneficiaryService struct {
	repo repositories.BeneficiaryRepository
}


func NewBeneficiaryService(repo repositories.BeneficiaryRepository) BeneficiaryService {
	return &beneficiaryService{repo: repo}
}

// Create implements BeneficiaryService.
func (s *beneficiaryService) Create(beneficiary *models.Beneficiary) (*models.Beneficiary, error) {
	newBeneficiary, err := s.repo.Create(beneficiary)
	if err != nil {
		return nil, err
	}

	return newBeneficiary, nil

}

func (s *beneficiaryService) GetBeneficiaries(userID string) ([]models.Beneficiary, error) {
	beneficiaries, err := s.repo.GetBeneficiaries(userID)
	if err != nil {
		return nil, err
	}
	return beneficiaries, nil
}

