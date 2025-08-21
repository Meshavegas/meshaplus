package entity

import (
	"time"

	"github.com/google/uuid"
)

// ==================== TASK RESPONSES ====================

// TaskListResponse représente une liste paginée de tâches
type TaskListResponse struct {
	Tasks      []Task `json:"tasks"`
	Total      int64  `json:"total" example:"100"`
	Page       int    `json:"page" example:"1"`
	PageSize   int    `json:"page_size" example:"10"`
	TotalPages int    `json:"total_pages" example:"10"`
}

// TaskStatsResponse représente les statistiques des tâches
type TaskStatsResponse struct {
	TotalTasks     int64   `json:"total_tasks" example:"50"`
	CompletedTasks int64   `json:"completed_tasks" example:"30"`
	PendingTasks   int64   `json:"pending_tasks" example:"20"`
	CompletionRate float64 `json:"completion_rate" example:"60.0"`
}

type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	AccountID   uuid.UUID `json:"account_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Type        string    `json:"type"` // income, expense
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Recurring   bool      `json:"recurring"`
}

type AccountResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	Balance  float64   `json:"balance"`
	Currency string    `json:"currency"`
}

type BudgetResponse struct {
	ID            uuid.UUID `json:"id"`
	CategoryID    uuid.UUID `json:"category_id"`
	AmountPlanned float64   `json:"amount_planned"`
	AmountSpent   float64   `json:"amount_spent"`
	Month         int       `json:"month"`
	Year          int       `json:"year"`
	Progress      float64   `json:"progress"` // AmountSpent / AmountPlanned * 100
}

type SavingGoalResponse struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	TargetAmount  float64    `json:"target_amount"`
	CurrentAmount float64    `json:"current_amount"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	IsAchieved    bool       `json:"is_achieved"`
	Progress      float64    `json:"progress"` // CurrentAmount / TargetAmount * 100
}

type CategoryResponse struct {
	ID       uuid.UUID          `json:"id"`
	Name     string             `json:"name"`
	Type     string             `json:"type"`
	ParentID *uuid.UUID         `json:"parent_id,omitempty"`
	Children []CategoryResponse `json:"children,omitempty"`
}

type ReminderResponse struct {
	ID        uuid.UUID  `json:"id"`
	TaskID    *uuid.UUID `json:"task_id,omitempty"`
	TransacID *uuid.UUID `json:"transac_id,omitempty"`
	Message   string     `json:"message"`
	TriggerAt time.Time  `json:"trigger_at"`
	IsSent    bool       `json:"is_sent"`
}

type NotificationResponse struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	IsRead  bool      `json:"is_read"`
}

type MotivationResponse struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
	Type    string    `json:"type"` // quote, tip, challenge
}

type LifeNoteResponse struct {
	ID            uuid.UUID  `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	RelatedGoalID *uuid.UUID `json:"related_goal_id,omitempty"`
}

// ==================== DASHBOARD RESPONSES ====================

// DashboardResponse représente le tableau de bord complet
type DashboardResponse struct {
	UserID          uuid.UUID           `json:"user_id"`
	TasksOverview   TasksOverview       `json:"tasks_overview"`
	FinanceOverview FinanceOverview     `json:"finance_overview"`
	MoodToday       *MoodResponse       `json:"mood_today,omitempty"`
	Motivation      *MotivationResponse `json:"motivation,omitempty"`
}

type MoodResponse struct {
	MoodAverage string `json:"mood_average"`
}

type TasksOverview struct {
	TotalTasks int `json:"total_tasks"`
	Completed  int `json:"completed"`
	Pending    int `json:"pending"`
	Incoming   int `json:"incoming"`
	Running    int `json:"running"`
}

type FinanceOverview struct {
	TotalRevenue  float64 `json:"total_revenue"`
	TotalExpense  float64 `json:"total_expense"`
	Balance       float64 `json:"balance"`
	MonthlyBudget float64 `json:"monthly_budget"`
	Savings       float64 `json:"savings"`
}

type WeeklySummaryResponse struct {
	WeekStartDate  time.Time `json:"week_start_date"`
	TasksCompleted int       `json:"tasks_completed"`
	TotalRevenue   float64   `json:"total_revenue"`
	TotalExpense   float64   `json:"total_expense"`
	Savings        float64   `json:"savings"`
	GoalProgress   string    `json:"goal_progress"`
	MoodAverage    float64   `json:"mood_average"`
	Notes          string    `json:"notes"`
}

type FinanceReportResponse struct {
	Month      int                   `json:"month"`
	Year       int                   `json:"year"`
	ByCategory []CategoryFinanceStat `json:"by_category"`
}

type CategoryFinanceStat struct {
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"` // income, expense
	Planned    float64   `json:"planned"`
	Spent      float64   `json:"spent"`
}
