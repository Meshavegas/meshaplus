package routes

import (
	"github.com/go-chi/chi/v5"
)

// SetupHealthRoutes configure les routes pour la santé de l'application
// TODO: Implémenter quand HealthHandler sera créé
func SetupHealthRoutes(r chi.Router, healthHandler interface{}) {
	// r.Route("/health", func(r chi.Router) {
	// 	r.Get("/", healthHandler.HealthCheck)
	// })
}
