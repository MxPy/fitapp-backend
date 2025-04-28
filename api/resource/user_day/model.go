package userday

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// Importuj model użytkownika, jeśli chcesz zdefiniować relację GORM
	// "fitapp-backend/api/resource/user"
)

// UserDay represents the structure of the 'user_days' table
type UserDay struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Indeks na DeletedAt jest często przydatny

	UserID        uuid.UUID `gorm:"type:uuid;not null;index:idx_userday_user_date"` // Klucz obcy + część indeksu złożonego
	UserDate      time.Time `gorm:"type:date;not null;index:idx_userday_user_date"` // Data + część indeksu złożonego
	DailyKcal     int       `gorm:"not null"`
	DailyProteins int       `gorm:"not null"`
	DailyCarbs    int       `gorm:"not null"`
	DailyFats     int       `gorm:"not null"`

	// Opcjonalna definicja relacji dla GORM (np. do Eager Loading)
	// User          user.User `gorm:"foreignKey:UserID"`
}

// UserDays is a slice of UserDay pointers
type UserDays []*UserDay

// DTO represents the data transfer object for a UserDay
type DTO struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	UserDate      string `json:"user_date"` // Format "YYYY-MM-DD"
	DailyKcal     int    `json:"daily_kcal"`
	DailyProteins int    `json:"daily_proteins"`
	DailyCarbs    int    `json:"daily_carbs"`
	DailyFats     int    `json:"daily_fats"`
	// Można dodać CreatedAt/UpdatedAt w razie potrzeby
}

// Form represents the data structure for creating/updating a UserDay
// Parsowanie stringów UserID i UserDate odbywa się w handlerze
type Form struct {
	UserID        string `json:"user_id"`   // Oczekiwany UUID jako string
	UserDate      string `json:"user_date"` // Oczekiwana data w formacie "YYYY-MM-DD"
	DailyKcal     int    `json:"daily_kcal"`
	DailyProteins int    `json:"daily_proteins"`
	DailyCarbs    int    `json:"daily_carbs"`
	DailyFats     int    `json:"daily_fats"`
}
