package routes

import (
	"backend/internal/handler"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupBudgetRoutes configure les routes pour les budgets
func SetupBudgetRoutes(r chi.Router, budgetHandler *handler.BudgetHandler, authMiddleware *middleware.AuthMiddleware) {
	// Groupe de routes pour les budgets (protégées par authentification)
	r.Route("/budgets", func(r chi.Router) {
		// Appliquer l'authentification à toutes les routes
		r.Use(authMiddleware.Authenticate)

		// Routes pour la gestion des budgets
		r.Post("/", budgetHandler.CreateBudget)       // POST /api/v1/budgets
		r.Get("/", budgetHandler.GetBudgets)          // GET /api/v1/budgets
		r.Get("/stats", budgetHandler.GetBudgetStats) // GET /api/v1/budgets/stats
		r.Get("/{id}", budgetHandler.GetBudget)       // GET /api/v1/budgets/{id}
		r.Put("/{id}", budgetHandler.UpdateBudget)    // PUT /api/v1/budgets/{id}
		r.Delete("/{id}", budgetHandler.DeleteBudget) // DELETE /api/v1/budgets/{id}
	})
}
