package repository

import (
	"errors"
	"log"

	//"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	db "methodi_razrabotki/internal/database"
	models "methodi_razrabotki/internal/models"
)

func GetCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := db.GetDB().Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func FindCategory(name string) (*models.Category, error) {
	if name == "" {
		return nil, emptyName
	}

	var existing models.Category
	err := db.GetDB().Where("name = ?", name).First(&existing).Error
	if err != nil {
		return nil, errorCategory2
	} else {
		return &existing, nil
	}
}

func CreateCategory(name string) (*models.Category, error) {
	if name == "" {
		return nil, emptyName
	}

	var existing models.Category
	err := db.GetDB().Where("name = ?", name).First(&existing).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		category := models.Category{
			Name:      name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.GetDB().Create(&category).Error; err != nil {
			return nil, err
		}
		return &category, nil

	case err != nil:
		return nil, err

	default:
		return nil, errorCategory
	}
}

func DeleteCategory(name string) error {
	var category models.Category
	if err := db.GetDB().Where("name = ?", name).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorDelete
		}
		return err
	}

	var productCount int64
	db.GetDB().Model(&models.Product{}).Where("category_id = ?", category.ID).Count(&productCount)
	if productCount > 0 {
		return errors.New("Невозможно удалить категорию: к ней привязаны товары")
	}

	result := db.GetDB().Delete(&category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateProduct(product *models.Product) (*models.Product, error) {
	var existing models.Product
	err := db.GetDB().Where("name = ?", product.Name).First(&existing).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		if err := db.GetDB().Create(product).Error; err != nil {
			return nil, err
		}
		return product, nil

	case err != nil:
		return nil, err

	default:
		return nil, errorProduct
	}
}

func GetProducts() ([]models.Product, error) {
	var products []models.Product
	result := db.GetDB().Preload("Category").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func DeleteProduct(name string) error {
	result := db.GetDB().Where("name = ?", name).Delete(&models.Product{})
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errorDelete
		} else {
			return result.Error
		}
	}
	return nil
}

func CreateStoreMan(req *models.RegisterRequest) error {
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		return empty
	}
	if strings.Contains(req.Email, "@") == false {
		return nonValidEmail
	}
	if len(req.Password) < 6 {
		return nonValidPassword
	}
	var exUser models.User
	result := db.GetDB().Where("email = ?", req.Email).First(&exUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			user := models.User{
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Email:     req.Email,
				Token:     string(hashedPassword),
				Role:      "store",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			result := db.GetDB().Create(&user)
			if result.Error != nil {
				log.Printf("Ошибка создания пользователя в БД: %v", result.Error)
				return result.Error
			}

			return nil
		}
		return result.Error
	} else {
		return usedEmail
	}
}

func UpdateProduct(req *models.Product) error {
	existing := &models.Product{}
	if err := db.GetDB().Where("id = ?", req.ID).First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Товар с таким ID не найден")
		}
		return err
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.CategoryID != 0 {
		existing.CategoryID = req.CategoryID
	}
	if req.Price != 0 {
		existing.Price = req.Price
	}
	if req.Stock != 0 {
		existing.Stock = req.Stock
	}

	result := db.GetDB().Model(existing).Updates(existing)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
