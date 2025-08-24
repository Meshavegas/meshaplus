package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"backend/pkg/logger"
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type CategoryRepository struct {
	db     *pg.DB
	logger logger.Logger
}

func NewCategoryRepository(db *pg.DB, logger logger.Logger) repository.CategoryRepository {
	return &CategoryRepository{
		db:     db,
		logger: logger,
	}
}

// Create crée une nouvelle catégorie
func (r *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	_, err := r.db.Model(category).Context(ctx).Insert()
	if err != nil {
		r.logger.Error("Erreur création catégorie", logger.Error(err))
		return err
	}

	return nil
}

// GetByID récupère une catégorie par son ID
func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error) {
	var category entity.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, color, created_at
		FROM categories 
		WHERE id = ?
	`

	_, err := r.db.QueryOneContext(ctx, &category, query, id)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Erreur récupération catégorie par ID", logger.Error(err))
		return nil, err
	}

	return &category, nil
}

// GetByUserID récupère toutes les catégories d'un utilisateur
func (r *CategoryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error) {
	var categories []*entity.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, color, created_at
		FROM categories 
		WHERE user_id = ?
		ORDER BY name ASC
	`

	_, err := r.db.QueryContext(ctx, &categories, query, userID)
	if err != nil {
		r.logger.Error("Erreur récupération catégories par utilisateur", logger.Error(err))
		return nil, err
	}

	return categories, nil
}

// Update met à jour une catégorie
func (r *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	_, err := r.db.Model(category).Context(ctx).WherePK().Update()
	if err != nil {
		r.logger.Error("Erreur mise à jour catégorie", logger.Error(err))
		return err
	}

	return nil
}

// Delete supprime une catégorie (soft delete)
func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE categories SET deleted_at = NOW() WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Erreur suppression catégorie", logger.Error(err))
		return err
	}

	return nil
}

// GetByType récupère toutes les catégories d'un utilisateur par type
func (r *CategoryRepository) GetByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]*entity.Category, error) {
	var categories []*entity.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, color, created_at
		FROM categories 
		WHERE user_id = ? AND type = ?
		ORDER BY name ASC
	`

	_, err := r.db.QueryContext(ctx, &categories, query, userID, categoryType)
	if err != nil {
		r.logger.Error("Erreur récupération catégories par type", logger.Error(err))
		return nil, err
	}

	return categories, nil
}

// GetChildren récupère les catégories enfants d'une catégorie parent
func (r *CategoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.Category, error) {
	var categories []*entity.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, color, created_at
		FROM categories 
		WHERE parent_id = ?
		ORDER BY name ASC
	`

	_, err := r.db.QueryContext(ctx, &categories, query, parentID)
	if err != nil {
		r.logger.Error("Erreur récupération catégories enfants", logger.Error(err))
		return nil, err
	}

	return categories, nil
}

// Méthodes utilitaires spécifiques à notre implémentation

// GetCategoryNamesByType récupère uniquement les noms des catégories d'un utilisateur par type
func (r *CategoryRepository) GetCategoryNamesByType(userID uuid.UUID, categoryType string) ([]string, error) {
	var categoryNames []string

	query := `
		SELECT name
		FROM categories 
		WHERE user_id = ? AND type = ?
		ORDER BY name ASC
	`

	_, err := r.db.Query(&categoryNames, query, userID, categoryType)
	if err != nil {
		r.logger.Error("Erreur récupération noms de catégories par type", logger.Error(err))
		return nil, err
	}

	return categoryNames, nil
}

// GetAll récupère toutes les catégories d'un utilisateur
func (r *CategoryRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error) {
	var categories []*entity.Category

	query := `
		SELECT id, user_id, name, type, parent_id, icon, color, created_at
		FROM categories 
		WHERE user_id = ?
		ORDER BY name ASC
	`

	_, err := r.db.QueryContext(ctx, &categories, query, userID)
	if err != nil {
		r.logger.Error("Erreur récupération catégories", logger.Error(err))
		return nil, err
	}

	return categories, nil
}
