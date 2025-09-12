package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	FirstName string `gorm:"not null"`
	LastName  string
	Email     string
	Token     string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RegisterRequest struct {
	FirstName string `json:"FirstName" binding:"required"`
	LastName  string `json:"LastName" binding:"required"`
	Email     string `json:"Email" binding:"required,email"`
	Password  string `json:"Password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"Email" binding:"required,email"`
	Password string `json:"Password" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type ForgotPasswordRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type ResetPasswordRequest struct {
	ResetToken  string `json:"resetToken" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `json:"user_id"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items"`
}

type Category struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Products []Product `gorm:"foreignkey:CategoryID"`
}

type Product struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"not null"`
	Description string
	Price       float64  `gorm:"not null"`
	CategoryID  uint     `gorm:"not null"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Stock       int      `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type Image struct {
	ID        uint
	ProductID uint
	URL       string
	Alt       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Product Product
}

type CartResponse struct {
	ID        uint `gorm:"primary_key"`
	UserID    uint
	Total     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []CartResponseItem
}

type CartItem struct {
	ID        uint `gorm:"primary_key"`
	CartID    uint
	ProductID uint
	Quantity  int
	Price     float64
	Total     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartResponseItem struct {
	ProductID uint
	Name      string
	Price     float64
	Quantity  int
	Total     float64
}
