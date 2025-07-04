package services

import (
	"errors"
	"pgpockets/internal/models"
	"pgpockets/internal/repositories"
	"pgpockets/internal/utils"
	"time"

	"github.com/google/uuid"
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
	Login(email, password, ipAddr, userAgent string) (string, string, error)
	Logout(userID uuid.UUID) error
}

type authService struct {
	userRepo  repositories.UserRepository
	walletRepo repositories.WalletRepository
	logger    *zap.Logger
	jwtSecret string
}

func NewAuthService(
	userRepo repositories.UserRepository,
	walletRepo repositories.WalletRepository,
	logger *zap.Logger,
	jwtSecret string,
	) *authService {
	return &authService{
		userRepo:  userRepo,
		walletRepo: walletRepo,
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
		Gender:      gender,
	}
	err = s.userRepo.CreateProfile(profile)
	if err != nil {
		s.logger.Error("Failed to create user profile", zap.Error(err))
		if rollbackErr := s.userRepo.DeleteUser(user.ID); rollbackErr != nil {
			s.logger.Error("Failed to rollback user creation after profile creation failure", zap.Error(rollbackErr))
		}
		return nil, err
	}

	// Create wallet
	if err := s.walletRepo.CreateWallet(user.ID); err != nil {
		s.logger.Error("Something went wrong while creating the wallet", zap.Error(err))
	}
	


	s.logger.Info("User registered successfully", zap.String("email", email))
	return user, nil
}

func (s *authService) Login(email, password, ipAddr, userAgent string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warn("Login failed: user not found", zap.String("email", email))
			return "", "", ErrInvalidCredentials
		}
		s.logger.Error("Error retrieving user", zap.Error(err))
		return "", "", err
	}
	// Check if password matches
	if user == nil {
		s.logger.Warn("Login failed: user not found", zap.String("email", email))
		return "", "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Warn("Login failed: invalid password", zap.String("email", email))
		return "", "", ErrInvalidCredentials
	}

	// Generate JWT tokens for access and refresh token
	accessToken, err := utils.GenerateJWTToken(
		user.ID,
		s.jwtSecret,
		time.Now().Add(15*time.Minute),
	)
	if err != nil {
		s.logger.Error("Failed to generate JWT token", zap.Error(err))
		return "", "", ErrTokenGenerationFailed
	}

	refreshToken, err := utils.GenerateJWTToken(
		user.ID,
		s.jwtSecret,
		time.Now().Add(15*time.Minute),
	)
	if err != nil {
		s.logger.Error("Failed to generate JWT token", zap.Error(err))
		return "", "", ErrTokenGenerationFailed
	}

	// Create a new session for the user
	session := &models.Session{
		UserID:                user.ID,
		ClientIP:              ipAddr,
		UserAgent:             userAgent,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  time.Now().Add(15 * time.Minute),
		RefreshTokenExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		RefreshToken:          refreshToken,
		IsActive:              true,
	}

	err = s.userRepo.CreateSession(session)
	if err != nil {
		s.logger.Error("Failed to create user session", zap.Error(err))
		return "", "", ErrTokenGenerationFailed
	}
	s.logger.Info("User logged in successfully", zap.String("email", email))
	return accessToken, refreshToken, nil
}

func (s *authService) Logout(userID uuid.UUID) error {
	err := s.userRepo.DeleteSession(userID)

	if err != nil {
		s.logger.Error("Failed to delete user sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return err
	}

	s.logger.Info("User logged out successfully!")
	return nil

}
