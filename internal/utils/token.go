package utils

import (
	"time"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrSigningToken = errors.New("Error while signing token")
)

// Generates a JWT token
func GenerateJWTToken(
	userID uuid.UUID,
	secret string,
	expiryDate time.Time,
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     expiryDate.Unix(),
	})
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", ErrSigningToken
	}


	return signedToken, err
}
