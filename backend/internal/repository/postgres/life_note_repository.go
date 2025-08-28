package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// LifeNoteRepository implémente repository.LifeNoteRepository
type LifeNoteRepository struct {
	db *pg.DB
}

// NewLifeNoteRepository crée une nouvelle instance de LifeNoteRepository
func NewLifeNoteRepository(db *pg.DB) repository.LifeNoteRepository {
	return &LifeNoteRepository{db: db}
}

// Create crée une nouvelle note de vie
func (r *LifeNoteRepository) Create(ctx context.Context, note *entity.LifeNote) error {
	_, err := r.db.WithContext(ctx).Model(note).Insert()
	if err != nil {
		return fmt.Errorf("erreur création note de vie: %w", err)
	}
	return nil
}

// GetByID récupère une note de vie par son ID
func (r *LifeNoteRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.LifeNote, error) {
	note := &entity.LifeNote{}
	err := r.db.WithContext(ctx).Model(note).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("note de vie non trouvée")
		}
		return nil, fmt.Errorf("erreur récupération note de vie: %w", err)
	}
	return note, nil
}

// GetByUserID récupère toutes les notes de vie d'un utilisateur
func (r *LifeNoteRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.LifeNote, error) {
	var notes []*entity.LifeNote
	err := r.db.WithContext(ctx).Model(&notes).Where("user_id = ?", userID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération notes de vie utilisateur: %w", err)
	}
	return notes, nil
}

// GetByGoalID récupère les notes de vie liées à un objectif
func (r *LifeNoteRepository) GetByGoalID(ctx context.Context, userID uuid.UUID, goalID uuid.UUID) ([]*entity.LifeNote, error) {
	var notes []*entity.LifeNote
	err := r.db.WithContext(ctx).Model(&notes).Where("user_id = ? AND related_goal_id = ?", userID, goalID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération notes de vie par objectif: %w", err)
	}
	return notes, nil
}

// GetRecent récupère les notes de vie récentes d'un utilisateur
func (r *LifeNoteRepository) GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.LifeNote, error) {
	var notes []*entity.LifeNote
	err := r.db.WithContext(ctx).Model(&notes).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération notes de vie récentes: %w", err)
	}
	return notes, nil
}

// Update met à jour une note de vie
func (r *LifeNoteRepository) Update(ctx context.Context, note *entity.LifeNote) error {
	_, err := r.db.WithContext(ctx).Model(note).Where("id = ?", note.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour note de vie: %w", err)
	}
	return nil
}

// Delete supprime une note de vie
func (r *LifeNoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.LifeNote{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression note de vie: %w", err)
	}
	return nil
}

// SearchByTitle recherche des notes de vie par titre
func (r *LifeNoteRepository) SearchByTitle(ctx context.Context, userID uuid.UUID, query string) ([]*entity.LifeNote, error) {
	var notes []*entity.LifeNote
	err := r.db.WithContext(ctx).Model(&notes).Where("user_id = ? AND title ILIKE ?", userID, "%"+query+"%").Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur recherche notes de vie par titre: %w", err)
	}
	return notes, nil
}

// SearchByContent recherche des notes de vie par contenu
func (r *LifeNoteRepository) SearchByContent(ctx context.Context, userID uuid.UUID, query string) ([]*entity.LifeNote, error) {
	var notes []*entity.LifeNote
	err := r.db.WithContext(ctx).Model(&notes).Where("user_id = ? AND content ILIKE ?", userID, "%"+query+"%").Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur recherche notes de vie par contenu: %w", err)
	}
	return notes, nil
}

// Search recherche des notes de vie par titre et contenu
func (r *LifeNoteRepository) Search(ctx context.Context, userID uuid.UUID, query string) ([]*entity.LifeNote, error) {
	var notes []*entity.LifeNote
	err := r.db.WithContext(ctx).Model(&notes).Where("user_id = ? AND (title ILIKE ? OR content ILIKE ?)", userID, "%"+query+"%", "%"+query+"%").Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur recherche notes de vie: %w", err)
	}
	return notes, nil
}
