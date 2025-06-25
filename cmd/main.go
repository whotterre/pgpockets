package main

import (
	"log"
	"pgpockets/internal/config"
	"pgpockets/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/zap"
)

func main() {
	app := fiber.New()
	var err error
	// Twelve factor apps need logs for 011y
	// Initialize Zap logger
	appLogger, err := zap.NewDevelopment() //TODO: change to newprod l8r
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger %v", err)
	}
	defer appLogger.Sync()

	app.Use(logger.New())
	
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	
	// Connect to the database
	db, err := database.ConnectToDatabase(config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	SetupRoutes(app, config, appLogger, db)
	
	log.Fatal(app.Listen(config.ServerAddr))
}
