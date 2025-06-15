package userday

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	// Assuming a shared error handling package exists
	// "fitapp-backend/internal/err"
)

// API holds the dependencies for the userday handlers
type API struct {
	repository *Repository
}

// New creates a new API instance for userday routes
func New(db *gorm.DB) *API {
	return &API{
		repository: NewRepository(db),
	}
}

// --- TODO: Implement proper error handling ---
// Placeholder for error response function
func handleErr(w http.ResponseWriter, status int, message string, err error) {
	// Log the error internally
	fmt.Printf("ERROR [%d]: %s - %v\n", status, message, err)
	// Basic error response - replace with your actual error handling
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// Simple JSON error response
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// --- TODO: Implement proper validation ---

// List godoc
//
//	@summary		List user days
//	@description	List all user day records (consider pagination)
//	@tags			user-days
//	@accept			json
//	@produce		json
//	@success		200	{array}		DTO
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/user-days [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userDays, err := a.repository.List()
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to retrieve user days", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(userDays) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(userDays.ToDto()); err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to encode response", err)
		return
	}
}

// Create godoc
//
//	@summary		Create user day
//	@description	Create a new user day record
//	@tags			user-days
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"User Day form"
//	@success		201	{object}	DTO "Returns the created user day"
//	@failure		400	{object}	map[string]string "Bad Request" // e.g., Invalid JSON, UUID, Date format
//	@failure		422	{object}	map[string]string "Unprocessable Entity" // e.g., Validation errors, Duplicate entry
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/user-days [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid request body JSON", err)
		return
	}

	// --- Parse and Validate Input ---
	userID, err := uuid.Parse(form.UserID)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid UserID format (must be UUID)", err)
		return
	}

	userDate, err := time.Parse(DateFormat, form.UserDate)
	if err != nil {
		handleErr(w, http.StatusBadRequest, fmt.Sprintf("Invalid UserDate format (must be %s)", DateFormat), err)
		return
	}

	// --- TODO: Add more validation if needed (e.g., kcal >= 0) ---
	// Example: if form.DailyKcal < 0 { handleErr(w, http.StatusUnprocessableEntity, "kcal cannot be negative", nil); return }

	newUserDay := &UserDay{
		ID:            uuid.New(), // Generate new primary key
		UserID:        userID,
		UserDate:      userDate,
		DailyKcal:     form.DailyKcal,
		DailyProteins: form.DailyProteins,
		DailyCarbs:    form.DailyCarbs,
		DailyFats:     form.DailyFats,
	}

	createdUserDay, err := a.repository.Create(newUserDay)
	if err != nil {
		// Check for specific DB errors like unique constraint violation
		// Example: if strings.Contains(err.Error(), "UNIQUE constraint failed") { ... }
		handleErr(w, http.StatusUnprocessableEntity, "Failed to create user day (possible duplicate UserID/Date?)", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUserDay.ToDto()); err != nil {
		// Log error, but the resource was created
		fmt.Printf("Error encoding created user_day response: %v\n", err)
	}
}

func (a *API) Create_from_product(uid string, kcal int, proteins int, carbs int, fats int) {

	fmt.Print("chujjj")
	userID, err := uuid.Parse(uid)
	if err != nil {
		return
	}

	toUpdate, err := a.repository.FindByUserAndDate(userID, time.Now())
	if err == nil {
		fmt.Print("Already exists")
		userDayToUpdate := &UserDay{
			ID:            toUpdate.ID,
			DailyKcal:     toUpdate.DailyKcal + kcal,
			DailyProteins: toUpdate.DailyProteins + proteins,
			DailyCarbs:    toUpdate.DailyCarbs + carbs,
			DailyFats:     toUpdate.DailyFats + fats,
			// UserID and UserDate are NOT updated via this method
		}

		rowsAffected, err := a.repository.Update(userDayToUpdate)
		if err != nil {
			fmt.Print("Failed to update user day")
			return
		}

		if rowsAffected == 0 {
			// Could be not found, or data was identical. Assume not found for simplicity.
			fmt.Print("User day record not found or no changes detected")
			return
		}
		return
	}

	newUserDay := &UserDay{
		ID:            uuid.New(), // Generate new primary key
		UserID:        userID,
		UserDate:      time.Now(),
		DailyKcal:     kcal,
		DailyProteins: proteins,
		DailyCarbs:    carbs,
		DailyFats:     fats,
	}
	if _, err := a.repository.Create(newUserDay); err != nil {
		// Log error, but the resource was created
		fmt.Printf("Error encoding created user_day response: %v\n", err)
	}

	fmt.Printf("%v", newUserDay)
}

// Read godoc
//
//	@summary		Read user day by ID
//	@description	Read a single user day record by its primary ID
//	@tags			user-days
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"UserDay Primary ID (UUID)"
//	@success		200	{object}	DTO
//	@failure		400	{object}	map[string]string "Bad Request" // e.g., Invalid UUID format
//	@failure		404	{object}	map[string]string "Not Found"
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/user-days/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	paramValue := chi.URLParam(r, "id")
	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid ID format (must be UUID)", err)
		return
	}

	userDay, err := a.repository.Read(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleErr(w, http.StatusNotFound, "User day record not found", err)
		} else {
			handleErr(w, http.StatusInternalServerError, "Failed to read user day", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userDay.ToDto()); err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to encode response", err)
		return
	}
}

// FindByUserAndDate godoc
//
//	@summary		Find user day by user and date
//	@description	Find a single user day record by UserID and Date
//	@tags			user-days
//	@accept			json
//	@produce		json
//	@param			userId	query		string	true	"User ID (UUID)" format(uuid)
//	@param			date	query		string	true	"Date (YYYY-MM-DD)" format(date)
//	@success		200	{object}	DTO
//	@failure		400	{object}	map[string]string "Bad Request" // e.g., Missing/Invalid query params or format
//	@failure		404	{object}	map[string]string "Not Found"
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/user-days/search [get] // Use a different path or rely on query params
func (a *API) FindByUserAndDate(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userIdParam := r.URL.Query().Get("userId")
	dateParam := r.URL.Query().Get("date")

	if userIdParam == "" || dateParam == "" {
		handleErr(w, http.StatusBadRequest, "Missing required query parameters: userId, date", nil)
		return
	}

	userID, err := uuid.Parse(userIdParam)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid userId format (must be UUID)", err)
		return
	}

	userDate, err := time.Parse(DateFormat, dateParam)
	if err != nil {
		handleErr(w, http.StatusBadRequest, fmt.Sprintf("Invalid date format (must be %s)", DateFormat), err)
		return
	}

	userDay, err := a.repository.FindByUserAndDate(userID, userDate)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleErr(w, http.StatusNotFound, "User day record not found for specified user and date", err)
		} else {
			handleErr(w, http.StatusInternalServerError, "Failed to find user day", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userDay.ToDto()); err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to encode response", err)
		return
	}
}

// Update godoc
//
//	@summary		Update user day
//	@description	Update nutritional info for an existing user day record by its primary ID
//	@tags			user-days
//	@accept			json
//	@produce		json
//	@param			id		path	string	true	"UserDay Primary ID (UUID)"
//	@param			body	body	Form	true	"User Day form (only nutritional fields are used)"
//	@success		200		{object}	map[string]string "Successfully updated"
//	@failure		400	{object}	map[string]string "Bad Request" // e.g., Invalid ID/JSON
//	@failure		404	{object}	map[string]string "Not Found"
//	@failure		422	{object}	map[string]string "Unprocessable Entity" // e.g., Validation errors
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/user-days/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	paramValue := chi.URLParam(r, "id")
	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid ID format (must be UUID)", err)
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid request body JSON", err)
		return
	}

	// --- TODO: Add validation for nutritional values (e.g., >= 0) ---

	// Create a model instance only with fields to be updated + ID for WHERE clause
	userDayToUpdate := &UserDay{
		ID:            id,
		DailyKcal:     form.DailyKcal,
		DailyProteins: form.DailyProteins,
		DailyCarbs:    form.DailyCarbs,
		DailyFats:     form.DailyFats,
		// UserID and UserDate are NOT updated via this method
	}

	rowsAffected, err := a.repository.Update(userDayToUpdate)
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to update user day", err)
		return
	}

	if rowsAffected == 0 {
		// Could be not found, or data was identical. Assume not found for simplicity.
		handleErr(w, http.StatusNotFound, "User day record not found or no changes detected", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User day updated successfully"})
}

// Delete godoc
//
//	@summary		Delete user day
//	@description	Soft delete a user day record by its primary ID
//	@tags			user-days
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"UserDay Primary ID (UUID)"
//	@success		200		{object}	map[string]string "Successfully deleted"
//	@failure		400	{object}	map[string]string "Bad Request" // e.g., Invalid ID format
//	@failure		404	{object}	map[string]string "Not Found"
//	@failure		500	{object}	map[string]string "Internal Server Error"
//	@router			/user-days/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	paramValue := chi.URLParam(r, "id")
	id, err := uuid.Parse(paramValue)
	if err != nil {
		handleErr(w, http.StatusBadRequest, "Invalid ID format (must be UUID)", err)
		return
	}

	rowsAffected, err := a.repository.Delete(id)
	if err != nil {
		handleErr(w, http.StatusInternalServerError, "Failed to delete user day", err)
		return
	}

	if rowsAffected == 0 {
		handleErr(w, http.StatusNotFound, "User day record not found or already deleted", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User day deleted successfully"})
}
