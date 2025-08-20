package postgres

import (
	"context"
	"fmt"
	"math"
	"strings"

	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// userRepository implémente repository.UserRepository avec go-pg
type userRepository struct {
	db *pg.DB
}

// NewUserRepository crée une nouvelle instance de userRepository
func NewUserRepository(db *pg.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create crée un nouvel utilisateur
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.db.ModelContext(ctx, user).Insert()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return entity.ErrUserAlreadyExists
		}
		return fmt.Errorf("erreur création utilisateur: %w", err)
	}
	return nil
}

// GetByID récupère un utilisateur par son ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.ModelContext(ctx, user).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("erreur récupération utilisateur: %w", err)
	}
	return user, nil
}

// GetByEmail récupère un utilisateur par son email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.ModelContext(ctx, user).Where("email = ?", email).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("erreur récupération utilisateur par email: %w", err)
	}
	return user, nil
}

// List récupère une liste paginée d'utilisateurs
func (r *userRepository) List(ctx context.Context, query *entity.UserQuery) (*entity.UserListResponse, error) {
	var users []entity.User
	var total int

	// Construire la requête de base
	dbQuery := r.db.ModelContext(ctx, &users)

	// Appliquer les filtres de recherche
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		dbQuery = dbQuery.Where(
			"LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
			searchTerm, searchTerm,
		)
	}

	// Compter le total
	total, err := dbQuery.Count()
	if err != nil {
		return nil, fmt.Errorf("erreur comptage utilisateurs: %w", err)
	}

	// Récupérer les utilisateurs avec pagination
	offset := (query.Page - 1) * query.PageSize
	err = dbQuery.Offset(offset).Limit(query.PageSize).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération utilisateurs: %w", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &entity.UserListResponse{
		Users:      users,
		Total:      int64(total),
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: totalPages,
	}, nil
}

// Update met à jour un utilisateur
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	result, err := r.db.ModelContext(ctx, user).WherePK().Update()
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return entity.ErrUserAlreadyExists
		}
		return fmt.Errorf("erreur mise à jour utilisateur: %w", err)
	}

	if result.RowsAffected() == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

// Delete supprime un utilisateur (soft delete)
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ModelContext(ctx, &entity.User{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression utilisateur: %w", err)
	}

	if result.RowsAffected() == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

// Exists vérifie si un utilisateur existe
func (r *userRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	count, err := r.db.ModelContext(ctx, &entity.User{}).Where("id = ?", id).Count()
	if err != nil {
		return false, fmt.Errorf("erreur vérification existence utilisateur: %w", err)
	}

	return count > 0, nil
}

// ExistsByEmail vérifie si un email existe déjà
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.db.ModelContext(ctx, &entity.User{}).Where("email = ?", email).Count()
	if err != nil {
		return false, fmt.Errorf("erreur vérification email: %w", err)
	}

	return count > 0, nil
}

// Search recherche des utilisateurs par nom ou email
func (r *userRepository) Search(ctx context.Context, searchTerm string, limit int) ([]*entity.User, error) {
	var users []*entity.User

	searchPattern := "%" + strings.ToLower(searchTerm) + "%"

	err := r.db.ModelContext(ctx, &users).
		Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
			searchPattern, searchPattern).
		Limit(limit).
		Order("created_at DESC").
		Select()

	if err != nil {
		return nil, fmt.Errorf("erreur recherche utilisateurs: %w", err)
	}

	return users, nil
}
