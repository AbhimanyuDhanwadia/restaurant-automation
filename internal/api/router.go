package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	"github.com/restaurantautomation/api/internal/api/handlers"
	appm "github.com/restaurantautomation/api/internal/api/middleware"
	"github.com/restaurantautomation/api/internal/config"
)

// NewRouter constructs the Chi router with the standard middleware stack
// and registers all API routes.
func NewRouter(cfg *config.Config, log zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	// 1. Basic Middleware
	r.Use(middleware.RealIP)
	r.Use(appm.RequestID)
	r.Use(appm.Logger(log))
	r.Use(middleware.Recoverer)

	// 2. CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// 3. System Routes (Unauthenticated)
	r.Get("/health", handlers.Health())
	r.Get("/ready", handlers.Ready())

	// 4. API Routes (Will be authenticated later)
	r.Route("/api/v1", func(r chi.Router) {
		// Placeholder for order routes (Milestone 4)
		r.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
		})
	})

	return r
}
