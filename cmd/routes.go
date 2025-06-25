package main

import (
	"pgpockets/internal/config"
	"pgpockets/internal/handlers"
	"pgpockets/internal/repositories"
	"pgpockets/internal/services"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, config config.Config, appLogger *zap.Logger, db *gorm.DB) {
	apiV1 := app.Group("/api/v1")
	appLogger.Info("Setting up routes...")
	// Auth routes
	authGroup := apiV1.Group("/auth")
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, appLogger, config.JWTSecret)
	authHandlers := handlers.NewAuthHandler(authService, appLogger)
	authGroup.Post("/register", authHandlers.RegisterUser)
	authGroup.Post("/login", authHandlers.LoginUser)
	authGroup.Delete("/logout", authHandlers.LogoutUser)
	

}