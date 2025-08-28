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
func (r *TransactionRepository) GetByCategoryID(ctx context.Context, userID uuid.UUID, categoryID *uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	query := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ?", userID)
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	} else {
		query = query.Where("category_id IS NULL")
	}
	err := query.Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions par catégorie: %w", err)
	}
	return transactions, nil
}

// GetByAccountID récupère les transactions d'un compte
func (r *TransactionRepository) GetByAccountID(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	query := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ?", userID)
	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	} else {
		query = query.Where("account_id IS NULL")
	}
	err := query.Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions par compte: %w", err)
	}
	return transactions, nil
}

// GetBySavingGoalID récupère les transactions d'un objectif d'épargne
func (r *TransactionRepository) GetBySavingGoalID(ctx context.Context, userID uuid.UUID, savingGoalID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ? AND saving_goal_id = ?", userID, savingGoalID).Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions par objectif d'épargne: %w", err)
	}
	return transactions, nil
}

// GetTransfers récupère les transferts entre comptes
func (r *TransactionRepository) GetTransfers(ctx context.Context, userID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ? AND type = ? AND to_account_id IS NOT NULL", userID, "transfer").Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transferts: %w", err)
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

// get all transactions by user id and account id order ber  created_at desc
func (r *TransactionRepository) GetAllTransactionsByUserIDAndAccountID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ? AND account_id = ?", userID, accountID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions par utilisateur et compte: %w", err)
	}
	return transactions, nil
}

// get all transactions by user id and saving goal id order by created_at desc
func (r *TransactionRepository) GetAllTransactionsBySavingGoalID(ctx context.Context, userID uuid.UUID, savingGoalID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ? AND saving_goal_id = ?", userID, savingGoalID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions par utilisateur et objectif d'épargne: %w", err)
	}
	return transactions, nil
}

// GetByAccountIDWithCategoryDetails récupère les transactions d'un compte avec les détails complets de la catégorie
func (r *TransactionRepository) GetByAccountIDWithCategoryDetails(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]*entity.TransactionWithDetails, error) {
	// D'abord, récupérer les transactions
	var transactions []*entity.Transaction
	query := r.db.WithContext(ctx).Model(&transactions).Where("user_id = ?", userID)

	if accountID != nil {
		query = query.Where("account_id = ?", *accountID)
	} else {
		query = query.Where("account_id IS NULL")
	}

	err := query.Order("date DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération transactions: %w", err)
	}

	// Ensuite, récupérer les catégories pour ces transactions
	var categoryIDs []uuid.UUID
	for _, tx := range transactions {
		if tx.CategoryID != nil {
			categoryIDs = append(categoryIDs, *tx.CategoryID)
		}
	}

	// Récupérer toutes les catégories en une seule requête avec SQL brut
	var categories []*entity.Category
	if len(categoryIDs) > 0 {
		query := `
			SELECT id, user_id, name, type, parent_id, icon, color, created_at
			FROM categories 
			WHERE id = ANY(?)
		`

		_, err = r.db.WithContext(ctx).Query(&categories, query, pg.Array(categoryIDs))
		if err != nil {
			return nil, fmt.Errorf("erreur récupération catégories: %w", err)
		}

	}

	// Créer un map pour un accès rapide aux catégories
	categoryMap := make(map[uuid.UUID]*entity.Category)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat
	}

	// Construire le résultat final
	var result []*entity.TransactionWithDetails
	for _, tx := range transactions {
		txWithDetails := &entity.TransactionWithDetails{
			ID:           tx.ID,
			UserID:       tx.UserID,
			AccountID:    tx.AccountID,
			CategoryID:   tx.CategoryID,
			Type:         tx.Type,
			ToAccountID:  tx.ToAccountID,
			SavingGoalID: tx.SavingGoalID,
			Amount:       tx.Amount,
			Description:  tx.Description,
			Date:         tx.Date,
			Recurring:    tx.Recurring,
			CreatedAt:    tx.CreatedAt,
			UpdatedAt:    tx.UpdatedAt,
		}

		// Ajouter les détails de la catégorie si elle existe
		if tx.CategoryID != nil {
			if category, exists := categoryMap[*tx.CategoryID]; exists {
				txWithDetails.Category = category
			}
		}

		result = append(result, txWithDetails)
	}

	return result, nil
}
