package routes

import (
	"backend/internal/handler"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupSavingGoalRoutes configure les routes pour les objectifs d'épargne
func SetupSavingGoalRoutes(r chi.Router, savingGoalHandler *handler.SavingGoalHandler, authMiddleware *middleware.AuthMiddleware) {
	// Groupe de routes pour les objectifs d'épargne (protégées par authentification)
	r.Route("/saving-goals", func(r chi.Router) {
		// Appliquer l'authentification à toutes les routes
		r.Use(authMiddleware.Authenticate)

		// Routes pour la gestion des objectifs d'épargne
		r.Post("/", savingGoalHandler.CreateSavingGoal)       // POST /api/v1/saving-goals
		r.Get("/", savingGoalHandler.GetSavingGoals)          // GET /api/v1/saving-goals
		r.Get("/{id}", savingGoalHandler.GetSavingGoal)       // GET /api/v1/saving-goals/{id}
		r.Put("/{id}", savingGoalHandler.UpdateSavingGoal)    // PUT /api/v1/saving-goals/{id}
		r.Delete("/{id}", savingGoalHandler.DeleteSavingGoal) // DELETE /api/v1/saving-goals/{id}
	})
}
