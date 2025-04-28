package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10" // Import validator
	"github.com/google/uuid"
	"gorm.io/gorm"
	// Assuming a shared error handling package exists
	// "fitapp-backend/internal/err"
	// "fitapp-backend/pkg/web" // For helper functions like Respond, RespondError
)

// API holds the dependencies for the user handlers
type API struct {
	repository *Repository
	validate   *validator.Validate // Add validator instance
}

// New creates a new API instance for user routes
func New(db *gorm.DB) *API {
	return &API{
		repository: NewRepository(db),
		validate:   validator.New(), // Initialize validator
	}
}

// --- Simple Error Handling Placeholder ---
// Replace with your application's standard error response mechanism
func handleErr(w http.ResponseWriter, status int, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	response := map[string]interface{}{"error": message}
	if details != nil {
		response["details"] = details
	}
	json.NewEncoder(w).Encode(response)
	// Log the error as well
	fmt.Printf("ERROR [%d]: %s - Details: %v\n", status, message, details)
}

// Helper to format validation errors
func formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fmt.Sprintf("failed validation on '%s'", fieldErr.Tag())
		}
	}
	return errors
}

// List godoc
//
//	@summary		List users
//	@description	List all non-deleted users
//	@tags			users
//	@accept			json
//	@produce		json
//	@success		200	{array}		DTO
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/users [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	users, err := a.repository.List()
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to retrieve users", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(users) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(users.ToDto()); err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to encode users response", err.Error())
		return
	}
}

// Create godoc
//
//	@summary		Create user
//	@description	Create a new user with profile details
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"User creation form"
//	@success		201	{object}	DTO "Returns the created user"
//	@failure		400	{object}	map[string]interface{} "Bad Request (Invalid JSON or Validation Errors)"
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/users [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid request body JSON", err.Error())
		return
	}

	// Validate the form using the validator
	if err := a.validate.Struct(form); err != nil {
		handleErr(w, http.StatusBadRequest, "Validation failed", formatValidationErrors(err))
		return
	}
	// Explicitly check Sex pointer after basic validation
	if form.Sex == nil {
		handleErr(w, http.StatusBadRequest, "Validation failed", map[string]string{"Sex": "field is required"})
		return
	}

	// Map form data to the model
	newUser := &User{
		ID:       uuid.New(), // Generate new primary key
		Username: form.Username,
		FullName: form.FullName,
		Sex:      *form.Sex, // Dereference pointer after validation ensures it's not nil
		Height:   form.Height,
		Weight:   form.Weight,
		Age:      form.Age,
	}

	createdUser, err := a.repository.Create(newUser)
	if err != nil {
		// Could check for specific DB errors like unique username constraint if applicable
		handleErr(w, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser.ToDto()); err != nil {
		// Log error, but resource was created
		fmt.Printf("Error encoding created user response: %v\n", err)
	}
}

// Read godoc
//
//	@summary		Read user
//	@description	Read a single user by ID
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"User ID (UUID)"
//	@success		200	{object}	DTO
//	@failure		400	{object}	map[string]string "Bad Request (Invalid UUID)"
//	@failure		404	{object}	map[string]string "Not Found"
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/users/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	paramValue := chi.URLParam(r, "id")
	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid user ID format", err.Error())
		return
	}

	user, err := a.repository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleErr(w, http.StatusNotFound, "User not found", nil)
		} else {
			handleErr(w, http.StatusInternalServerError, "Failed to read user", err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user.ToDto()); err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to encode user response", err.Error())
		return
	}
}

// Update godoc
//
//	@summary		Update user
//	@description	Update an existing user's profile details by ID
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id		path	string	true	"User ID (UUID)"
//	@param			body	body	Form	true	"User update form"
//	@success		200		{object}	DTO "Returns the updated user"
//	@failure		400	{object}	map[string]interface{} "Bad Request (Invalid ID, JSON or Validation Errors)"
//	@failure		404	{object}	map[string]string "Not Found"
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/users/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	paramValue := chi.URLParam(r, "id")
	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid user ID format", err.Error())
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid request body JSON", err.Error())
		return
	}

	// Validate the form
	if err := a.validate.Struct(form); err != nil {
		handleErr(w, http.StatusBadRequest, "Validation failed", formatValidationErrors(err))
		return
	}
	if form.Sex == nil {
		handleErr(w, http.StatusBadRequest, "Validation failed", map[string]string{"Sex": "field is required"})
		return
	}

	// Create the model instance for update
	userToUpdate := &User{
		ID:       id, // Set ID from path parameter
		Username: form.Username,
		FullName: form.FullName,
		Sex:      *form.Sex,
		Height:   form.Height,
		Weight:   form.Weight,
		Age:      form.Age,
	}

	rowsAffected, err := a.repository.Update(userToUpdate)
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to update user", err.Error())
		return
	}

	if rowsAffected == 0 {
		// Could be not found, or data was identical. Check existence first for better 404.
		// Read the user again to confirm existence and return correct status/message.
		_, readErr := a.repository.Read(id)
		if readErr == gorm.ErrRecordNotFound {
			handleErr(w, http.StatusNotFound, "User not found", nil)
		} else if readErr != nil {
			// Handle unexpected error during read check
			handleErr(w, http.StatusInternalServerError, "Failed to verify user update", readErr.Error())
		} else {
			// User exists, but no rows were affected (likely data was the same)
			// Still return the user's current state as success.
			updatedUser, _ := a.repository.Read(id) // Read again to get potentially auto-updated UpdatedAt
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(updatedUser.ToDto())
		}
		return // Exit after handling 0 rows affected case
	}

	// If update was successful (rowsAffected > 0), read the updated user to return it
	updatedUser, err := a.repository.Read(id)
	if err != nil {
		// This shouldn't happen if rowsAffected > 0, but handle defensively
		handleErr(w, http.StatusInternalServerError, "Failed to read user after update", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser.ToDto())
}

// Delete godoc
//
//	@summary		Delete user
//	@description	Soft delete a user by ID
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"User ID (UUID)"
//	@success		200		{object}	map[string]string "Successfully deleted"
//	@failure		400	{object}	map[string]string "Bad Request (Invalid ID)"
//	@failure		404	{object}	map[string]string "Not Found"
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/users/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	paramValue := chi.URLParam(r, "id")
	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid user ID format", err.Error())
		return
	}

	rowsAffected, err := a.repository.Delete(id)
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	if rowsAffected == 0 {
		handleErr(w, http.StatusNotFound, "User not found or already deleted", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
