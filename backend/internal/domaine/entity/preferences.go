package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserPreferences représente les préférences utilisateur du wizard
type UserPreferences struct {
	tableName struct{} `pg:"preferences"`

	ID        uuid.UUID          `json:"id" db:"id"`
	UserID    uuid.UUID          `json:"user_id" db:"user_id"`
	Income    IncomePreferences  `json:"income" db:"income"`
	Expenses  ExpensePreferences `json:"expenses" db:"expenses"`
	Goals     GoalPreferences    `json:"goals" db:"goals"`
	Habits    HabitPreferences   `json:"habits" db:"habits"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" db:"updated_at"`
}

// IncomePreferences représente les préférences de revenus
type IncomePreferences struct {
	Sources      []string `json:"sources" db:"sources"`
	MonthlyTotal float64  `json:"monthly_total" db:"monthly_total"`
	Accounts     []string `json:"accounts" db:"accounts"`
	HasDebt      bool     `json:"has_debt" db:"has_debt"`
	DebtAmount   float64  `json:"debt_amount" db:"debt_amount"`
}

// ExpensePreferences représente les préférences de dépenses
type ExpensePreferences struct {
	TopCategories []string `json:"top_categories" db:"top_categories"`
	Food          float64  `json:"food" db:"food"`
	Transport     float64  `json:"transport" db:"transport"`
	Housing       float64  `json:"housing" db:"housing"`
	Subscriptions float64  `json:"subscriptions" db:"subscriptions"`
	AlertsEnabled bool     `json:"alerts_enabled" db:"alerts_enabled"`
	AutoBudget    bool     `json:"auto_budget" db:"auto_budget"`
}

// GoalPreferences représente les préférences d'objectifs
type GoalPreferences struct {
	MainGoal      string  `json:"main_goal" db:"main_goal"`
	SecondaryGoal string  `json:"secondary_goal" db:"secondary_goal"`
	SavingsTarget float64 `json:"savings_target" db:"savings_target"`
	Deadline      string  `json:"deadline" db:"deadline"`
	AdviceEnabled bool    `json:"advice_enabled" db:"advice_enabled"`
}

// HabitPreferences représente les préférences d'habitudes
type HabitPreferences struct {
	PlanningTime   string `json:"planning_time" db:"planning_time"`
	DailyFocusTime string `json:"daily_focus_time" db:"daily_focus_time"`
	CustomHabit    string `json:"custom_habit" db:"custom_habit"`
	SummaryType    string `json:"summary_type" db:"summary_type"`
}

// CreatePreferencesRequest représente la requête pour créer des préférences
type CreatePreferencesRequest struct {
	Income   IncomePreferences  `json:"income" validate:"required"`
	Expenses ExpensePreferences `json:"expenses" validate:"required"`
	Goals    GoalPreferences    `json:"goals" validate:"required"`
	Habits   HabitPreferences   `json:"habits" validate:"required"`
}

// UpdatePreferencesRequest représente la requête pour mettre à jour des préférences
type UpdatePreferencesRequest struct {
	Income   *IncomePreferences  `json:"income,omitempty"`
	Expenses *ExpensePreferences `json:"expenses,omitempty"`
	Goals    *GoalPreferences    `json:"goals,omitempty"`
	Habits   *HabitPreferences   `json:"habits,omitempty"`
}
