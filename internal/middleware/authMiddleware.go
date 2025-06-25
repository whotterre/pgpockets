package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	jwtSecret []byte
	logger *zap.Logger
}

func NewAuthMiddleware(jwtSecret string, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		logger: logger,
	}
}

func (a *AuthMiddleware) RequireAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
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

        // Extract the token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // Parse and validate the token
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

        // Check if token is valid and extract claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            // Check if token is expired
            if exp, ok := claims["exp"].(float64); ok {
                if time.Now().Unix() > int64(exp) {
                    a.logger.Warn("Token expired")
                    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                        "error": "Token expired",
                    })
                }
            }

            // Extract user ID
            userIDStr, ok := claims["user_id"].(string)
            if !ok {
                a.logger.Warn("Invalid user ID in token")
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "Invalid token claims",
                })
            }

            // Parse user ID to UUID
            userID, err := uuid.Parse(userIDStr)
            if err != nil {
                a.logger.Warn("Invalid user ID format", zap.Error(err))
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "Invalid user ID format",
                })
            }

            c.Locals("userID", userID)
            c.Locals("access_token", tokenString)

            return c.Next()
        }

        a.logger.Warn("Invalid token claims")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }
}