package usecase

import (
	"context"
	"fmt"

	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"backend/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserUsecase gère la logique métier des utilisateurs
type UserUsecase struct {
	userRepo repository.UserRepository
	logger   logger.Logger
}

// NewUserUsecase crée une nouvelle instance de UserUsecase
func NewUserUsecase(userRepo repository.UserRepository, logger logger.Logger) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// hashPassword hashe un mot de passe avec bcrypt
func (uc *UserUsecase) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("erreur hashage mot de passe: %w", err)
	}
	return string(bytes), nil
}

// CreateUser crée un nouvel utilisateur
func (uc *UserUsecase) CreateUser(ctx context.Context, req entity.CreateUserRequest) (*entity.User, error) {
	// Validation métier
	if req.Name == "" || req.Email == "" || req.Password == "" {
		uc.logger.Warn("Données utilisateur invalides",
			logger.String("email", req.Email),
			logger.String("name", req.Name),
		)
		return nil, entity.ErrInvalidUserData
	}

	// Vérifier si l'email existe déjà
	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		uc.logger.Error("Erreur vérification email",
			logger.Error(err),
			logger.String("email", req.Email),
		)
		return nil, fmt.Errorf("erreur vérification email: %w", err)
	}

	if exists {
		uc.logger.Warn("Tentative création utilisateur avec email existant",
			logger.String("email", req.Email),
		)
		return nil, entity.ErrUserAlreadyExists
	}

	// Hasher le mot de passe
	passwordHash, err := uc.hashPassword(req.Password)
	if err != nil {
		uc.logger.Error("Erreur hashage mot de passe",
			logger.Error(err),
			logger.String("email", req.Email),
		)
		return nil, fmt.Errorf("erreur hashage mot de passe: %w", err)
	}

	// Créer l'entité utilisateur
	user := &entity.User{
		Name:         req.Name,
		Email:        req.Email,
		Avatar:       req.Avatar,
		PasswordHash: passwordHash,
	}

	// Validation avant création
	if err := user.IsValidForCreation(); err != nil {
		return nil, err
	}

	// Persister en base
	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("Erreur création utilisateur",
			logger.Error(err),
			logger.String("email", req.Email),
		)
		return nil, fmt.Errorf("erreur création utilisateur: %w", err)
	}

	uc.logger.Info("Utilisateur créé avec succès",
		logger.String("id", user.ID.String()),
		logger.String("email", user.Email),
	)

	return user, nil
}

// GetUserByID récupère un utilisateur par son ID
func (uc *UserUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if id == uuid.Nil {
		return nil, entity.ErrInvalidUserData
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Erreur récupération utilisateur",
			logger.Error(err),
			logger.String("id", id.String()),
		)
		return nil, err
	}

	return user, nil
}

// GetUsers récupère une liste paginée d'utilisateurs
func (uc *UserUsecase) GetUsers(ctx context.Context, query entity.UserQuery) (*entity.UserListResponse, error) {
	// Validation des paramètres
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 10
	}

	response, err := uc.userRepo.List(ctx, &query)
	if err != nil {
		uc.logger.Error("Erreur récupération liste utilisateurs",
			logger.Error(err),
			logger.Int("page", query.Page),
			logger.Int("pageSize", query.PageSize),
		)
		return nil, fmt.Errorf("erreur récupération utilisateurs: %w", err)
	}

	uc.logger.Debug("Liste utilisateurs récupérée",
		logger.Int("count", len(response.Users)),
		logger.Int64("total", response.Total),
	)

	return response, nil
}

// UpdateUser met à jour un utilisateur
func (uc *UserUsecase) UpdateUser(ctx context.Context, id uuid.UUID, req entity.UpdateUserRequest) (*entity.User, error) {
	if id == uuid.Nil {
		return nil, entity.ErrInvalidUserData
	}

	// Récupérer l'utilisateur existant
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Mettre à jour les champs modifiés
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		// Vérifier que le nouvel email n'est pas déjà utilisé
		if *req.Email != user.Email {
			exists, err := uc.userRepo.ExistsByEmail(ctx, *req.Email)
			if err != nil {
				return nil, fmt.Errorf("erreur vérification email: %w", err)
			}
			if exists {
				return nil, entity.ErrUserAlreadyExists
			}
		}
		user.Email = *req.Email
	}
	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Password != nil {
		// Hasher le nouveau mot de passe
		passwordHash, err := uc.hashPassword(*req.Password)
		if err != nil {
			uc.logger.Error("Erreur hashage nouveau mot de passe",
				logger.Error(err),
				logger.String("id", id.String()),
			)
			return nil, fmt.Errorf("erreur hashage mot de passe: %w", err)
		}
		user.PasswordHash = passwordHash
	}

	// Persister les modifications
	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Error("Erreur mise à jour utilisateur",
			logger.Error(err),
			logger.String("id", id.String()),
		)
		return nil, fmt.Errorf("erreur mise à jour utilisateur: %w", err)
	}

	uc.logger.Info("Utilisateur mis à jour avec succès",
		logger.String("id", user.ID.String()),
		logger.String("email", user.Email),
	)

	return user, nil
}

// DeleteUser supprime un utilisateur (soft delete)
func (uc *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return entity.ErrInvalidUserData
	}

	// Vérifier que l'utilisateur existe
	exists, err := uc.userRepo.Exists(ctx, id)
	if err != nil {
		return fmt.Errorf("erreur vérification existence utilisateur: %w", err)
	}
	if !exists {
		return entity.ErrUserNotFound
	}

	// Supprimer l'utilisateur
	if err := uc.userRepo.Delete(ctx, id); err != nil {
		uc.logger.Error("Erreur suppression utilisateur",
			logger.Error(err),
			logger.String("id", id.String()),
		)
		return fmt.Errorf("erreur suppression utilisateur: %w", err)
	}

	uc.logger.Info("Utilisateur supprimé avec succès",
		logger.String("id", id.String()),
	)

	return nil
}

// SearchUsers recherche des utilisateurs
func (uc *UserUsecase) SearchUsers(ctx context.Context, searchTerm string, limit int) ([]*entity.User, error) {
	if searchTerm == "" {
		return []*entity.User{}, nil
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	users, err := uc.userRepo.Search(ctx, searchTerm, limit)
	if err != nil {
		uc.logger.Error("Erreur recherche utilisateurs",
			logger.Error(err),
			logger.String("search", searchTerm),
		)
		return nil, fmt.Errorf("erreur recherche utilisateurs: %w", err)
	}

	return users, nil
}
