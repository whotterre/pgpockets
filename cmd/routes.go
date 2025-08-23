package main

import (
	"pgpockets/internal/config"
	"pgpockets/internal/handlers"
	"pgpockets/internal/middleware"
	"pgpockets/internal/repositories"
	"pgpockets/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
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
	// Initialize reusable rate limiter concern
	rateLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 2 * time.Second,
	})
	/* Protected routes */
	apiV1.Use(authMiddleware.RequireAuth())

	authGroup.Delete("/logout", authHandlers.LogoutUser)

	// Dashboard routes
	dashboardRepo := repositories.NewDashboardRepository(db)
	dashboardService := services.NewDashboardService(dashboardRepo, appLogger, config.ExchangeRatesAPIKey)
	dashboardHandlers := handlers.NewDashboardHandler(dashboardService, appLogger)
	dashboardGroup := apiV1.Group("/dashboard")
	dashboardGroup.Use(rateLimiter)
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
	walletService := services.NewWalletService(walletRepo, appLogger, db)
	walletHandlers := handlers.NewWalletHandler(walletService, appLogger, config.ExchangeRatesAPIKey)
	walletGroup := apiV1.Group("/wallets")
	// Rate limiting
	walletGroup.Use(rateLimiter)
	walletGroup.Get("/balance", walletHandlers.GetWalletBalance)
	walletGroup.Get("/balances", walletHandlers.GetBalancesForAllWallets)
	walletGroup.Patch("/currency/:desiredCurrency", walletHandlers.ChangeWalletCurrency)
	// Transaction routes
	txnRepo := repositories.NewTransactionRepository(db)
	txnService := services.NewTransactionService(txnRepo, appLogger, walletRepo, db)
	txnHandlers := handlers.NewTransactionHandler(txnService, appLogger)
	txnGroup := apiV1.Group("/transaction")
	txnGroup.Use(rateLimiter)
	txnGroup.Patch("/make-transfer", txnHandlers.TransferFunds)
	txnGroup.Get("/history", txnHandlers.GetUserTransactionHistory)
	txnGroup.Get("/history/date-range", txnHandlers.GetTransactionsInDateRange)
	txnGroup.Get("/transaction/:txnID", txnHandlers.GetTransactionByID)

	// Profile routes
	profileRepo := repositories.NewProfileRepository(db)
	profileService := services.NewProfileService(profileRepo, appLogger)
	profileHandlers := handlers.NewProfileHandler(profileService, appLogger)
	profileGroup := apiV1.Group("/profile")
	profileGroup.Get("/", profileHandlers.GetProfile)
	profileGroup.Put("/", profileHandlers.UpdateProfile)
	// Beneficiary routes
	beneficiaryRepo := repositories.NewBeneficiaryRepository(db)
	beneficiaryService := services.NewBeneficiaryService(beneficiaryRepo, appLogger)
	beneficiaryHandlers := handlers.NewBeneficiaryHandler(beneficiaryService, appLogger)
	beneficiaryGroup := apiV1.Group("/beneficiaries")
	beneficiaryGroup.Get("/", beneficiaryHandlers.GetBeneficiaries)
	beneficiaryGroup.Post("/", beneficiaryHandlers.AddBeneficiary)
	beneficiaryGroup.Delete("/beneficiary/:beneID", beneficiaryHandlers.DeleteBeneficiary)
	// Notification Routes
	notifRepo := repositories.NewNotifRepo(db)
	notifServices := services.NewNotificationService(notifRepo)
	notifHandlers := handlers.NewNotificationHandlers(notifServices, appLogger)
	notifGroup := apiV1.Group("/notifications")
	notifGroup.Get("/count", notifHandlers.GetNotificationCount)
	notifGroup.Get("/unread-count", notifHandlers.GetUnreadNotificationCount)
	notifGroup.Get("/", notifHandlers.GetNotifications)
	notifGroup.Get("/unread", notifHandlers.GetUnreadNotifications)
	notifGroup.Get("/:id", notifHandlers.GetNotification)
	notifGroup.Patch("/:id/read", notifHandlers.MarkAsRead)
	notifGroup.Patch("/:id/unread", notifHandlers.MarkAsUnread)
	notifGroup.Patch("/read-all", notifHandlers.MarkAllAsRead)
	notifGroup.Delete("/:id", notifHandlers.DeleteNotification)
	notifGroup.Delete("/", notifHandlers.DeleteAllNotifications)
	notifGroup.Delete("/read", notifHandlers.DeleteAllReadNotifications)

}
