package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository handles database operations for products
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new product repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// List retrieves all non-deleted products
func (r *Repository) List() (Products, error) {
	products := make([]*Product, 0)
	// GORM automatically handles the DeletedAt field for soft deletes
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// Create inserts a new product into the database
func (r *Repository) Create(product *Product) (*Product, error) {
	// GORM automatically handles CreatedAt and UpdatedAt
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// Read retrieves a single product by its ID
func (r *Repository) Read(id uuid.UUID) (*Product, error) {
	product := &Product{}
	// First will automatically add WHERE deleted_at IS NULL
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		// err could be gorm.ErrRecordNotFound
		return nil, err
	}
	return product, nil
}

// Update modifies an existing product in the database
func (r *Repository) Update(product *Product) (int64, error) {
	// GORM automatically handles UpdatedAt
	// Select specifies which fields are allowed to be updated
	result := r.db.Model(&Product{}).
		Select("ProductName", "Kcal", "Proteins", "Carbs", "Fats", "UpdatedAt").
		Where("id = ?", product.ID).
		Updates(product) // Pass the product struct with new values

	// result.Error could be gorm.ErrRecordNotFound if ID doesn't exist
	// result.RowsAffected will be 0 if no record found or if data is the same
	return result.RowsAffected, result.Error
}

// Delete performs a soft delete on a product by its ID
func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	// GORM's Delete performs a soft delete if gorm.DeletedAt field exists
	result := r.db.Where("id = ?", id).Delete(&Product{})

	// result.Error can occur
	// result.RowsAffected will be 0 if the record was already deleted or not found
	return result.RowsAffected, result.Error
}

// Helper to explicitly tell GORM the table name if needed (usually inferred)
// func (Product) TableName() string {
//  return "products"
// }
