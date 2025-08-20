package routes

import (
	"backend/internal/handler"
	"backend/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// SetupTaskRoutes configure les routes pour les tâches
func SetupTaskRoutes(r chi.Router, taskHandler *handler.TaskHandler, authMiddleware *middleware.AuthMiddleware) {
	// Groupe de routes pour les tâches (protégées par authentification)
	r.Route("/tasks", func(r chi.Router) {
		// Appliquer l'authentification à toutes les routes
		r.Use(authMiddleware.Authenticate)

		// Routes pour la gestion des tâches
		r.Post("/", taskHandler.CreateTask)                // POST /api/v1/tasks
		r.Get("/", taskHandler.GetUserTasks)               // GET /api/v1/tasks
		r.Get("/stats", taskHandler.GetTaskStats)          // GET /api/v1/tasks/stats
		r.Get("/{id}", taskHandler.GetTask)                // GET /api/v1/tasks/{id}
		r.Put("/{id}", taskHandler.UpdateTask)             // PUT /api/v1/tasks/{id}
		r.Delete("/{id}", taskHandler.DeleteTask)          // DELETE /api/v1/tasks/{id}
		r.Post("/{id}/complete", taskHandler.CompleteTask) // POST /api/v1/tasks/{id}/complete
	})
}
