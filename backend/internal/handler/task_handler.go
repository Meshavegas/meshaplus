package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/domaine/entity"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TaskHandler gère les requêtes HTTP pour les tâches
type TaskHandler struct {
	taskService *service.TaskService
	logger      logger.Logger
}

// NewTaskHandler crée une nouvelle instance de TaskHandler
func NewTaskHandler(taskService *service.TaskService, logger logger.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		logger:      logger,
	}
}

// CreateTask crée une nouvelle tâche
// @Summary Créer une nouvelle tâche
// @Description Crée une nouvelle tâche pour l'utilisateur authentifié
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body entity.CreateTaskRequest true "Données de la tâche"
// @Success 201 {object} response.Response "Tâche créée"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	var req entity.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	task, err := h.taskService.CreateTask(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("Erreur création tâche", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur création tâche", err)
		return
	}

	response.Success(w, http.StatusCreated, "Tâche créée avec succès", task)
}

// GetTask récupère une tâche par son ID
// @Summary Récupérer une tâche
// @Description Récupère une tâche spécifique par son ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la tâche"
// @Success 200 {object} response.Response "Tâche récupérée"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 403 {object} response.ErrorResponse "Accès non autorisé"
// @Failure 404 {object} response.ErrorResponse "Tâche non trouvée"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID de la tâche depuis l'URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de tâche invalide", err)
		return
	}

	task, err := h.taskService.GetTask(r.Context(), userID, taskID)
	if err != nil {
		switch err.Error() {
		case "tâche non trouvée":
			response.Error(w, http.StatusNotFound, "Tâche non trouvée", err)
		case "accès non autorisé":
			response.Error(w, http.StatusForbidden, "Accès non autorisé", err)
		default:
			h.logger.Error("Erreur récupération tâche", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur récupération tâche", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Tâche récupérée avec succès", task)
}

// GetUserTasks récupère toutes les tâches de l'utilisateur
// @Summary Récupérer toutes les tâches
// @Description Récupère toutes les tâches de l'utilisateur authentifié
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Tâches récupérées"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /tasks [get]
func (h *TaskHandler) GetUserTasks(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	tasks, err := h.taskService.GetUserTasks(r.Context(), userID)
	if err != nil {
		h.logger.Error("Erreur récupération tâches", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération tâches", err)
		return
	}

	response.Success(w, http.StatusOK, "Tâches récupérées avec succès", tasks)
}

// UpdateTask met à jour une tâche
// @Summary Mettre à jour une tâche
// @Description Met à jour une tâche existante
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la tâche"
// @Param task body entity.UpdateTaskRequest true "Données de mise à jour"
// @Success 200 {object} response.Response "Tâche mise à jour"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 403 {object} response.ErrorResponse "Accès non autorisé"
// @Failure 404 {object} response.ErrorResponse "Tâche non trouvée"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID de la tâche depuis l'URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de tâche invalide", err)
		return
	}

	var req entity.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	task, err := h.taskService.UpdateTask(r.Context(), userID, taskID, req)
	if err != nil {
		switch err.Error() {
		case "tâche non trouvée":
			response.Error(w, http.StatusNotFound, "Tâche non trouvée", err)
		case "accès non autorisé":
			response.Error(w, http.StatusForbidden, "Accès non autorisé", err)
		default:
			h.logger.Error("Erreur mise à jour tâche", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur mise à jour tâche", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Tâche mise à jour avec succès", task)
}

// DeleteTask supprime une tâche
// @Summary Supprimer une tâche
// @Description Supprime une tâche existante
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la tâche"
// @Success 200 {object} response.Response "Tâche supprimée"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 403 {object} response.ErrorResponse "Accès non autorisé"
// @Failure 404 {object} response.ErrorResponse "Tâche non trouvée"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID de la tâche depuis l'URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de tâche invalide", err)
		return
	}

	err = h.taskService.DeleteTask(r.Context(), userID, taskID)
	if err != nil {
		switch err.Error() {
		case "tâche non trouvée":
			response.Error(w, http.StatusNotFound, "Tâche non trouvée", err)
		case "accès non autorisé":
			response.Error(w, http.StatusForbidden, "Accès non autorisé", err)
		default:
			h.logger.Error("Erreur suppression tâche", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur suppression tâche", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Tâche supprimée avec succès", nil)
}

// CompleteTask marque une tâche comme complétée
// @Summary Marquer une tâche comme complétée
// @Description Marque une tâche comme complétée
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la tâche"
// @Success 200 {object} response.Response "Tâche marquée comme complétée"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 403 {object} response.ErrorResponse "Accès non autorisé"
// @Failure 404 {object} response.ErrorResponse "Tâche non trouvée"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /tasks/{id}/complete [post]
func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID de la tâche depuis l'URL
	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de tâche invalide", err)
		return
	}

	task, err := h.taskService.CompleteTask(r.Context(), userID, taskID)
	if err != nil {
		switch err.Error() {
		case "tâche non trouvée":
			response.Error(w, http.StatusNotFound, "Tâche non trouvée", err)
		case "accès non autorisé":
			response.Error(w, http.StatusForbidden, "Accès non autorisé", err)
		default:
			h.logger.Error("Erreur complétion tâche", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur complétion tâche", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Tâche marquée comme complétée", task)
}

// GetTaskStats récupère les statistiques des tâches
// @Summary Récupérer les statistiques des tâches
// @Description Récupère les statistiques des tâches de l'utilisateur
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Statistiques récupérées"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /tasks/stats [get]
func (h *TaskHandler) GetTaskStats(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	stats, err := h.taskService.GetTaskStats(r.Context(), userID)
	if err != nil {
		h.logger.Error("Erreur récupération statistiques", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération statistiques", err)
		return
	}

	response.Success(w, http.StatusOK, "Statistiques récupérées avec succès", stats)
}
