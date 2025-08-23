package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/response"

	"github.com/google/uuid"
)

// CategoryHandler gère les requêtes HTTP pour la catégorisation
type CategoryHandler struct {
	categoryService *service.CategoryService
	logger          logger.Logger
}

// NewCategoryHandler crée une nouvelle instance de CategoryHandler
func NewCategoryHandler(categoryService *service.CategoryService, logger logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		logger:          logger,
	}
}

// CategorizeItem catégorise un item en utilisant l'IA
// @Summary Catégoriser un item
// @Description Utilise l'IA pour catégoriser un item en choisissant une catégorie existante ou en créant une nouvelle
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CategorizeRequest true "Données pour la catégorisation"
// @Success 200 {object} response.Response "Item catégorisé"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /categories/categorize [post]
func (h *CategoryHandler) CategorizeItem(w http.ResponseWriter, r *http.Request) {
	var req CategorizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Validation des données
	if req.Item == "" {
		response.Error(w, http.StatusBadRequest, "L'item à catégoriser est requis", nil)
		return
	}

	// Appel au service de catégorisation
	categoryResponse, err := h.categoryService.CategorizeItem(r.Context(), userID, req.Item, req.CategoryType)
	if err != nil {
		h.logger.Error("Erreur catégorisation", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la catégorisation", err)
		return
	}

	response.Success(w, http.StatusOK, "Item catégorisé avec succès", categoryResponse)
}

// GetCategories récupère toutes les catégories by type
// @Summary Récupère toutes les catégories
// @Description Récupère toutes les catégories
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Catégories récupérées"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /categories [get]
func (h *CategoryHandler) GetCategoriesByType(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	categoryType := r.URL.Query().Get("categoryType")
	if categoryType == "" {
		caterorieResponse, err := h.categoryService.GetAllCategories(r.Context(), userID)
		if err != nil {
			h.logger.Error("Erreur récupération catégories", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur lors de la récupération des catégories", err)
			return
		}
		response.Success(w, http.StatusOK, "Catégories récupérées avec succès", caterorieResponse)
		return
	}

	categories, err := h.categoryService.GetCategoriesByType(r.Context(), userID, categoryType)
	if err != nil {
		h.logger.Error("Erreur récupération catégories", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la récupération des catégories", err)
		return
	}

	response.Success(w, http.StatusOK, "Catégories récupérées avec succès", categories)
}

// CategorizeRequest représente la requête pour catégoriser un item
type CategorizeRequest struct {
	Item         string `json:"item" validate:"required"`
	CategoryType string `json:"categoryType" validate:"required,oneof=expense revenue task"`
}
