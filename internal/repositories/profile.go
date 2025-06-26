package repositories

import (
	"pgpockets/internal/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetProfileByUserID(userID string) (*models.Profile, error)
	UpdateProfile(profile *models.Profile) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{
		db: db,
	}
}

func (r *profileRepository) GetProfileByUserID(userID string) (*models.Profile, error) {
	var profile models.Profile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) UpdateProfile(profile *models.Profile) error {
	if err := r.db.Save(profile).Error; err != nil {
		return err
	}
	return nil
}