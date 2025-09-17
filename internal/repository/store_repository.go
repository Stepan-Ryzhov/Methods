package repository

import (
	// "errors"
	// "log"

	//"strconv"
	// "strings"
	// "time"

	// "golang.org/x/crypto/bcrypt"
	// "gorm.io/gorm"

	db "methodi_razrabotki/internal/database"
	models "methodi_razrabotki/internal/models"
)

func GetOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := db.GetDB().Preload("Items").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
