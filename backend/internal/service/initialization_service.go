package service

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"backend/pkg/logger"
	"context"

	"github.com/google/uuid"
)

// InitializationService gère l'initialisation des données pour les nouveaux utilisateurs
type InitializationService struct {
	categoryRepo repository.CategoryRepository
	accountRepo  repository.AccountRepository
	logger       logger.Logger
}

// NewInitializationService crée une nouvelle instance de InitializationService
func NewInitializationService(categoryRepo repository.CategoryRepository, accountRepo repository.AccountRepository, logger logger.Logger) *InitializationService {
	return &InitializationService{
		categoryRepo: categoryRepo,
		accountRepo:  accountRepo,
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
