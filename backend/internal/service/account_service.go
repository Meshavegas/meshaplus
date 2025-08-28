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

// AccountService gère la logique métier des comptes bancaires
type AccountService struct {
	accountRepo     repository.AccountRepository
	logger          logger.Logger
	transactionRepo repository.TransactionRepository
	savingGoalRepo  repository.SavingGoalRepository
}

// NewAccountService crée une nouvelle instance de AccountService
func NewAccountService(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	savingGoalRepo repository.SavingGoalRepository,
	logger logger.Logger,
) *AccountService {
	return &AccountService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		savingGoalRepo:  savingGoalRepo,
		logger:          logger,
	}
}

// CreateAccount crée un nouveau compte bancaire
func (s *AccountService) CreateAccount(ctx context.Context, userID uuid.UUID, req entity.CreateAccountRequest) (*entity.Account, error) {
	// Validation des données
	if req.Name == "" {
		return nil, fmt.Errorf("le nom du compte est requis")
	}

	if req.Type != "checking" && req.Type != "savings" && req.Type != "mobile_money" {
		return nil, fmt.Errorf("le type doit être 'checking', 'savings' ou 'mobile_money'")
	}

	if len(req.Currency) != 3 {
		return nil, fmt.Errorf("la devise doit avoir 3 caractères")
	}

	// Création du compte
	account := &entity.Account{
		UserID:        userID,
		Name:          req.Name,
		Type:          req.Type,
		Icon:          req.Icon,
		AccountNumber: req.AccountNumber,
		Color:         req.Color,
		Balance:       req.Balance,
		Currency:      req.Currency,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		s.logger.Error("Erreur création compte", logger.Error(err))
		return nil, fmt.Errorf("erreur création compte: %w", err)
	}

	s.logger.Info("Compte créé avec succès",
		logger.String("account_id", account.ID.String()),
		logger.String("user_id", userID.String()),
		logger.String("name", account.Name),
		logger.String("type", account.Type),
	)

	return account, nil
}

// GetAccount récupère un compte par son ID
func (s *AccountService) GetAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*entity.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		s.logger.Error("Erreur récupération compte", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération compte: %w", err)
	}

	// Vérifier que le compte appartient à l'utilisateur
	if account.UserID != userID {
		s.logger.Warn("Tentative d'accès non autorisé à un compte",
			logger.String("user_id", userID.String()),
			logger.String("account_id", accountID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	return account, nil
}

// GetAccounts récupère tous les comptes d'un utilisateur
func (s *AccountService) GetAccounts(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.Account, int64, error) {
	// Pour l'instant, récupérer tous les comptes de l'utilisateur
	// TODO: Implémenter la pagination quand le repository sera créé
	accounts, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération comptes", logger.Error(err))
		return nil, 0, fmt.Errorf("erreur récupération comptes: %w", err)
	}

	// Calculer le total
	total := int64(len(accounts))

	// Appliquer la pagination basique
	start := (page - 1) * limit
	end := start + limit
	if int64(start) >= total {
		return []*entity.Account{}, total, nil
	}
	if int64(end) > total {
		end = int(total)
	}

	return accounts[start:end], total, nil
}

// UpdateAccount met à jour un compte
func (s *AccountService) UpdateAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, req entity.UpdateAccountRequest) (*entity.Account, error) {
	// Récupérer le compte existant
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		s.logger.Error("Erreur récupération compte pour mise à jour", logger.Error(err))
		return nil, fmt.Errorf("compte non trouvé")
	}

	// Vérifier que le compte appartient à l'utilisateur
	if account.UserID != userID {
		s.logger.Warn("Tentative de mise à jour non autorisée d'un compte",
			logger.String("user_id", userID.String()),
			logger.String("account_id", accountID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	// Mettre à jour les champs
	if req.Name != nil {
		account.Name = *req.Name
	}

	if req.Type != nil {
		if *req.Type != "checking" && *req.Type != "savings" && *req.Type != "mobile_money" {
			return nil, fmt.Errorf("le type doit être 'checking', 'savings' ou 'mobile_money'")
		}
		account.Type = *req.Type
	}

	if req.Balance != nil {
		if *req.Balance < 0 {
			return nil, fmt.Errorf("le solde ne peut pas être négatif")
		}
		account.Balance = *req.Balance
	}

	if req.Currency != nil {
		if len(*req.Currency) != 3 {
			return nil, fmt.Errorf("la devise doit avoir 3 caractères")
		}
		account.Currency = *req.Currency
	}

	account.UpdatedAt = time.Now()

	// Sauvegarder les modifications
	if err := s.accountRepo.Update(ctx, account); err != nil {
		s.logger.Error("Erreur mise à jour compte", logger.Error(err))
		return nil, fmt.Errorf("erreur mise à jour compte: %w", err)
	}

	s.logger.Info("Compte mis à jour avec succès",
		logger.String("account_id", account.ID.String()),
		logger.String("user_id", userID.String()),
	)

	return account, nil
}

// DeleteAccount supprime un compte
func (s *AccountService) DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	// Récupérer le compte existant
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		s.logger.Error("Erreur récupération compte pour suppression", logger.Error(err))
		return fmt.Errorf("compte non trouvé")
	}

	// Vérifier que le compte appartient à l'utilisateur
	if account.UserID != userID {
		s.logger.Warn("Tentative de suppression non autorisée d'un compte",
			logger.String("user_id", userID.String()),
			logger.String("account_id", accountID.String()),
		)
		return fmt.Errorf("accès non autorisé")
	}

	// Supprimer le compte
	if err := s.accountRepo.Delete(ctx, accountID); err != nil {
		s.logger.Error("Erreur suppression compte", logger.Error(err))
		return fmt.Errorf("erreur suppression compte: %w", err)
	}

	s.logger.Info("Compte supprimé avec succès",
		logger.String("account_id", accountID.String()),
		logger.String("user_id", userID.String()),
	)

	return nil
}

// GetAccountBalance récupère le solde d'un compte
func (s *AccountService) GetAccountBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (float64, error) {
	// Récupérer le compte
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		s.logger.Error("Erreur récupération compte pour solde", logger.Error(err))
		return 0, fmt.Errorf("compte non trouvé")
	}

	// Vérifier que le compte appartient à l'utilisateur
	if account.UserID != userID {
		s.logger.Warn("Tentative d'accès non autorisé au solde d'un compte",
			logger.String("user_id", userID.String()),
			logger.String("account_id", accountID.String()),
		)
		return 0, fmt.Errorf("accès non autorisé")
	}

	return account.Balance, nil
}

// GetAccountsByUserID récupère tous les comptes d'un utilisateur
func (s *AccountService) GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error) {
	accounts, err := s.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération comptes utilisateur", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération comptes: %w", err)
	}

	return accounts, nil
}

// check and update account balance
func (s *AccountService) CheckAndUpdateAccountBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*entity.Account, error) {

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		s.logger.Error("Erreur récupération compte", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération compte: %w", err)
	}

	transaction, err := s.transactionRepo.GetAllTransactionsByUserIDAndAccountID(ctx, userID, accountID)
	if err != nil {
		s.logger.Error("Erreur récupération transactions", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération transactions: %w", err)
	}
	if len(transaction) == 0 {
		s.logger.Info("Aucune transaction trouvée", logger.String("account_id", accountID.String()))
		return account, nil
	}

	savingGoal, err := s.savingGoalRepo.GetAllSavingGoalsByUserIDAndAccountID(ctx, userID, accountID)
	if err != nil {
		s.logger.Error("Erreur récupération transactions", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération transactions: %w", err)
	}
	// sum all transactions amount
	totalAmount := 0.0
	for _, t := range transaction {
		if t.Type == "income" {
			totalAmount += t.Amount
		} else {
			totalAmount -= t.Amount
		}
	}
	// sum all saving goals amount
	for _, v := range savingGoal {
		transactionBySavingGoal, err := s.transactionRepo.GetAllTransactionsBySavingGoalID(ctx, userID, v.ID)
		for _, t := range transactionBySavingGoal {
			totalAmount += t.Amount
		}
		if err != nil {
			s.logger.Error("Erreur récupération transactions par objectif d'épargne", logger.Error(err))
			return nil, fmt.Errorf("erreur récupération transactions par objectif d'épargne: %w", err)
		}
	}
	// update account balance
	account.Balance = totalAmount
	if err := s.accountRepo.Update(ctx, account); err != nil {
		s.logger.Error("Erreur mise à jour compte", logger.Error(err))
		return nil, fmt.Errorf("erreur mise à jour compte: %w", err)
	}

	return account, nil

}
