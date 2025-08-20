package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// SavingGoalRepository implémente repository.SavingGoalRepository
type SavingGoalRepository struct {
	db *pg.DB
}

// NewSavingGoalRepository crée une nouvelle instance de SavingGoalRepository
func NewSavingGoalRepository(db *pg.DB) repository.SavingGoalRepository {
	return &SavingGoalRepository{db: db}
}

// Create crée un nouvel objectif d'épargne
func (r *SavingGoalRepository) Create(ctx context.Context, goal *entity.SavingGoal) error {
	_, err := r.db.WithContext(ctx).Model(goal).Insert()
	if err != nil {
		return fmt.Errorf("erreur création objectif d'épargne: %w", err)
	}
	return nil
}

// GetByID récupère un objectif d'épargne par son ID
func (r *SavingGoalRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.SavingGoal, error) {
	goal := &entity.SavingGoal{}
	err := r.db.WithContext(ctx).Model(goal).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("objectif d'épargne non trouvé")
		}
		return nil, fmt.Errorf("erreur récupération objectif d'épargne: %w", err)
	}
	return goal, nil
}

// GetByUserID récupère tous les objectifs d'épargne d'un utilisateur
func (r *SavingGoalRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavingGoal, error) {
	var goals []*entity.SavingGoal
	err := r.db.WithContext(ctx).Model(&goals).Where("user_id = ?", userID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération objectifs d'épargne utilisateur: %w", err)
	}
	return goals, nil
}

// Update met à jour un objectif d'épargne
func (r *SavingGoalRepository) Update(ctx context.Context, goal *entity.SavingGoal) error {
	_, err := r.db.WithContext(ctx).Model(goal).Where("id = ?", goal.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour objectif d'épargne: %w", err)
	}
	return nil
}

// Delete supprime un objectif d'épargne
func (r *SavingGoalRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.SavingGoal{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression objectif d'épargne: %w", err)
	}
	return nil
}

// GetAchievedByUserID récupère les objectifs d'épargne atteints d'un utilisateur
func (r *SavingGoalRepository) GetAchievedByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavingGoal, error) {
	var goals []*entity.SavingGoal
	err := r.db.WithContext(ctx).Model(&goals).Where("user_id = ? AND is_achieved = ?", userID, true).Order("updated_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération objectifs d'épargne atteints: %w", err)
	}
	return goals, nil
}

// GetActiveByUserID récupère les objectifs d'épargne actifs d'un utilisateur
func (r *SavingGoalRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavingGoal, error) {
	var goals []*entity.SavingGoal
	err := r.db.WithContext(ctx).Model(&goals).Where("user_id = ? AND is_achieved = ?", userID, false).Order("deadline ASC NULLS LAST, created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération objectifs d'épargne actifs: %w", err)
	}
	return goals, nil
}
