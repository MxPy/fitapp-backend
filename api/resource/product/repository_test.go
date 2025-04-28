package product_test

import (
	"regexp" // Use regexp for more flexible SQL matching
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"fitapp-backend/api/resource/product" // Adjust import path as needed
	mockDB "fitapp-backend/mock/db"       // Adjust import path as needed
	testUtil "fitapp-backend/util/test"   // Adjust import path as needed
)

// Common columns expected in queries (adjust based on GORM's behavior)
// Note: GORM often selects specific columns or `*`. Adapt regex accordingly.
var productColumns = []string{"id", "created_at", "updated_at", "deleted_at", "product_name", "kcal", "proteins", "carbs", "fats"}

func TestRepository_List(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := product.NewRepository(db)

	// Define rows mock data should return
	mockRows := sqlmock.NewRows(productColumns).
		AddRow(uuid.New(), time.Now(), time.Now(), gorm.DeletedAt{}, "Product A", 100, 10, 20, 5).
		AddRow(uuid.New(), time.Now(), time.Now(), gorm.DeletedAt{}, "Product B", 250, 25, 30, 8)

	// Expect a SELECT query matching the pattern for finding products
	// GORM's Find usually generates `SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL`
	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL`)
	mock.ExpectQuery(expectedSQL).WillReturnRows(mockRows)

	products, err := repo.List()
	testUtil.NoError(t, err)
	testUtil.Equal(t, 2, len(products)) // Check if two products were returned

	// Ensure all expectations were met
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := product.NewRepository(db)

	// Prepare the product data to be created
	id := uuid.New()
	newProduct := &product.Product{
		ID:          id, // ID is set before calling Create in the handler
		ProductName: "New Product",
		Kcal:        150,
		Proteins:    15,
		Carbs:       10,
		Fats:        5,
		// CreatedAt, UpdatedAt, DeletedAt are handled by GORM/DB
	}

	// Expect transaction begin
	mock.ExpectBegin()
	// Expect an INSERT statement
	// The exact columns and placeholders depend on GORM version and configuration.
	// This regex assumes GORM inserts all non-zero fields + auto fields.
	expectedSQL := regexp.QuoteMeta(`INSERT INTO "products" ("id","created_at","updated_at","deleted_at","product_name","kcal","proteins","carbs","fats") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`) // Adjust based on actual query
	mock.ExpectExec(expectedSQL).
		WithArgs(
			newProduct.ID,
			mockDB.AnyTime{}, // CreatedAt
			mockDB.AnyTime{}, // UpdatedAt
			nil,              // DeletedAt should be NULL
			newProduct.ProductName,
			newProduct.Kcal,
			newProduct.Proteins,
			newProduct.Carbs,
			newProduct.Fats,
		).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate 1 row inserted
	// Expect transaction commit
	mock.ExpectCommit()

	createdProduct, err := repo.Create(newProduct)
	testUtil.NoError(t, err)
	testUtil.NotNil(t, createdProduct)       // Check if a product object was returned
	testUtil.Equal(t, id, createdProduct.ID) // Check if the ID matches

	// Ensure all expectations were met
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Read(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := product.NewRepository(db)

	id := uuid.New()
	expectedProductName := "Read Product"

	// Define the row mock data should return for the specific ID
	mockRows := sqlmock.NewRows(productColumns).
		AddRow(id, time.Now(), time.Now(), gorm.DeletedAt{}, expectedProductName, 200, 20, 25, 10)

	// Expect a SELECT query with a WHERE clause for the ID
	// GORM's First usually generates `SELECT * FROM "products" WHERE id = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`
	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(id).            // Expect the ID as argument
		WillReturnRows(mockRows) // Return the mock row

	foundProduct, err := repo.Read(id)
	testUtil.NoError(t, err)
	testUtil.NotNil(t, foundProduct)                                 // Check if a product was found
	testUtil.Equal(t, id, foundProduct.ID)                           // Verify the ID
	testUtil.Equal(t, expectedProductName, foundProduct.ProductName) // Verify the name

	// Ensure all expectations were met
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Read_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := product.NewRepository(db)
	id := uuid.New()

	// Expect the query but return gorm.ErrRecordNotFound
	expectedSQL := regexp.QuoteMeta(`SELECT * FROM "products" WHERE id = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT 1`)
	mock.ExpectQuery(expectedSQL).
		WithArgs(id).
		WillReturnError(gorm.ErrRecordNotFound) // Simulate record not found

	_, err = repo.Read(id)
	testUtil.ErrorIs(t, err, gorm.ErrRecordNotFound) // Check if the correct error is returned

	// Ensure all expectations were met
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := product.NewRepository(db)

	id := uuid.New()
	productToUpdate := &product.Product{
		ID:          id,
		ProductName: "Updated Product",
		Kcal:        300,
		Proteins:    30,
		Carbs:       35,
		Fats:        15,
		// UpdatedAt is handled by GORM
	}

	// Expect transaction begin
	mock.ExpectBegin()
	// Expect an UPDATE statement
	// GORM's Updates with Select generates specific SET clauses
	expectedSQL := regexp.QuoteMeta(`UPDATE "products" SET "product_name"=$1,"kcal"=$2,"proteins"=$3,"carbs"=$4,"fats"=$5,"updated_at"=$6 WHERE id = $7`)
	mock.ExpectExec(expectedSQL).
		WithArgs(
			productToUpdate.ProductName,
			productToUpdate.Kcal,
			productToUpdate.Proteins,
			productToUpdate.Carbs,
			productToUpdate.Fats,
			mockDB.AnyTime{}, // UpdatedAt
			id,               // WHERE clause ID
		).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulate 1 row affected
	// Expect transaction commit
	mock.ExpectCommit()

	rowsAffected, err := repo.Update(productToUpdate)
	testUtil.NoError(t, err)
	testUtil.Equal(t, int64(1), rowsAffected) // Check if 1 row was affected

	// Ensure all expectations were met
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := product.NewRepository(db)

	id := uuid.New()

	// Expect transaction begin
	mock.ExpectBegin()
	// Expect an UPDATE statement for soft delete
	// GORM's Delete sets the deleted_at field
	expectedSQL := regexp.QuoteMeta(`UPDATE "products" SET "deleted_at"=$1 WHERE id = $2 AND "products"."deleted_at" IS NULL`)
	mock.ExpectExec(expectedSQL).
		WithArgs(mockDB.AnyTime{}, id).           // Expect timestamp for deleted_at and the ID
		WillReturnResult(sqlmock.NewResult(0, 1)) // Simulate 1 row affected
	// Expect transaction commit
	mock.ExpectCommit()

	rowsAffected, err := repo.Delete(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, int64(1), rowsAffected) // Check if 1 row was affected

	// Ensure all expectations were met
	testUtil.NoError(t, mock.ExpectationsWereMet())
}
