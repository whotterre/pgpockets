package main

import (
	"pgpockets/internal/config"
	"pgpockets/internal/handlers"
	"pgpockets/internal/middleware"
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
	walletRepo := repositories.NewWalletRepository(db)
	authService := services.NewAuthService(userRepo, walletRepo, appLogger, config.JWTSecret)
	authHandlers := handlers.NewAuthHandler(authService, appLogger)
	authGroup.Post("/register", authHandlers.RegisterUser)
	authGroup.Post("/login", authHandlers.LoginUser)

	// Initialize authMiddleware
	authMiddleware := middleware.NewAuthMiddleware(config.JWTSecret, appLogger, userRepo)

	/* Protected routes */
	apiV1.Use(authMiddleware.RequireAuth())
	authGroup.Delete("/logout", authHandlers.LogoutUser)

	// Dashboard routes
	dashboardRepo := repositories.NewDashboardRepository(db)
	dashboardService := services.NewDashboardService(dashboardRepo, appLogger, config.ExchangeRatesAPIKey)
	dashboardHandlers := handlers.NewDashboardHandler(dashboardService, appLogger)
	dashboardGroup := apiV1.Group("/dashboard")
	dashboardGroup.Get("/exchange-rates", dashboardHandlers.GetExchangeRates)
	// Card routes
	cardRepo := repositories.NewCardRepository(db)
	cardService := services.NewCardService(cardRepo, appLogger)
	cardHandlers := handlers.NewCardHandler(cardService, appLogger)
	cardGroup := apiV1.Group("/cards")
	cardGroup.Post("/", cardHandlers.CreateCard)
	cardGroup.Get("/cards", cardHandlers.RetrieveAllCards)
	cardGroup.Get("/card/:cardID", cardHandlers.GetCardByID)
	cardGroup.Delete("/card/:cardID", cardHandlers.DeleteCard)
	// Wallet routes
	
}


