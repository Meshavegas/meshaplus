package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// AccountRepository implémente repository.AccountRepository
type AccountRepository struct {
	db *pg.DB
}

// NewAccountRepository crée une nouvelle instance de AccountRepository
func NewAccountRepository(db *pg.DB) repository.AccountRepository {
	return &AccountRepository{db: db}
}

// Create crée un nouveau compte
func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) error {
	_, err := r.db.WithContext(ctx).Model(account).Insert()
	if err != nil {
		return fmt.Errorf("erreur création compte: %w", err)
	}
	return nil
}

// GetByID récupère un compte par son ID
func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	account := &entity.Account{}
	err := r.db.WithContext(ctx).Model(account).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("compte non trouvé")
		}
		return nil, fmt.Errorf("erreur récupération compte: %w", err)
	}
	return account, nil
}

// GetByUserID récupère tous les comptes d'un utilisateur
func (r *AccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error) {
	var accounts []*entity.Account
	err := r.db.WithContext(ctx).Model(&accounts).Where("user_id = ?", userID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération comptes utilisateur: %w", err)
	}
	return accounts, nil
}

// Update met à jour un compte
func (r *AccountRepository) Update(ctx context.Context, account *entity.Account) error {
	_, err := r.db.WithContext(ctx).Model(account).Where("id = ?", account.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour compte: %w", err)
	}
	return nil
}

// Delete supprime un compte
func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Account{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression compte: %w", err)
	}
	return nil
}

// GetBalanceByUserID récupère le solde total d'un utilisateur
func (r *AccountRepository) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&entity.Account{}).Column("balance").Where("user_id = ?", userID).Select(&total)
	if err != nil {
		return 0, fmt.Errorf("erreur récupération solde total: %w", err)
	}
	return total, nil
} 