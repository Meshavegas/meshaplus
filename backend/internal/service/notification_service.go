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

// NotificationService gère la logique métier des notifications
type NotificationService struct {
	notificationRepo repository.NotificationRepository
	logger           logger.Logger
}

// NewNotificationService crée une nouvelle instance de NotificationService
func NewNotificationService(notificationRepo repository.NotificationRepository, logger logger.Logger) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		logger:           logger,
	}
}

// CreateNotification crée une nouvelle notification
func (s *NotificationService) CreateNotification(ctx context.Context, userID uuid.UUID, req entity.CreateNotificationRequest) (*entity.Notification, error) {
	// Validation des données
	if req.Title == "" {
		return nil, fmt.Errorf("le titre est requis")
	}

	if req.Message == "" {
		return nil, fmt.Errorf("le message est requis")
	}

	// Création de la notification
	notification := &entity.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     req.Title,
		Message:   req.Message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		s.logger.Error("Erreur lors de la création de la notification", logger.Error(err))
		return nil, fmt.Errorf("erreur lors de la création de la notification: %w", err)
	}

	s.logger.Info("Notification créée avec succès", logger.String("notification_id", notification.ID.String()))
	return notification, nil
}

// GetNotificationByID récupère une notification par son ID
func (s *NotificationService) GetNotificationByID(ctx context.Context, userID, notificationID uuid.UUID) (*entity.Notification, error) {
	notification, err := s.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la notification: %w", err)
	}

	// Vérifier que la notification appartient à l'utilisateur
	if notification.UserID != userID {
		return nil, fmt.Errorf("notification non trouvée")
	}

	return notification, nil
}

// GetNotificationsByUserID récupère toutes les notifications d'un utilisateur
func (s *NotificationService) GetNotificationsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	notifications, err := s.notificationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des notifications: %w", err)
	}

	return notifications, nil
}

// GetUnreadNotifications récupère les notifications non lues d'un utilisateur
func (s *NotificationService) GetUnreadNotifications(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	notifications, err := s.notificationRepo.GetUnread(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des notifications non lues: %w", err)
	}

	return notifications, nil
}

// GetReadNotifications récupère les notifications lues d'un utilisateur
func (s *NotificationService) GetReadNotifications(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	notifications, err := s.notificationRepo.GetRead(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des notifications lues: %w", err)
	}

	return notifications, nil
}

// CountUnreadNotifications compte les notifications non lues d'un utilisateur
func (s *NotificationService) CountUnreadNotifications(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := s.notificationRepo.CountUnread(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("erreur lors du comptage des notifications non lues: %w", err)
	}

	return count, nil
}

// MarkNotificationAsRead marque une notification comme lue
func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	// Vérifier que la notification existe et appartient à l'utilisateur
	_, err := s.GetNotificationByID(ctx, userID, notificationID)
	if err != nil {
		return err
	}

	err = s.notificationRepo.MarkAsRead(ctx, notificationID)
	if err != nil {
		s.logger.Error("Erreur lors du marquage de la notification comme lue", logger.Error(err))
		return fmt.Errorf("erreur lors du marquage de la notification comme lue: %w", err)
	}

	s.logger.Info("Notification marquée comme lue", logger.String("notification_id", notificationID.String()))
	return nil
}

// MarkAllNotificationsAsRead marque toutes les notifications d'un utilisateur comme lues
func (s *NotificationService) MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error {
	err := s.notificationRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		s.logger.Error("Erreur lors du marquage de toutes les notifications comme lues", logger.Error(err))
		return fmt.Errorf("erreur lors du marquage de toutes les notifications comme lues: %w", err)
	}

	s.logger.Info("Toutes les notifications marquées comme lues", logger.String("user_id", userID.String()))
	return nil
}

// DeleteNotification supprime une notification
func (s *NotificationService) DeleteNotification(ctx context.Context, userID, notificationID uuid.UUID) error {
	// Vérifier que la notification existe et appartient à l'utilisateur
	_, err := s.GetNotificationByID(ctx, userID, notificationID)
	if err != nil {
		return err
	}

	err = s.notificationRepo.Delete(ctx, notificationID)
	if err != nil {
		s.logger.Error("Erreur lors de la suppression de la notification", logger.Error(err))
		return fmt.Errorf("erreur lors de la suppression de la notification: %w", err)
	}

	s.logger.Info("Notification supprimée avec succès", logger.String("notification_id", notificationID.String()))
	return nil
}

// SendNotificationToUser envoie une notification à un utilisateur
func (s *NotificationService) SendNotificationToUser(ctx context.Context, userID uuid.UUID, title, message string) error {
	req := entity.CreateNotificationRequest{
		Title:   title,
		Message: message,
	}

	_, err := s.CreateNotification(ctx, userID, req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la notification: %w", err)
	}

	return nil
}

// SendTaskReminder envoie un rappel de tâche
func (s *NotificationService) SendTaskReminder(ctx context.Context, userID uuid.UUID, taskTitle string) error {
	title := "Rappel de tâche"
	message := fmt.Sprintf("N'oubliez pas votre tâche: %s", taskTitle)

	return s.SendNotificationToUser(ctx, userID, title, message)
}

// SendBudgetAlert envoie une alerte de budget
func (s *NotificationService) SendBudgetAlert(ctx context.Context, userID uuid.UUID, budgetName string, percentage float64) error {
	title := "Alerte Budget"
	message := fmt.Sprintf("Attention: Vous avez dépensé %.1f%% de votre budget '%s'", percentage, budgetName)

	return s.SendNotificationToUser(ctx, userID, title, message)
}

// SendGoalAchievement envoie une notification d'objectif atteint
func (s *NotificationService) SendGoalAchievement(ctx context.Context, userID uuid.UUID, goalTitle string) error {
	title := "Objectif Atteint!"
	message := fmt.Sprintf("Félicitations! Vous avez atteint votre objectif: %s", goalTitle)

	return s.SendNotificationToUser(ctx, userID, title, message)
}
