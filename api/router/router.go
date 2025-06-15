package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"gorm.io/gorm"

	"fitapp-backend/api/resource/health"
	"fitapp-backend/api/resource/product"
	"fitapp-backend/api/resource/user"
	userday "fitapp-backend/api/resource/user_day"
)

func New(db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)
	r.Get("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		// Ustaw nagłówki HTTP, aby zapobiec cache'owaniu przez przeglądarkę
		// Te nagłówki dają silne instrukcje, aby nie cache'ować
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies

		// Użyj http.ServeFile, aby wysłać zawartość pliku
		// Upewnij się, że ścieżka "docs/swagger.yaml" jest poprawna względem katalogu roboczego aplikacji
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	// 2. Handler Swagger UI (bez zmian)
	r.Get("/swagger/*", httpSwagger.Handler(
		// Upewnij się, że URL jest poprawny (schemat, host, port, ścieżka)
		httpSwagger.URL("http://localhost:8080/swagger/doc.yaml"),
	))

	r.Route("/v1", func(r chi.Router) {
		userAPI := user.New(db)
		r.Get("/users", userAPI.List)
		r.Post("/users", userAPI.Create)
		r.Get("/users/{id}", userAPI.Read)
		r.Put("/users/{id}", userAPI.Update)
		r.Delete("/users/{id}", userAPI.Delete)
		userdayAPI := userday.New(db)
		r.Get("/user-days", userdayAPI.List)
		r.Post("/user-days", userdayAPI.Create)
		r.Get("/user-days/{id}", userdayAPI.Read)
		r.Get("/user-days/search", userdayAPI.FindByUserAndDate)
		r.Put("/user-days/{id}", userdayAPI.Update)
		r.Delete("/user-days/{id}", userdayAPI.Delete)
		productAPI := product.New(db, userdayAPI)
		r.Get("/products", productAPI.List)
		r.Post("/products", productAPI.Create)
		r.Get("/products/{id}", productAPI.Read)
		r.Put("/products/{id}", productAPI.Update)
		r.Delete("/products/{id}", productAPI.Delete)

	})
	return r
}
