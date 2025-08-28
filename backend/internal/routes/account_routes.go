package routes

import (
	"backend/internal/handler"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupAccountRoutes configure les routes pour les comptes
func SetupAccountRoutes(r chi.Router, accountHandler *handler.AccountHandler, authMiddleware *middleware.AuthMiddleware) {
	// Groupe de routes pour les comptes (protégées par authentification)
	r.Route("/accounts", func(r chi.Router) {
		// Appliquer l'authentification à toutes les routes
		r.Use(authMiddleware.Authenticate)

		// Routes pour la gestion des comptes
		r.Post("/", accountHandler.CreateAccount)                // POST /api/v1/accounts
		r.Get("/", accountHandler.GetAccounts)                   // GET /api/v1/accounts
		r.Get("/{id}", accountHandler.GetAccount)                // GET /api/v1/accounts/{id}
		r.Put("/{id}", accountHandler.UpdateAccount)             // PUT /api/v1/accounts/{id}
		r.Delete("/{id}", accountHandler.DeleteAccount)          // DELETE /api/v1/accounts/{id}
		r.Get("/{id}/balance", accountHandler.GetAccountBalance) // GET /api/v1/accounts/{id}/balance
		r.Get("/{id}/details", accountHandler.GetAccountDetails) // GET /api/v1/accounts/{id}/details
	})
}
