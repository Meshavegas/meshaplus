package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/internal/domaine/entity"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/response"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TransactionHandler gère les requêtes HTTP pour les transactions
type TransactionHandler struct {
	transactionService *service.TransactionService
	logger             logger.Logger
}

// NewTransactionHandler crée une nouvelle instance de TransactionHandler
func NewTransactionHandler(transactionService *service.TransactionService, logger logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		logger:             logger,
	}
}

// CreateTransaction crée une nouvelle transaction
// @Summary Créer une nouvelle transaction
// @Description Crée une nouvelle transaction financière pour l'utilisateur authentifié
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction body entity.CreateTransactionRequest true "Données de la transaction"
// @Success 201 {object} response.Response "Transaction créée"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	var req entity.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	transaction, err := h.transactionService.CreateTransaction(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("Erreur création transaction", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur création transaction", err)
		return
	}

	response.Success(w, http.StatusCreated, "Transaction créée avec succès", transaction)
}

// GetTransaction récupère une transaction par son ID
// @Summary Récupérer une transaction
// @Description Récupère une transaction spécifique par son ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la transaction"
// @Success 200 {object} response.Response "Transaction récupérée"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Transaction non trouvée"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID de la transaction depuis l'URL
	transactionIDStr := chi.URLParam(r, "id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de transaction invalide", err)
		return
	}

	transaction, err := h.transactionService.GetTransaction(r.Context(), userID, transactionID)
	if err != nil {
		h.logger.Error("Erreur récupération transaction", logger.Error(err))
		response.Error(w, http.StatusNotFound, "Transaction non trouvée", err)
		return
	}

	response.Success(w, http.StatusOK, "Transaction récupérée avec succès", transaction)
}

// GetTransactions récupère toutes les transactions de l'utilisateur
// @Summary Récupérer toutes les transactions
// @Description Récupère toutes les transactions de l'utilisateur authentifié
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Type de transaction (income/expense)"
// @Param category_id query string false "ID de la catégorie"
// @Param account_id query string false "ID du compte"
// @Param start_date query string false "Date de début (YYYY-MM-DD)"
// @Param end_date query string false "Date de fin (YYYY-MM-DD)"
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(10)
// @Success 200 {object} response.Response "Transactions récupérées"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
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

	// Paramètres de filtrage
	txType := query.Get("type")
	categoryIDStr := query.Get("category_id")
	accountIDStr := query.Get("account_id")
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	// Parser les UUIDs optionnels
	var categoryID, accountID *uuid.UUID
	if categoryIDStr != "" {
		if parsed, err := uuid.Parse(categoryIDStr); err == nil {
			categoryID = &parsed
		}
	}
	if accountIDStr != "" {
		if parsed, err := uuid.Parse(accountIDStr); err == nil {
			accountID = &parsed
		}
	}

	transactions, total, err := h.transactionService.GetTransactions(r.Context(), userID, service.TransactionQuery{
		Type:       txType,
		CategoryID: categoryID,
		AccountID:  accountID,
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		h.logger.Error("Erreur récupération transactions", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération transactions", err)
		return
	}

	response.Success(w, http.StatusOK, "Transactions récupérées avec succès", map[string]interface{}{
		"transactions": transactions,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateTransaction met à jour une transaction
// @Summary Mettre à jour une transaction
// @Description Met à jour une transaction existante
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la transaction"
// @Param transaction body entity.UpdateTransactionRequest true "Données de mise à jour"
// @Success 200 {object} response.Response "Transaction mise à jour"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Transaction non trouvée"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID de la transaction depuis l'URL
	transactionIDStr := chi.URLParam(r, "id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de transaction invalide", err)
		return
	}

	var req entity.UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	transaction, err := h.transactionService.UpdateTransaction(r.Context(), userID, transactionID, req)
	if err != nil {
		h.logger.Error("Erreur mise à jour transaction", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur mise à jour transaction", err)
		return
	}

	response.Success(w, http.StatusOK, "Transaction mise à jour avec succès", transaction)
}

// DeleteTransaction supprime une transaction
// @Summary Supprimer une transaction
// @Description Supprime une transaction existante
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la transaction"
// @Success 200 {object} response.Response "Transaction supprimée"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Transaction non trouvée"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID de la transaction depuis l'URL
	transactionIDStr := chi.URLParam(r, "id")
	transactionID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de transaction invalide", err)
		return
	}

	err = h.transactionService.DeleteTransaction(r.Context(), userID, transactionID)
	if err != nil {
		h.logger.Error("Erreur suppression transaction", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur suppression transaction", err)
		return
	}

	response.Success(w, http.StatusOK, "Transaction supprimée avec succès", nil)
}

// GetTransactionStats récupère les statistiques des transactions
// @Summary Récupérer les statistiques des transactions
// @Description Récupère les statistiques des transactions de l'utilisateur
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param period query string false "Période (week/month/year)" default(month)
// @Success 200 {object} response.Response "Statistiques récupérées"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /transactions/stats [get]
func (h *TransactionHandler) GetTransactionStats(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer la période depuis les paramètres de requête
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "month"
	}

	stats, err := h.transactionService.GetTransactionStats(r.Context(), userID, period)
	if err != nil {
		h.logger.Error("Erreur récupération statistiques", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération statistiques", err)
		return
	}

	response.Success(w, http.StatusOK, "Statistiques récupérées avec succès", stats)
}
