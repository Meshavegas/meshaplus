package service

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"backend/internal/service/ai"
	"backend/pkg/logger"
	"context"

	"github.com/google/uuid"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
	aiService    *ai.AIService
	logger       logger.Logger
}

func NewCategoryService(categoryRepo repository.CategoryRepository, aiService *ai.AIService, logger logger.Logger) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		aiService:    aiService,
		logger:       logger,
	}
}

// CategorizeItem catégorise un item en utilisant l'IA et la base de données
func (s *CategoryService) CategorizeItem(ctx context.Context, userID uuid.UUID, item string, categoryType string) (*ai.CategoryResponse, error) {
	// Récupérer les catégories existantes depuis la base de données
	existingCategories, err := s.categoryRepo.GetByType(ctx, userID, categoryType)
	if err != nil {
		s.logger.Error("Erreur récupération catégories existantes", logger.Error(err))
		return nil, err
	}

	// Extraire les noms des catégories
	var categoryNames []string
	for _, category := range existingCategories {
		categoryNames = append(categoryNames, category.Name)
	}

	// Utiliser l'IA pour catégoriser
	categoryResponse, err := s.aiService.GenerateCatherorie(categoryNames, item)
	if err != nil {
		s.logger.Error("Erreur catégorisation IA", logger.Error(err))
		return nil, err
	}

	// Si c'est une nouvelle catégorie, la créer dans la base de données
	if categoryResponse.IsNewCategory {
		newCategory := &entity.Category{
			UserID: userID,
			Name:   categoryResponse.CategoryName,
			Type:   categoryType,
		}

		if err := s.categoryRepo.Create(ctx, newCategory); err != nil {
			s.logger.Error("Erreur création nouvelle catégorie", logger.Error(err))
			return nil, err
		}

		s.logger.Info("Nouvelle catégorie créée",
			logger.String("name", categoryResponse.CategoryName),
			logger.String("type", categoryType))
	}

	return categoryResponse, nil
}

// GetCategoriesByType récupère les catégories d'un utilisateur par type
func (s *CategoryService) GetCategoriesByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]*entity.Category, error) {
	return s.categoryRepo.GetByType(ctx, userID, categoryType)
}

// get all categories
func (s *CategoryService) GetAllCategories(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error) {
	categories, err := s.categoryRepo.GetAll(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération catégories", logger.Error(err))
		return nil, err
	}

	// Récupérer les catégories parentes
	var parentCategories []*entity.Category
	for _, category := range categories {
		if category.ParentID == nil {
			parentCategories = append(parentCategories, category)
		}
	}

	return parentCategories, nil
}

// CreateCategory crée une nouvelle catégorie manuellement
func (s *CategoryService) CreateCategory(ctx context.Context, userID uuid.UUID, req *entity.CreateCategoryRequest) (*entity.Category, error) {
	category := &entity.Category{
		UserID:   userID,
		Name:     req.Name,
		Type:     req.Type,
		ParentID: req.ParentID,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		s.logger.Error("Erreur création catégorie", logger.Error(err))
		return nil, err
	}

	return category, nil
}
