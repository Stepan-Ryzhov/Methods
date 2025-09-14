package database

import (
	"fmt"
	"log"

	models "methodi_razrabotki/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Автомиграция моделей
	if err := db.AutoMigrate(
		&models.Cart{},
		&models.CartItem{},
		&models.Category{},
		&models.Product{},
		&models.User{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderList{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	log.Println("Successfully connected to SQLite database")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
