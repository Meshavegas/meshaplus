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

// MoodRepository implémente repository.MoodRepository
type MoodRepository struct {
	db *pg.DB
}

// NewMoodRepository crée une nouvelle instance de MoodRepository
func NewMoodRepository(db *pg.DB) repository.MoodRepository {
	return &MoodRepository{db: db}
}

// Create crée un nouveau mood
func (r *MoodRepository) Create(ctx context.Context, mood *entity.Mood) error {
	_, err := r.db.WithContext(ctx).Model(mood).Insert()
	if err != nil {
		return fmt.Errorf("erreur création mood: %w", err)
	}
	return nil
}

// GetByID récupère un mood par son ID
func (r *MoodRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Mood, error) {
	mood := &entity.Mood{}
	err := r.db.WithContext(ctx).Model(mood).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("mood non trouvé")
		}
		return nil, fmt.Errorf("erreur récupération mood: %w", err)
	}
	return mood, nil
}

// GetByUserID récupère tous les moods d'un utilisateur
func (r *MoodRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Mood, error) {
	var moods []*entity.Mood
	err := r.db.WithContext(ctx).Model(&moods).Where("user_id = ?", userID).Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération moods utilisateur: %w", err)
	}
	return moods, nil
}

// GetByDate récupère le mood d'un utilisateur pour une date donnée
func (r *MoodRepository) GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) (*entity.Mood, error) {
	mood := &entity.Mood{}
	err := r.db.WithContext(ctx).Model(mood).Where("user_id = ? AND date = ?", userID, date).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("mood non trouvé pour cette date")
		}
		return nil, fmt.Errorf("erreur récupération mood par date: %w", err)
	}
	return mood, nil
}

// GetByDateRange récupère les moods dans une plage de dates
func (r *MoodRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*entity.Mood, error) {
	var moods []*entity.Mood
	err := r.db.WithContext(ctx).Model(&moods).Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération moods par plage de dates: %w", err)
	}
	return moods, nil
}

// Update met à jour un mood
func (r *MoodRepository) Update(ctx context.Context, mood *entity.Mood) error {
	_, err := r.db.WithContext(ctx).Model(mood).Where("id = ?", mood.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour mood: %w", err)
	}
	return nil
}

// Delete supprime un mood
func (r *MoodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Mood{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression mood: %w", err)
	}
	return nil
}

// GetByFeeling récupère les moods par sentiment
func (r *MoodRepository) GetByFeeling(ctx context.Context, userID uuid.UUID, feeling string) ([]*entity.Mood, error) {
	var moods []*entity.Mood
	err := r.db.WithContext(ctx).Model(&moods).Where("user_id = ? AND feeling = ?", userID, feeling).Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération moods par sentiment: %w", err)
	}
	return moods, nil
}
