package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the structure of the 'users' table
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Map struct fields to potentially different DB column names
	Username string `gorm:"column:user_username;not null"`
	FullName string `gorm:"column:user_full_name;not null"`
	Sex      bool   `gorm:"column:user_sex;not null"`    // true = male, false = female (by convention)
	Height   int    `gorm:"column:user_height;not null"` // in cm
	Weight   int    `gorm:"column:user_weight;not null"` // in kg (or other unit)
	Age      int    `gorm:"column:user_age;not null"`
}

// Users is a slice of User pointers
type Users []*User

// DTO represents the data transfer object for a User
type DTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Sex      string `json:"sex"` // "male" or "female" or "unknown"
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Age      int    `json:"age"`
	// Optionally add CreatedAt/UpdatedAt strings if needed
}

// Form represents the data structure for creating/updating a User
type Form struct {
	Username string `json:"username" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Sex      *bool  `json:"sex" validate:"required"`         // Use pointer to distinguish false from nil (not provided)
	Height   int    `json:"height" validate:"required,gt=0"` // Height must be positive
	Weight   int    `json:"weight" validate:"required,gt=0"` // Weight must be positive
	Age      int    `json:"age" validate:"required,gt=0"`    // Age must be positive
}

// --- TableName (Optional) ---
// Uncomment if GORM has trouble inferring the table name
// func (User) TableName() string {
// 	return "users"
// }
