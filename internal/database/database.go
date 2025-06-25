package database

import (
	"pgpockets/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectToDatabase(dbSource string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	db.AutoMigrate(
		&models.Card{},
		&models.User{},
		&models.Transaction{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.Profile{},
		&models.Session{},
		&models.Wallet{},
	)

	return db, nil
}