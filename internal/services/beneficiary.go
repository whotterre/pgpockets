package services

import (
	"fmt"
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BeneficiaryService interface {
	Create(beneficiary *models.Beneficiary) (*models.Beneficiary, error)
	GetBeneficiaries(userID string) ([]models.Beneficiary, error)
	DeleteBeneficiary(userIDStr, beneIDStr string) error
}

type beneficiaryService struct {
	repo repositories.BeneficiaryRepository
	logger *zap.Logger
}


func NewBeneficiaryService(repo repositories.BeneficiaryRepository, appLogger *zap.Logger) BeneficiaryService {
	return &beneficiaryService{
		repo: repo,
		logger: appLogger,
	}
}

// Create implements BeneficiaryService.
func (s *beneficiaryService) Create(beneficiary *models.Beneficiary) (*models.Beneficiary, error) {
	newBeneficiary, err := s.repo.Create(beneficiary)
	if err != nil {
		s.logger.Error("Failed to create beneficiary")
		return nil, err
	}
	s.logger.Info("Successfully created new beneficiary")
	return newBeneficiary, nil

}

func (s *beneficiaryService) GetBeneficiaries(userID string) ([]models.Beneficiary, error) {
	beneficiaries, err := s.repo.GetBeneficiaries(userID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to get beneficiaries for user with id %s", userID))
		return nil, err
	}
	s.logger.Error(fmt.Sprintf("Successfully got beneficiaries for user with id %s", userID))
	return beneficiaries, nil
}

func (s *beneficiaryService) DeleteBeneficiary(userIDStr, beneIDStr string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.logger.Error("Failed to get parse user ID string to uuid format")
		return err 
	}

	beneID, err := uuid.Parse(beneIDStr)
	if err != nil {
		s.logger.Error("Failed to get parse beneficiary ID string to uuid format")
		return err 
	}
	err = s.repo.DeleteBeneficiary(beneID, userID)
	if err != nil {
		s.logger.Error("Failed to delete beneficiary")
		return err
	}
	return nil 
}