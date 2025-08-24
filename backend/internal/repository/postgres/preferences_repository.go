package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// PreferencesRepository implémente repository.PreferencesRepository
type PreferencesRepository struct {
	db *pg.DB
}

// NewPreferencesRepository crée une nouvelle instance de PreferencesRepository
func NewPreferencesRepository(db *pg.DB) repository.PreferencesRepository {
	return &PreferencesRepository{db: db}
}

// Create crée de nouvelles préférences utilisateur
func (r *PreferencesRepository) Create(ctx context.Context, preferences *entity.UserPreferences) error {
	_, err := r.db.WithContext(ctx).Model(preferences).Insert()
	if err != nil {
		return fmt.Errorf("erreur création préférences: %w", err)
	}
	return nil
}

// GetByUserID récupère les préférences d'un utilisateur
func (r *PreferencesRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserPreferences, error) {
	preferences := &entity.UserPreferences{}
	err := r.db.WithContext(ctx).Model(preferences).Where("user_id = ?", userID).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil // Pas de préférences trouvées
		}
		return nil, fmt.Errorf("erreur récupération préférences: %w", err)
	}
	return preferences, nil
}

// Update met à jour les préférences d'un utilisateur
func (r *PreferencesRepository) Update(ctx context.Context, preferences *entity.UserPreferences) error {
	_, err := r.db.WithContext(ctx).Model(preferences).Where("user_id = ?", preferences.UserID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour préférences: %w", err)
	}
	return nil
}

// Delete supprime les préférences d'un utilisateur
func (r *PreferencesRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.UserPreferences{}).Where("user_id = ?", userID).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression préférences: %w", err)
	}
	return nil
}
