package repository

import (
	"context"
	"time"

	"backend/internal/domaine/entity"

	"github.com/google/uuid"
)

// USER
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, query *entity.UserQuery) (*entity.UserListResponse, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Search(ctx context.Context, searchTerm string, limit int) ([]*entity.User, error)
}

// TASK
type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetCompletedByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error)
	GetPendingByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error)
	GetIncomingByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error)
}

// HABIT (basé sur Task avec Recurrence)
type HabitRepository interface {
	Create(ctx context.Context, habit *entity.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error)
	Update(ctx context.Context, habit *entity.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Task, error)
}

// ACCOUNT
type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (float64, error)
}

// TRANSACTION
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Transaction, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCategoryID(ctx context.Context, userID uuid.UUID, categoryID *uuid.UUID) ([]*entity.Transaction, error)
	GetByAccountID(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]*entity.Transaction, error)
	GetByAccountIDWithCategoryDetails(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID) ([]*entity.TransactionWithDetails, error)
	GetBySavingGoalID(ctx context.Context, userID uuid.UUID, savingGoalID uuid.UUID) ([]*entity.Transaction, error)
	GetTransfers(ctx context.Context, userID uuid.UUID) ([]*entity.Transaction, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*entity.Transaction, error)
	GetByType(ctx context.Context, userID uuid.UUID, txType string) ([]*entity.Transaction, error)
	GetTotalByUserID(ctx context.Context, userID uuid.UUID) (float64, error)
	GetAllTransactionsByUserIDAndAccountID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]*entity.Transaction, error)
	//get All transaction by saving goal id
	GetAllTransactionsBySavingGoalID(ctx context.Context, userID uuid.UUID, savingGoalID uuid.UUID) ([]*entity.Transaction, error)
}

// CATEGORY
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entity.Category, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]*entity.Category, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.Category, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
}

type PreferencesRepository interface {
	Create(ctx context.Context, preferences *entity.UserPreferences) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserPreferences, error)
	Update(ctx context.Context, preferences *entity.UserPreferences) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

// BUDGET
type BudgetRepository interface {
	Create(ctx context.Context, budget *entity.Budget) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Budget, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Budget, error)
	Update(ctx context.Context, budget *entity.Budget) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]*entity.Budget, error)
	GetByPeriod(ctx context.Context, userID uuid.UUID, period string) ([]*entity.Budget, error)
	GetByName(ctx context.Context, userID uuid.UUID, name string) ([]*entity.Budget, error)
}

// SAVING GOAL
type SavingGoalRepository interface {
	Create(ctx context.Context, goal *entity.SavingGoal) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.SavingGoal, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavingGoal, error)
	Update(ctx context.Context, goal *entity.SavingGoal) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAchievedByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavingGoal, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SavingGoal, error)
	GetByAccountID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]*entity.SavingGoal, error)
	GetByFrequency(ctx context.Context, userID uuid.UUID, frequency string) ([]*entity.SavingGoal, error)
	GetAllSavingGoalsByUserIDAndAccountID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]*entity.SavingGoal, error)
}

// REMINDER
type ReminderRepository interface {
	Create(ctx context.Context, reminder *entity.Reminder) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Reminder, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Reminder, error)
	Update(ctx context.Context, reminder *entity.Reminder) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByTaskID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) ([]*entity.Reminder, error)
	GetByTransactionID(ctx context.Context, userID uuid.UUID, transacID uuid.UUID) ([]*entity.Reminder, error)
	GetPending(ctx context.Context, userID uuid.UUID) ([]*entity.Reminder, error)
	GetDue(ctx context.Context, now time.Time) ([]*entity.Reminder, error)
	MarkAsSent(ctx context.Context, id uuid.UUID) error
}

// NOTIFICATION
type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
	Update(ctx context.Context, notification *entity.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUnread(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
	GetRead(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	CountUnread(ctx context.Context, userID uuid.UUID) (int, error)
}

// MOOD
type MoodRepository interface {
	Create(ctx context.Context, mood *entity.Mood) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Mood, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Mood, error)
	Update(ctx context.Context, mood *entity.Mood) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByDate(ctx context.Context, userID uuid.UUID, date time.Time) (*entity.Mood, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*entity.Mood, error)
	GetByFeeling(ctx context.Context, userID uuid.UUID, feeling string) ([]*entity.Mood, error)
}

// MOTIVATION
type MotivationRepository interface {
	Create(ctx context.Context, motivation *entity.Motivation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Motivation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Motivation, error)
	Update(ctx context.Context, motivation *entity.Motivation) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByType(ctx context.Context, userID uuid.UUID, motivationType string) ([]*entity.Motivation, error)
	GetRandom(ctx context.Context, userID uuid.UUID) (*entity.Motivation, error)
	GetRandomByType(ctx context.Context, userID uuid.UUID, motivationType string) (*entity.Motivation, error)
	Search(ctx context.Context, userID uuid.UUID, query string) ([]*entity.Motivation, error)
}

// LIFE NOTE
type LifeNoteRepository interface {
	Create(ctx context.Context, note *entity.LifeNote) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.LifeNote, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.LifeNote, error)
	Update(ctx context.Context, note *entity.LifeNote) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByGoalID(ctx context.Context, userID uuid.UUID, goalID uuid.UUID) ([]*entity.LifeNote, error)
	GetRecent(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.LifeNote, error)
	SearchByTitle(ctx context.Context, userID uuid.UUID, query string) ([]*entity.LifeNote, error)
	SearchByContent(ctx context.Context, userID uuid.UUID, query string) ([]*entity.LifeNote, error)
	Search(ctx context.Context, userID uuid.UUID, query string) ([]*entity.LifeNote, error)
}

// WEEKLY SUMMARY
type WeeklySummaryRepository interface {
	Create(ctx context.Context, summary *entity.WeeklySummary) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.WeeklySummary, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.WeeklySummary, error)
	Update(ctx context.Context, summary *entity.WeeklySummary) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByWeek(ctx context.Context, userID uuid.UUID, weekStartDate time.Time) (*entity.WeeklySummary, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*entity.WeeklySummary, error)
	GetLatest(ctx context.Context, userID uuid.UUID) (*entity.WeeklySummary, error)
	Upsert(ctx context.Context, summary *entity.WeeklySummary) error
}
