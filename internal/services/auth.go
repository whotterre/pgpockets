package services

import (
	"errors"
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists     = errors.New("user with this email already exists")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrPasswordHashFailed    = errors.New("failed to process password securely")
	ErrTokenGenerationFailed = errors.New("failed to generate authentication token")
)

type AuthService interface {
	Register(
		email string,
		password string,
		firstName string,
		lastName string,
		dob string,
		phoneNumber string,
		address string,
		gender string,
	) (*models.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	logger    *zap.Logger
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, logger *zap.Logger, jwtSecret string) *authService {
	return &authService{
		userRepo:  userRepo,
		logger:    logger,
		jwtSecret: jwtSecret,
	}
}

// Business logic for user registration
func (s *authService) Register(
	email string,
	password string,
	firstName string,
	lastName string,
	dob string,
	phoneNumber string,
	address string,
	gender string,
) (*models.User, error) {
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		s.logger.Warn("User already exists", zap.String("email", email))
		return nil, ErrUserAlreadyExists
	}
	if err != gorm.ErrRecordNotFound {
		s.logger.Error("Error checking user existence", zap.Error(err))
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, ErrPasswordHashFailed
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	parsedDOB, err := time.Parse("02-01-2006", dob)
    if err != nil {
        s.logger.Error("Failed to parse date of birth", zap.Error(err))
        return nil, errors.New("invalid date of birth format, expected DD-MM-YYYY")
    }

	
	err = s.userRepo.CreateUser(user)
	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, err
	}
	profile := &models.Profile{
		UserID:      user.ID,
		FirstName:   firstName,
		LastName:    lastName,
		DateOfBirth: &parsedDOB,
		PhoneNumber: phoneNumber,
		Address:     address,
	}
	err = s.userRepo.CreateProfile(profile)
	if err != nil {
		s.logger.Error("Failed to create user profile", zap.Error(err))
		if rollbackErr := s.userRepo.DeleteUser(user.ID); rollbackErr != nil {
			s.logger.Error("Failed to rollback user creation after profile creation failure", zap.Error(rollbackErr))
		}
		return nil, err
	}
	s.logger.Info("User registered successfully", zap.String("email", email))
	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	return "", nil
}
