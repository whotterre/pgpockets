package main

import (
	"log"
	"pgpockets/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main(){
	app := fiber.New()
	// Twelve factor apps need logs
	app.Use(logger.New())
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	log.Fatal(app.Listen(config.ServerAddr))
}