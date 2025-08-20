package service

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"backend/pkg/logger"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TransactionQuery représente les paramètres de requête pour les transactions
type TransactionQuery struct {
	Type       string
	CategoryID *uuid.UUID
	AccountID  *uuid.UUID
	StartDate  string
	EndDate    string
	Page       int
	Limit      int
}

// TransactionService gère la logique métier des transactions
type TransactionService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	logger          logger.Logger
}

// NewTransactionService crée une nouvelle instance de TransactionService
func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	logger logger.Logger,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		logger:          logger,
	}
}

// CreateTransaction crée une nouvelle transaction
func (s *TransactionService) CreateTransaction(ctx context.Context, userID uuid.UUID, req entity.CreateTransactionRequest) (*entity.Transaction, error) {
	// Validation des données
	if req.Amount <= 0 {
		return nil, fmt.Errorf("le montant doit être positif")
	}

	if req.Type != "income" && req.Type != "expense" {
		return nil, fmt.Errorf("le type doit être 'income' ou 'expense'")
	}

	// Vérifier que le compte existe et appartient à l'utilisateur
	account, err := s.accountRepo.GetByID(ctx, req.AccountID)
	if err != nil {
		return nil, fmt.Errorf("compte non trouvé")
	}
	if account.UserID != userID {
		return nil, fmt.Errorf("accès non autorisé au compte")
	}

	// Création de la transaction
	transaction := &entity.Transaction{
		UserID:      userID,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Type:        req.Type,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
		Recurring:   req.Recurring,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		s.logger.Error("Erreur création transaction", logger.Error(err))
		return nil, fmt.Errorf("erreur création transaction: %w", err)
	}

	s.logger.Info("Transaction créée avec succès",
		logger.String("transaction_id", transaction.ID.String()),
		logger.String("user_id", userID.String()),
		logger.String("type", transaction.Type),
		logger.Float64("amount", transaction.Amount),
	)

	return transaction, nil
}

// GetTransaction récupère une transaction par son ID
func (s *TransactionService) GetTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*entity.Transaction, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		s.logger.Error("Erreur récupération transaction", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération transaction: %w", err)
	}

	// Vérifier que la transaction appartient à l'utilisateur
	if transaction.UserID != userID {
		s.logger.Warn("Tentative d'accès non autorisé à une transaction",
			logger.String("user_id", userID.String()),
			logger.String("transaction_id", transactionID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	return transaction, nil
}

// GetTransactions récupère les transactions avec filtres et pagination
func (s *TransactionService) GetTransactions(ctx context.Context, userID uuid.UUID, query TransactionQuery) ([]*entity.Transaction, int64, error) {
	// Pour l'instant, récupérer toutes les transactions de l'utilisateur
	// TODO: Implémenter les filtres quand le repository sera créé
	transactions, err := s.transactionRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération transactions", logger.Error(err))
		return nil, 0, fmt.Errorf("erreur récupération transactions: %w", err)
	}

	// Calculer le total
	total := int64(len(transactions))

	// Appliquer la pagination basique
	start := (query.Page - 1) * query.Limit
	end := start + query.Limit
	if int64(start) >= total {
		return []*entity.Transaction{}, total, nil
	}
	if int64(end) > total {
		end = int(total)
	}

	return transactions[start:end], total, nil
}

// UpdateTransaction met à jour une transaction
func (s *TransactionService) UpdateTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID, req entity.UpdateTransactionRequest) (*entity.Transaction, error) {
	// Récupérer la transaction existante
	transaction, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		s.logger.Error("Erreur récupération transaction pour mise à jour", logger.Error(err))
		return nil, fmt.Errorf("transaction non trouvée")
	}

	// Vérifier que la transaction appartient à l'utilisateur
	if transaction.UserID != userID {
		s.logger.Warn("Tentative de mise à jour non autorisée d'une transaction",
			logger.String("user_id", userID.String()),
			logger.String("transaction_id", transactionID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	// Mettre à jour les champs
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, fmt.Errorf("le montant doit être positif")
		}
		transaction.Amount = *req.Amount
	}

	if req.Type != nil {
		if *req.Type != "income" && *req.Type != "expense" {
			return nil, fmt.Errorf("le type doit être 'income' ou 'expense'")
		}
		transaction.Type = *req.Type
	}

	if req.Description != nil {
		transaction.Description = *req.Description
	}

	if req.Date != nil {
		transaction.Date = *req.Date
	}

	if req.CategoryID != nil {
		transaction.CategoryID = *req.CategoryID
	}

	if req.AccountID != nil {
		// Vérifier que le nouveau compte appartient à l'utilisateur
		account, err := s.accountRepo.GetByID(ctx, *req.AccountID)
		if err != nil {
			return nil, fmt.Errorf("compte non trouvé")
		}
		if account.UserID != userID {
			return nil, fmt.Errorf("accès non autorisé au compte")
		}
		transaction.AccountID = *req.AccountID
	}

	transaction.UpdatedAt = time.Now()

	// Sauvegarder les modifications
	if err := s.transactionRepo.Update(ctx, transaction); err != nil {
		s.logger.Error("Erreur mise à jour transaction", logger.Error(err))
		return nil, fmt.Errorf("erreur mise à jour transaction: %w", err)
	}

	s.logger.Info("Transaction mise à jour avec succès",
		logger.String("transaction_id", transaction.ID.String()),
		logger.String("user_id", userID.String()),
	)

	return transaction, nil
}

// DeleteTransaction supprime une transaction
func (s *TransactionService) DeleteTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) error {
	// Récupérer la transaction existante
	transaction, err := s.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		s.logger.Error("Erreur récupération transaction pour suppression", logger.Error(err))
		return fmt.Errorf("transaction non trouvée")
	}

	// Vérifier que la transaction appartient à l'utilisateur
	if transaction.UserID != userID {
		s.logger.Warn("Tentative de suppression non autorisée d'une transaction",
			logger.String("user_id", userID.String()),
			logger.String("transaction_id", transactionID.String()),
		)
		return fmt.Errorf("accès non autorisé")
	}

	// Supprimer la transaction
	if err := s.transactionRepo.Delete(ctx, transactionID); err != nil {
		s.logger.Error("Erreur suppression transaction", logger.Error(err))
		return fmt.Errorf("erreur suppression transaction: %w", err)
	}

	s.logger.Info("Transaction supprimée avec succès",
		logger.String("transaction_id", transactionID.String()),
		logger.String("user_id", userID.String()),
	)

	return nil
}

// GetTransactionStats récupère les statistiques des transactions
func (s *TransactionService) GetTransactionStats(ctx context.Context, userID uuid.UUID, period string) (map[string]interface{}, error) {
	// Validation de la période
	if period != "week" && period != "month" && period != "year" {
		period = "month" // Valeur par défaut
	}

	// Pour l'instant, retourner des statistiques basiques
	// TODO: Implémenter les vraies statistiques quand le repository sera créé
	transactions, err := s.transactionRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération statistiques transactions", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération statistiques: %w", err)
	}

	// Calculer les statistiques basiques
	var totalIncome, totalExpense float64
	for _, tx := range transactions {
		if tx.Type == "income" {
			totalIncome += tx.Amount
		} else if tx.Type == "expense" {
			totalExpense += tx.Amount
		}
	}

	stats := map[string]interface{}{
		"period":        period,
		"total_income":  totalIncome,
		"total_expense": totalExpense,
		"net_amount":    totalIncome - totalExpense,
		"count":         len(transactions),
	}

	return stats, nil
}
