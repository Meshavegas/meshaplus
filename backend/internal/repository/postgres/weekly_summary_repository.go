package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// WeeklySummaryRepository implémente repository.WeeklySummaryRepository
type WeeklySummaryRepository struct {
	db *pg.DB
}

// NewWeeklySummaryRepository crée une nouvelle instance de WeeklySummaryRepository
func NewWeeklySummaryRepository(db *pg.DB) repository.WeeklySummaryRepository {
	return &WeeklySummaryRepository{db: db}
}

// Create crée un nouveau résumé hebdomadaire
func (r *WeeklySummaryRepository) Create(ctx context.Context, summary *entity.WeeklySummary) error {
	_, err := r.db.WithContext(ctx).Model(summary).Insert()
	if err != nil {
		return fmt.Errorf("erreur création résumé hebdomadaire: %w", err)
	}
	return nil
}

// GetByID récupère un résumé hebdomadaire par son ID
func (r *WeeklySummaryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.WeeklySummary, error) {
	summary := &entity.WeeklySummary{}
	err := r.db.WithContext(ctx).Model(summary).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("résumé hebdomadaire non trouvé")
		}
		return nil, fmt.Errorf("erreur récupération résumé hebdomadaire: %w", err)
	}
	return summary, nil
}

// GetByUserID récupère tous les résumés hebdomadaires d'un utilisateur
func (r *WeeklySummaryRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.WeeklySummary, error) {
	var summaries []*entity.WeeklySummary
	err := r.db.WithContext(ctx).Model(&summaries).Where("user_id = ?", userID).Order("week_start_date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération résumés hebdomadaires utilisateur: %w", err)
	}
	return summaries, nil
}

// GetByWeek récupère le résumé d'une semaine spécifique
func (r *WeeklySummaryRepository) GetByWeek(ctx context.Context, userID uuid.UUID, weekStartDate time.Time) (*entity.WeeklySummary, error) {
	summary := &entity.WeeklySummary{}
	err := r.db.WithContext(ctx).Model(summary).Where("user_id = ? AND week_start_date = ?", userID, weekStartDate).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("résumé hebdomadaire non trouvé pour cette semaine")
		}
		return nil, fmt.Errorf("erreur récupération résumé hebdomadaire par semaine: %w", err)
	}
	return summary, nil
}

// GetByDateRange récupère les résumés dans une plage de dates
func (r *WeeklySummaryRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*entity.WeeklySummary, error) {
	var summaries []*entity.WeeklySummary
	err := r.db.WithContext(ctx).Model(&summaries).Where("user_id = ? AND week_start_date >= ? AND week_start_date <= ?", userID, startDate, endDate).Order("week_start_date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération résumés hebdomadaires par plage de dates: %w", err)
	}
	return summaries, nil
}

// GetLatest récupère le dernier résumé hebdomadaire d'un utilisateur
func (r *WeeklySummaryRepository) GetLatest(ctx context.Context, userID uuid.UUID) (*entity.WeeklySummary, error) {
	summary := &entity.WeeklySummary{}
	err := r.db.WithContext(ctx).Model(summary).Where("user_id = ?", userID).Order("week_start_date DESC").First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("aucun résumé hebdomadaire trouvé")
		}
		return nil, fmt.Errorf("erreur récupération dernier résumé hebdomadaire: %w", err)
	}
	return summary, nil
}

// Update met à jour un résumé hebdomadaire
func (r *WeeklySummaryRepository) Update(ctx context.Context, summary *entity.WeeklySummary) error {
	_, err := r.db.WithContext(ctx).Model(summary).Where("id = ?", summary.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour résumé hebdomadaire: %w", err)
	}
	return nil
}

// Delete supprime un résumé hebdomadaire
func (r *WeeklySummaryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.WeeklySummary{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression résumé hebdomadaire: %w", err)
	}
	return nil
}

// Upsert crée ou met à jour un résumé hebdomadaire
func (r *WeeklySummaryRepository) Upsert(ctx context.Context, summary *entity.WeeklySummary) error {
	_, err := r.db.WithContext(ctx).Model(summary).
		OnConflict("(user_id, week_start_date) DO UPDATE").
		Insert()
	if err != nil {
		return fmt.Errorf("erreur upsert résumé hebdomadaire: %w", err)
	}
	return nil
}
