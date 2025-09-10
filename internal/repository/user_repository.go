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

func GetProductsUser() ([]models.Product, error) {
	var products []models.Product
	result := db.GetDB().Preload("Category").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
