package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles database operations for users
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// List retrieves all non-deleted users
func (r *Repository) List() (Users, error) {
	users := make([]*User, 0)
	// GORM automatically uses column mappings defined in the User struct
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Create inserts a new user into the database
func (r *Repository) Create(user *User) (*User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Read retrieves a single user by its ID
func (r *Repository) Read(id uuid.UUID) (*User, error) {
	user := &User{}
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err // Can be gorm.ErrRecordNotFound
	}
	return user, nil
}

// Update modifies an existing user in the database
func (r *Repository) Update(user *User) (int64, error) {
	// Specify fields allowed to be updated using GORM struct field names
	result := r.db.Model(&User{}).
		Select("Username", "FullName", "Sex", "Height", "Weight", "Age", "UpdatedAt").
		Where("id = ?", user.ID).
		Updates(user) // GORM handles mapping to correct DB columns

	return result.RowsAffected, result.Error
}

// Delete performs a soft delete on a user by its ID
func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&User{})
	return result.RowsAffected, result.Error
}
