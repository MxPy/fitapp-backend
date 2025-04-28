package userday

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles database operations for user_days
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new userday repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// List retrieves all non-deleted user_days (consider pagination for large datasets)
func (r *Repository) List() (UserDays, error) {
	userDays := make([]*UserDay, 0)
	if err := r.db.Find(&userDays).Error; err != nil {
		return nil, err
	}
	return userDays, nil
}

// Create inserts a new user_day record into the database
func (r *Repository) Create(userDay *UserDay) (*UserDay, error) {
	if err := r.db.Create(userDay).Error; err != nil {
		// Handle potential constraint violations (e.g., duplicate user_id/date if unique index exists)
		return nil, err
	}
	return userDay, nil
}

// Read retrieves a single user_day record by its primary ID
func (r *Repository) Read(id uuid.UUID) (*UserDay, error) {
	userDay := &UserDay{}
	if err := r.db.Where("id = ?", id).First(&userDay).Error; err != nil {
		return nil, err // Can be gorm.ErrRecordNotFound
	}
	return userDay, nil
}

// FindByUserAndDate retrieves a single user_day record by user ID and date
func (r *Repository) FindByUserAndDate(userID uuid.UUID, date time.Time) (*UserDay, error) {
	userDay := &UserDay{}
	// Format time.Time to "YYYY-MM-DD" string for WHERE clause on DATE type column
	dateStr := date.Format(DateFormat)
	if err := r.db.Where("user_id = ? AND user_date = ?", userID, dateStr).First(&userDay).Error; err != nil {
		return nil, err // Can be gorm.ErrRecordNotFound
	}
	return userDay, nil
}

// Update modifies an existing user_day record in the database.
// Only updates nutritional values and UpdatedAt timestamp.
func (r *Repository) Update(userDay *UserDay) (int64, error) {
	result := r.db.Model(&UserDay{}).
		// Select only the fields allowed to be updated
		Select("DailyKcal", "DailyProteins", "DailyCarbs", "DailyFats", "UpdatedAt").
		Where("id = ?", userDay.ID).
		Updates(userDay) // Pass the userDay struct with new values

	return result.RowsAffected, result.Error
}

// Delete performs a soft delete on a user_day record by its primary ID
func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&UserDay{})
	return result.RowsAffected, result.Error
}

// Helper to explicitly tell GORM the table name if needed (usually inferred)
// func (UserDay) TableName() string {
// 	return "user_days"
// }
