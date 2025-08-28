package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// MotivationRepository implémente repository.MotivationRepository
type MotivationRepository struct {
	db *pg.DB
}

// NewMotivationRepository crée une nouvelle instance de MotivationRepository
func NewMotivationRepository(db *pg.DB) repository.MotivationRepository {
	return &MotivationRepository{db: db}
}

// Create crée une nouvelle motivation
func (r *MotivationRepository) Create(ctx context.Context, motivation *entity.Motivation) error {
	_, err := r.db.WithContext(ctx).Model(motivation).Insert()
	if err != nil {
		return fmt.Errorf("erreur création motivation: %w", err)
	}
	return nil
}

// GetByID récupère une motivation par son ID
func (r *MotivationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Motivation, error) {
	motivation := &entity.Motivation{}
	err := r.db.WithContext(ctx).Model(motivation).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("motivation non trouvée")
		}
		return nil, fmt.Errorf("erreur récupération motivation: %w", err)
	}
	return motivation, nil
}

// GetByUserID récupère toutes les motivations d'un utilisateur
func (r *MotivationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Motivation, error) {
	var motivations []*entity.Motivation
	err := r.db.WithContext(ctx).Model(&motivations).Where("user_id = ?", userID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération motivations utilisateur: %w", err)
	}
	return motivations, nil
}

// GetByType récupère les motivations par type
func (r *MotivationRepository) GetByType(ctx context.Context, userID uuid.UUID, motivationType string) ([]*entity.Motivation, error) {
	var motivations []*entity.Motivation
	err := r.db.WithContext(ctx).Model(&motivations).Where("user_id = ? AND type = ?", userID, motivationType).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération motivations par type: %w", err)
	}
	return motivations, nil
}

// GetRandom récupère une motivation aléatoire d'un utilisateur
func (r *MotivationRepository) GetRandom(ctx context.Context, userID uuid.UUID) (*entity.Motivation, error) {
	motivation := &entity.Motivation{}
	err := r.db.WithContext(ctx).Model(motivation).Where("user_id = ?", userID).OrderExpr("RANDOM()").First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("aucune motivation trouvée")
		}
		return nil, fmt.Errorf("erreur récupération motivation aléatoire: %w", err)
	}
	return motivation, nil
}

// GetRandomByType récupère une motivation aléatoire d'un type spécifique
func (r *MotivationRepository) GetRandomByType(ctx context.Context, userID uuid.UUID, motivationType string) (*entity.Motivation, error) {
	motivation := &entity.Motivation{}
	err := r.db.WithContext(ctx).Model(motivation).Where("user_id = ? AND type = ?", userID, motivationType).OrderExpr("RANDOM()").First()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("aucune motivation de ce type trouvée")
		}
		return nil, fmt.Errorf("erreur récupération motivation aléatoire par type: %w", err)
	}
	return motivation, nil
}

// Update met à jour une motivation
func (r *MotivationRepository) Update(ctx context.Context, motivation *entity.Motivation) error {
	_, err := r.db.WithContext(ctx).Model(motivation).Where("id = ?", motivation.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour motivation: %w", err)
	}
	return nil
}

// Delete supprime une motivation
func (r *MotivationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Motivation{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression motivation: %w", err)
	}
	return nil
}

// Search recherche des motivations par contenu
func (r *MotivationRepository) Search(ctx context.Context, userID uuid.UUID, query string) ([]*entity.Motivation, error) {
	var motivations []*entity.Motivation
	err := r.db.WithContext(ctx).Model(&motivations).Where("user_id = ? AND content ILIKE ?", userID, "%"+query+"%").Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur recherche motivations: %w", err)
	}
	return motivations, nil
}
