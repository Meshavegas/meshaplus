package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"backend/pkg/config"
	"backend/pkg/logger"

	"google.golang.org/genai"
)

type AIService struct {
	config config.AIConfig
	logger logger.Logger
}

// CategoryResponse représente la réponse structurée de l'IA
type CategoryResponse struct {
	CategoryName  string `json:"categoryName"`
	Icon          string `json:"icon"`
	Color         string `json:"color"`
	IsNewCategory bool   `json:"isNewCategory"`
	Confidence    int    `json:"confidence"`
	Reasoning     string `json:"reasoning"`
}

func NewAIService(config config.AIConfig, logger logger.Logger) *AIService {
	return &AIService{
		config: config,
		logger: logger,
	}
}

func (j *AIService) GenerateCatherorie(existingsCat []string, currentItem string) (*CategoryResponse, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		j.logger.Error("Failed to create Gemini client: " + err.Error())
		return nil, err
	}

	// Construire le prompt avec les catégories existantes
	prompt := fmt.Sprintf(`Analysez l'item suivant et choisissez la catégorie la plus appropriée.

Item à catégoriser: "%s"

Catégories existantes: %v

Instructions:
1. Si une catégorie existante correspond bien à l'item, utilisez-la
2. Si aucune catégorie existante ne convient, proposez une nouvelle catégorie logique et descriptive
3. Le nom de la catégorie doit être en français, court et clair
4. Pour chaque catégorie (nouvelle ou existante), proposez une icône appropriée au format "library:iconname"
5. Utilisez les bibliothèques d'icônes suivantes:
   - fad: FontAwesome Duotone (icônes avec deux couleurs)
   - fas: FontAwesome Solid (icônes pleines)
   - far: FontAwesome Regular (icônes régulières)
   - fab: FontAwesome Brands (marques)
   - io: Ionicons
   - md: Material Icons
   - mc: Material Community Icons
6. Pour chaque catégorie, proposez une couleur hexadécimale appropriée (ex: #FF6B6B pour rouge, #4ECDC4 pour turquoise)
7. Expliquez votre raisonnement

Exemples d'icônes: "fad:utensils", "fas:shopping-cart", "io:restaurant", "md:food-bank"
Exemples de couleurs: "#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7"

Répondez avec la structure JSON demandée.`, currentItem, existingsCat)

	// Configuration pour la sortie structurée JSON
	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"categoryName": {
					Type: genai.TypeString,
				},
				"icon": {
					Type: genai.TypeString,
				},
				"color": {
					Type: genai.TypeString,
				},
				"isNewCategory": {
					Type: genai.TypeBoolean,
				},
				"confidence": {
					Type: genai.TypeInteger,
				},
				"reasoning": {
					Type: genai.TypeString,
				},
			},
			PropertyOrdering: []string{"categoryName", "icon", "color", "isNewCategory", "confidence", "reasoning"},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		j.logger.Error("Failed to generate content: " + err.Error())
		return nil, err
	}

	// Parser la réponse JSON
	var categoryResponse CategoryResponse
	if err := json.Unmarshal([]byte(result.Text()), &categoryResponse); err != nil {
		j.logger.Error("Failed to parse JSON response: " + err.Error())
		return nil, err
	}

	j.logger.Info(fmt.Sprintf("Catégorie générée: %s (nouvelle: %t, confiance: %d%%)",
		categoryResponse.CategoryName,
		categoryResponse.IsNewCategory,
		categoryResponse.Confidence))

	return &categoryResponse, nil
}
