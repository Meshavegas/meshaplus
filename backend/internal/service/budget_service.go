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

// BudgetService gère la logique métier des budgets
type BudgetService struct {
	budgetRepo repository.BudgetRepository
	logger     logger.Logger
}

// NewBudgetService crée une nouvelle instance de BudgetService
func NewBudgetService(budgetRepo repository.BudgetRepository, logger logger.Logger) *BudgetService {
	return &BudgetService{
		budgetRepo: budgetRepo,
		logger:     logger,
	}
}

// CreateBudget crée un nouveau budget
func (s *BudgetService) CreateBudget(ctx context.Context, userID uuid.UUID, req entity.CreateBudgetRequest) (*entity.Budget, error) {
	// Validation des données
	if req.AmountPlanned <= 0 {
		return nil, fmt.Errorf("le montant planifié doit être positif")
	}

	if req.Month < 1 || req.Month > 12 {
		return nil, fmt.Errorf("le mois doit être entre 1 et 12")
	}

	if req.Year < 2020 {
		return nil, fmt.Errorf("l'année doit être 2020 ou plus")
	}

	// TODO: Vérifier qu'il n'y a pas déjà un budget pour cette catégorie et cette période
	// Cette vérification sera implémentée quand le repository sera créé

	// Création du budget
	budget := &entity.Budget{
		UserID:        userID,
		CategoryID:    req.CategoryID,
		AmountPlanned: req.AmountPlanned,
		AmountSpent:   0, // Initialiser à 0
		Month:         req.Month,
		Year:          req.Year,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.budgetRepo.Create(ctx, budget); err != nil {
		s.logger.Error("Erreur création budget", logger.Error(err))
		return nil, fmt.Errorf("erreur création budget: %w", err)
	}

	s.logger.Info("Budget créé avec succès",
		logger.String("budget_id", budget.ID.String()),
		logger.String("user_id", userID.String()),
		logger.String("category_id", req.CategoryID.String()),
		logger.Int("month", req.Month),
		logger.Int("year", req.Year),
	)

	return budget, nil
}

// GetBudget récupère un budget par son ID
func (s *BudgetService) GetBudget(ctx context.Context, userID uuid.UUID, budgetID uuid.UUID) (*entity.Budget, error) {
	budget, err := s.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		s.logger.Error("Erreur récupération budget", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération budget: %w", err)
	}

	// Vérifier que le budget appartient à l'utilisateur
	if budget.UserID != userID {
		s.logger.Warn("Tentative d'accès non autorisé à un budget",
			logger.String("user_id", userID.String()),
			logger.String("budget_id", budgetID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	return budget, nil
}

// GetBudgets récupère tous les budgets d'un utilisateur pour une période donnée
func (s *BudgetService) GetBudgets(ctx context.Context, userID uuid.UUID, month, year int, page, limit int) ([]*entity.Budget, int64, error) {
	// Pour l'instant, récupérer tous les budgets de l'utilisateur
	// TODO: Implémenter le filtrage par période et la pagination quand le repository sera créé
	budgets, err := s.budgetRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération budgets", logger.Error(err))
		return nil, 0, fmt.Errorf("erreur récupération budgets: %w", err)
	}

	// Filtrer par période
	var filteredBudgets []*entity.Budget
	for _, budget := range budgets {
		if budget.Month == month && budget.Year == year {
			filteredBudgets = append(filteredBudgets, budget)
		}
	}

	// Calculer le total
	total := int64(len(filteredBudgets))

	// Appliquer la pagination basique
	start := (page - 1) * limit
	end := start + limit
	if int64(start) >= total {
		return []*entity.Budget{}, total, nil
	}
	if int64(end) > total {
		end = int(total)
	}

	return filteredBudgets[start:end], total, nil
}

// UpdateBudget met à jour un budget
func (s *BudgetService) UpdateBudget(ctx context.Context, userID uuid.UUID, budgetID uuid.UUID, req entity.UpdateBudgetRequest) (*entity.Budget, error) {
	// Récupérer le budget existant
	budget, err := s.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		s.logger.Error("Erreur récupération budget pour mise à jour", logger.Error(err))
		return nil, fmt.Errorf("budget non trouvé")
	}

	// Vérifier que le budget appartient à l'utilisateur
	if budget.UserID != userID {
		s.logger.Warn("Tentative de mise à jour non autorisée d'un budget",
			logger.String("user_id", userID.String()),
			logger.String("budget_id", budgetID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	// Mettre à jour les champs
	if req.AmountPlanned != nil {
		if *req.AmountPlanned <= 0 {
			return nil, fmt.Errorf("le montant planifié doit être positif")
		}
		budget.AmountPlanned = *req.AmountPlanned
	}

	if req.AmountSpent != nil {
		if *req.AmountSpent < 0 {
			return nil, fmt.Errorf("le montant dépensé ne peut pas être négatif")
		}
		budget.AmountSpent = *req.AmountSpent
	}

	budget.UpdatedAt = time.Now()

	// Sauvegarder les modifications
	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		s.logger.Error("Erreur mise à jour budget", logger.Error(err))
		return nil, fmt.Errorf("erreur mise à jour budget: %w", err)
	}

	s.logger.Info("Budget mis à jour avec succès",
		logger.String("budget_id", budget.ID.String()),
		logger.String("user_id", userID.String()),
	)

	return budget, nil
}

// DeleteBudget supprime un budget
func (s *BudgetService) DeleteBudget(ctx context.Context, userID uuid.UUID, budgetID uuid.UUID) error {
	// Récupérer le budget existant
	budget, err := s.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		s.logger.Error("Erreur récupération budget pour suppression", logger.Error(err))
		return fmt.Errorf("budget non trouvé")
	}

	// Vérifier que le budget appartient à l'utilisateur
	if budget.UserID != userID {
		s.logger.Warn("Tentative de suppression non autorisée d'un budget",
			logger.String("user_id", userID.String()),
			logger.String("budget_id", budgetID.String()),
		)
		return fmt.Errorf("accès non autorisé")
	}

	// Supprimer le budget
	if err := s.budgetRepo.Delete(ctx, budgetID); err != nil {
		s.logger.Error("Erreur suppression budget", logger.Error(err))
		return fmt.Errorf("erreur suppression budget: %w", err)
	}

	s.logger.Info("Budget supprimé avec succès",
		logger.String("budget_id", budgetID.String()),
		logger.String("user_id", userID.String()),
	)

	return nil
}

// GetBudgetStats récupère les statistiques des budgets pour une période donnée
func (s *BudgetService) GetBudgetStats(ctx context.Context, userID uuid.UUID, month, year int) (map[string]interface{}, error) {
	// Récupérer tous les budgets de l'utilisateur pour la période
	budgets, err := s.budgetRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération statistiques budgets", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération statistiques: %w", err)
	}

	// Filtrer par période et calculer les statistiques
	var totalPlanned, totalSpent float64
	var budgetCount int
	var overBudgetCount int

	for _, budget := range budgets {
		if budget.Month == month && budget.Year == year {
			totalPlanned += budget.AmountPlanned
			totalSpent += budget.AmountSpent
			budgetCount++

			if budget.AmountSpent > budget.AmountPlanned {
				overBudgetCount++
			}
		}
	}

	stats := map[string]interface{}{
		"period":            fmt.Sprintf("%d/%d", month, year),
		"total_planned":     totalPlanned,
		"total_spent":       totalSpent,
		"remaining":         totalPlanned - totalSpent,
		"budget_count":      budgetCount,
		"over_budget_count": overBudgetCount,
		"utilization_rate":  0.0,
	}

	if totalPlanned > 0 {
		stats["utilization_rate"] = (totalSpent / totalPlanned) * 100
	}

	return stats, nil
}
