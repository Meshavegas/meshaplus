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

// InitializationService gère l'initialisation des données pour les nouveaux utilisateurs
type InitializationService struct {
	categoryRepo repository.CategoryRepository
	accountRepo  repository.AccountRepository
	budgetRepo   repository.BudgetRepository
	logger       logger.Logger
}

// NewInitializationService crée une nouvelle instance de InitializationService
func NewInitializationService(categoryRepo repository.CategoryRepository, accountRepo repository.AccountRepository, budgetRepo repository.BudgetRepository, logger logger.Logger) *InitializationService {
	return &InitializationService{
		categoryRepo: categoryRepo,
		accountRepo:  accountRepo,
		budgetRepo:   budgetRepo,
		logger:       logger,
	}
}

// InitializeUserData initialise les données par défaut pour un nouvel utilisateur
func (s *InitializationService) InitializeUserData(ctx context.Context, userID uuid.UUID) error {
	s.logger.Info("Initialisation des données pour le nouvel utilisateur", logger.String("user_id", userID.String()))

	// Créer les catégories par défaut
	if err := s.createDefaultCategories(ctx, userID); err != nil {
		s.logger.Error("Erreur création catégories par défaut", logger.Error(err))
		return err
	}

	// Créer les comptes par défaut
	if err := s.createDefaultAccounts(ctx, userID); err != nil {
		s.logger.Error("Erreur création comptes par défaut", logger.Error(err))
		return err
	}

	s.logger.Info("Initialisation terminée avec succès", logger.String("user_id", userID.String()))
	return nil
}

// CreateBudgetsFromPreferences crée des budgets basés sur les préférences utilisateur
func (s *InitializationService) CreateBudgetsFromPreferences(ctx context.Context, userID uuid.UUID, preferences *entity.UserPreferences) error {
	s.logger.Info("Création budgets automatiques basés sur les préférences", logger.String("user_id", userID.String()))

	// Récupérer les catégories de l'utilisateur pour mapper les budgets
	categories, err := s.categoryRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("erreur récupération catégories: %w", err)
	}

	// Créer une map pour rapidement trouver les catégories par nom
	categoryMap := make(map[string]*entity.Category)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat
	}

	// Obtenir le mois et l'année actuels
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	// Créer des budgets basés sur les préférences
	budgetsToCreate := []struct {
		categoryName string
		amount       float64
	}{
		{"Nourriture", preferences.Expenses.Food},
		{"Transport", preferences.Expenses.Transport},
		{"Logement", preferences.Expenses.Housing},
		{"Abonnements", preferences.Expenses.Subscriptions},
	}

	for _, budgetData := range budgetsToCreate {
		if budgetData.amount > 0 {
			category, exists := categoryMap[budgetData.categoryName]
			if !exists {
				s.logger.Warn("Catégorie non trouvée pour le budget",
					logger.String("category", budgetData.categoryName),
					logger.String("user_id", userID.String()))
				continue
			}

			budget := &entity.Budget{
				ID:            uuid.New(),
				UserID:        userID,
				CategoryID:    category.ID,
				AmountPlanned: budgetData.amount,
				AmountSpent:   0, // Nouveau budget, pas encore de dépenses
				Month:         currentMonth,
				Year:          currentYear,
				CreatedAt:     now,
				UpdatedAt:     now,
			}

			if err := s.budgetRepo.Create(ctx, budget); err != nil {
				s.logger.Error("Erreur création budget automatique",
					logger.Error(err),
					logger.String("category", budgetData.categoryName),
					logger.Float64("amount", budgetData.amount))
				continue
			}

			s.logger.Info("Budget automatique créé",
				logger.String("category", budgetData.categoryName),
				logger.Float64("amount", budgetData.amount),
				logger.String("user_id", userID.String()))
		}
	}

	return nil
}

// createDefaultAccounts crée les comptes par défaut pour un utilisateur
func (s *InitializationService) createDefaultAccounts(ctx context.Context, userID uuid.UUID) error {
	s.logger.Info("Début création des comptes par défaut", logger.String("user_id", userID.String()))

	var hasErrors bool
	var firstError error

	for i, defaultAccount := range entity.DefaultAccounts {
		s.logger.Info("Tentative de création du compte",
			logger.Int("index", i),
			logger.String("name", defaultAccount.Name),
			logger.String("type", defaultAccount.Type))

		account := &entity.Account{
			UserID:        userID,
			Name:          defaultAccount.Name,
			Type:          defaultAccount.Type,
			Currency:      defaultAccount.Currency,
			Icon:          defaultAccount.Icon,
			Color:         defaultAccount.Color,
			Balance:       defaultAccount.Balance,
			AccountNumber: defaultAccount.AccountNumber,
		}

		if err := s.accountRepo.Create(ctx, account); err != nil {
			s.logger.Error("Erreur création compte par défaut",
				logger.String("name", defaultAccount.Name),
				logger.String("type", defaultAccount.Type),
				logger.String("currency", defaultAccount.Currency),
				logger.String("icon", defaultAccount.Icon),
				logger.String("color", defaultAccount.Color),
				logger.Error(err))
			hasErrors = true
			if firstError == nil {
				firstError = err
			}
			continue // Continue avec les autres comptes
		}

		s.logger.Info("Compte par défaut créé avec succès",
			logger.String("name", defaultAccount.Name),
			logger.String("type", defaultAccount.Type),
			logger.String("currency", defaultAccount.Currency),
			logger.String("icon", defaultAccount.Icon),
			logger.String("color", defaultAccount.Color))
	}

	if hasErrors {
		s.logger.Error("Des erreurs sont survenues lors de la création des comptes par défaut",
			logger.String("user_id", userID.String()),
			logger.Error(firstError))
		return firstError
	}

	s.logger.Info("Création des comptes par défaut terminée avec succès", logger.String("user_id", userID.String()))
	return nil
}

// createDefaultCategories crée les catégories par défaut pour un utilisateur
func (s *InitializationService) createDefaultCategories(ctx context.Context, userID uuid.UUID) error {
	for _, defaultCategory := range entity.DefaultCategories {
		category := &entity.Category{
			UserID: userID,
			Name:   defaultCategory.Name,
			Type:   defaultCategory.Type,
			Icon:   defaultCategory.Icon,
			Color:  defaultCategory.Color,
		}

		if err := s.categoryRepo.Create(ctx, category); err != nil {
			s.logger.Error("Erreur création catégorie par défaut",
				logger.String("name", defaultCategory.Name),
				logger.String("type", defaultCategory.Type),
				logger.Error(err))
			return err
		}

		s.logger.Info("Catégorie par défaut créée",
			logger.String("name", defaultCategory.Name),
			logger.String("type", defaultCategory.Type),
			logger.String("icon", defaultCategory.Icon),
			logger.String("color", defaultCategory.Color))
	}

	return nil
}
