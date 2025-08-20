package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// TaskRepository implémente repository.TaskRepository
type TaskRepository struct {
	db *pg.DB
}

// NewTaskRepository crée une nouvelle instance de TaskRepository
func NewTaskRepository(db *pg.DB) repository.TaskRepository {
	return &TaskRepository{db: db}
}

// Create crée une nouvelle tâche
func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	_, err := r.db.WithContext(ctx).Model(task).Insert()
	if err != nil {
		return fmt.Errorf("erreur création tâche: %w", err)
	}
	return nil
}

// GetIncomingByUserID récupère les tâches entrantes d'un utilisateur
func (r *TaskRepository) GetIncomingByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	var tasks []*entity.Task
	err := r.db.WithContext(ctx).Model(&tasks).Where("user_id = ? AND status = ?", userID, "incoming").Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération tâches entrantes: %w", err)
	}
	return tasks, nil
}

// GetByID récupère une tâche par son ID
func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	task := &entity.Task{}
	err := r.db.WithContext(ctx).Model(task).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("tâche non trouvée")
		}
		return nil, fmt.Errorf("erreur récupération tâche: %w", err)
	}
	return task, nil
}

// GetByUserID récupère toutes les tâches d'un utilisateur
func (r *TaskRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	var tasks []*entity.Task
	err := r.db.WithContext(ctx).Model(&tasks).Where("user_id = ?", userID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération tâches utilisateur: %w", err)
	}
	return tasks, nil
}

// Update met à jour une tâche
func (r *TaskRepository) Update(ctx context.Context, task *entity.Task) error {
	_, err := r.db.WithContext(ctx).Model(task).Where("id = ?", task.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour tâche: %w", err)
	}
	return nil
}

// Delete supprime une tâche
func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression tâche: %w", err)
	}
	return nil
}

// GetCompletedByUserID récupère les tâches complétées d'un utilisateur
func (r *TaskRepository) GetCompletedByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	var tasks []*entity.Task
	err := r.db.WithContext(ctx).Model(&tasks).Where("user_id = ? AND is_completed = ?", userID, true).Order("updated_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération tâches complétées: %w", err)
	}
	return tasks, nil
}

// GetPendingByUserID récupère les tâches en attente d'un utilisateur
func (r *TaskRepository) GetPendingByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error) {
	var tasks []*entity.Task
	err := r.db.WithContext(ctx).Model(&tasks).Where("user_id = ? AND is_completed = ?", userID, false).Order("due_date ASC NULLS LAST, created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération tâches en attente: %w", err)
	}
	return tasks, nil
}
