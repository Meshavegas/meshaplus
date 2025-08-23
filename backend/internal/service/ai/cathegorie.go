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
4. Expliquez votre raisonnement

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
			PropertyOrdering: []string{"categoryName", "isNewCategory", "confidence", "reasoning"},
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
