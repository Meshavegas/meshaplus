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

// SavingGoalService gère la logique métier des objectifs d'épargne
type SavingGoalService struct {
	savingGoalRepo repository.SavingGoalRepository
	logger         logger.Logger
}

// NewSavingGoalService crée une nouvelle instance de SavingGoalService
func NewSavingGoalService(savingGoalRepo repository.SavingGoalRepository, logger logger.Logger) *SavingGoalService {
	return &SavingGoalService{
		savingGoalRepo: savingGoalRepo,
		logger:         logger,
	}
}

// CreateSavingGoal crée un nouvel objectif d'épargne
func (s *SavingGoalService) CreateSavingGoal(ctx context.Context, userID uuid.UUID, req entity.CreateSavingGoalRequest) (*entity.SavingGoal, error) {
	// Validation des données
	if req.Title == "" {
		return nil, fmt.Errorf("le titre est requis")
	}

	if req.TargetAmount <= 0 {
		return nil, fmt.Errorf("le montant cible doit être positif")
	}

	// Validation de la fréquence si fournie
	if req.Frequency != nil {
		if *req.Frequency != "weekly" && *req.Frequency != "monthly" && *req.Frequency != "yearly" {
			return nil, fmt.Errorf("la fréquence doit être 'weekly', 'monthly' ou 'yearly'")
		}
	}

	// Création de l'objectif d'épargne
	savingGoal := &entity.SavingGoal{
		UserID:        userID,
		Title:         req.Title,
		TargetAmount:  req.TargetAmount,
		CurrentAmount: 0, // Initialiser à 0
		Deadline:      req.Deadline,
		IsAchieved:    false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Gérer la fréquence si fournie
	if req.Frequency != nil {
		savingGoal.Frequency = *req.Frequency
	}

	if err := s.savingGoalRepo.Create(ctx, savingGoal); err != nil {
		s.logger.Error("Erreur création objectif d'épargne", logger.Error(err))
		return nil, fmt.Errorf("erreur création objectif d'épargne: %w", err)
	}

	s.logger.Info("Objectif d'épargne créé avec succès",
		logger.String("saving_goal_id", savingGoal.ID.String()),
		logger.String("user_id", userID.String()),
		logger.String("title", savingGoal.Title),
		logger.Float64("target_amount", savingGoal.TargetAmount),
	)

	return savingGoal, nil
}

// GetSavingGoal récupère un objectif d'épargne par son ID
func (s *SavingGoalService) GetSavingGoal(ctx context.Context, userID uuid.UUID, goalID uuid.UUID) (*entity.SavingGoal, error) {
	savingGoal, err := s.savingGoalRepo.GetByID(ctx, goalID)
	if err != nil {
		s.logger.Error("Erreur récupération objectif d'épargne", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération objectif d'épargne: %w", err)
	}

	// Vérifier que l'objectif appartient à l'utilisateur
	if savingGoal.UserID != userID {
		s.logger.Warn("Tentative d'accès non autorisé à un objectif d'épargne",
			logger.String("user_id", userID.String()),
			logger.String("goal_id", goalID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	return savingGoal, nil
}

// GetSavingGoals récupère tous les objectifs d'épargne d'un utilisateur
func (s *SavingGoalService) GetSavingGoals(ctx context.Context, userID uuid.UUID, page, limit int) ([]*entity.SavingGoal, int64, error) {
	// Pour l'instant, récupérer tous les objectifs de l'utilisateur
	// TODO: Implémenter la pagination quand le repository sera créé
	savingGoals, err := s.savingGoalRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération objectifs d'épargne", logger.Error(err))
		return nil, 0, fmt.Errorf("erreur récupération objectifs d'épargne: %w", err)
	}

	// Calculer le total
	total := int64(len(savingGoals))

	// Appliquer la pagination basique
	start := (page - 1) * limit
	end := start + limit
	if int64(start) >= total {
		return []*entity.SavingGoal{}, total, nil
	}
	if int64(end) > total {
		end = int(total)
	}

	return savingGoals[start:end], total, nil
}

// UpdateSavingGoal met à jour un objectif d'épargne
func (s *SavingGoalService) UpdateSavingGoal(ctx context.Context, userID uuid.UUID, goalID uuid.UUID, req entity.UpdateSavingGoalRequest) (*entity.SavingGoal, error) {
	// Récupérer l'objectif existant
	savingGoal, err := s.savingGoalRepo.GetByID(ctx, goalID)
	if err != nil {
		s.logger.Error("Erreur récupération objectif d'épargne pour mise à jour", logger.Error(err))
		return nil, fmt.Errorf("objectif d'épargne non trouvé")
	}

	// Vérifier que l'objectif appartient à l'utilisateur
	if savingGoal.UserID != userID {
		s.logger.Warn("Tentative de mise à jour non autorisée d'un objectif d'épargne",
			logger.String("user_id", userID.String()),
			logger.String("goal_id", goalID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	// Mettre à jour les champs
	if req.Title != nil {
		savingGoal.Title = *req.Title
	}

	if req.TargetAmount != nil {
		if *req.TargetAmount <= 0 {
			return nil, fmt.Errorf("le montant cible doit être positif")
		}
		savingGoal.TargetAmount = *req.TargetAmount
	}

	if req.CurrentAmount != nil {
		if *req.CurrentAmount < 0 {
			return nil, fmt.Errorf("le montant actuel ne peut pas être négatif")
		}
		savingGoal.CurrentAmount = *req.CurrentAmount
	}

	if req.Deadline != nil {
		savingGoal.Deadline = req.Deadline
	}

	if req.IsAchieved != nil {
		savingGoal.IsAchieved = *req.IsAchieved
	}

	if req.Frequency != nil {
		if *req.Frequency != "weekly" && *req.Frequency != "monthly" && *req.Frequency != "yearly" {
			return nil, fmt.Errorf("la fréquence doit être 'weekly', 'monthly' ou 'yearly'")
		}
		savingGoal.Frequency = *req.Frequency
	}

	// Vérifier si l'objectif est atteint
	if savingGoal.CurrentAmount >= savingGoal.TargetAmount {
		savingGoal.IsAchieved = true
	}

	savingGoal.UpdatedAt = time.Now()

	// Sauvegarder les modifications
	if err := s.savingGoalRepo.Update(ctx, savingGoal); err != nil {
		s.logger.Error("Erreur mise à jour objectif d'épargne", logger.Error(err))
		return nil, fmt.Errorf("erreur mise à jour objectif d'épargne: %w", err)
	}

	s.logger.Info("Objectif d'épargne mis à jour avec succès",
		logger.String("goal_id", savingGoal.ID.String()),
		logger.String("user_id", userID.String()),
	)

	return savingGoal, nil
}

// DeleteSavingGoal supprime un objectif d'épargne
func (s *SavingGoalService) DeleteSavingGoal(ctx context.Context, userID uuid.UUID, goalID uuid.UUID) error {
	// Récupérer l'objectif existant
	savingGoal, err := s.savingGoalRepo.GetByID(ctx, goalID)
	if err != nil {
		s.logger.Error("Erreur récupération objectif d'épargne pour suppression", logger.Error(err))
		return fmt.Errorf("objectif d'épargne non trouvé")
	}

	// Vérifier que l'objectif appartient à l'utilisateur
	if savingGoal.UserID != userID {
		s.logger.Warn("Tentative de suppression non autorisée d'un objectif d'épargne",
			logger.String("user_id", userID.String()),
			logger.String("goal_id", goalID.String()),
		)
		return fmt.Errorf("accès non autorisé")
	}

	// Supprimer l'objectif
	if err := s.savingGoalRepo.Delete(ctx, goalID); err != nil {
		s.logger.Error("Erreur suppression objectif d'épargne", logger.Error(err))
		return fmt.Errorf("erreur suppression objectif d'épargne: %w", err)
	}

	s.logger.Info("Objectif d'épargne supprimé avec succès",
		logger.String("goal_id", goalID.String()),
		logger.String("user_id", userID.String()),
	)

	return nil
}
