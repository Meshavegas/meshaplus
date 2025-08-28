package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// ReminderRepository implémente repository.ReminderRepository
type ReminderRepository struct {
	db *pg.DB
}

// NewReminderRepository crée une nouvelle instance de ReminderRepository
func NewReminderRepository(db *pg.DB) repository.ReminderRepository {
	return &ReminderRepository{db: db}
}

// Create crée un nouveau reminder
func (r *ReminderRepository) Create(ctx context.Context, reminder *entity.Reminder) error {
	_, err := r.db.WithContext(ctx).Model(reminder).Insert()
	if err != nil {
		return fmt.Errorf("erreur création reminder: %w", err)
	}
	return nil
}

// GetByID récupère un reminder par son ID
func (r *ReminderRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Reminder, error) {
	reminder := &entity.Reminder{}
	err := r.db.WithContext(ctx).Model(reminder).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("reminder non trouvé")
		}
		return nil, fmt.Errorf("erreur récupération reminder: %w", err)
	}
	return reminder, nil
}

// GetByUserID récupère tous les reminders d'un utilisateur
func (r *ReminderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Reminder, error) {
	var reminders []*entity.Reminder
	err := r.db.WithContext(ctx).Model(&reminders).Where("user_id = ?", userID).Order("trigger_at ASC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération reminders utilisateur: %w", err)
	}
	return reminders, nil
}

// GetByTaskID récupère les reminders d'une tâche
func (r *ReminderRepository) GetByTaskID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) ([]*entity.Reminder, error) {
	var reminders []*entity.Reminder
	err := r.db.WithContext(ctx).Model(&reminders).Where("user_id = ? AND task_id = ?", userID, taskID).Order("trigger_at ASC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération reminders par tâche: %w", err)
	}
	return reminders, nil
}

// GetByTransactionID récupère les reminders d'une transaction
func (r *ReminderRepository) GetByTransactionID(ctx context.Context, userID uuid.UUID, transacID uuid.UUID) ([]*entity.Reminder, error) {
	var reminders []*entity.Reminder
	err := r.db.WithContext(ctx).Model(&reminders).Where("user_id = ? AND transac_id = ?", userID, transacID).Order("trigger_at ASC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération reminders par transaction: %w", err)
	}
	return reminders, nil
}

// GetPending récupère les reminders non envoyés
func (r *ReminderRepository) GetPending(ctx context.Context, userID uuid.UUID) ([]*entity.Reminder, error) {
	var reminders []*entity.Reminder
	err := r.db.WithContext(ctx).Model(&reminders).Where("user_id = ? AND is_sent = ?", userID, false).Order("trigger_at ASC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération reminders en attente: %w", err)
	}
	return reminders, nil
}

// GetDue récupère les reminders qui doivent être déclenchés
func (r *ReminderRepository) GetDue(ctx context.Context, now time.Time) ([]*entity.Reminder, error) {
	var reminders []*entity.Reminder
	err := r.db.WithContext(ctx).Model(&reminders).Where("is_sent = ? AND trigger_at <= ?", false, now).Order("trigger_at ASC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération reminders à déclencher: %w", err)
	}
	return reminders, nil
}

// Update met à jour un reminder
func (r *ReminderRepository) Update(ctx context.Context, reminder *entity.Reminder) error {
	_, err := r.db.WithContext(ctx).Model(reminder).Where("id = ?", reminder.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour reminder: %w", err)
	}
	return nil
}

// Delete supprime un reminder
func (r *ReminderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Reminder{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression reminder: %w", err)
	}
	return nil
}

// MarkAsSent marque un reminder comme envoyé
func (r *ReminderRepository) MarkAsSent(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Reminder{}).Set("is_sent = ?", true).Where("id = ?", id).Update()
	if err != nil {
		return fmt.Errorf("erreur marquage reminder comme envoyé: %w", err)
	}
	return nil
}
