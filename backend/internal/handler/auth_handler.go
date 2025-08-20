package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/response"
)

// AuthHandler gère les requêtes HTTP pour l'authentification
type AuthHandler struct {
	authService *service.AuthService
	logger      logger.Logger
}

// NewAuthHandler crée une nouvelle instance de AuthHandler
func NewAuthHandler(authService *service.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register inscrit un nouvel utilisateur
// @Summary Inscription d'un utilisateur
// @Description Inscrit un nouvel utilisateur dans le système
// @Tags auth
// @Accept json
// @Produce json
// @Param user body service.RegisterRequest true "Données d'inscription"
// @Success 201 {object} response.Response "Utilisateur inscrit"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 409 {object} response.ErrorResponse "Email déjà utilisé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	result, err := h.authService.Register(r.Context(), req)
	if err != nil {
		switch err.Error() {
		case "email déjà utilisé":
			response.Error(w, http.StatusConflict, "Email déjà utilisé", err)
		default:
			h.logger.Error("Erreur inscription", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur interne", err)
		}
		return
	}

	response.Success(w, http.StatusCreated, "Utilisateur inscrit avec succès", result)
}

// Login authentifie un utilisateur
// @Summary Connexion d'un utilisateur
// @Description Authentifie un utilisateur et retourne les tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body service.LoginRequest true "Données de connexion"
// @Success 200 {object} response.Response "Connexion réussie"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Email ou mot de passe incorrect"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	result, err := h.authService.Login(r.Context(), req)
	if err != nil {
		switch err.Error() {
		case "email ou mot de passe incorrect":
			response.Error(w, http.StatusUnauthorized, "Email ou mot de passe incorrect", err)
		default:
			h.logger.Error("Erreur connexion", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur interne", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Connexion réussie", result)
}

// RefreshToken rafraîchit un token d'accès
// @Summary Rafraîchir un token
// @Description Rafraîchit un token d'accès avec un token de rafraîchissement
// @Tags auth
// @Accept json
// @Produce json
// @Param token body service.RefreshTokenRequest true "Token de rafraîchissement"
// @Success 200 {object} response.Response "Token rafraîchi"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 401 {object} response.ErrorResponse "Token invalide"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req service.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	result, err := h.authService.RefreshToken(r.Context(), req)
	if err != nil {
		switch err.Error() {
		case "token de rafraîchissement invalide", "token n'est pas un token de rafraîchissement":
			response.Error(w, http.StatusUnauthorized, "Token invalide", err)
		default:
			h.logger.Error("Erreur rafraîchissement token", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur interne", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Token rafraîchi avec succès", result)
}

// GetCurrentUser récupère l'utilisateur actuel
// @Summary Récupérer l'utilisateur actuel
// @Description Récupère les informations de l'utilisateur authentifié
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "Utilisateur actuel"
// @Failure 401 {object} response.ErrorResponse "Non authentifié"
// @Failure 404 {object} response.ErrorResponse "Utilisateur non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.GetCurrentUser(r.Context())
	if err != nil {
		switch err.Error() {
		case "utilisateur non authentifié":
			response.Error(w, http.StatusUnauthorized, "Non authentifié", err)
		case "utilisateur non trouvé":
			response.Error(w, http.StatusNotFound, "Utilisateur non trouvé", err)
		default:
			h.logger.Error("Erreur récupération utilisateur actuel", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur interne", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Utilisateur récupéré avec succès", user)
}
