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

// AccountHandler gère les requêtes HTTP pour les comptes bancaires
type AccountHandler struct {
	accountService     AccountService
	transactionService TransactionService
	logger             logger.Logger
}

type AccountDetails struct {
	Account      *entity.Account       `json:"account"`
	Transactions []*entity.Transaction `json:"transactions"`
}

// AccountService interface pour les services de comptes
type AccountService interface {
	CreateAccount(ctx context.Context, userID uuid.UUID, req entity.CreateAccountRequest) (*entity.Account, error)
	GetAccount(ctx context.Context, userID, accountID uuid.UUID) (*entity.Account, error)
	GetAccounts(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Account, int64, error)
	UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, req entity.UpdateAccountRequest) (*entity.Account, error)
	DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error
	GetAccountBalance(ctx context.Context, userID, accountID uuid.UUID) (float64, error)
}

type TransactionService interface {
	GetTransactionsByAccountID(ctx context.Context, userID, accountID uuid.UUID) ([]*entity.TransactionWithDetails, error)
}

// NewAccountHandler crée une nouvelle instance de AccountHandler
func NewAccountHandler(accountService AccountService, transactionService TransactionService, logger logger.Logger) *AccountHandler {
	return &AccountHandler{
		accountService:     accountService,
		transactionService: transactionService,
		logger:             logger,
	}
}

// CreateAccount crée un nouveau compte bancaire
// @Summary Créer un nouveau compte
// @Description Crée un nouveau compte bancaire pour l'utilisateur authentifié
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param account body entity.CreateAccountRequest true "Données du compte"
// @Success 201 {object} response.Response "Compte créé"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	var req entity.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	account, err := h.accountService.CreateAccount(r.Context(), userID, req)
	if err != nil {
		h.logger.Error("Erreur création compte", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur création compte", err)
		return
	}

	response.Success(w, http.StatusCreated, "Compte créé avec succès", account)
}

// GetAccount récupère un compte par son ID
// @Summary Récupérer un compte
// @Description Récupère un compte spécifique par son ID
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID du compte"
// @Success 200 {object} response.Response "Compte récupéré"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Compte non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID du compte depuis l'URL
	accountIDStr := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de compte invalide", err)
		return
	}

	account, err := h.accountService.GetAccount(r.Context(), userID, accountID)
	if err != nil {
		h.logger.Error("Erreur récupération compte", logger.Error(err))
		response.Error(w, http.StatusNotFound, "Compte non trouvé", err)
		return
	}

	response.Success(w, http.StatusOK, "Compte récupéré avec succès", account)
}

// GetAccounts récupère tous les comptes de l'utilisateur
// @Summary Récupérer tous les comptes
// @Description Récupère tous les comptes de l'utilisateur authentifié
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Type de compte (checking/savings/mobile_money)"
// @Param page query int false "Numéro de page" default(1)
// @Param limit query int false "Nombre d'éléments par page" default(10)
// @Success 200 {object} response.Response "Comptes récupérés"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /accounts [get]
func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
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

	accounts, total, err := h.accountService.GetAccounts(r.Context(), userID, page, limit)
	if err != nil {
		h.logger.Error("Erreur récupération comptes", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération comptes", err)
		return
	}

	response.Success(w, http.StatusOK, "Comptes récupérés avec succès", map[string]interface{}{
		"accounts": accounts,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateAccount met à jour un compte
// @Summary Mettre à jour un compte
// @Description Met à jour un compte existant
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID du compte"
// @Param account body entity.UpdateAccountRequest true "Données de mise à jour"
// @Success 200 {object} response.Response "Compte mis à jour"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Compte non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID du compte depuis l'URL
	accountIDStr := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de compte invalide", err)
		return
	}

	var req entity.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	account, err := h.accountService.UpdateAccount(r.Context(), userID, accountID, req)
	if err != nil {
		h.logger.Error("Erreur mise à jour compte", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur mise à jour compte", err)
		return
	}

	response.Success(w, http.StatusOK, "Compte mis à jour avec succès", account)
}

// DeleteAccount supprime un compte
// @Summary Supprimer un compte
// @Description Supprime un compte existant
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID du compte"
// @Success 200 {object} response.Response "Compte supprimé"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Compte non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID du compte depuis l'URL
	accountIDStr := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de compte invalide", err)
		return
	}

	err = h.accountService.DeleteAccount(r.Context(), userID, accountID)
	if err != nil {
		h.logger.Error("Erreur suppression compte", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur suppression compte", err)
		return
	}

	response.Success(w, http.StatusOK, "Compte supprimé avec succès", nil)
}

// GetAccountBalance récupère le solde d'un compte
// @Summary Récupérer le solde d'un compte
// @Description Récupère le solde actuel d'un compte spécifique
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID du compte"
// @Success 200 {object} response.Response "Solde récupéré"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Compte non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /accounts/{id}/balance [get]
func (h *AccountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID du compte depuis l'URL
	accountIDStr := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de compte invalide", err)
		return
	}

	balance, err := h.accountService.GetAccountBalance(r.Context(), userID, accountID)
	if err != nil {
		h.logger.Error("Erreur récupération solde", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération solde", err)
		return
	}

	response.Success(w, http.StatusOK, "Solde récupéré avec succès", map[string]interface{}{
		"balance": balance,
	})
}

func (h *AccountHandler) GetAccountDetails(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID utilisateur du contexte
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	// Récupérer l'ID du compte depuis l'URL
	accountIDStr := chi.URLParam(r, "id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID de compte invalide", err)
		return
	}

	accountDetails, err := h.accountService.GetAccount(r.Context(), userID, accountID)
	if err != nil {
		h.logger.Error("Erreur récupération détails du compte", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération détails du compte", err)
		return
	}

	transaction, err := h.transactionService.GetTransactionsByAccountID(r.Context(), userID, accountID)
	if err != nil {
		h.logger.Error("Erreur récupération transactions", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération transactions", err)
		return
	}

	response.Success(w, http.StatusOK, "Détails du compte récupérés avec succès", map[string]interface{}{
		"account":      accountDetails,
		"transactions": transaction,
	})
}
