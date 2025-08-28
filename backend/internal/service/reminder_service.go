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

// ReminderService gère la logique métier des reminders
type ReminderService struct {
	reminderRepo repository.ReminderRepository
	logger       logger.Logger
}

// NewReminderService crée une nouvelle instance de ReminderService
func NewReminderService(reminderRepo repository.ReminderRepository, logger logger.Logger) *ReminderService {
	return &ReminderService{
		reminderRepo: reminderRepo,
		logger:       logger,
	}
}

// CreateReminder crée un nouveau reminder
func (s *ReminderService) CreateReminder(ctx context.Context, userID uuid.UUID, req entity.CreateReminderRequest) (*entity.Reminder, error) {
	// Validation des données
	if req.Message == "" {
		return nil, fmt.Errorf("le message est requis")
	}

	if req.TriggerAt.Before(time.Now()) {
		return nil, fmt.Errorf("la date de déclenchement doit être dans le futur")
	}

	// Création du reminder
	reminder := &entity.Reminder{
		ID:        uuid.New(),
		UserID:    userID,
		TaskID:    req.TaskID,
		TransacID: req.TransacID,
		Message:   req.Message,
		TriggerAt: req.TriggerAt,
		IsSent:    false,
		CreatedAt: time.Now(),
	}

	err := s.reminderRepo.Create(ctx, reminder)
	if err != nil {
		s.logger.Error("Erreur lors de la création du reminder", logger.Error(err))
		return nil, fmt.Errorf("erreur lors de la création du reminder: %w", err)
	}

	s.logger.Info("Reminder créé avec succès", logger.String("reminder_id", reminder.ID.String()))
	return reminder, nil
}

// GetReminderByID récupère un reminder par son ID
func (s *ReminderService) GetReminderByID(ctx context.Context, userID, reminderID uuid.UUID) (*entity.Reminder, error) {
	reminder, err := s.reminderRepo.GetByID(ctx, reminderID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du reminder: %w", err)
	}

	// Vérifier que le reminder appartient à l'utilisateur
	if reminder.UserID != userID {
		return nil, fmt.Errorf("reminder non trouvé")
	}

	return reminder, nil
}

// GetRemindersByUserID récupère tous les reminders d'un utilisateur
func (s *ReminderService) GetRemindersByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Reminder, error) {
	reminders, err := s.reminderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des reminders: %w", err)
	}

	return reminders, nil
}

// GetRemindersByTaskID récupère les reminders d'une tâche
func (s *ReminderService) GetRemindersByTaskID(ctx context.Context, userID, taskID uuid.UUID) ([]*entity.Reminder, error) {
	reminders, err := s.reminderRepo.GetByTaskID(ctx, userID, taskID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des reminders par tâche: %w", err)
	}

	return reminders, nil
}

// GetRemindersByTransactionID récupère les reminders d'une transaction
func (s *ReminderService) GetRemindersByTransactionID(ctx context.Context, userID, transacID uuid.UUID) ([]*entity.Reminder, error) {
	reminders, err := s.reminderRepo.GetByTransactionID(ctx, userID, transacID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des reminders par transaction: %w", err)
	}

	return reminders, nil
}

// GetPendingReminders récupère les reminders non envoyés
func (s *ReminderService) GetPendingReminders(ctx context.Context, userID uuid.UUID) ([]*entity.Reminder, error) {
	reminders, err := s.reminderRepo.GetPending(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des reminders en attente: %w", err)
	}

	return reminders, nil
}

// GetDueReminders récupère les reminders qui doivent être déclenchés
func (s *ReminderService) GetDueReminders(ctx context.Context) ([]*entity.Reminder, error) {
	now := time.Now()
	reminders, err := s.reminderRepo.GetDue(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des reminders à déclencher: %w", err)
	}

	return reminders, nil
}

// UpdateReminder met à jour un reminder
func (s *ReminderService) UpdateReminder(ctx context.Context, userID, reminderID uuid.UUID, req entity.UpdateReminderRequest) (*entity.Reminder, error) {
	// Récupérer le reminder existant
	reminder, err := s.GetReminderByID(ctx, userID, reminderID)
	if err != nil {
		return nil, err
	}

	// Vérifier que le reminder n'a pas encore été envoyé
	if reminder.IsSent {
		return nil, fmt.Errorf("impossible de modifier un reminder déjà envoyé")
	}

	// Mettre à jour les champs
	if req.Message != nil {
		reminder.Message = *req.Message
	}
	if req.TriggerAt != nil {
		if req.TriggerAt.Before(time.Now()) {
			return nil, fmt.Errorf("la date de déclenchement doit être dans le futur")
		}
		reminder.TriggerAt = *req.TriggerAt
	}

	err = s.reminderRepo.Update(ctx, reminder)
	if err != nil {
		s.logger.Error("Erreur lors de la mise à jour du reminder", logger.Error(err))
		return nil, fmt.Errorf("erreur lors de la mise à jour du reminder: %w", err)
	}

	s.logger.Info("Reminder mis à jour avec succès", logger.String("reminder_id", reminder.ID.String()))
	return reminder, nil
}

// DeleteReminder supprime un reminder
func (s *ReminderService) DeleteReminder(ctx context.Context, userID, reminderID uuid.UUID) error {
	// Vérifier que le reminder existe et appartient à l'utilisateur
	_, err := s.GetReminderByID(ctx, userID, reminderID)
	if err != nil {
		return err
	}

	err = s.reminderRepo.Delete(ctx, reminderID)
	if err != nil {
		s.logger.Error("Erreur lors de la suppression du reminder", logger.Error(err))
		return fmt.Errorf("erreur lors de la suppression du reminder: %w", err)
	}

	s.logger.Info("Reminder supprimé avec succès", logger.String("reminder_id", reminderID.String()))
	return nil
}

// MarkReminderAsSent marque un reminder comme envoyé
func (s *ReminderService) MarkReminderAsSent(ctx context.Context, reminderID uuid.UUID) error {
	err := s.reminderRepo.MarkAsSent(ctx, reminderID)
	if err != nil {
		s.logger.Error("Erreur lors du marquage du reminder comme envoyé", logger.Error(err))
		return fmt.Errorf("erreur lors du marquage du reminder comme envoyé: %w", err)
	}

	s.logger.Info("Reminder marqué comme envoyé", logger.String("reminder_id", reminderID.String()))
	return nil
}

// ProcessDueReminders traite les reminders qui doivent être déclenchés
func (s *ReminderService) ProcessDueReminders(ctx context.Context) error {
	dueReminders, err := s.GetDueReminders(ctx)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des reminders à traiter: %w", err)
	}

	for _, reminder := range dueReminders {
		// Ici vous pouvez ajouter la logique pour envoyer la notification
		// Par exemple, envoyer un email, une notification push, etc.
		s.logger.Info("Traitement du reminder",
			logger.String("reminder_id", reminder.ID.String()),
			logger.String("message", reminder.Message))

		// Marquer comme envoyé
		err := s.MarkReminderAsSent(ctx, reminder.ID)
		if err != nil {
			s.logger.Error("Erreur lors du marquage du reminder comme envoyé",
				logger.String("reminder_id", reminder.ID.String()),
				logger.Error(err))
		}
	}

	s.logger.Info("Traitement des reminders terminé", logger.Int("count", len(dueReminders)))
	return nil
}
