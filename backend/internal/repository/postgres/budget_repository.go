package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// BudgetRepository implémente repository.BudgetRepository
type BudgetRepository struct {
	db *pg.DB
}

// NewBudgetRepository crée une nouvelle instance de BudgetRepository
func NewBudgetRepository(db *pg.DB) repository.BudgetRepository {
	return &BudgetRepository{db: db}
}

// Create crée un nouveau budget
func (r *BudgetRepository) Create(ctx context.Context, budget *entity.Budget) error {
	_, err := r.db.WithContext(ctx).Model(budget).Insert()
	if err != nil {
		return fmt.Errorf("erreur création budget: %w", err)
	}
	return nil
}

// GetByID récupère un budget par son ID
func (r *BudgetRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Budget, error) {
	budget := &entity.Budget{}
	err := r.db.WithContext(ctx).Model(budget).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("budget non trouvé")
		}
		return nil, fmt.Errorf("erreur récupération budget: %w", err)
	}
	return budget, nil
}

// GetByUserID récupère tous les budgets d'un utilisateur
func (r *BudgetRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Budget, error) {
	var budgets []*entity.Budget
	err := r.db.WithContext(ctx).Model(&budgets).Where("user_id = ?", userID).Order("year DESC").Order("month DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération budgets utilisateur: %w", err)
	}
	return budgets, nil
}

// Update met à jour un budget
func (r *BudgetRepository) Update(ctx context.Context, budget *entity.Budget) error {
	_, err := r.db.WithContext(ctx).Model(budget).Where("id = ?", budget.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour budget: %w", err)
	}
	return nil
}

// Delete supprime un budget
func (r *BudgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Budget{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression budget: %w", err)
	}
	return nil
}

// GetByCategoryID récupère les budgets d'une catégorie
func (r *BudgetRepository) GetByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]*entity.Budget, error) {
	var budgets []*entity.Budget
	err := r.db.WithContext(ctx).Model(&budgets).Where("user_id = ? AND category_id = ?", userID, categoryID).Order("year DESC").Order("month DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération budgets par catégorie: %w", err)
	}
	return budgets, nil
}

// GetByMonth récupère les budgets d'un mois spécifique
func (r *BudgetRepository) GetByMonth(ctx context.Context, userID uuid.UUID, month, year int) ([]*entity.Budget, error) {
	var budgets []*entity.Budget
	err := r.db.WithContext(ctx).Model(&budgets).Where("user_id = ? AND month = ? AND year = ?", userID, month, year).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération budgets par mois: %w", err)
	}
	return budgets, nil
}
