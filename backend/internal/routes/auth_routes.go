package routes

import (
	"backend/internal/handler"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupAuthRoutes configure les routes pour l'authentification
func SetupAuthRoutes(r chi.Router, authService *service.AuthService, authMiddleware *middleware.AuthMiddleware, logger logger.Logger) {
	authHandler := handler.NewAuthHandler(authService, logger)

	r.Route("/auth", func(r chi.Router) {
		// Routes publiques (sans authentification)
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.RefreshToken)

		// Routes protégées (avec authentification)
		r.With(authMiddleware.Authenticate).Get("/me", authHandler.GetCurrentUser)
	})
}
