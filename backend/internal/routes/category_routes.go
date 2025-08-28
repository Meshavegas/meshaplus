package routes

import (
	"backend/internal/handler"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupCategoryRoutes configure les routes pour la catégorisation
func SetupCategoryRoutes(r chi.Router, categoryHandler *handler.CategoryHandler, authMiddleware *middleware.AuthMiddleware) {
	// Groupe de routes pour les catégories (protégées par authentification)
	r.Route("/categories", func(r chi.Router) {
		// Appliquer l'authentification à toutes les routes
		r.Use(authMiddleware.Authenticate)

		// Routes pour la catégorisation
		r.Post("/categorize", categoryHandler.CategorizeItem) // POST /api/v1/categories/categorize
		r.Post("/", categoryHandler.CreateCategory)           // POST /api/v1/categories
		r.Get("/", categoryHandler.GetCategoriesByType)       // GET /api/v1/categories?categoryType=expense
		r.Get("/{id}", categoryHandler.GetCategoryByID)       // GET /api/v1/categories/{id}
	})
}
