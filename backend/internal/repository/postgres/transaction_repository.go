package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// TransactionRepository implémente repository.TransactionRepository
type TransactionRepository struct {
	db *pg.DB
}

// NewTransactionRepository crée une nouvelle instance de TransactionRepository
func NewTransactionRepository(db *pg.DB) repository.TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create crée une nouvelle transaction
func (r *TransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	_, err := r.db.WithContext(ctx).Model(transaction).Insert()
	if err != nil {
		return fmt.Errorf("erreur création transaction: %w", err)
	}
	return nil
}

// GetByID récupère une transaction par son ID
func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	transaction := &entity.Transaction{}
	err := r.db.WithContext(ctx).Model(transaction).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("transaction non trouvée")
		}
		return nil, fmt.Errorf("erreur récupération transaction: %w", err)
	}
	return transaction, nil
}

// GetByUserID récupère toutes les transactions d'un utilisateur
func (r *TransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ?", userID).Order("date DESC").Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions utilisateur: %w", err)
	}
	return transactions, nil
}

// Update met à jour une transaction
func (r *TransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	_, err := r.db.WithContext(ctx).Model(transaction).Where("id = ?", transaction.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour transaction: %w", err)
	}
	return nil
}

// Delete supprime une transaction
func (r *TransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Transaction{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression transaction: %w", err)
	}
	return nil
}

// GetByCategoryID récupère les transactions d'une catégorie
func (r *TransactionRepository) GetByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ? AND category_id = ?", userID, categoryID).Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions par catégorie: %w", err)
	}
	return transactions, nil
}

// GetByDateRange récupère les transactions dans une plage de dates
func (r *TransactionRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions par plage de dates: %w", err)
	}
	return transactions, nil
}

// GetByType récupère les transactions par type
func (r *TransactionRepository) GetByType(ctx context.Context, userID uuid.UUID, txType string) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ? AND type = ?", userID, txType).Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions par type: %w", err)
	}
	return transactions, nil
}

// GetTotalByUserID récupère le total des transactions d'un utilisateur
func (r *TransactionRepository) GetTotalByUserID(ctx context.Context, userID uuid.UUID) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).Model(&entity.Transaction{}).Column("amount").Where("user_id = ?", userID).Select(&total)
	if err != nil {
		return 0, fmt.Errorf("erreur récupération total transactions: %w", err)
	}
	return total, nil
}
