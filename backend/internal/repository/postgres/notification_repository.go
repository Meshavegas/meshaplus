package postgres

import (
	"backend/internal/domaine/entity"
	"backend/internal/domaine/repository"
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

// NotificationRepository implémente repository.NotificationRepository
type NotificationRepository struct {
	db *pg.DB
}

// NewNotificationRepository crée une nouvelle instance de NotificationRepository
func NewNotificationRepository(db *pg.DB) repository.NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create crée une nouvelle notification
func (r *NotificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	_, err := r.db.WithContext(ctx).Model(notification).Insert()
	if err != nil {
		return fmt.Errorf("erreur création notification: %w", err)
	}
	return nil
}

// GetByID récupère une notification par son ID
func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	notification := &entity.Notification{}
	err := r.db.WithContext(ctx).Model(notification).Where("id = ?", id).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("notification non trouvée")
		}
		return nil, fmt.Errorf("erreur récupération notification: %w", err)
	}
	return notification, nil
}

// GetByUserID récupère toutes les notifications d'un utilisateur
func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	err := r.db.WithContext(ctx).Model(&notifications).Where("user_id = ?", userID).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération notifications utilisateur: %w", err)
	}
	return notifications, nil
}

// GetUnread récupère les notifications non lues d'un utilisateur
func (r *NotificationRepository) GetUnread(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	err := r.db.WithContext(ctx).Model(&notifications).Where("user_id = ? AND is_read = ?", userID, false).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération notifications non lues: %w", err)
	}
	return notifications, nil
}

// GetRead récupère les notifications lues d'un utilisateur
func (r *NotificationRepository) GetRead(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	err := r.db.WithContext(ctx).Model(&notifications).Where("user_id = ? AND is_read = ?", userID, true).Order("created_at DESC").Select()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération notifications lues: %w", err)
	}
	return notifications, nil
}

// Update met à jour une notification
func (r *NotificationRepository) Update(ctx context.Context, notification *entity.Notification) error {
	_, err := r.db.WithContext(ctx).Model(notification).Where("id = ?", notification.ID).Update()
	if err != nil {
		return fmt.Errorf("erreur mise à jour notification: %w", err)
	}
	return nil
}

// Delete supprime une notification
func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("id = ?", id).Delete()
	if err != nil {
		return fmt.Errorf("erreur suppression notification: %w", err)
	}
	return nil
}

// MarkAsRead marque une notification comme lue
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Notification{}).Set("is_read = ?", true).Where("id = ?", id).Update()
	if err != nil {
		return fmt.Errorf("erreur marquage notification comme lue: %w", err)
	}
	return nil
}

// MarkAllAsRead marque toutes les notifications d'un utilisateur comme lues
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.WithContext(ctx).Model(&entity.Notification{}).Set("is_read = ?", true).Where("user_id = ? AND is_read = ?", userID, false).Update()
	if err != nil {
		return fmt.Errorf("erreur marquage toutes notifications comme lues: %w", err)
	}
	return nil
}

// CountUnread compte les notifications non lues d'un utilisateur
func (r *NotificationRepository) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := r.db.WithContext(ctx).Model(&entity.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count()
	if err != nil {
		return 0, fmt.Errorf("erreur comptage notifications non lues: %w", err)
	}
	return count, nil
}
