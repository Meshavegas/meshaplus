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

// TaskService gère la logique métier des tâches
type TaskService struct {
	taskRepo repository.TaskRepository
	logger   logger.Logger
}

// NewTaskService crée une nouvelle instance de TaskService
func NewTaskService(taskRepo repository.TaskRepository, logger logger.Logger) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		logger:   logger,
	}
}

// CreateTask crée une nouvelle tâche
func (s *TaskService) CreateTask(ctx context.Context, userID uuid.UUID, req entity.CreateTaskRequest) (*entity.Task, error) {
	// Validation des données
	if req.Title == "" {
		return nil, fmt.Errorf("le titre est requis")
	}

	// Création de la tâche
	task := &entity.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		CategoryID:  req.CategoryID,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		s.logger.Error("Erreur création tâche", logger.Error(err))
		return nil, fmt.Errorf("erreur création tâche: %w", err)
	}

	s.logger.Info("Tâche créée avec succès",
		logger.String("task_id", task.ID.String()),
		logger.String("user_id", userID.String()),
		logger.String("title", task.Title),
	)

	return task, nil
}

// GetTask récupère une tâche par son ID
func (s *TaskService) GetTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*entity.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		s.logger.Error("Erreur récupération tâche", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération tâche: %w", err)
	}

	// Vérifier que la tâche appartient à l'utilisateur
	if task.UserID != userID {
		s.logger.Warn("Tentative d'accès non autorisé à une tâche",
			logger.String("user_id", userID.String()),
			logger.String("task_id", taskID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	return task, nil
}

// GetUserTasks récupère toutes les tâches d'un utilisateur
func (s *TaskService) GetUserTasks(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	tasks, err := s.taskRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération tâches utilisateur", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération tâches: %w", err)
	}

	return tasks, nil
}

// UpdateTask met à jour une tâche
func (s *TaskService) UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, req entity.UpdateTaskRequest) (*entity.Task, error) {
	// Récupérer la tâche existante
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		s.logger.Error("Erreur récupération tâche pour mise à jour", logger.Error(err))
		return nil, fmt.Errorf("tâche non trouvée")
	}

	// Vérifier que la tâche appartient à l'utilisateur
	if task.UserID != userID {
		s.logger.Warn("Tentative de mise à jour non autorisée d'une tâche",
			logger.String("user_id", userID.String()),
			logger.String("task_id", taskID.String()),
		)
		return nil, fmt.Errorf("accès non autorisé")
	}

	// Mettre à jour les champs fournis
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.CategoryID != nil {
		task.CategoryID = req.CategoryID
	}

	task.UpdatedAt = time.Now()

	// Sauvegarder les modifications
	if err := s.taskRepo.Update(ctx, task); err != nil {
		s.logger.Error("Erreur mise à jour tâche", logger.Error(err))
		return nil, fmt.Errorf("erreur mise à jour tâche: %w", err)
	}

	s.logger.Info("Tâche mise à jour avec succès",
		logger.String("task_id", task.ID.String()),
		logger.String("user_id", userID.String()),
	)

	return task, nil
}

// DeleteTask supprime une tâche
func (s *TaskService) DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error {
	// Récupérer la tâche pour vérifier la propriété
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		s.logger.Error("Erreur récupération tâche pour suppression", logger.Error(err))
		return fmt.Errorf("tâche non trouvée")
	}

	// Vérifier que la tâche appartient à l'utilisateur
	if task.UserID != userID {
		s.logger.Warn("Tentative de suppression non autorisée d'une tâche",
			logger.String("user_id", userID.String()),
			logger.String("task_id", taskID.String()),
		)
		return fmt.Errorf("accès non autorisé")
	}

	// Supprimer la tâche
	if err := s.taskRepo.Delete(ctx, taskID); err != nil {
		s.logger.Error("Erreur suppression tâche", logger.Error(err))
		return fmt.Errorf("erreur suppression tâche: %w", err)
	}

	s.logger.Info("Tâche supprimée avec succès",
		logger.String("task_id", taskID.String()),
		logger.String("user_id", userID.String()),
	)

	return nil
}

// GetCompletedTasks récupère les tâches complétées d'un utilisateur
func (s *TaskService) GetCompletedTasks(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	tasks, err := s.taskRepo.GetCompletedByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération tâches complétées", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération tâches complétées: %w", err)
	}

	return tasks, nil
}

// GetPendingTasks récupère les tâches en attente d'un utilisateur
func (s *TaskService) GetPendingTasks(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	tasks, err := s.taskRepo.GetPendingByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération tâches en attente", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération tâches en attente: %w", err)
	}

	return tasks, nil
}

// CompleteTask marque une tâche comme complétée
func (s *TaskService) CompleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*entity.Task, error) {
	req := entity.UpdateTaskRequest{
		Status: &[]string{"completed"}[0],
	}
	return s.UpdateTask(ctx, userID, taskID, req)
}

// GetTaskStats récupère les statistiques des tâches d'un utilisateur
func (s *TaskService) GetTaskStats(ctx context.Context, userID uuid.UUID) (*entity.TaskStatsResponse, error) {
	// Récupérer toutes les tâches
	tasks, err := s.taskRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur récupération tâches pour statistiques", logger.Error(err))
		return nil, fmt.Errorf("erreur récupération statistiques: %w", err)
	}

	// Calculer les statistiques
	totalTasks := int64(len(tasks))
	completedTasks := int64(0)
	pendingTasks := int64(0)

	for _, task := range tasks {
		if task.Status == "completed" {
			completedTasks++
		} else {
			pendingTasks++
		}
	}

	var completionRate float64
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}

	stats := &entity.TaskStatsResponse{
		TotalTasks:     totalTasks,
		CompletedTasks: completedTasks,
		PendingTasks:   pendingTasks,
		CompletionRate: completionRate,
	}

	return stats, nil
}
