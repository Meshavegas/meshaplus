package routes

import (
	"backend/internal/handler"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupFinanceRoutes configure les routes pour le tableau de bord financier
func SetupFinanceRoutes(r chi.Router, financeDashboardHandler *handler.FinanceDashboardHandler, authMiddleware *middleware.AuthMiddleware) {
	r.Route("/finance", func(r chi.Router) {
		// Appliquer le middleware d'authentification à toutes les routes
		r.Use(authMiddleware.Authenticate)

		// Route pour le tableau de bord financier
		r.Get("/dashboard", financeDashboardHandler.GetFinanceDashboard)
	})
}
