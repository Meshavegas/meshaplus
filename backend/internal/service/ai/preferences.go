package ai

import (
	"backend/internal/domaine/entity"
	"backend/pkg/logger"
	"context"
	"fmt"

	"google.golang.org/genai"
)

// PreferencesAIService gère l'analyse et la génération de contenu basé sur les préférences utilisateur
type PreferencesAIService struct {
	logger logger.Logger
}

// NewPreferencesAIService crée une nouvelle instance de PreferencesAIService
func NewPreferencesAIService(apiKey string, logger logger.Logger) (*PreferencesAIService, error) {
	return &PreferencesAIService{
		logger: logger,
	}, nil
}

// GenerateDefaultPreferences génère des préférences par défaut basées sur les informations utilisateur
func (s *PreferencesAIService) GenerateDefaultPreferences(ctx context.Context, userInfo *entity.User) (*entity.CreatePreferencesRequest, error) {
	// Pour l'instant, on utilise des préférences statiques par défaut
	// TODO: Réactiver l'AI quand l'API key sera configurée
	s.logger.Info("Génération de préférences par défaut statiques",
		logger.String("user_id", userInfo.ID.String()),
		logger.String("user_name", userInfo.Name),
	)

	return &entity.CreatePreferencesRequest{
		Income: entity.IncomePreferences{
			Sources:      []string{"Salaire"},
			MonthlyTotal: 150000,
			Accounts:     []string{"Compte Bancaire", "Mobile Money"},
			HasDebt:      false,
			DebtAmount:   0,
		},
		Expenses: entity.ExpensePreferences{
			TopCategories: []string{"Nourriture", "Transport", "Logement"},
			Food:          60000,
			Transport:     20000,
			Housing:       80000,
			Subscriptions: 15000,
			AlertsEnabled: true,
			AutoBudget:    true, // Important pour générer les budgets automatiquement
		},
		Goals: entity.GoalPreferences{
			MainGoal:      "Épargne",
			SecondaryGoal: "Investissement",
			SavingsTarget: 50000,
			Deadline:      "6 mois",
			AdviceEnabled: true,
		},
		Habits: entity.HabitPreferences{
			PlanningTime:   "Matin",
			DailyFocusTime: "30min",
			CustomHabit:    "Vérification quotidienne des dépenses",
			SummaryType:    "Hebdomadaire",
		},
	}, nil
}

// GeneratePersonalizedAdvice génère des conseils personnalisés basés sur les préférences
func (s *PreferencesAIService) GeneratePersonalizedAdvice(ctx context.Context, preferences *entity.UserPreferences, adviceType string) (string, error) {
	var prompt string

	switch adviceType {
	case "budget":
		prompt = fmt.Sprintf(`
Basé sur ces préférences de dépenses :
- Catégories principales: %v
- Nourriture: %.0f XAF
- Transport: %.0f XAF
- Logement: %.0f XAF
- Abonnements: %.0f XAF

Génère un conseil budgétaire personnalisé et pratique en français (max 150 mots).`,
			preferences.Expenses.TopCategories,
			preferences.Expenses.Food,
			preferences.Expenses.Transport,
			preferences.Expenses.Housing,
			preferences.Expenses.Subscriptions)

	case "savings":
		prompt = fmt.Sprintf(`
Basé sur ces objectifs :
- Objectif principal: %s
- Objectif secondaire: %s
- Cible d'épargne: %.0f XAF
- Échéance: %s

Génère un conseil d'épargne personnalisé et motivant en français (max 150 mots).`,
			preferences.Goals.MainGoal,
			preferences.Goals.SecondaryGoal,
			preferences.Goals.SavingsTarget,
			preferences.Goals.Deadline)

	case "habits":
		prompt = fmt.Sprintf(`
Basé sur ces habitudes :
- Temps de planification: %s
- Temps de focus quotidien: %s
- Habitude personnalisée: %s

Génère un conseil pour améliorer les habitudes financières en français (max 150 mots).`,
			preferences.Habits.PlanningTime,
			preferences.Habits.DailyFocusTime,
			preferences.Habits.CustomHabit)

	default:
		return "", fmt.Errorf("type de conseil non supporté: %s", adviceType)
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		s.logger.Error("Failed to create Gemini client: " + err.Error())
		return "", err
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		s.logger.Error("Failed to generate content: " + err.Error())
		return "", err
	}

	advice := result.Text()

	s.logger.Info("Conseil personnalisé généré",
		logger.String("type", adviceType),
		logger.String("user_id", preferences.UserID.String()),
	)

	return advice, nil
}

// Close ferme la connexion au client AI (plus nécessaire avec la nouvelle API)
func (s *PreferencesAIService) Close() error {
	return nil
}
