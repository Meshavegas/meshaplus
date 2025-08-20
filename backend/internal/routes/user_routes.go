package routes

import (
	"backend/internal/handler"
	"backend/internal/usecase"
	"backend/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// SetupUserRoutes configure les routes pour les utilisateurs
func SetupUserRoutes(r chi.Router, userUsecase *usecase.UserUsecase, logger logger.Logger) {
	userHandler := handler.NewUserHandler(userUsecase, logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.GetUsers)
		r.Get("/{id}", userHandler.GetUserByID)
		r.Put("/{id}", userHandler.UpdateUser)
		r.Delete("/{id}", userHandler.DeleteUser)
	})
}
