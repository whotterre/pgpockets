package repositories

import (
	"errors"
	"pgpockets/internal/models"
	"time"

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
	DeleteSession(userID uuid.UUID) error
	GetActiveSessionByToken(accessToken string) (*models.Session, error)
	RefreshSession(refreshToken string) (*models.Session, error)
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

func (r *gormUserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *gormUserRepository) DeleteSession(userID uuid.UUID) error {
    return r.db.Where("user_id = ?", userID).Delete(&models.Session{}).Error
}

func (r *gormUserRepository) GetActiveSessionByToken(accessToken string) (*models.Session, error) {
    var session models.Session
    err := r.db.Where("access_token = ? AND is_active = ? AND access_token_expires_at > ?", 
        accessToken, true, time.Now()).First(&session).Error
    if err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *gormUserRepository) RefreshSession(refreshToken string) (*models.Session, error) {
    var session models.Session
    err := r.db.Where("refresh_token = ? AND is_active = ? AND refresh_token_expires_at > ?", 
        refreshToken, true, time.Now()).First(&session).Error
    if err != nil {
        return nil, err
    }
    return &session, nil
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}
