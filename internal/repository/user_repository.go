package repository

import (
	"errors"
	// "log"

	//"strconv"
	// "strings"
	// "time"

	// "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

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

func GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	result := db.GetDB().Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil

}

func CreateCart(cart *models.Cart) error {
	result := db.GetDB().Create(cart)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func AddItemToCart(item *models.CartItem) error {
	result := db.GetDB().Create(item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetCartID(id uint) (uint, error) {
	var cart models.Cart
	result := db.GetDB().Where("user_id = ?", id).First(&cart)
	if result.Error != nil {
		return 0, result.Error
	}
	return cart.ID, nil
}

func GetCartByID(id uint) (*models.Cart, error) {
	var cart models.Cart
	result := db.GetDB().Preload("Items").Where("id = ?", id).First(&cart)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("Корзина не найдена")
		}
		return nil, result.Error
	}
	return &cart, nil
}

func ClearCart(cart_id uint) error {
	result := db.GetDB().Where("cart_id = ?", cart_id).Delete(&models.CartItem{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
