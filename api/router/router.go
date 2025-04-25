package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"gorm.io/gorm"

	"fitapp-backend/api/resource/health"
	"fitapp-backend/api/resource/user"
)

func New(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)
	r.Get("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		// Użyj http.ServeFile, aby wysłać zawartość pliku .swagger/doc.yaml
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	// 2. Dodaj handler Swagger UI.
	// Ten handler serwuje pliki HTML/JS/CSS dla Swagger UI.
	// Parametr httpSwagger.URL wskazuje, gdzie UI ma pobrać definicję API (czyli nasz doc.yaml)
	// Ścieżka musi kończyć się na "/*" dla poprawnego routingu w chi
	r.Get("/swagger/*", httpSwagger.Handler(
		// Wskazujesz pełny URL (schemat + host + port + ścieżka) do pliku doc.yaml,
		// który serwujesz w punkcie 1.
		httpSwagger.URL("http://localhost:8080/swagger/doc.yaml"),
	))
	r.Route("/v1", func(r chi.Router) {
		userAPI := user.New(db)
		r.Get("/users", userAPI.List)
		r.Post("/users", userAPI.Create)
		r.Get("/users/{id}", userAPI.Read)
		r.Put("/users/{id}", userAPI.Update)
		r.Delete("/users/{id}", userAPI.Delete)
	})

	return r
}
