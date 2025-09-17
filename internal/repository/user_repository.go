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

func IncrementItem(cartID uint, productID uint) error {
	var item models.CartItem
	result := db.GetDB().Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item)
	if result.Error != nil {
		return result.Error
	}
	item.Quantity += 1
	item.Total = item.Price * float64(item.Quantity)
	result = db.GetDB().Save(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DecrementItem(cartID uint, productID uint) error {
	var item models.CartItem
	result := db.GetDB().Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item)
	if result.Error != nil {
		return result.Error
	}
	if item.Quantity <= 1 {
		return RemoveFromCart(cartID, productID)
	}
	item.Quantity -= 1
	item.Total = item.Price * float64(item.Quantity)
	result = db.GetDB().Save(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func RemoveFromCart(cartID uint, productID uint) error {
	result := db.GetDB().Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&models.CartItem{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateOrder(order *models.Order) error {
	if err := db.GetDB().Create(order).Error; err != nil {
		return err
	}
	return nil
}

func GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	result := db.GetDB().Where("user_id = ?", userID).Preload("Items").Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func DeleteOrder(orderID uint, userID uint) error {
	result := db.GetDB().Where("id = ? AND user_id = ?", orderID, userID).Delete(&models.Order{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateOrderStatus(orderID uint, status string) error {
	if status == "" {
		return errors.New("Статус не может быть пустым")
	} else {
		result := db.GetDB().Model(&models.Order{}).Where("id = ?", orderID).Update("status", status)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
}
