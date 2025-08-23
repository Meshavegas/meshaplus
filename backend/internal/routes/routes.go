package routes

import (
	"backend/internal/handler"
	"backend/internal/service"
	"backend/internal/usecase"
	"backend/pkg/logger"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes configure toutes les routes de l'application
func SetupRoutes(
	r chi.Router,
	userUsecase *usecase.UserUsecase,
	authService *service.AuthService,
	authMiddleware *middleware.AuthMiddleware,
	taskHandler *handler.TaskHandler,
	transactionHandler *handler.TransactionHandler,
	accountHandler *handler.AccountHandler,
	budgetHandler *handler.BudgetHandler,
	savingGoalHandler *handler.SavingGoalHandler,
	categoryHandler *handler.CategoryHandler,
	logger logger.Logger,
) {
	// Routes pour la documentation Swagger (publiques) - à la racine
	SetupSwaggerRoutes(r)

	// Préfixe API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Routes d'authentification (publiques et protégées)
		SetupAuthRoutes(r, authService, authMiddleware, logger)

		// Routes pour les tâches (protégées)
		SetupTaskRoutes(r, taskHandler, authMiddleware)

		// Routes pour les transactions (protégées)
		SetupTransactionRoutes(r, transactionHandler, authMiddleware)

		// Routes pour les comptes (protégées)
		SetupAccountRoutes(r, accountHandler, authMiddleware)

		// Routes pour les budgets (protégées)
		SetupBudgetRoutes(r, budgetHandler, authMiddleware)

		// Routes pour les objectifs d'épargne (protégées)
		SetupSavingGoalRoutes(r, savingGoalHandler, authMiddleware)

		// Routes pour la catégorisation (protégées)
		SetupCategoryRoutes(r, categoryHandler, authMiddleware)

		// TODO: Ajouter d'autres routes selon les besoins
		// SetupUserRoutes(r, userHandler, authMiddleware)
		// SetupFileRoutes(r, fileHandler, authMiddleware)
		// SetupHealthRoutes(r, healthHandler)
	})
}
