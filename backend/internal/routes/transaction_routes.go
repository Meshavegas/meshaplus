package routes

import (
	"backend/internal/handler"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupTransactionRoutes configure les routes pour les transactions
func SetupTransactionRoutes(r chi.Router, transactionHandler *handler.TransactionHandler, authMiddleware *middleware.AuthMiddleware) {
	// Groupe de routes pour les transactions (protégées par authentification)
	r.Route("/transactions", func(r chi.Router) {
		// Appliquer l'authentification à toutes les routes
		r.Use(authMiddleware.Authenticate)

		// Routes pour la gestion des transactions
		r.Post("/", transactionHandler.CreateTransaction)       // POST /api/v1/transactions
		r.Get("/", transactionHandler.GetTransactions)          // GET /api/v1/transactions
		r.Get("/stats", transactionHandler.GetTransactionStats) // GET /api/v1/transactions/stats
		r.Get("/{id}", transactionHandler.GetTransaction)       // GET /api/v1/transactions/{id}
		r.Put("/{id}", transactionHandler.UpdateTransaction)    // PUT /api/v1/transactions/{id}
		r.Delete("/{id}", transactionHandler.DeleteTransaction) // DELETE /api/v1/transactions/{id}
	})
}
