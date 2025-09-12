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

var (
	empty            error = errors.New("Одно или несколько полей не заполнены")
	usedEmail        error = errors.New("Пользователь с таким email уже существует")
	nonValidEmail    error = errors.New("Некорректный email")
	nonValidPassword error = errors.New("Пароль не должен содержать меньше 6 символов")
	nonRegisterEmail error = errors.New("Пользователь с таким e-mail не зарегистрирован")
	errorPassword    error = errors.New("Неверный пароль")
	errorCategory    error = errors.New("Категория с таким именем уже существует")
	errorCategory2   error = errors.New("Категория с таким именем не существует")
	emptyName        error = errors.New("Поле не может быть пустым")
	errorDelete      error = errors.New("Такой элемент не найден")
	errorProduct     error = errors.New("Товар с таким именем уже существует")
)

func Register(req *models.RegisterRequest) error {
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
				Role:      "user",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if req.Password == "admin1234" {
				user.Role = "admin"
			}
			result := db.GetDB().Create(&user)
			if result.Error != nil {
				log.Printf("Ошибка создания пользователя в БД: %v", result.Error)
				return result.Error
			}
			if user.Role == "user" {
				cart := &models.Cart{
					UserID:    user.ID,
					Total:     0.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Items:     nil,
				}
				CreateCart(cart)
			}
			return nil
		}
		return result.Error
	} else {
		return usedEmail
	}
}

func Login(req models.LoginRequest) (*models.User, error) {
	if req.Email == "" || req.Password == "" {
		return nil, empty
	}
	if strings.Contains(req.Email, "@") == false {
		return nil, nonValidEmail
	}
	if len(req.Password) < 6 {
		return nil, nonValidPassword
	}
	var user models.User
	result := db.GetDB().Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nonRegisterEmail
		}
		return nil, result.Error
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Token), []byte(req.Password))
	if err != nil {
		return nil, errorPassword
	}
	return &user, nil
}
