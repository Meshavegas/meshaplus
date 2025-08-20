package routes

import (
	"backend/pkg/logger"

	"github.com/go-chi/chi/v5"
)

// SetupFileRoutes configure les routes pour les fichiers
// TODO: Implémenter quand FileUsecase et FileHandler seront créés
func SetupFileRoutes(r chi.Router, fileUsecase interface{}, logger logger.Logger) {
	// fileHandler := handler.NewFileHandler(fileUsecase, logger)

	// r.Route("/files", func(r chi.Router) {
	// 	r.Post("/upload", fileHandler.UploadFile)
	// 	r.Get("/{id}", fileHandler.GetFile)
	// 	r.Delete("/{id}", fileHandler.DeleteFile)
	// })
}
