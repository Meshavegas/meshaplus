package service

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"backend/internal/service/ai"
	"backend/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// PreferencesService gère les préférences utilisateur
type PreferencesService struct {
	preferencesRepo repository.PreferencesRepository
	aiService       *ai.AIService
	logger          logger.Logger
}

// NewPreferencesService crée une nouvelle instance de PreferencesService
func NewPreferencesService(preferencesRepo repository.PreferencesRepository, aiService *ai.AIService, logger logger.Logger) *PreferencesService {
	return &PreferencesService{
		preferencesRepo: preferencesRepo,
		aiService:       aiService,
		logger:          logger,
	}
}

// CreatePreferences crée les préférences utilisateur
func (s *PreferencesService) CreatePreferences(ctx context.Context, userID uuid.UUID, req *entity.CreatePreferencesRequest) (*entity.UserPreferences, error) {
	preferences := &entity.UserPreferences{
		UserID:   userID,
		Income:   req.Income,
		Expenses: req.Expenses,
		Goals:    req.Goals,
		Habits:   req.Habits,
	}

	if err := s.preferencesRepo.Create(ctx, preferences); err != nil {
		s.logger.Error("Erreur création préférences", logger.Error(err))
		return nil, err
	}

	s.logger.Info("Préférences créées avec succès", logger.String("user_id", userID.String()))
	return preferences, nil
}

// GetPreferences récupère les préférences d'un utilisateur
func (s *PreferencesService) GetPreferences(ctx context.Context, userID uuid.UUID) (*entity.UserPreferences, error) {
	preferences, err := s.preferencesRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération préférences", logger.Error(err))
		return nil, err
	}

	return preferences, nil
}

// UpdatePreferences met à jour les préférences d'un utilisateur
func (s *PreferencesService) UpdatePreferences(ctx context.Context, userID uuid.UUID, req *entity.UpdatePreferencesRequest) (*entity.UserPreferences, error) {
	// Récupérer les préférences existantes
	preferences, err := s.preferencesRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération préférences existantes", logger.Error(err))
		return nil, err
	}

	if preferences == nil {
		// Créer de nouvelles préférences si elles n'existent pas
		if req.Income != nil && req.Expenses != nil && req.Goals != nil && req.Habits != nil {
			createReq := &entity.CreatePreferencesRequest{
				Income:   *req.Income,
				Expenses: *req.Expenses,
				Goals:    *req.Goals,
				Habits:   *req.Habits,
			}
			return s.CreatePreferences(ctx, userID, createReq)
		}
		return nil, fmt.Errorf("préférences non trouvées")
	}

	// Mettre à jour les champs fournis
	if req.Income != nil {
		preferences.Income = *req.Income
	}
	if req.Expenses != nil {
		preferences.Expenses = *req.Expenses
	}
	if req.Goals != nil {
		preferences.Goals = *req.Goals
	}
	if req.Habits != nil {
		preferences.Habits = *req.Habits
	}

	if err := s.preferencesRepo.Update(ctx, preferences); err != nil {
		s.logger.Error("Erreur mise à jour préférences", logger.Error(err))
		return nil, err
	}

	s.logger.Info("Préférences mises à jour avec succès", logger.String("user_id", userID.String()))
	return preferences, nil
}

// DeletePreferences supprime les préférences d'un utilisateur
func (s *PreferencesService) DeletePreferences(ctx context.Context, userID uuid.UUID) error {
	if err := s.preferencesRepo.Delete(ctx, userID); err != nil {
		s.logger.Error("Erreur suppression préférences", logger.Error(err))
		return err
	}

	s.logger.Info("Préférences supprimées avec succès", logger.String("user_id", userID.String()))
	return nil
}

// GeneratePersonalizedContent génère du contenu personnalisé basé sur les préférences
func (s *PreferencesService) GeneratePersonalizedContent(ctx context.Context, userID uuid.UUID, contentType string) (string, error) {
	preferences, err := s.preferencesRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération préférences pour génération de contenu", logger.Error(err))
		return "", err
	}

	if preferences == nil {
		return "", fmt.Errorf("préférences non trouvées")
	}

	// Générer du contenu personnalisé basé sur les préférences
	// Cette fonction peut être étendue pour utiliser l'IA
	content := s.generateContentFromPreferences(preferences, contentType)

	return content, nil
}

// generateContentFromPreferences génère du contenu basé sur les préférences
func (s *PreferencesService) generateContentFromPreferences(preferences *entity.UserPreferences, contentType string) string {
	switch contentType {
	case "budget_tip":
		return s.generateBudgetTip(preferences)
	case "savings_advice":
		return s.generateSavingsAdvice(preferences)
	case "habit_reminder":
		return s.generateHabitReminder(preferences)
	default:
		return "Contenu personnalisé basé sur vos préférences."
	}
}

// generateBudgetTip génère un conseil de budget personnalisé
func (s *PreferencesService) generateBudgetTip(preferences *entity.UserPreferences) string {
	if preferences.Expenses.AutoBudget {
		return "Votre budget automatique 50/30/20 est activé. Assurez-vous de respecter vos limites mensuelles."
	}
	return "Considérez activer le budget automatique pour mieux gérer vos dépenses."
}

// generateSavingsAdvice génère un conseil d'épargne personnalisé
func (s *PreferencesService) generateSavingsAdvice(preferences *entity.UserPreferences) string {
	if preferences.Goals.SavingsTarget > 0 {
		return fmt.Sprintf("Objectif d'épargne : %.0f XAF par mois. Vous êtes sur la bonne voie !", preferences.Goals.SavingsTarget)
	}
	return "Définissez un objectif d'épargne pour commencer à épargner régulièrement."
}

// generateHabitReminder génère un rappel d'habitude personnalisé
func (s *PreferencesService) generateHabitReminder(preferences *entity.UserPreferences) string {
	if preferences.Habits.CustomHabit != "" {
		return fmt.Sprintf("N'oubliez pas votre habitude : %s", preferences.Habits.CustomHabit)
	}
	return "Prenez le temps de planifier votre journée selon vos préférences."
}
