package services

import (
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProfileService interface {
	GetProfileByUserID(userID string) (*models.Profile, error)
	UpdateProfile(profile *models.Profile, userID uuid.UUID) error
}

type profileService struct {
	profileRepo repositories.ProfileRepository
	appLogger   *zap.Logger
}

func NewProfileService(profileRepo repositories.ProfileRepository, appLogger *zap.Logger) ProfileService {
	return &profileService{
		profileRepo: profileRepo,
		appLogger:   appLogger,
	}
}

func (s *profileService) GetProfileByUserID(userID string) (*models.Profile, error) {
	profile, err := s.profileRepo.GetProfileByUserID(userID)
	if err != nil {
		s.appLogger.Error("Failed to get profile by user ID", zap.String("userID", userID), zap.Error(err))
		return nil, err
	}
	return profile, nil
}



func (s *profileService) UpdateProfile(profile *models.Profile, userID uuid.UUID) error {
	if err := s.profileRepo.UpdateProfile(profile); err != nil {
		s.appLogger.Error("Failed to update profile", zap.String("userID", userID.String()), zap.Error(err))
		return err
	}
	return nil
}

