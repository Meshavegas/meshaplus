package entity

import (
	"time"

	"github.com/google/uuid"
)

// Task représente une tâche utilisateur
type Task struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty" db:"category_id"`
	Title           string     `json:"title" db:"title"`
	Description     string     `json:"description" db:"description"`
	Priority        string     `json:"priority" db:"priority"` // low, medium, high
	DueDate         *time.Time `json:"due_date,omitempty" db:"due_date"`
	DurationPlanned int        `json:"duration_planned" db:"duration_planned"`
	DurationSpent   int        `json:"duration_spent" db:"duration_spent"`
	Status          string     `json:"status" db:"status"`                             // expired, done, icommming, runninng
	RecurrenceRule  *string    `json:"recurrence_rule,omitempty" db:"recurrence_rule"` // daily, weekly, cron-like
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// Mood permet de suivre l'état d'esprit de l'utilisateur
type Mood struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Date      time.Time `json:"date" db:"date"`
	Feeling   string    `json:"feeling" db:"feeling"` // happy, neutral, sad, etc.
	Note      string    `json:"note,omitempty" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Account struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`       // ex: "Bancaire", "Cash"
	Type      string    `json:"type" db:"type"`       // checking, savings, mobile_money
	Balance   float64   `json:"balance" db:"balance"` // solde actuel
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	AccountID   uuid.UUID  `json:"account_id" db:"account_id"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty" db:"category_id"`
	Type        string     `json:"type" db:"type"` // income, expense
	Amount      float64    `json:"amount" db:"amount"`
	Description string     `json:"description" db:"description"`
	Date        time.Time  `json:"date" db:"date"`
	Recurring   bool       `json:"recurring" db:"recurring"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Reminder représente un rappel ou notification intelligente
type Reminder struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	TaskID    *uuid.UUID `json:"task_id,omitempty" db:"task_id"`
	TransacID *uuid.UUID `json:"transac_id,omitempty" db:"transac_id"`
	Message   string     `json:"message" db:"message"`
	TriggerAt time.Time  `json:"trigger_at" db:"trigger_at"`
	IsSent    bool       `json:"is_sent" db:"is_sent"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// Category représente une catégorie hiérarchique
type Category struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	Name      string     `json:"name" db:"name"`
	Type      string     `json:"type" db:"type"`                     // expense, revenue, task
	ParentID  *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"` // sous-catégorie
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// SavingGoal représente un objectif d'épargne
type SavingGoal struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	Title         string     `json:"title" db:"title"`
	TargetAmount  float64    `json:"target_amount" db:"target_amount"`
	CurrentAmount float64    `json:"current_amount" db:"current_amount"`
	Deadline      *time.Time `json:"deadline,omitempty" db:"deadline"`
	IsAchieved    bool       `json:"is_achieved" db:"is_achieved"`
	Frequency     string     `json:"frequency" db:"frequency"` // weekly, monthly, yearly, cron
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// Budget représente un budget mensuel ou annuel
type Budget struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	CategoryID    uuid.UUID `json:"category_id" db:"category_id"`
	AmountPlanned float64   `json:"amount_planned" db:"amount_planned"`
	AmountSpent   float64   `json:"amount_spent" db:"amount_spent"`
	Month         int       `json:"month" db:"month"`
	Year          int       `json:"year" db:"year"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Motivation représente une motivation
type Motivation struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Type      string    `json:"type" db:"type"` // quote, tip, challenge
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// LifeNote représente une note de vie
type LifeNote struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	Title         string     `json:"title" db:"title"`
	Content       string     `json:"content" db:"content"`
	RelatedGoalID *uuid.UUID `json:"related_goal_id,omitempty" db:"related_goal_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// Notification représente une notification
type Notification struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Message   string    `json:"message" db:"message"`
	IsRead    bool      `json:"is_read" db:"is_read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// WeeklySummary représente un résumé hebdomadaire
type WeeklySummary struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	UserID              uuid.UUID `json:"user_id" db:"user_id"`
	WeekStartDate       time.Time `json:"week_start_date" db:"week_start_date"`
	TasksCompleted      int       `json:"tasks_completed" db:"tasks_completed"`
	TotalRevenue        float64   `json:"total_revenue" db:"total_revenue"`
	TotalExpense        float64   `json:"total_expense" db:"total_expense"`
	Savings             float64   `json:"savings" db:"savings"`
	GoalProgressSummary string    `json:"goal_progress_summary" db:"goal_progress_summary"`
	MoodAverage         float64   `json:"mood_average" db:"mood_average"`
	Notes               string    `json:"notes" db:"notes"`
}
