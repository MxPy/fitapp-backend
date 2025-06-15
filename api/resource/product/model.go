package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents the structure of the 'products' table
type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ProductName string         `gorm:"not null"`
	Kcal        int            `gorm:"not null"`
	Proteins    int            `gorm:"not null"`
	Carbs       int            `gorm:"not null"`
	Fats        int            `gorm:"not null"`
}

// Products is a slice of Product pointers
type Products []*Product

// DTO represents the data transfer object for a Product
type DTO struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name"`
	Kcal        int    `json:"kcal"`
	Proteins    int    `json:"proteins"`
	Carbs       int    `json:"carbs"`
	Fats        int    `json:"fats"`
}

// Form represents the data structure for creating/updating a Product
type Form struct {
	UserID      string `json:"user_id"`
	ProductName string `json:"product_name"`
	Grams       int    `json:"grams"`
	Kcal        int    `json:"kcal"`
	Proteins    int    `json:"proteins"`
	Carbs       int    `json:"carbs"`
	Fats        int    `json:"fats"`
}
