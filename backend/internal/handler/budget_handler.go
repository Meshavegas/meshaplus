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

// BudgetHandler gère les requêtes HTTP pour les budgets
type BudgetHandler struct {
	budgetService BudgetService
	logger        logger.Logger
}

// BudgetService interface pour les services de budgets
type BudgetService interface {
	CreateBudget(ctx context.Context, userID uuid.UUID, req entity.CreateBudgetRequest) (*entity.Budget, error)
	GetBudget(ctx context.Context, userID, budgetID uuid.UUID) (*entity.Budget, error)
	GetBudgets(ctx context.Context, userID uuid.UUID, month, year int, page, limit int) ([]*entity.Budget, int64, error)
	UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, req entity.UpdateBudgetRequest) (*entity.Budget, error)
	DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error
	GetBudgetStats(ctx context.Context, userID uuid.UUID, month, year int) (map[string]interface{}, error)
}

// NewBudgetHandler crée une nouvelle instance de BudgetHandler
func NewBudgetHandler(budgetService BudgetService, logger logger.Logger) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
		logger:        logger,
	}
}

// CreateBudget crée un nouveau budget
// @Summary Créer un nouveau budget
// @Description Crée un nouveau budget pour l'utilisateur authentifié
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param budget body entity.CreateBudgetRequest true "Données du budget"
// @Success 201 {object} response.Response "Budget créé"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /budgets [post]
func (h *BudgetHandler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	var req entity.CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	budget, err := h.budgetService.CreateBudget(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("Erreur création budget", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur création budget", err)
		return
	}

	response.Success(w, http.StatusCreated, "Budget créé avec succès", budget)
}

// GetBudget récupère un budget par son ID
// @Summary Récupérer un budget
// @Description Récupère un budget spécifique par son ID
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID du budget"
// @Success 200 {object} response.Response "Budget récupéré"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Budget non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /budgets/{id} [get]
func (h *BudgetHandler) GetBudget(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID du budget depuis l'URL
	budgetIDStr := chi.URLParam(r, "id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de budget invalide", err)
		return
	}

	budget, err := h.budgetService.GetBudget(r.Context(), userID, budgetID)
	if err != nil {
		h.logger.Error("Erreur récupération budget", logger.Error(err))
		response.Error(w, http.StatusNotFound, "Budget non trouvé", err)
		return
	}

	response.Success(w, http.StatusOK, "Budget récupéré avec succès", budget)
}

// GetBudgets récupère tous les budgets de l'utilisateur
// @Summary Récupérer tous les budgets
// @Description Récupère tous les budgets de l'utilisateur authentifié
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param month query int false "Mois (1-12)" default(1)
// @Param year query int false "Année" default(2024)
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(10)
// @Success 200 {object} response.Response "Budgets récupérés"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /budgets [get]
func (h *BudgetHandler) GetBudgets(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer les paramètres de requête
	query := r.URL.Query()

	// Paramètres de pagination
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

	// Paramètres de filtrage par période
	month, _ := strconv.Atoi(query.Get("month"))
	if month <= 0 || month > 12 {
		month = 1
	}

	year, _ := strconv.Atoi(query.Get("year"))
	if year <= 0 {
		year = 2024
	}

	budgets, total, err := h.budgetService.GetBudgets(r.Context(), userID, month, year, page, limit)
	if err != nil {
		h.logger.Error("Erreur récupération budgets", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération budgets", err)
		return
	}

	response.Success(w, http.StatusOK, "Budgets récupérés avec succès", map[string]interface{}{
		"budgets": budgets,
		"period": map[string]interface{}{
			"month": month,
			"year":  year,
		},
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateBudget met à jour un budget
// @Summary Mettre à jour un budget
// @Description Met à jour un budget existant
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID du budget"
// @Param budget body entity.UpdateBudgetRequest true "Données de mise à jour"
// @Success 200 {object} response.Response "Budget mis à jour"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Budget non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /budgets/{id} [put]
func (h *BudgetHandler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID du budget depuis l'URL
	budgetIDStr := chi.URLParam(r, "id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de budget invalide", err)
		return
	}

	var req entity.UpdateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	budget, err := h.budgetService.UpdateBudget(r.Context(), userID, budgetID, req)
	if err != nil {
		h.logger.Error("Erreur mise à jour budget", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur mise à jour budget", err)
		return
	}

	response.Success(w, http.StatusOK, "Budget mis à jour avec succès", budget)
}

// DeleteBudget supprime un budget
// @Summary Supprimer un budget
// @Description Supprime un budget existant
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID du budget"
// @Success 200 {object} response.Response "Budget supprimé"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Budget non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /budgets/{id} [delete]
func (h *BudgetHandler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID du budget depuis l'URL
	budgetIDStr := chi.URLParam(r, "id")
	budgetID, err := uuid.Parse(budgetIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de budget invalide", err)
		return
	}

	err = h.budgetService.DeleteBudget(r.Context(), userID, budgetID)
	if err != nil {
		h.logger.Error("Erreur suppression budget", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur suppression budget", err)
		return
	}

	response.Success(w, http.StatusOK, "Budget supprimé avec succès", nil)
}

// GetBudgetStats récupère les statistiques des budgets
// @Summary Récupérer les statistiques des budgets
// @Description Récupère les statistiques des budgets pour une période donnée
// @Tags budgets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param month query int false "Mois (1-12)" default(1)
// @Param year query int false "Année" default(2024)
// @Success 200 {object} response.Response "Statistiques récupérées"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /budgets/stats [get]
func (h *BudgetHandler) GetBudgetStats(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer les paramètres de requête
	query := r.URL.Query()

	month, _ := strconv.Atoi(query.Get("month"))
	if month <= 0 || month > 12 {
		month = 1
	}

	year, _ := strconv.Atoi(query.Get("year"))
	if year <= 0 {
		year = 2024
	}

	stats, err := h.budgetService.GetBudgetStats(r.Context(), userID, month, year)
	if err != nil {
		h.logger.Error("Erreur récupération statistiques", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération statistiques", err)
		return
	}

	response.Success(w, http.StatusOK, "Statistiques récupérées avec succès", stats)
}
