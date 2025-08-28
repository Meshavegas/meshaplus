package service

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"backend/internal/service/ai"
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
	categoryRepo    repository.CategoryRepository
	aiService       *ai.AIService
	logger          logger.Logger
	accountService  *AccountService
}

// NewTransactionService crée une nouvelle instance de TransactionService
func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
	aiService *ai.AIService,
	accountService *AccountService,
	logger logger.Logger,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		categoryRepo:    categoryRepo,
		aiService:       aiService,
		accountService:  accountService,
		logger:          logger,
	}
}

// CreateTransaction crée une nouvelle transaction
func (s *TransactionService) CreateTransaction(ctx context.Context, userID uuid.UUID, req entity.CreateTransactionRequest) (*entity.Transaction, error) {
	// Validation des données
	if req.Amount <= 0 {
		return nil, fmt.Errorf("le montant doit être positif")
	}

	validTypes := []string{"income", "expense", "transfer", "saving", "refund"}
	isValidType := false
	for _, validType := range validTypes {
		if req.Type == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return nil, fmt.Errorf("le type doit être l'un des suivants: %v", validTypes)
	}

	// Vérifier que le compte existe et appartient à l'utilisateur (si spécifié)
	if req.AccountID != nil {
		account, err := s.accountRepo.GetByID(ctx, *req.AccountID)

		s.logger.Info("Compte récupéré", logger.String("account_id", account.ID.String()), logger.String("user_id", userID.String()), logger.Float64("account_name", account.Balance))
		if err != nil {
			return nil, fmt.Errorf("compte non trouvé")
		}
		if account.UserID != userID {
			return nil, fmt.Errorf("accès non autorisé au compte")
		}
	}

	// Gestion de la catégorie
	var categoryID *uuid.UUID = req.CategoryID

	// Si aucune catégorie n'est spécifiée, utiliser l'IA pour en créer une automatiquement
	if categoryID == nil && req.Description != "" {
		// Déterminer le type de catégorie basé sur le type de transaction
		categoryType := "expense"
		if req.Type == "income" {
			categoryType = "revenue"
		}

		// Utiliser l'IA pour catégoriser automatiquement
		categoryResponse, err := s.aiService.GenerateCatherorie([]string{}, req.Description)
		if err != nil {
			s.logger.Warn("Erreur catégorisation automatique, transaction créée sans catégorie", logger.Error(err))
		} else {
			// Log de debug pour voir la réponse de l'IA
			s.logger.Info("Réponse IA catégorisation",
				logger.String("category_name", categoryResponse.CategoryName),
				logger.String("icon", categoryResponse.Icon),
				logger.String("color", categoryResponse.Color),
				logger.Bool("is_new", categoryResponse.IsNewCategory),
				logger.Int("confidence", categoryResponse.Confidence))
			// Créer la nouvelle catégorie dans la base de données
			// Utiliser les valeurs de l'IA ou des valeurs par défaut si elles sont vides
			icon := categoryResponse.Icon
			if icon == "" {
				icon = "md:category" // Icône par défaut
			}

			color := categoryResponse.Color
			if color == "" {
				color = "#6C5CE7" // Couleur par défaut
			}

			newCategory := &entity.Category{
				UserID: userID,
				Name:   categoryResponse.CategoryName,
				Type:   categoryType,
				Icon:   icon,
				Color:  color,
			}

			if err := s.categoryRepo.Create(ctx, newCategory); err != nil {
				s.logger.Warn("Erreur création catégorie automatique", logger.Error(err))
			} else {
				categoryID = &newCategory.ID
				s.logger.Info("Catégorie créée automatiquement",
					logger.String("category_name", categoryResponse.CategoryName),
					logger.String("category_id", newCategory.ID.String()),
					logger.Int("confidence", categoryResponse.Confidence),
				)
			}
		}
	} else if categoryID != nil {
		// Vérifier que la catégorie existe et appartient à l'utilisateur
		category, err := s.categoryRepo.GetByID(ctx, userID, *categoryID)
		if err != nil {
			return nil, fmt.Errorf("catégorie non trouvée")
		}
		if category.UserID != userID {
			return nil, fmt.Errorf("accès non autorisé à la catégorie")
		}
	}

	if req.Type == "transfer" {
		account, err := s.accountRepo.GetByID(ctx, *req.AccountID)
		if err != nil {
			s.logger.Error("Erreur récupération compte", logger.Error(err))
			return nil, fmt.Errorf("erreur récupération compte: %w", err)
		}

		//expense accountID transaction
		transaction := &entity.Transaction{
			ID:          uuid.New(),
			UserID:      userID,
			AccountID:   req.AccountID,
			Type:        "expense",
			Amount:      req.Amount,
			Description: req.Description,
			Date:        req.Date,
			Recurring:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.transactionRepo.Create(ctx, transaction); err != nil {
			s.logger.Error("Erreur création transaction", logger.Error(err))
			return nil, fmt.Errorf("erreur création transaction: %w", err)
		}

		// Mettre à jour la balance du compte source (expense)
		account.Balance -= req.Amount
		if err := s.accountRepo.Update(ctx, account); err != nil {
			s.logger.Error("Erreur mise à jour balance du compte source", logger.Error(err))
		}

		//income toAccountID transaction
		toAccount, err := s.accountRepo.GetByID(ctx, *req.ToAccountID)
		if err != nil {
			s.logger.Error("Erreur récupération compte destination", logger.Error(err))
			return nil, fmt.Errorf("erreur récupération compte destination: %w", err)
		}

		transaction2 := &entity.Transaction{
			ID:          uuid.New(),
			UserID:      userID,
			AccountID:   req.ToAccountID,
			Type:        "income",
			Amount:      req.Amount,
			Description: req.Description,
			Date:        req.Date,
			Recurring:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.transactionRepo.Create(ctx, transaction2); err != nil {
			s.logger.Error("Erreur création transaction destination", logger.Error(err))
			return nil, fmt.Errorf("erreur création transaction destination: %w", err)
		}

		// Mettre à jour la balance du compte destination (income)
		toAccount.Balance += req.Amount
		if err := s.accountRepo.Update(ctx, toAccount); err != nil {
			s.logger.Error("Erreur mise à jour balance du compte destination", logger.Error(err))
		}

		return transaction, nil
	}

	// Création de la transaction
	account, err := s.accountRepo.GetByID(ctx, *req.AccountID)
	if err != nil {
		s.logger.Error("Erreur récupération compte", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération compte: %w", err)
	}

	transaction := &entity.Transaction{
		ID:           uuid.New(),
		UserID:       userID,
		AccountID:    req.AccountID,
		CategoryID:   categoryID,
		Type:         req.Type,
		SavingGoalID: req.SavingGoalID,
		Amount:       req.Amount,
		Description:  req.Description,
		Date:         req.Date,
		Recurring:    req.Recurring,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		s.logger.Error("Erreur création transaction", logger.Error(err))
		return nil, fmt.Errorf("erreur création transaction: %w", err)
	}

	// Mettre à jour la balance du compte selon le type de transaction
	switch req.Type {
	case "income", "refund":
		account.Balance += req.Amount
	case "expense", "saving":
		account.Balance -= req.Amount
	}

	// Sauvegarder la mise à jour du compte
	if err := s.accountRepo.Update(ctx, account); err != nil {
		s.logger.Error("Erreur mise à jour balance du compte", logger.Error(err))
		// Note: On ne fait pas échouer la création de transaction si la mise à jour de balance échoue
		// Mais on log l'erreur pour le debugging
	}

	s.logger.Info("Transaction créée avec succès",
		logger.String("transaction_id", transaction.ID.String()),
		logger.String("user_id", userID.String()),
		logger.String("type", transaction.Type),
		logger.Float64("amount", transaction.Amount),
		logger.String("category_id", func() string {
			if categoryID != nil {
				return categoryID.String()
			}
			return "none"
		}()),
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
		transaction.CategoryID = req.CategoryID
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
		transaction.AccountID = req.AccountID
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

// GetTransactionsByDateRange récupère les transactions dans une plage de dates
func (s *TransactionService) GetTransactionsByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*entity.Transaction, error) {
	transactions, err := s.transactionRepo.GetByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		s.logger.Error("Erreur récupération transactions par plage de dates", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionsByUserID récupère toutes les transactions d'un utilisateur
func (s *TransactionService) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Transaction, error) {
	transactions, err := s.transactionRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération transactions utilisateur", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionsByAccountID récupère toutes les transactions d'un compte spécifique avec les détails de la catégorie
func (s *TransactionService) GetTransactionsByAccountID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]*entity.TransactionWithDetails, error) {
	transactions, err := s.transactionRepo.GetByAccountIDWithCategoryDetails(ctx, userID, &accountID)
	if err != nil {
		s.logger.Error("Erreur récupération transactions par compte", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération transactions: %w", err)
	}

	return transactions, nil
}
