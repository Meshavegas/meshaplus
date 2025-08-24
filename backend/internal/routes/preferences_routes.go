package routes

import (
	"backend/internal/handler"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupPreferencesRoutes configure les routes pour les préférences utilisateur
func SetupPreferencesRoutes(r chi.Router, preferencesHandler *handler.PreferencesHandler, authMiddleware *middleware.AuthMiddleware) {
	// Groupe de routes pour les préférences (protégées par authentification)
	r.Route("/preferences", func(r chi.Router) {
		// Appliquer l'authentification à toutes les routes
		r.Use(authMiddleware.Authenticate)

		// Routes pour les préférences
		r.Post("/", preferencesHandler.CreatePreferences)     // POST /api/v1/preferences
		r.Get("/", preferencesHandler.GetPreferences)         // GET /api/v1/preferences
		r.Put("/", preferencesHandler.UpdatePreferences)      // PUT /api/v1/preferences
		r.Delete("/", preferencesHandler.DeletePreferences)   // DELETE /api/v1/preferences
		r.Get("/content", preferencesHandler.GenerateContent) // GET /api/v1/preferences/content?contentType=budget_tip
	})
}
