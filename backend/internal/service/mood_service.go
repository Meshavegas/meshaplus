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

// MoodService gère la logique métier des moods
type MoodService struct {
	moodRepo repository.MoodRepository
	logger   logger.Logger
}

// NewMoodService crée une nouvelle instance de MoodService
func NewMoodService(moodRepo repository.MoodRepository, logger logger.Logger) *MoodService {
	return &MoodService{
		moodRepo: moodRepo,
		logger:   logger,
	}
}

// CreateMood crée un nouveau mood
func (s *MoodService) CreateMood(ctx context.Context, userID uuid.UUID, req entity.CreateMoodRequest) (*entity.Mood, error) {
	// Validation des données
	if req.Feeling == "" {
		return nil, fmt.Errorf("le sentiment est requis")
	}

	// Vérifier si un mood existe déjà pour cette date
	existingMood, err := s.moodRepo.GetByDate(ctx, userID, req.Date)
	if err == nil && existingMood != nil {
		return nil, fmt.Errorf("un mood existe déjà pour cette date")
	}

	// Création du mood
	mood := &entity.Mood{
		ID:        uuid.New(),
		UserID:    userID,
		Date:      req.Date,
		Feeling:   req.Feeling,
		Note:      req.Note,
		CreatedAt: time.Now(),
	}

	err = s.moodRepo.Create(ctx, mood)
	if err != nil {
		s.logger.Error("Erreur lors de la création du mood", logger.Error(err))
		return nil, fmt.Errorf("erreur lors de la création du mood: %w", err)
	}

	s.logger.Info("Mood créé avec succès", logger.String("mood_id", mood.ID.String()))
	return mood, nil
}

// GetMoodByID récupère un mood par son ID
func (s *MoodService) GetMoodByID(ctx context.Context, userID, moodID uuid.UUID) (*entity.Mood, error) {
	mood, err := s.moodRepo.GetByID(ctx, moodID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du mood: %w", err)
	}

	// Vérifier que le mood appartient à l'utilisateur
	if mood.UserID != userID {
		return nil, fmt.Errorf("mood non trouvé")
	}

	return mood, nil
}

// GetMoodsByUserID récupère tous les moods d'un utilisateur
func (s *MoodService) GetMoodsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Mood, error) {
	moods, err := s.moodRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des moods: %w", err)
	}

	return moods, nil
}

// GetMoodByDate récupère le mood d'une date spécifique
func (s *MoodService) GetMoodByDate(ctx context.Context, userID uuid.UUID, date time.Time) (*entity.Mood, error) {
	mood, err := s.moodRepo.GetByDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du mood par date: %w", err)
	}

	return mood, nil
}

// GetMoodsByDateRange récupère les moods dans une plage de dates
func (s *MoodService) GetMoodsByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*entity.Mood, error) {
	moods, err := s.moodRepo.GetByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des moods par plage de dates: %w", err)
	}

	return moods, nil
}

// GetMoodsByFeeling récupère les moods par sentiment
func (s *MoodService) GetMoodsByFeeling(ctx context.Context, userID uuid.UUID, feeling string) ([]*entity.Mood, error) {
	moods, err := s.moodRepo.GetByFeeling(ctx, userID, feeling)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des moods par sentiment: %w", err)
	}

	return moods, nil
}

// UpdateMood met à jour un mood
func (s *MoodService) UpdateMood(ctx context.Context, userID, moodID uuid.UUID, req entity.UpdateMoodRequest) (*entity.Mood, error) {
	// Récupérer le mood existant
	mood, err := s.GetMoodByID(ctx, userID, moodID)
	if err != nil {
		return nil, err
	}

	// Mettre à jour les champs
	if req.Feeling != nil && *req.Feeling != "" {
		mood.Feeling = *req.Feeling
	}
	if req.Note != nil {
		mood.Note = *req.Note
	}

	err = s.moodRepo.Update(ctx, mood)
	if err != nil {
		s.logger.Error("Erreur lors de la mise à jour du mood", logger.Error(err))
		return nil, fmt.Errorf("erreur lors de la mise à jour du mood: %w", err)
	}

	s.logger.Info("Mood mis à jour avec succès", logger.String("mood_id", mood.ID.String()))
	return mood, nil
}

// DeleteMood supprime un mood
func (s *MoodService) DeleteMood(ctx context.Context, userID, moodID uuid.UUID) error {
	// Vérifier que le mood existe et appartient à l'utilisateur
	_, err := s.GetMoodByID(ctx, userID, moodID)
	if err != nil {
		return err
	}

	err = s.moodRepo.Delete(ctx, moodID)
	if err != nil {
		s.logger.Error("Erreur lors de la suppression du mood", logger.Error(err))
		return fmt.Errorf("erreur lors de la suppression du mood: %w", err)
	}

	s.logger.Info("Mood supprimé avec succès", logger.String("mood_id", moodID.String()))
	return nil
}

// UpsertMood crée ou met à jour un mood pour une date donnée
func (s *MoodService) UpsertMood(ctx context.Context, userID uuid.UUID, req entity.CreateMoodRequest) (*entity.Mood, error) {
	// Vérifier si un mood existe déjà pour cette date
	existingMood, err := s.moodRepo.GetByDate(ctx, userID, req.Date)
	if err != nil {
		// Aucun mood existant, créer un nouveau
		return s.CreateMood(ctx, userID, req)
	}

	// Mettre à jour le mood existant
	updateReq := entity.UpdateMoodRequest{
		Feeling: &req.Feeling,
		Note:    &req.Note,
	}

	return s.UpdateMood(ctx, userID, existingMood.ID, updateReq)
}
