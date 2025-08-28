package ai

import (
	"backend/internal/domaine/entity"
	"backend/pkg/logger"
	"context"
	"encoding/json"
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
	// Essayer d'abord avec l'AI
	aiPreferences, err := s.generatePreferencesWithAI(ctx, userInfo)
	if err != nil {
		s.logger.Warn("Échec génération AI, utilisation du fallback statique",
			logger.String("user_id", userInfo.ID.String()),
			logger.Error(err))
		return s.generateStaticPreferences(userInfo), nil
	}

	s.logger.Info("Préférences par défaut générées par AI",
		logger.String("user_id", userInfo.ID.String()),
		logger.String("user_name", userInfo.Name),
	)

	return aiPreferences, nil
}

// generatePreferencesWithAI génère des préférences par défaut avec l'API Google Generative AI
func (s *PreferencesAIService) generatePreferencesWithAI(ctx context.Context, userInfo *entity.User) (*entity.CreatePreferencesRequest, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erreur création client AI: %w", err)
	}

	prompt := fmt.Sprintf(`Tu es un expert en finances personnelles qui génère des préférences par défaut personnalisées.

INFORMATIONS UTILISATEUR :
- Nom: %s
- Email: %s

CONTEXTE CAMEROUNAIS :
- Devise: XAF (Francs CFA)
- Salaire moyen: 150,000-300,000 XAF/mois
- Services populaires: Mobile Money (MOMO, Orange Money), Comptes bancaires
- Coût de vie: Nourriture (40-60k), Transport (15-25k), Logement (50-100k), Abonnements (10-20k)

INSTRUCTIONS DE GÉNÉRATION :

💰 REVENUS :
- Sources: 1-2 sources réalistes (Salaire, Business, Freelance, etc.)
- Montant: Entre 150,000 et 300,000 XAF selon le profil
- Comptes: 2-3 comptes (Bancaire, Mobile Money, Cash)
- Dettes: Probabilité faible (20%%), montant réaliste si applicable

💸 DÉPENSES :
- Catégories: 3-4 catégories principales réalistes
- Nourriture: 40-60%% du revenu
- Transport: 10-20%% du revenu  
- Logement: 25-40%% du revenu
- Abonnements: 5-10%% du revenu
- Total: 80-90%% du revenu (laisser marge pour épargne)

🎯 OBJECTIFS :
- Principal: Épargne, Investissement, Voyage, ou autre objectif courant
- Secondaire: Complémentaire au principal
- Cible: 10-30%% du revenu mensuel
- Échéance: 3-12 mois

⏰ HABITUDES :
- Planification: Matin, Midi, ou Soir
- Focus: 15min, 30min, 1h, ou +1h
- Habitude: Spécifique et réalisable
- Résumé: Quotidien, Hebdomadaire, ou Aucun

Génère des préférences cohérentes et équilibrées pour un utilisateur camerounais.
Les montants doivent être réalistes et les proportions logiques.`, userInfo.Name, userInfo.Email)

	// Configuration pour la sortie structurée JSON
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"income": {
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"sources": {
							Type:  genai.TypeArray,
							Items: &genai.Schema{Type: genai.TypeString},
						},
						"monthly_total": {
							Type: genai.TypeInteger,
						},
						"accounts": {
							Type:  genai.TypeArray,
							Items: &genai.Schema{Type: genai.TypeString},
						},
						"has_debt": {
							Type: genai.TypeBoolean,
						},
						"debt_amount": {
							Type: genai.TypeInteger,
						},
					},
				},
				"expenses": {
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"top_categories": {
							Type:  genai.TypeArray,
							Items: &genai.Schema{Type: genai.TypeString},
						},
						"food": {
							Type: genai.TypeInteger,
						},
						"transport": {
							Type: genai.TypeInteger,
						},
						"housing": {
							Type: genai.TypeInteger,
						},
						"subscriptions": {
							Type: genai.TypeInteger,
						},
						"alerts_enabled": {
							Type: genai.TypeBoolean,
						},
						"auto_budget": {
							Type: genai.TypeBoolean,
						},
					},
				},
				"goals": {
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"main_goal": {
							Type: genai.TypeString,
						},
						"secondary_goal": {
							Type: genai.TypeString,
						},
						"savings_target": {
							Type: genai.TypeInteger,
						},
						"deadline": {
							Type: genai.TypeString,
						},
						"advice_enabled": {
							Type: genai.TypeBoolean,
						},
					},
				},
				"habits": {
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"planning_time": {
							Type: genai.TypeString,
						},
						"daily_focus_time": {
							Type: genai.TypeString,
						},
						"custom_habit": {
							Type: genai.TypeString,
						},
						"summary_type": {
							Type: genai.TypeString,
						},
					},
				},
			},
		},
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return nil, fmt.Errorf("erreur génération AI: %w", err)
	}

	// Parser la réponse JSON
	var preferencesResponse entity.CreatePreferencesRequest
	if err := json.Unmarshal([]byte(result.Text()), &preferencesResponse); err != nil {
		s.logger.Error("Failed to parse JSON response: " + err.Error())
		return nil, err
	}

	return &preferencesResponse, nil
}

// generateStaticPreferences génère des préférences statiques en fallback
func (s *PreferencesAIService) generateStaticPreferences(userInfo *entity.User) *entity.CreatePreferencesRequest {
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
	}
}

// GeneratePersonalizedAdvice génère des conseils personnalisés basés sur les préférences
func (s *PreferencesAIService) GeneratePersonalizedAdvice(ctx context.Context, preferences *entity.UserPreferences, adviceType string) (string, error) {
	// Essayer d'abord avec l'AI
	aiAdvice, err := s.generateAdviceWithAI(ctx, preferences, adviceType)
	if err != nil {
		s.logger.Warn("Échec génération AI, utilisation du fallback statique",
			logger.String("type", adviceType),
			logger.Error(err))
		return s.generateStaticAdvice(preferences, adviceType), nil
	}

	s.logger.Info("Conseil personnalisé généré par AI",
		logger.String("type", adviceType),
		logger.String("user_id", preferences.UserID.String()),
	)

	return aiAdvice, nil
}

// generateAdviceWithAI génère un conseil avec l'API Google Generative AI
func (s *PreferencesAIService) generateAdviceWithAI(ctx context.Context, preferences *entity.UserPreferences, adviceType string) (string, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("erreur création client AI: %w", err)
	}

	var prompt string

	switch adviceType {
	case "budget":
		prompt = fmt.Sprintf(`
Tu es un expert financier qui donne des conseils personnalisés et pratiques.

PROFIL FINANCIER COMPLET DE L'UTILISATEUR :

💰 REVENUS :
- Sources: %v
- Total mensuel: %.0f XAF
- Comptes: %v
- Dettes: %s (%.0f XAF)

💸 DÉPENSES ACTUELLES :
- Catégories principales: %v
- Nourriture: %.0f XAF (%.1f%% du revenu)
- Transport: %.0f XAF (%.1f%% du revenu)
- Logement: %.0f XAF (%.1f%% du revenu)
- Abonnements: %.0f XAF (%.1f%% du revenu)
- Total dépenses: %.0f XAF (%.1f%% du revenu)

🎯 OBJECTIFS :
- Principal: %s
- Secondaire: %s
- Cible d'épargne: %.0f XAF
- Échéance: %s

⏰ HABITUDES :
- Planification: %s
- Focus quotidien: %s
- Habitude personnalisée: %s

Génère un conseil budgétaire personnalisé et pratique en français (max 250 mots).
Analyse la situation actuelle et propose des améliorations concrètes.
Utilise des emojis et sois motivant.`,
			preferences.Income.Sources,
			preferences.Income.MonthlyTotal,
			preferences.Income.Accounts,
			func() string {
				if preferences.Income.HasDebt {
					return "Oui"
				} else {
					return "Non"
				}
			}(),
			preferences.Income.DebtAmount,
			preferences.Expenses.TopCategories,
			preferences.Expenses.Food,
			(preferences.Expenses.Food/preferences.Income.MonthlyTotal)*100,
			preferences.Expenses.Transport,
			(preferences.Expenses.Transport/preferences.Income.MonthlyTotal)*100,
			preferences.Expenses.Housing,
			(preferences.Expenses.Housing/preferences.Income.MonthlyTotal)*100,
			preferences.Expenses.Subscriptions,
			(preferences.Expenses.Subscriptions/preferences.Income.MonthlyTotal)*100,
			preferences.Expenses.Food+preferences.Expenses.Transport+preferences.Expenses.Housing+preferences.Expenses.Subscriptions,
			((preferences.Expenses.Food+preferences.Expenses.Transport+preferences.Expenses.Housing+preferences.Expenses.Subscriptions)/preferences.Income.MonthlyTotal)*100,
			preferences.Goals.MainGoal,
			preferences.Goals.SecondaryGoal,
			preferences.Goals.SavingsTarget,
			preferences.Goals.Deadline,
			preferences.Habits.PlanningTime,
			preferences.Habits.DailyFocusTime,
			preferences.Habits.CustomHabit)

	case "savings":
		prompt = fmt.Sprintf(`
Tu es un expert financier qui donne des conseils d'épargne personnalisés.

PROFIL FINANCIER COMPLET DE L'UTILISATEUR :

💰 REVENUS :
- Sources: %v
- Total mensuel: %.0f XAF
- Comptes: %v

💸 DÉPENSES ACTUELLES :
- Total dépenses: %.0f XAF
- Marge disponible: %.0f XAF (%.1f%% du revenu)

🎯 OBJECTIFS ACTUELS :
- Principal: %s
- Secondaire: %s
- Cible d'épargne: %.0f XAF
- Échéance: %s
- Épargne mensuelle nécessaire: %.0f XAF

⏰ HABITUDES :
- Planification: %s
- Focus quotidien: %s
- Habitude personnalisée: %s

Génère un conseil d'épargne personnalisé et motivant en français (max 250 mots).
Analyse la faisabilité de l'objectif et propose des stratégies concrètes.
Utilise des emojis et sois encourageant.`,
			preferences.Income.Sources,
			preferences.Income.MonthlyTotal,
			preferences.Income.Accounts,
			preferences.Expenses.Food+preferences.Expenses.Transport+preferences.Expenses.Housing+preferences.Expenses.Subscriptions,
			preferences.Income.MonthlyTotal-(preferences.Expenses.Food+preferences.Expenses.Transport+preferences.Expenses.Housing+preferences.Expenses.Subscriptions),
			((preferences.Income.MonthlyTotal-(preferences.Expenses.Food+preferences.Expenses.Transport+preferences.Expenses.Housing+preferences.Expenses.Subscriptions))/preferences.Income.MonthlyTotal)*100,
			preferences.Goals.MainGoal,
			preferences.Goals.SecondaryGoal,
			preferences.Goals.SavingsTarget,
			preferences.Goals.Deadline,
			preferences.Goals.SavingsTarget/6, // Estimation mensuelle
			preferences.Habits.PlanningTime,
			preferences.Habits.DailyFocusTime,
			preferences.Habits.CustomHabit)

	case "habits":
		prompt = fmt.Sprintf(`
Tu es un expert en habitudes financières qui donne des conseils personnalisés.

PROFIL COMPLET DE L'UTILISATEUR :

💰 SITUATION FINANCIÈRE :
- Revenus: %.0f XAF/mois
- Dépenses: %.0f XAF/mois
- Marge: %.0f XAF/mois
- Objectif: %s (%.0f XAF en %s)

⏰ HABITUDES ACTUELLES :
- Temps de planification: %s
- Temps de focus quotidien: %s
- Habitude personnalisée: %s
- Type de résumé: %s

🎯 OBJECTIFS :
- Principal: %s
- Secondaire: %s

Génère un conseil pour améliorer les habitudes financières en français (max 250 mots).
Propose des améliorations concrètes basées sur le profil actuel.
Utilise des emojis et sois motivant.`,
			preferences.Income.MonthlyTotal,
			preferences.Expenses.Food+preferences.Expenses.Transport+preferences.Expenses.Housing+preferences.Expenses.Subscriptions,
			preferences.Income.MonthlyTotal-(preferences.Expenses.Food+preferences.Expenses.Transport+preferences.Expenses.Housing+preferences.Expenses.Subscriptions),
			preferences.Goals.MainGoal,
			preferences.Goals.SavingsTarget,
			preferences.Goals.Deadline,
			preferences.Habits.PlanningTime,
			preferences.Habits.DailyFocusTime,
			preferences.Habits.CustomHabit,
			preferences.Habits.SummaryType,
			preferences.Goals.MainGoal,
			preferences.Goals.SecondaryGoal)

	default:
		return "", fmt.Errorf("type de conseil non supporté: %s", adviceType)
	}

	// Configuration pour la génération de contenu
	temp := float32(0.7)
	topK := float32(40)
	topP := float32(0.8)
	maxTokens := int32(1000)

	config := &genai.GenerateContentConfig{
		Temperature:     &temp, // Créativité modérée
		TopK:            &topK,
		TopP:            &topP,
		MaxOutputTokens: maxTokens, // Limite pour éviter des réponses trop longues
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		return "", fmt.Errorf("erreur génération AI: %w", err)
	}

	return result.Text(), nil
}

// generateStaticAdvice génère un conseil statique en fallback
func (s *PreferencesAIService) generateStaticAdvice(preferences *entity.UserPreferences, adviceType string) string {
	switch adviceType {
	case "budget":
		return s.generateBudgetAdvice(preferences)
	case "savings":
		return s.generateSavingsAdvice(preferences)
	case "habits":
		return s.generateHabitsAdvice(preferences)
	default:
		return "Conseil personnalisé basé sur vos préférences."
	}
}

// generateBudgetAdvice génère un conseil budgétaire personnalisé
func (s *PreferencesAIService) generateBudgetAdvice(preferences *entity.UserPreferences) string {
	totalExpenses := preferences.Expenses.Food + preferences.Expenses.Transport + preferences.Expenses.Housing + preferences.Expenses.Subscriptions

	return fmt.Sprintf(`Conseil budgétaire personnalisé :

Basé sur vos dépenses mensuelles de %.0f XAF, voici mes recommandations :

🍽️ Nourriture (%.0f XAF) : Planifiez vos repas à l'avance pour éviter les dépenses imprévues. Privilégiez les achats en gros pour les produits de base.

🚗 Transport (%.0f XAF) : Considérez le covoiturage ou les transports en commun pour réduire vos frais de transport.

🏠 Logement (%.0f XAF) : C'est votre plus gros poste de dépenses. Évaluez si vous pouvez négocier votre loyer ou chercher un logement moins cher.

📱 Abonnements (%.0f XAF) : Vérifiez régulièrement vos abonnements et annulez ceux que vous n'utilisez plus.

💡 Conseil : Activez les alertes de budget pour être notifié quand vous approchez de vos limites.`,
		totalExpenses,
		preferences.Expenses.Food,
		preferences.Expenses.Transport,
		preferences.Expenses.Housing,
		preferences.Expenses.Subscriptions)
}

// generateSavingsAdvice génère un conseil d'épargne personnalisé
func (s *PreferencesAIService) generateSavingsAdvice(preferences *entity.UserPreferences) string {
	return fmt.Sprintf(`Conseil d'épargne personnalisé :

🎯 Objectif principal : %s
🎯 Objectif secondaire : %s
💰 Cible d'épargne : %.0f XAF
⏰ Échéance : %s

Pour atteindre votre objectif de %.0f XAF en %s, je recommande :

1. Épargne automatique : Configurez un virement automatique de %.0f XAF par mois
2. Règle 50/30/20 : 50%% pour les besoins, 30%% pour les envies, 20%% pour l'épargne
3. Compte séparé : Ouvrez un compte d'épargne dédié à vos objectifs
4. Suivi régulier : Vérifiez vos progrès chaque semaine

💪 Vous êtes sur la bonne voie ! Chaque petit montant compte vers votre objectif.`,
		preferences.Goals.MainGoal,
		preferences.Goals.SecondaryGoal,
		preferences.Goals.SavingsTarget,
		preferences.Goals.Deadline,
		preferences.Goals.SavingsTarget,
		preferences.Goals.Deadline,
		preferences.Goals.SavingsTarget/6) // Estimation mensuelle
}

// generateHabitsAdvice génère un conseil d'habitudes personnalisé
func (s *PreferencesAIService) generateHabitsAdvice(preferences *entity.UserPreferences) string {
	return fmt.Sprintf(`Conseil d'habitudes personnalisé :

⏰ Temps de planification : %s
⏱️ Temps de focus quotidien : %s
🎯 Habitude personnalisée : %s

Recommandations pour optimiser votre routine :

1. Planification %s : Prenez %s chaque jour pour organiser vos finances
2. Focus quotidien : Utilisez vos %s pour des tâches financières importantes
3. Habitude spécifique : Maintenez votre habitude "%s"
4. Suivi : Consultez vos résumés %s pour rester motivé

💡 Astuce : Commencez par de petites habitudes et augmentez progressivement. La constance est plus importante que la perfection !`,
		preferences.Habits.PlanningTime,
		preferences.Habits.DailyFocusTime,
		preferences.Habits.CustomHabit,
		preferences.Habits.PlanningTime,
		preferences.Habits.DailyFocusTime,
		preferences.Habits.DailyFocusTime,
		preferences.Habits.CustomHabit,
		preferences.Habits.SummaryType)
}

// Close ferme la connexion au client AI (plus nécessaire avec la nouvelle API)
func (s *PreferencesAIService) Close() error {
	return nil
}
