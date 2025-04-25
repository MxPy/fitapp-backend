package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Username  string
	Full_name string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Users []*User

type DTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Full_name string `json:"full_name"`
}

type Form struct {
	Username  string `json:"username"`
	Full_name string `json:"full_name"`
}
