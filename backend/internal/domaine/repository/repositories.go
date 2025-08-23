package repository

import (
	"context"

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
	GetByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]*entity.Transaction, error)
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*entity.Transaction, error)
	GetByType(ctx context.Context, userID uuid.UUID, txType string) ([]*entity.Transaction, error)
	GetTotalByUserID(ctx context.Context, userID uuid.UUID) (float64, error)
}

// CATEGORY
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]*entity.Category, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]*entity.Category, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]*entity.Category, error)
}

// BUDGET
type BudgetRepository interface {
	Create(ctx context.Context, budget *entity.Budget) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Budget, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Budget, error)
	Update(ctx context.Context, budget *entity.Budget) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]*entity.Budget, error)
	GetByMonth(ctx context.Context, userID uuid.UUID, month, year int) ([]*entity.Budget, error)
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
}

// REMINDER
type ReminderRepository interface {
	Create(ctx context.Context, reminder *entity.Reminder) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Reminder, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Reminder, error)
	Update(ctx context.Context, reminder *entity.Reminder) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUpcomingByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Reminder, error)
}

// NOTIFICATION
type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
	Update(ctx context.Context, notification *entity.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	GetUnreadByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
}

// MOOD
type MoodRepository interface {
	Create(ctx context.Context, mood *entity.Mood) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Mood, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Mood, error)
	Update(ctx context.Context, mood *entity.Mood) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]*entity.Mood, error)
	GetAverageByUserID(ctx context.Context, userID uuid.UUID) (float64, error)
}

// MOTIVATION
type MotivationRepository interface {
	Create(ctx context.Context, motivation *entity.Motivation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Motivation, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Motivation, error)
	Update(ctx context.Context, motivation *entity.Motivation) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetRandomByType(ctx context.Context, userID uuid.UUID, motivationType string) (*entity.Motivation, error)
}

// LIFE NOTE
type LifeNoteRepository interface {
	Create(ctx context.Context, note *entity.LifeNote) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.LifeNote, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.LifeNote, error)
	Update(ctx context.Context, note *entity.LifeNote) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByGoalID(ctx context.Context, goalID uuid.UUID) ([]*entity.LifeNote, error)
}

// WEEKLY SUMMARY
type WeeklySummaryRepository interface {
	Create(ctx context.Context, summary *entity.WeeklySummary) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.WeeklySummary, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.WeeklySummary, error)
	Update(ctx context.Context, summary *entity.WeeklySummary) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.WeeklySummary, error)
}
