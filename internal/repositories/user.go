package repositories

import (
	"errors"
	"pgpockets/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists     = errors.New("user with this email already exists")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrPasswordHashFailed    = errors.New("failed to process password securely")
	ErrTokenGenerationFailed = errors.New("failed to generate authentication token")
	ErrUserNotFound          = errors.New("user not found")
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	DeleteUser(id uuid.UUID) error
	CreateProfile(profile *models.Profile) error

	CreateSession(session *models.Session) error
}

func (u *gormUserRepository) CreateProfile(profile *models.Profile) error {
	if err := u.db.Create(profile).Error; err != nil {
		return err
	}
	return nil
}

func (u *gormUserRepository) DeleteUser(d uuid.UUID) error {
	if err := u.db.Where("id = ?", d).Delete(&models.User{}).Error; err != nil {
		return err
	}
	return nil
}


// CreateUser creates a new user in the database.
func (r *gormUserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByEmail retrieves a user by their email from the database.
func (r *gormUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}
