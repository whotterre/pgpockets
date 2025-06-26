package middleware

import (
	"pgpockets/internal/repositories"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
    userRepo repositories.UserRepository
	jwtSecret []byte
	logger    *zap.Logger
}

func NewAuthMiddleware(jwtSecret string, logger *zap.Logger, userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
        userRepo: userRepo,
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

func (a *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			a.logger.Warn("Missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			a.logger.Warn("Invalid authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return a.jwtSecret, nil
		})

		if err != nil {
			a.logger.Warn("Invalid token", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			a.logger.Warn("Invalid token claims")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			a.logger.Warn("Invalid user ID in token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			a.logger.Warn("Invalid user ID format", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID format",
			})
		}

		// **Key Addition**: Verify token exists in database
		session, err := a.userRepo.GetActiveSessionByToken(tokenString)
		if err != nil {
			a.logger.Warn("Session not found or expired", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session not found or expired",
			})
		}

		if !session.IsActive {
			a.logger.Warn("Session has been revoked")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session has been revoked",
			})
		}

		// Store in context
		c.Locals("userID", userID)
		c.Locals("sessionID", session.ID)
		c.Locals("accessToken", tokenString)

		return c.Next()
	}
}
