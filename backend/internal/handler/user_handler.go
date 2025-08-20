package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"backend/internal/domaine/entity"
	"backend/internal/usecase"
	"backend/pkg/logger"
	"backend/pkg/response"
)

// UserHandler gère les requêtes HTTP pour les utilisateurs
type UserHandler struct {
	userUsecase *usecase.UserUsecase
	logger      logger.Logger
}

// NewUserHandler crée une nouvelle instance de UserHandler
func NewUserHandler(userUsecase *usecase.UserUsecase, logger logger.Logger) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		logger:      logger,
	}
}

// CreateUser crée un nouvel utilisateur
// @Summary Créer un utilisateur
// @Description Crée un nouvel utilisateur dans le système
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.CreateUserRequest true "Données de l'utilisateur"
// @Success 201 {object} response.Response{data=entity.User} "Utilisateur créé"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 409 {object} response.ErrorResponse "Email déjà existant"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req entity.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	user, err := h.userUsecase.CreateUser(r.Context(), req)
	if err != nil {
		switch err {
		case entity.ErrInvalidUserData:
			response.Error(w, http.StatusBadRequest, "Données utilisateur invalides", err)
		case entity.ErrUserAlreadyExists:
			response.Error(w, http.StatusConflict, "Email déjà existant", err)
		default:
			h.logger.Error("Erreur création utilisateur", logger.Error(err))
			response.Error(w, http.StatusInternalServerError, "Erreur interne", err)
		}
		return
	}

	response.Success(w, http.StatusCreated, "Utilisateur créé avec succès", user)
}

// GetUsers récupère la liste des utilisateurs
// @Summary Lister les utilisateurs
// @Description Récupère une liste paginée d'utilisateurs
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Numéro de page" default(1)
// @Param page_size query int false "Taille de page" default(10)
// @Param search query string false "Terme de recherche"

// @Success 200 {object} response.Response{data=entity.UserListResponse} "Liste des utilisateurs"
// @Failure 400 {object} response.ErrorResponse "Paramètres invalides"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := entity.UserQuery{
		Page:     1,
		PageSize: 10,
	}

	// Parser les paramètres de query
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			query.Page = page
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 && pageSize <= 100 {
			query.PageSize = pageSize
		}
	}

	if search := r.URL.Query().Get("search"); search != "" {
		query.Search = search
	}

	users, err := h.userUsecase.GetUsers(r.Context(), query)
	if err != nil {
		h.logger.Error("Erreur récupération utilisateurs", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur récupération utilisateurs", err)
		return
	}

	response.Success(w, http.StatusOK, "Utilisateurs récupérés avec succès", users)
}

// GetUserByID récupère un utilisateur par son ID
// @Summary Récupérer un utilisateur
// @Description Récupère un utilisateur par son identifiant
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID de l'utilisateur" format(uuid)
// @Success 200 {object} response.Response{data=entity.User} "Utilisateur trouvé"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 404 {object} response.ErrorResponse "Utilisateur non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "ID manquant", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID invalide", err)
		return
	}

	user, err := h.userUsecase.GetUserByID(r.Context(), id)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			response.Error(w, http.StatusNotFound, "Utilisateur non trouvé", err)
		default:
			h.logger.Error("Erreur récupération utilisateur", logger.Error(err), logger.String("id", idStr))
			response.Error(w, http.StatusInternalServerError, "Erreur interne", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Utilisateur récupéré avec succès", user)
}

// UpdateUser met à jour un utilisateur
// @Summary Mettre à jour un utilisateur
// @Description Met à jour les informations d'un utilisateur
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID de l'utilisateur" format(uuid)
// @Param user body entity.UpdateUserRequest true "Données à mettre à jour"
// @Success 200 {object} response.Response{data=entity.User} "Utilisateur mis à jour"
// @Failure 400 {object} response.ErrorResponse "Données invalides"
// @Failure 404 {object} response.ErrorResponse "Utilisateur non trouvé"
// @Failure 409 {object} response.ErrorResponse "Email déjà existant"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "ID manquant", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID invalide", err)
		return
	}

	var req entity.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Erreur décodage JSON", logger.Error(err))
		response.Error(w, http.StatusBadRequest, "Données JSON invalides", err)
		return
	}

	user, err := h.userUsecase.UpdateUser(r.Context(), id, req)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			response.Error(w, http.StatusNotFound, "Utilisateur non trouvé", err)
		case entity.ErrInvalidUserData:
			response.Error(w, http.StatusBadRequest, "Données utilisateur invalides", err)
		case entity.ErrUserAlreadyExists:
			response.Error(w, http.StatusConflict, "Email déjà existant", err)
		default:
			h.logger.Error("Erreur mise à jour utilisateur", logger.Error(err), logger.String("id", idStr))
			response.Error(w, http.StatusInternalServerError, "Erreur interne", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Utilisateur mis à jour avec succès", user)
}

// DeleteUser supprime un utilisateur
// @Summary Supprimer un utilisateur
// @Description Supprime un utilisateur (soft delete)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID de l'utilisateur" format(uuid)
// @Success 200 {object} response.Response "Utilisateur supprimé"
// @Failure 400 {object} response.ErrorResponse "ID invalide"
// @Failure 404 {object} response.ErrorResponse "Utilisateur non trouvé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "ID manquant", nil)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "ID invalide", err)
		return
	}

	err = h.userUsecase.DeleteUser(r.Context(), id)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			response.Error(w, http.StatusNotFound, "Utilisateur non trouvé", err)
		case entity.ErrInvalidUserData:
			response.Error(w, http.StatusBadRequest, "ID invalide", err)
		default:
			h.logger.Error("Erreur suppression utilisateur", logger.Error(err), logger.String("id", idStr))
			response.Error(w, http.StatusInternalServerError, "Erreur interne", err)
		}
		return
	}

	response.Success(w, http.StatusOK, "Utilisateur supprimé avec succès", nil)
}
