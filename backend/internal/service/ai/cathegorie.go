package ai

import (
	"context"
	"fmt"
	"log"

	"backend/pkg/config"
	"backend/pkg/logger"

	"google.golang.org/genai"
)

type AIService struct {
	config config.AIConfig

	logger logger.Logger
}

func NewAIService(config config.AIConfig, logger logger.Logger) *AIService {
	return &AIService{
		config: config,
		logger: logger,
	}
}

func (j *AIService) GenerateCatherorie(existingsCat []string, currentItem string) *genai.GenerateContentResponse {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)

	if err != nil {
		j.logger.Error(err.Error())
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text("Explain how AI works in a few words"),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())
	j.logger.Info(result.Text())
	return result
}
