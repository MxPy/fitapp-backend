package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type API struct {
	repository *Repository
}

func New(db *gorm.DB) *API {
	return &API{
		repository: NewRepository(db),
	}
}

// List godoc
//
//	@summary        List user
//	@description    List user
//	@tags           users
//	@accept         json
//	@produce        json
//	@success        200 {array}     DTO
//	@failure        500 {object}    err.Error
//	@router         /users [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	users, err := a.repository.List()
	if err != nil {
		// handle later
		return
	}

	if len(users) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(users.ToDto()); err != nil {
		// handle later
		return
	}
}

// Create godoc
//
//	@summary        Create user
//	@description    Create user
//	@tags           users
//	@accept         json
//	@produce        json
//	@param          body    body    Form    true    "User form"
//	@success        201
//	@failure        400 {object}    err.Error
//	@failure        422 {object}    err.Errors
//	@failure        500 {object}    err.Error
//	@router         /users [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		// handle later
		return
	}

	newUser := form.ToModel()
	newUser.ID = uuid.New()

	_, err := a.repository.Create(newUser)
	if err != nil {
		// handle later
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary        Read user
//	@description    Read user
//	@tags           users
//	@accept         json
//	@produce        json
//	@param          id	path        string  true    "User ID"
//	@success        200 {object}    DTO
//	@failure        400 {object}    err.Error
//	@failure        404
//	@failure        500 {object}    err.Error
//	@router         /users/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("========== READ REQUEST ==========\n")
	fmt.Printf("Request URL: %s\n", r.URL.String())
	fmt.Printf("Request Method: %s\n", r.Method)

	paramValue := chi.URLParam(r, "id")
	fmt.Printf("URL parameter 'id': '%s', length: %d\n", paramValue, len(paramValue))

	if paramValue == "" {
		fmt.Printf("ERROR: Empty ID parameter\n")
		http.Error(w, "Missing required path parameter: id", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(paramValue)
	if err != nil {
		fmt.Printf("ERROR parsing UUID: %v\n", err)
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	fmt.Printf("Parsed UUID: %v\n", id)

	user, err := a.repository.Read(id)
	if err != nil {
		fmt.Printf("ERROR reading from repository: %v\n", err)
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("Record not found\n")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("User found: %+v\n", user)

	dto := user.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		fmt.Printf("ERROR encoding response: %v\n", err)
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Response sent successfully\n")
	fmt.Printf("======== END READ REQUEST ========\n")
}

// Update godoc
//
//	@summary        Update user
//	@description    Update user
//	@tags           users
//	@accept         json
//	@produce        json
//	@param          id      path    string  true    "User ID"
//	@param          body    body    Form    true    "User form"
//	@success        200
//	@failure        400 {object}    err.Error
//	@failure        404
//	@failure        422 {object}    err.Errors
//	@failure        500 {object}    err.Error
//	@router         /users/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		// handle later
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		// handle later
		return
	}

	user := form.ToModel()
	user.ID = id

	rows, err := a.repository.Update(user)
	if err != nil {
		// handle later
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// Delete godoc
//
//	@summary        Delete user
//	@description    Delete user
//	@tags           users
//	@accept         json
//	@produce        json
//	@param          id  path    string  true    "User ID"
//	@success        200
//	@failure        400 {object}    err.Error
//	@failure        404
//	@failure        500 {object}    err.Error
//	@router         /users/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		// handle later
		return
	}

	rows, err := a.repository.Delete(id)
	if err != nil {
		// handle later
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
