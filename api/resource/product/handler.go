package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	// Assuming a shared error handling package exists
	// "fitapp-backend/internal/err"
	userday "fitapp-backend/api/resource/user_day"
)

// API holds the dependencies for the product handlers
type API struct {
	repository   *Repository
	user_day_api *userday.API
}

// New creates a new API instance for product routes
func New(db *gorm.DB, user_day_api *userday.API) *API {
	return &API{
		repository:   NewRepository(db),
		user_day_api: user_day_api,
	}
}

// --- TODO: Implement proper error handling ---
// Placeholder for error response function
func handleErr(w http.ResponseWriter, status int, message string, err error) {
	// Log the error internally
	fmt.Printf("ERROR [%d]: %s - %v\n", status, message, err)
	// Basic error response - replace with your actual error handling
	http.Error(w, message, status)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// --- TODO: Implement proper validation ---

// List godoc
//
//	@summary		List products
//	@description	List all available products
//	@tags			products
//	@accept			json
//	@produce		json
//	@success		200	{array}		DTO
//	@failure		500	{object}	string "Internal Server Error"
//	@router			/products [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	products, err := a.repository.List()
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to retrieve products", err)
		return
	}

	// Return empty array explicitly if no products, otherwise GORM might return null
	if len(products) == 0 {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "[]")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products.ToDto()); err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to encode products response", err)
		return
	}
}

// Create godoc
//
//	@summary		Create product
//	@description	Create a new product
//	@tags			products
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"Product form"
//	@success		201
//	@failure		400	{object}	string "Bad Request" // e.g., Invalid JSON
//	@failure		422	{object}	string "Unprocessable Entity" // e.g., Validation errors
//	@failure		500	{object}	string "Internal Server Error"
//	@router			/products [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if form.Grams == 0 {
		http.Error(w, "Grams must be greater than zero", http.StatusBadRequest)
		return
	}

	// --- TODO: Add validation for the form data ---
	// Example: if form.ProductName == "" { handleErr(...) return }

	newProduct := form.ToModel()
	newProduct.ID = uuid.New()

	// // Obliczenie wartości w przeliczeniu na 100g z użyciem float64
	// kcal := int(float64(form.Kcal) * 100 / float64(form.Grams))
	// proteins := int(float64(form.Proteins) * 100 / float64(form.Grams))
	// carbs := int(float64(form.Carbs) * 100 / float64(form.Grams))
	// fats := int(float64(form.Fats) * 100 / float64(form.Grams))

	// // Wstawienie do user_day
	// a.user_day_api.Create_from_product(form.UserID, kcal, proteins, carbs, fats)

	_, err := a.repository.Create(newProduct)
	if err != nil {
		// Consider more specific errors, e.g., duplicate product name if constraint exists
		handleErr(w, http.StatusInternalServerError, "Failed to create product", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read product
//	@description	Read a single product by ID
//	@tags			products
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"Product ID (UUID)"
//	@success		200	{object}	DTO
//	@failure		400	{object}	string "Bad Request" // e.g., Invalid UUID format
//	@failure		404	{object}	string "Not Found"
//	@failure		500	{object}	string "Internal Server Error"
//	@router			/products/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	paramValue := chi.URLParam(r, "id")
	if paramValue == "" {
		handleErr(w, http.StatusBadRequest, "Missing required path parameter: id", nil)
		return
	}

	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid UUID format for id parameter", err)
		return
	}

	product, err := a.repository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleErr(w, http.StatusNotFound, "Product not found", err)
		} else {
			handleErr(w, http.StatusInternalServerError, "Failed to read product", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product.ToDto()); err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to encode product response", err)
		return
	}
}

// Update godoc
//
//	@summary		Update product
//	@description	Update an existing product by ID
//	@tags			products
//	@accept			json
//	@produce		json
//	@param			id		path	string	true	"Product ID (UUID)"
//	@param			body	body	Form	true	"Product form"
//	@success		200		"Successfully updated" // Indicate success, maybe return updated object?
//	@failure		400	{object}	string "Bad Request" // e.g., Invalid UUID format or JSON
//	@failure		404	{object}	string "Not Found" // Product ID does not exist
//	@failure		422	{object}	string "Unprocessable Entity" // e.g., Validation errors
//	@failure		500	{object}	string "Internal Server Error"
//	@router			/products/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	paramValue := chi.URLParam(r, "id")
	if paramValue == "" {
		handleErr(w, http.StatusBadRequest, "Missing required path parameter: id", nil)
		return
	}

	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid UUID format for id parameter", err)
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// --- TODO: Add validation for the form data ---

	product := form.ToModel()
	product.ID = id // Set the ID from the path parameter

	rowsAffected, err := a.repository.Update(product)
	if err != nil {
		// GORM Update doesn't return ErrRecordNotFound directly on WHERE fail, check RowsAffected
		handleErr(w, http.StatusInternalServerError, "Failed to update product", err)
		return
	}

	if rowsAffected == 0 {
		// Could be not found, or data was identical. Check if exists first?
		// For simplicity, assume not found if 0 rows affected by update.
		// A Read call first could distinguish.
		handleErr(w, http.StatusNotFound, "Product not found or no changes made", nil)
		return
	}

	// Optionally return the updated product DTO or just status OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Product updated successfully") // Or return JSON
}

// Delete godoc
//
//	@summary		Delete product
//	@description	Soft delete a product by ID
//	@tags			products
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"Product ID (UUID)"
//	@success		200		"Successfully deleted" // Indicate success
//	@failure		400	{object}	string "Bad Request" // e.g., Invalid UUID format
//	@failure		404	{object}	string "Not Found" // Product ID does not exist or already deleted
//	@failure		500	{object}	string "Internal Server Error"
//	@router			/products/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	paramValue := chi.URLParam(r, "id")
	if paramValue == "" {
		handleErr(w, http.StatusBadRequest, "Missing required path parameter: id", nil)
		return
	}

	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid UUID format for id parameter", err)
		return
	}

	rowsAffected, err := a.repository.Delete(id)
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	if rowsAffected == 0 {
		// GORM soft delete returns 0 rows affected if already deleted or not found
		handleErr(w, http.StatusNotFound, "Product not found or already deleted", nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Product deleted successfully") // Or return JSON
}
