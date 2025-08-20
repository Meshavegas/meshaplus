package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/domaine/entity"
	"backend/pkg/logger"
	"backend/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// SavingGoalHandler gère les requêtes HTTP pour les objectifs d'épargne
type SavingGoalHandler struct {
	savingGoalService SavingGoalService
	logger            logger.Logger
}

// SavingGoalService interface pour les services d'objectifs d'épargne
type SavingGoalService interface {
	CreateSavingGoal(ctx context.Context, userID uuid.UUID, req entity.CreateSavingGoalRequest) (*entity.SavingGoal, error)
	GetSavingGoal(ctx context.Context, userID, goalID uuid.UUID) (*entity.SavingGoal, error)
	GetSavingGoals(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.SavingGoal, int64, error)
	UpdateSavingGoal(ctx context.Context, userID, goalID uuid.UUID, req entity.UpdateSavingGoalRequest) (*entity.SavingGoal, error)
	DeleteSavingGoal(ctx context.Context, userID, goalID uuid.UUID) error
}

// NewSavingGoalHandler crée une nouvelle instance de SavingGoalHandler
func NewSavingGoalHandler(savingGoalService SavingGoalService, logger logger.Logger) *SavingGoalHandler {
	return &SavingGoalHandler{
		savingGoalService: savingGoalService,
		logger:            logger,
	}
}

// CreateSavingGoal crée un nouvel objectif d'épargne
// @Summary Créer un objectif d'épargne
// @Description Crée un nouvel objectif d'épargne pour l'utilisateur authentifié
// @Tags saving-goals
// @Accept json
// @Produce json
// @Param goal body entity.CreateSavingGoalRequest true "Données de l'objectif d'épargne"
// @Success 201 {object} response.Response "Objectif d'épargne créé"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /saving-goals [post]
func (h *SavingGoalHandler) CreateSavingGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	var req entity.CreateSavingGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	goal, err := h.savingGoalService.CreateSavingGoal(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("Erreur création objectif d'épargne", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur création objectif d'épargne", err)
		return
	}

	response.Success(w, http.StatusCreated, "Objectif d'épargne créé avec succès", goal)
}

// GetSavingGoal récupère un objectif d'épargne par son ID
// @Summary Récupérer un objectif d'épargne
// @Description Récupère un objectif d'épargne par son ID
// @Tags saving-goals
// @Accept json
// @Produce json
// @Param id path string true "ID de l'objectif d'épargne"
// @Success 200 {object} response.Response "Objectif d'épargne récupéré"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Objectif d'épargne non trouvé"
// @Router /saving-goals/{id} [get]
func (h *SavingGoalHandler) GetSavingGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	goalIDStr := chi.URLParam(r, "id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID d'objectif invalide", err)
		return
	}

	goal, err := h.savingGoalService.GetSavingGoal(r.Context(), userID, goalID)
	if err != nil {
		h.logger.Error("Erreur récupération objectif d'épargne", logger.Error(err))
		response.Error(w, http.StatusNotFound, "Objectif d'épargne non trouvé", err)
		return
	}

	response.Success(w, http.StatusOK, "Objectif d'épargne récupéré avec succès", goal)
}

// GetSavingGoals récupère tous les objectifs d'épargne de l'utilisateur
// @Summary Lister les objectifs d'épargne
// @Description Récupère tous les objectifs d'épargne de l'utilisateur avec pagination
// @Tags saving-goals
// @Accept json
// @Produce json
// @Param page query int false "Numéro de page (défaut: 1)"
// @Param limit query int false "Nombre d'éléments par page (défaut: 10)"
// @Success 200 {object} response.Response "Liste des objectifs d'épargne"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Router /saving-goals [get]
func (h *SavingGoalHandler) GetSavingGoals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	goals, total, err := h.savingGoalService.GetSavingGoals(r.Context(), userID, page, limit)
	if err != nil {
		h.logger.Error("Erreur récupération objectifs d'épargne", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération objectifs d'épargne", err)
		return
	}

	response.Success(w, http.StatusOK, "Objectifs d'épargne récupérés avec succès", map[string]interface{}{
		"goals": goals,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateSavingGoal met à jour un objectif d'épargne
// @Summary Mettre à jour un objectif d'épargne
// @Description Met à jour un objectif d'épargne existant
// @Tags saving-goals
// @Accept json
// @Produce json
// @Param id path string true "ID de l'objectif d'épargne"
// @Param goal body entity.UpdateSavingGoalRequest true "Données de mise à jour"
// @Success 200 {object} response.Response "Objectif d'épargne mis à jour"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Objectif d'épargne non trouvé"
// @Router /saving-goals/{id} [put]
func (h *SavingGoalHandler) UpdateSavingGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	goalIDStr := chi.URLParam(r, "id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID d'objectif invalide", err)
		return
	}

	var req entity.UpdateSavingGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	goal, err := h.savingGoalService.UpdateSavingGoal(r.Context(), userID, goalID, req)
	if err != nil {
		h.logger.Error("Erreur mise à jour objectif d'épargne", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur mise à jour objectif d'épargne", err)
		return
	}

	response.Success(w, http.StatusOK, "Objectif d'épargne mis à jour avec succès", goal)
}

// DeleteSavingGoal supprime un objectif d'épargne
// @Summary Supprimer un objectif d'épargne
// @Description Supprime un objectif d'épargne existant
// @Tags saving-goals
// @Accept json
// @Produce json
// @Param id path string true "ID de l'objectif d'épargne"
// @Success 200 {object} response.Response "Objectif d'épargne supprimé"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Objectif d'épargne non trouvé"
// @Router /saving-goals/{id} [delete]
func (h *SavingGoalHandler) DeleteSavingGoal(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	goalIDStr := chi.URLParam(r, "id")
	goalID, err := uuid.Parse(goalIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID d'objectif invalide", err)
		return
	}

	err = h.savingGoalService.DeleteSavingGoal(r.Context(), userID, goalID)
	if err != nil {
		h.logger.Error("Erreur suppression objectif d'épargne", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur suppression objectif d'épargne", err)
		return
	}

	response.Success(w, http.StatusOK, "Objectif d'épargne supprimé avec succès", nil)
}
