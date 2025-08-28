package entity

import (
	"time"

	"github.com/google/uuid"
)

// ==================== TASK REQUESTS ====================

// CreateTaskRequest représente la requête pour créer une tâche
type CreateTaskRequest struct {
	CategoryID      *uuid.UUID `json:"category_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title           string     `json:"title" validate:"required,min=1,max=255" example:"Faire les courses"`
	Description     string     `json:"description" validate:"max=1000" example:"Acheter du pain, du lait et des fruits"`
	Priority        string     `json:"priority" validate:"required,oneof=low medium high" example:"medium"`
	DueDate         *time.Time `json:"due_date,omitempty" validate:"omitempty,gt=now" example:"2024-12-31T23:59:59Z"`
	DurationPlanned int        `json:"duration_planned" validate:"min=0,max=1440" example:"60"` // en minutes
	RecurrenceRule  *string    `json:"recurrence_rule,omitempty" validate:"omitempty,max=255" example:"daily"`
}

// UpdateTaskRequest représente la requête pour mettre à jour une tâche
type UpdateTaskRequest struct {
	CategoryID      *uuid.UUID `json:"category_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title           *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255" example:"Faire les courses"`
	Description     *string    `json:"description,omitempty" validate:"omitempty,max=1000" example:"Acheter du pain, du lait et des fruits"`
	Priority        *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high" example:"high"`
	DueDate         *time.Time `json:"due_date,omitempty" validate:"omitempty,gt=now" example:"2024-12-31T23:59:59Z"`
	DurationPlanned *int       `json:"duration_planned,omitempty" validate:"omitempty,min=0,max=1440" example:"90"`
	DurationSpent   *int       `json:"duration_spent,omitempty" validate:"omitempty,min=0,max=1440" example:"45"`
	Status          *string    `json:"status,omitempty" validate:"omitempty,oneof=expired done incoming running" example:"running"`
	RecurrenceRule  *string    `json:"recurrence_rule,omitempty" validate:"omitempty,max=255" example:"weekly"`
}

// ==================== TRANSACTION REQUESTS ====================

// CreateTransactionRequest représente la requête pour créer une transaction
type CreateTransactionRequest struct {
	AccountID    *uuid.UUID `json:"account_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Type         string     `json:"type" validate:"required,oneof=income expense transfer saving refund" example:"expense"`
	ToAccountID  *uuid.UUID `json:"to_account_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	SavingGoalID *uuid.UUID `json:"saving_goal_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Amount       float64    `json:"amount" validate:"required" example:"25.50"`
	Description  string     `json:"description" validate:"required,min=1,max=255" example:"Achat alimentaire"`
	Date         time.Time  `json:"date" validate:"required" example:"2024-01-15T00:00:00Z"`
	Recurring    bool       `json:"recurring" example:"false"`
}

// UpdateTransactionRequest représente la requête pour mettre à jour une transaction
type UpdateTransactionRequest struct {
	AccountID    *uuid.UUID `json:"account_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	CategoryID   *uuid.UUID `json:"category_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Type         *string    `json:"type,omitempty" validate:"omitempty,oneof=income expense transfer saving refund" example:"income"`
	ToAccountID  *uuid.UUID `json:"to_account_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	SavingGoalID *uuid.UUID `json:"saving_goal_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Amount       *float64   `json:"amount,omitempty" validate:"omitempty" example:"100.00"`
	Description  *string    `json:"description,omitempty" validate:"omitempty,min=1,max=255" example:"Salaire mensuel"`
	Date         *time.Time `json:"date,omitempty" validate:"omitempty" example:"2024-01-15T00:00:00Z"`
	Recurring    *bool      `json:"recurring,omitempty" example:"true"`
}

// ==================== BUDGET REQUESTS ====================

// CreateBudgetRequest représente la requête pour créer un budget
type CreateBudgetRequest struct {
	CategoryID    uuid.UUID `json:"category_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name          string    `json:"name" validate:"required,min=1,max=255" example:"Budget Alimentation Janvier"`
	AmountPlanned float64   `json:"amount_planned" validate:"required,gt=0" example:"500.00"`
	Period        string    `json:"period" validate:"required,oneof=monthly yearly weekly daily" example:"monthly"`
}

// UpdateBudgetRequest représente la requête pour mettre à jour un budget
type UpdateBudgetRequest struct {
	Name          *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255" example:"Budget Alimentation Février"`
	AmountPlanned *float64 `json:"amount_planned,omitempty" validate:"omitempty,gt=0" example:"600.00"`
	AmountSpent   *float64 `json:"amount_spent,omitempty" validate:"omitempty,gte=0" example:"450.00"`
	Period        *string  `json:"period,omitempty" validate:"omitempty,oneof=monthly yearly weekly daily" example:"monthly"`
}

// ==================== SAVING GOAL REQUESTS ====================

// CreateSavingGoalRequest représente la requête pour créer un objectif d'épargne
type CreateSavingGoalRequest struct {
	Title        string     `json:"title" validate:"required,min=1,max=255" example:"Vacances d'été"`
	TargetAmount float64    `json:"target_amount" validate:"required,gt=0" example:"2000.00"`
	Deadline     *time.Time `json:"deadline,omitempty" validate:"omitempty,gt=now" example:"2024-06-30T00:00:00Z"`
	Frequency    *string    `json:"frequency,omitempty" validate:"omitempty,oneof=weekly monthly yearly" example:"monthly"`
}

// UpdateSavingGoalRequest représente la requête pour mettre à jour un objectif d'épargne
type UpdateSavingGoalRequest struct {
	Title         *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255" example:"Vacances d'été 2024"`
	TargetAmount  *float64   `json:"target_amount,omitempty" validate:"omitempty,gt=0" example:"2500.00"`
	CurrentAmount *float64   `json:"current_amount,omitempty" validate:"omitempty,gte=0" example:"1500.00"`
	Deadline      *time.Time `json:"deadline,omitempty" validate:"omitempty,gt=now" example:"2024-07-31T00:00:00Z"`
	IsAchieved    *bool      `json:"is_achieved,omitempty" example:"false"`
	Frequency     *string    `json:"frequency,omitempty" validate:"omitempty,oneof=weekly monthly yearly" example:"monthly"`
}

// ==================== REMINDER REQUESTS ====================

// CreateReminderRequest représente la requête pour créer un rappel
type CreateReminderRequest struct {
	TaskID    *uuid.UUID `json:"task_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	TransacID *uuid.UUID `json:"transac_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Message   string     `json:"message" validate:"required,min=1,max=500" example:"N'oubliez pas de payer la facture d'électricité"`
	TriggerAt time.Time  `json:"trigger_at" validate:"required,gt=now" example:"2024-01-20T09:00:00Z"`
}

// UpdateReminderRequest représente la requête pour mettre à jour un rappel
type UpdateReminderRequest struct {
	Message   *string    `json:"message,omitempty" validate:"omitempty,min=1,max=500" example:"Rappel: Payer la facture d'électricité"`
	TriggerAt *time.Time `json:"trigger_at,omitempty" validate:"omitempty,gt=now" example:"2024-01-21T09:00:00Z"`
	IsSent    *bool      `json:"is_sent,omitempty" example:"false"`
}

// ==================== CATEGORY REQUESTS ====================

// CreateCategoryRequest représente la requête pour créer une catégorie
type CreateCategoryRequest struct {
	Name     string     `json:"name" validate:"required,min=1,max=100" example:"Alimentation"`
	Type     string     `json:"type" validate:"required,oneof=expense revenue task" example:"expense"`
	ParentID *uuid.UUID `json:"parent_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Icon     string     `json:"icon,omitempty" validate:"omitempty,max=50" example:"fad:utensils"`
	Color    string     `json:"color,omitempty" validate:"omitempty,max=20" example:"#FF6B6B"`
}

// UpdateCategoryRequest représente la requête pour mettre à jour une catégorie
type UpdateCategoryRequest struct {
	Name     *string    `json:"name,omitempty" validate:"omitempty,min=1,max=100" example:"Alimentation et courses"`
	Type     *string    `json:"type,omitempty" validate:"omitempty,oneof=expense revenue task" example:"expense"`
	ParentID *uuid.UUID `json:"parent_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Icon     *string    `json:"icon,omitempty" validate:"omitempty,max=50" example:"fad:utensils"`
	Color    *string    `json:"color,omitempty" validate:"omitempty,max=20" example:"#FF6B6B"`
}

// ==================== ACCOUNT REQUESTS ====================

// CreateAccountRequest représente la requête pour créer un compte
type CreateAccountRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=255" example:"Compte principal"`
	Type          string  `json:"type" validate:"required,oneof=checking savings mobile_money" example:"checking"`
	Balance       float64 `json:"balance" validate:"gte=0" example:"1500.00"`
	Currency      string  `json:"currency" validate:"required,len=3" example:"EUR"`
	Icon          string  `json:"icon" validate:"omitempty,max=50" example:"fad:utensils"`
	Color         string  `json:"color" validate:"omitempty,max=20" example:"#FF6B6B"`
	AccountNumber *string `json:"account_number,omitempty" validate:"omitempty,max=50" example:"1234567890"`
}

// UpdateAccountRequest représente la requête pour mettre à jour un compte
type UpdateAccountRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255" example:"Compte principal BNP"`
	Type     *string  `json:"type,omitempty" validate:"omitempty,oneof=checking savings mobile_money" example:"checking"`
	Balance  *float64 `json:"balance,omitempty" validate:"omitempty,gte=0" example:"2000.00"`
	Currency *string  `json:"currency,omitempty" validate:"omitempty,len=3" example:"EUR"`
}

// ==================== MOOD REQUESTS ====================

// CreateMoodRequest représente la requête pour créer une humeur
type CreateMoodRequest struct {
	Feeling string    `json:"feeling" validate:"required,min=1,max=50" example:"happy"`
	Note    string    `json:"note,omitempty" validate:"omitempty,max=500" example:"Journée productive aujourd'hui"`
	Date    time.Time `json:"date" validate:"required" example:"2024-01-15T00:00:00Z"`
}

// UpdateMoodRequest représente la requête pour mettre à jour une humeur
type UpdateMoodRequest struct {
	Feeling *string `json:"feeling,omitempty" validate:"omitempty,min=1,max=50" example:"excited"`
	Note    *string `json:"note,omitempty" validate:"omitempty,max=500" example:"Projet terminé avec succès"`
}

// ==================== NOTIFICATION REQUESTS ====================

// CreateNotificationRequest représente la requête pour créer une notification
type CreateNotificationRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=255" example:"Nouveau message"`
	Message string `json:"message" validate:"required,min=1,max=1000" example:"Vous avez reçu un nouveau message important"`
}

// ==================== MOTIVATION REQUESTS ====================

// CreateMotivationRequest représente la requête pour créer une motivation
type CreateMotivationRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000" example:"La persistance est le chemin du succès"`
	Type    string `json:"type" validate:"required,oneof=quote tip challenge" example:"quote"`
}

// UpdateMotivationRequest représente la requête pour mettre à jour une motivation
type UpdateMotivationRequest struct {
	Content *string `json:"content,omitempty" validate:"omitempty,min=1,max=1000" example:"La persistance mène toujours au succès"`
	Type    *string `json:"type,omitempty" validate:"omitempty,oneof=quote tip challenge" example:"quote"`
}

// ==================== LIFE NOTE REQUESTS ====================

// CreateLifeNoteRequest représente la requête pour créer une note de vie
type CreateLifeNoteRequest struct {
	Title         string     `json:"title" validate:"required,min=1,max=255" example:"Réflexions du jour"`
	Content       string     `json:"content" validate:"required,min=1,max=5000" example:"Aujourd'hui j'ai appris que..."`
	RelatedGoalID *uuid.UUID `json:"related_goal_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// UpdateLifeNoteRequest représente la requête pour mettre à jour une note de vie
type UpdateLifeNoteRequest struct {
	Title         *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255" example:"Nouvelles réflexions"`
	Content       *string    `json:"content,omitempty" validate:"omitempty,min=1,max=5000" example:"Mise à jour de mes pensées..."`
	RelatedGoalID *uuid.UUID `json:"related_goal_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// ==================== WEEKLY SUMMARY REQUESTS ====================

// CreateWeeklySummaryRequest représente la requête pour créer un résumé hebdomadaire
type CreateWeeklySummaryRequest struct {
	WeekStartDate       time.Time `json:"week_start_date" validate:"required" example:"2024-01-15T00:00:00Z"`
	TasksCompleted      int       `json:"tasks_completed" validate:"gte=0" example:"12"`
	TotalRevenue        float64   `json:"total_revenue" validate:"gte=0" example:"2500.00"`
	TotalExpense        float64   `json:"total_expense" validate:"gte=0" example:"1800.00"`
	Savings             float64   `json:"savings" validate:"gte=0" example:"700.00"`
	GoalProgressSummary string    `json:"goal_progress_summary" validate:"max=1000" example:"3 objectifs en cours, 1 atteint cette semaine"`
	MoodAverage         float64   `json:"mood_average" validate:"gte=0,lte=5" example:"4.2"`
	Notes               string    `json:"notes,omitempty" validate:"omitempty,max=2000" example:"Bonne semaine productive"`
}

// UpdateWeeklySummaryRequest représente la requête pour mettre à jour un résumé hebdomadaire
type UpdateWeeklySummaryRequest struct {
	TasksCompleted      *int     `json:"tasks_completed,omitempty" validate:"omitempty,gte=0" example:"15"`
	TotalRevenue        *float64 `json:"total_revenue,omitempty" validate:"omitempty,gte=0" example:"2600.00"`
	TotalExpense        *float64 `json:"total_expense,omitempty" validate:"omitempty,gte=0" example:"1900.00"`
	Savings             *float64 `json:"savings,omitempty" validate:"omitempty,gte=0" example:"700.00"`
	GoalProgressSummary *string  `json:"goal_progress_summary,omitempty" validate:"omitempty,max=1000" example:"4 objectifs en cours"`
	MoodAverage         *float64 `json:"mood_average,omitempty" validate:"omitempty,gte=0,lte=5" example:"4.5"`
	Notes               *string  `json:"notes,omitempty" validate:"omitempty,max=2000" example:"Excellente semaine"`
}

// ==================== BULK OPERATION REQUESTS ====================

// BulkDeleteRequest représente la requête pour supprimer plusieurs éléments
type BulkDeleteRequest struct {
	IDs []uuid.UUID `json:"ids" validate:"required,min=1,max=100,dive,uuid" example:"['123e4567-e89b-12d3-a456-426614174000','987fcdeb-51a2-43d1-b789-123456789abc']"`
}

// BulkUpdateRequest représente la requête pour mettre à jour plusieurs éléments
type BulkUpdateRequest struct {
	IDs    []uuid.UUID `json:"ids" validate:"required,min=1,max=100,dive,uuid" example:"['123e4567-e89b-12d3-a456-426614174000']"`
	Status *string     `json:"status,omitempty" validate:"omitempty,oneof=expired done incoming running" example:"done"`
}
