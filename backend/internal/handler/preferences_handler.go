package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/domaine/entity"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/response"

	"github.com/google/uuid"
)

// PreferencesHandler gère les requêtes HTTP pour les préférences utilisateur
type PreferencesHandler struct {
	preferencesService *service.PreferencesService
	logger             logger.Logger
}

// NewPreferencesHandler crée une nouvelle instance de PreferencesHandler
func NewPreferencesHandler(preferencesService *service.PreferencesService, logger logger.Logger) *PreferencesHandler {
	return &PreferencesHandler{
		preferencesService: preferencesService,
		logger:             logger,
	}
}

// CreatePreferences crée les préférences utilisateur
// @Summary Créer les préférences utilisateur
// @Description Crée les préférences utilisateur basées sur le wizard
// @Tags preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.CreatePreferencesRequest true "Données des préférences"
// @Success 201 {object} response.Response "Préférences créées"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /preferences [post]
func (h *PreferencesHandler) CreatePreferences(w http.ResponseWriter, r *http.Request) {
	var req entity.CreatePreferencesRequest
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

	// Appel au service de création de préférences
	preferences, err := h.preferencesService.CreatePreferences(r.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Erreur création préférences", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la création des préférences", err)
		return
	}

	response.Success(w, http.StatusCreated, "Préférences créées avec succès", preferences)
}

// GetPreferences récupère les préférences utilisateur
// @Summary Récupérer les préférences utilisateur
// @Description Récupère les préférences utilisateur
// @Tags preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Préférences récupérées"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Préférences non trouvées"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /preferences [get]
func (h *PreferencesHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Appel au service de récupération de préférences
	preferences, err := h.preferencesService.GetPreferences(r.Context(), userID)
	if err != nil {
		h.logger.Error("Erreur récupération préférences", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la récupération des préférences", err)
		return
	}

	if preferences == nil {
		response.Error(w, http.StatusNotFound, "Préférences non trouvées", nil)
		return
	}

	response.Success(w, http.StatusOK, "Préférences récupérées avec succès", preferences)
}

// UpdatePreferences met à jour les préférences utilisateur
// @Summary Mettre à jour les préférences utilisateur
// @Description Met à jour les préférences utilisateur
// @Tags preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body entity.UpdatePreferencesRequest true "Données de mise à jour"
// @Success 200 {object} response.Response "Préférences mises à jour"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Préférences non trouvées"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /preferences [put]
func (h *PreferencesHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	var req entity.UpdatePreferencesRequest
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

	// Appel au service de mise à jour de préférences
	preferences, err := h.preferencesService.UpdatePreferences(r.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Erreur mise à jour préférences", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la mise à jour des préférences", err)
		return
	}

	response.Success(w, http.StatusOK, "Préférences mises à jour avec succès", preferences)
}

// DeletePreferences supprime les préférences utilisateur
// @Summary Supprimer les préférences utilisateur
// @Description Supprime les préférences utilisateur
// @Tags preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Préférences supprimées"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /preferences [delete]
func (h *PreferencesHandler) DeletePreferences(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Appel au service de suppression de préférences
	if err := h.preferencesService.DeletePreferences(r.Context(), userID); err != nil {
		h.logger.Error("Erreur suppression préférences", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la suppression des préférences", err)
		return
	}

	response.Success(w, http.StatusOK, "Préférences supprimées avec succès", nil)
}

// GenerateContent génère du contenu personnalisé
// @Summary Générer du contenu personnalisé
// @Description Génère du contenu personnalisé basé sur les préférences
// @Tags preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param contentType query string true "Type de contenu (budget_tip, savings_advice, habit_reminder)"
// @Success 200 {object} response.Response "Contenu généré"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Préférences non trouvées"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /preferences/content [get]
func (h *PreferencesHandler) GenerateContent(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer le type de contenu depuis les paramètres de requête
	contentType := r.URL.Query().Get("contentType")
	if contentType == "" {
		response.Error(w, http.StatusBadRequest, "Type de contenu requis", nil)
		return
	}

	// Appel au service de génération de contenu
	content, err := h.preferencesService.GeneratePersonalizedContent(r.Context(), userID, contentType)
	if err != nil {
		h.logger.Error("Erreur génération contenu", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la génération du contenu", err)
		return
	}

	response.Success(w, http.StatusOK, "Contenu généré avec succès", map[string]string{
		"content": content,
		"type":    contentType,
	})
}
