package database

import (
	"backend/pkg/logger"
	"fmt"

	"github.com/go-pg/pg/v10"
)

// RunMigrations exécute toutes les migrations de base de données
func RunMigrations(db *pg.DB, loggerInstance logger.Logger) error {
	loggerInstance.Info("Démarrage des migrations de base de données")

	// Migration 1: Table users (déjà existante, mise à jour si nécessaire)
	if err := createUsersTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table users: %w", err)
	}

	// Migration 2: Table categories
	if err := createCategoriesTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table categories: %w", err)
	}

	// Migration 3: Table tasks
	if err := createTasksTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table tasks: %w", err)
	}

	// Migration 4: Table daily_routines
	if err := createDailyRoutinesTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table daily_routines: %w", err)
	}

	// Migration 5: Table exotic_tasks
	if err := createExoticTasksTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table exotic_tasks: %w", err)
	}

	// Migration 5.1: Ajouter la colonne type à exotic_tasks si elle n'existe pas
	if err := addTypeColumnToExoticTasks(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur ajout colonne type à exotic_tasks: %w", err)
	}

	// Migration 6: Table revenues
	if err := createRevenuesTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table revenues: %w", err)
	}

	// Migration 7: Table expenses
	if err := createExpensesTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table expenses: %w", err)
	}

	// Migration 8: Table goals
	if err := createGoalsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table goals: %w", err)
	}

	// Migration 9: Table goal_progress
	if err := createGoalProgressTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table goal_progress: %w", err)
	}

	// Migration 10: Table accounts
	if err := createAccountsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table accounts: %w", err)
	}

	// Migration 11: Table transactions
	if err := createTransactionsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table transactions: %w", err)
	}

	// Migration 12: Table budgets
	if err := createBudgetsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table budgets: %w", err)
	}

	// Migration 13: Table saving_goals
	if err := createSavingGoalsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table saving_goals: %w", err)
	}

	// Migration 14: Table saving_strategies
	if err := createSavingStrategiesTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table saving_strategies: %w", err)
	}

	// Migration 15: Table motivations
	if err := createMotivationsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table motivations: %w", err)
	}

	// Migration 16: Table life_notes
	if err := createLifeNotesTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table life_notes: %w", err)
	}

	// Migration 17: Table reminders
	if err := createRemindersTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table reminders: %w", err)
	}

	// Migration 18: Table notifications
	if err := createNotificationsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table notifications: %w", err)
	}

	// Migration 19: Table habits
	if err := createHabitsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table habits: %w", err)
	}

	// Migration 20: Table moods
	if err := createMoodsTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table moods: %w", err)
	}

	// Migration 21: Table weekly_summaries
	if err := createWeeklySummariesTable(db, loggerInstance); err != nil {
		return fmt.Errorf("erreur création table weekly_summaries: %w", err)
	}

	loggerInstance.Info("Toutes les migrations ont été exécutées avec succès")
	return nil
}

// createUsersTable crée la table users
func createUsersTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		avatar TEXT DEFAULT '',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE
	);
	
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table users", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table users créée/mise à jour")
	return nil
}

// createCategoriesTable crée la table categories
func createCategoriesTable(db *pg.DB, loggerInstance logger.Logger) error {
	// Créer la table de base sans parent_id d'abord
	createQuery := `
	CREATE TABLE IF NOT EXISTS categories (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(100) NOT NULL,
		type VARCHAR(50) NOT NULL CHECK (type IN ('expense', 'revenue', 'task')),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(user_id, name, type)
	);
	`

	_, err := db.Exec(createQuery)
	if err != nil {
		loggerInstance.Error("Erreur création table categories", logger.Error(err))
		return err
	}

	// Ajouter la colonne parent_id si elle n'existe pas
	migrationQuery := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'categories' AND column_name = 'parent_id') THEN
			ALTER TABLE categories ADD COLUMN parent_id UUID REFERENCES categories(id) ON DELETE CASCADE;
		END IF;
	END $$;
	`

	_, err = db.Exec(migrationQuery)
	if err != nil {
		loggerInstance.Error("Erreur ajout colonne parent_id", logger.Error(err))
		return err
	}

	// Ajouter la colonne icon si elle n'existe pas
	iconMigrationQuery := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'categories' AND column_name = 'icon') THEN
			ALTER TABLE categories ADD COLUMN icon VARCHAR(50);
		END IF;
	END $$;
	`

	_, err = db.Exec(iconMigrationQuery)
	if err != nil {
		loggerInstance.Error("Erreur ajout colonne icon", logger.Error(err))
		return err
	}

	// Ajouter la colonne color si elle n'existe pas
	colorMigrationQuery := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'categories' AND column_name = 'color') THEN
			ALTER TABLE categories ADD COLUMN color VARCHAR(20);
		END IF;
	END $$;
	`

	_, err = db.Exec(colorMigrationQuery)
	if err != nil {
		loggerInstance.Error("Erreur ajout colonne color", logger.Error(err))
		return err
	}

	// Créer les index
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_categories_type ON categories(type);`,
		`CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);`,
	}

	for _, indexQuery := range indexQueries {
		_, err := db.Exec(indexQuery)
		if err != nil {
			loggerInstance.Error("Erreur création index categories", logger.Error(err))
			return err
		}
	}

	loggerInstance.Info("Table categories créée/mise à jour")
	return nil
}

// createTasksTable crée la table tasks
func createTasksTable(db *pg.DB, loggerInstance logger.Logger) error {
	// Créer la table si elle n'existe pas
	createQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		priority VARCHAR(20) CHECK (priority IN ('low', 'medium', 'high')),
		due_date TIMESTAMP WITH TIME ZONE,
		duration_planned INTEGER DEFAULT 0,
		duration_spent INTEGER DEFAULT 0,
		status VARCHAR(20) CHECK (status IN ('expired', 'done', 'incoming', 'running')),
		recurrence_rule TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`

	_, err := db.Exec(createQuery)
	if err != nil {
		loggerInstance.Error("Erreur création table tasks", logger.Error(err))
		return err
	}

	// Migrer la structure existante si nécessaire
	migrationQueries := []string{
		// Ajouter la colonne priority si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tasks' AND column_name = 'priority') THEN
				ALTER TABLE tasks ADD COLUMN priority VARCHAR(20) CHECK (priority IN ('low', 'medium', 'high'));
			END IF;
		END $$;`,

		// Ajouter la colonne duration_planned si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tasks' AND column_name = 'duration_planned') THEN
				ALTER TABLE tasks ADD COLUMN duration_planned INTEGER DEFAULT 0;
			END IF;
		END $$;`,

		// Ajouter la colonne duration_spent si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tasks' AND column_name = 'duration_spent') THEN
				ALTER TABLE tasks ADD COLUMN duration_spent INTEGER DEFAULT 0;
			END IF;
		END $$;`,

		// Ajouter la colonne recurrence_rule si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tasks' AND column_name = 'recurrence_rule') THEN
				ALTER TABLE tasks ADD COLUMN recurrence_rule TEXT;
			END IF;
		END $$;`,

		// Migrer is_completed vers status
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tasks' AND column_name = 'is_completed') 
			   AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tasks' AND column_name = 'status') THEN
				-- Ajouter la colonne status
				ALTER TABLE tasks ADD COLUMN status VARCHAR(20) CHECK (status IN ('expired', 'done', 'incoming', 'running'));
				
				-- Migrer les données de is_completed vers status
				UPDATE tasks SET status = CASE 
					WHEN is_completed = true THEN 'done'
					ELSE 'incoming'
				END;
				
				-- Supprimer l'ancienne colonne is_completed
				ALTER TABLE tasks DROP COLUMN is_completed;
			END IF;
		END $$;`,
	}

	// Exécuter les migrations
	for _, query := range migrationQueries {
		_, err := db.Exec(query)
		if err != nil {
			loggerInstance.Error("Erreur migration table tasks", logger.Error(err))
			return err
		}
	}

	// Créer les index s'ils n'existent pas
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_category_id ON tasks(category_id);`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);`,
	}

	for _, query := range indexQueries {
		_, err := db.Exec(query)
		if err != nil {
			loggerInstance.Error("Erreur création index tasks", logger.Error(err))
			return err
		}
	}

	loggerInstance.Info("Table tasks créée/mise à jour")
	return nil
}

// createDailyRoutinesTable crée la table daily_routines
func createDailyRoutinesTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS daily_routines (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		task_title VARCHAR(255) NOT NULL,
		scheduled_time TIME NOT NULL,
		is_completed BOOLEAN DEFAULT FALSE,
		frequency VARCHAR(50) NOT NULL DEFAULT 'daily' CHECK (frequency IN ('daily', 'weekly', 'monthly')),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_daily_routines_user_id ON daily_routines(user_id);
	CREATE INDEX IF NOT EXISTS idx_daily_routines_frequency ON daily_routines(frequency);
	CREATE INDEX IF NOT EXISTS idx_daily_routines_scheduled_time ON daily_routines(scheduled_time);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table daily_routines", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table daily_routines créée")
	return nil
}

// createExoticTasksTable crée la table exotic_tasks
func createExoticTasksTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS exotic_tasks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		is_completed BOOLEAN DEFAULT FALSE,
		type VARCHAR(50) NOT NULL DEFAULT 'exotic',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_exotic_tasks_user_id ON exotic_tasks(user_id);
	CREATE INDEX IF NOT EXISTS idx_exotic_tasks_is_completed ON exotic_tasks(is_completed);
	CREATE INDEX IF NOT EXISTS idx_exotic_tasks_type ON exotic_tasks(type);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table exotic_tasks", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table exotic_tasks créée")
	return nil
}

// createRevenuesTable crée la table revenues
func createRevenuesTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS revenues (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
		source VARCHAR(255) NOT NULL,
		date DATE NOT NULL,
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_revenues_user_id ON revenues(user_id);
	CREATE INDEX IF NOT EXISTS idx_revenues_date ON revenues(date);
	CREATE INDEX IF NOT EXISTS idx_revenues_source ON revenues(source);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table revenues", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table revenues créée")
	return nil
}

// createExpensesTable crée la table expenses
func createExpensesTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS expenses (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
		description VARCHAR(255) NOT NULL,
		category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
		date DATE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
	CREATE INDEX IF NOT EXISTS idx_expenses_date ON expenses(date);
	CREATE INDEX IF NOT EXISTS idx_expenses_category_id ON expenses(category_id);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table expenses", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table expenses créée")
	return nil
}

// createGoalsTable crée la table goals
func createGoalsTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS goals (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		target_date DATE,
		is_achieved BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_goals_user_id ON goals(user_id);
	CREATE INDEX IF NOT EXISTS idx_goals_is_achieved ON goals(is_achieved);
	CREATE INDEX IF NOT EXISTS idx_goals_target_date ON goals(target_date);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table goals", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table goals créée")
	return nil
}

// createGoalProgressTable crée la table goal_progress
func createGoalProgressTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS goal_progress (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		goal_id UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
		step_description VARCHAR(255) NOT NULL,
		progress_percentage DECIMAL(5,2) NOT NULL CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
		completed BOOLEAN DEFAULT FALSE,
		date DATE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_goal_progress_goal_id ON goal_progress(goal_id);
	CREATE INDEX IF NOT EXISTS idx_goal_progress_date ON goal_progress(date);
	CREATE INDEX IF NOT EXISTS idx_goal_progress_completed ON goal_progress(completed);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table goal_progress", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table goal_progress créée")
	return nil
}

// createSavingStrategiesTable crée la table saving_strategies
func createSavingStrategiesTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS saving_strategies (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		strategy_name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL CHECK (type IN ('percentage', 'fixed', 'goal_based')),
		amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
		frequency VARCHAR(50) NOT NULL CHECK (frequency IN ('weekly', 'monthly', 'yearly')),
		target_goal_id UUID REFERENCES goals(id) ON DELETE SET NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_saving_strategies_user_id ON saving_strategies(user_id);
	CREATE INDEX IF NOT EXISTS idx_saving_strategies_type ON saving_strategies(type);
	CREATE INDEX IF NOT EXISTS idx_saving_strategies_target_goal_id ON saving_strategies(target_goal_id);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table saving_strategies", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table saving_strategies créée")
	return nil
}

// createMotivationsTable crée la table motivations
func createMotivationsTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS motivations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		content TEXT NOT NULL,
		type VARCHAR(50) NOT NULL CHECK (type IN ('quote', 'tip', 'challenge')),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_motivations_user_id ON motivations(user_id);
	CREATE INDEX IF NOT EXISTS idx_motivations_type ON motivations(type);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table motivations", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table motivations créée")
	return nil
}

// createLifeNotesTable crée la table life_notes
func createLifeNotesTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS life_notes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		related_goal_id UUID REFERENCES goals(id) ON DELETE SET NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_life_notes_user_id ON life_notes(user_id);
	CREATE INDEX IF NOT EXISTS idx_life_notes_related_goal_id ON life_notes(related_goal_id);
	CREATE INDEX IF NOT EXISTS idx_life_notes_created_at ON life_notes(created_at);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table life_notes", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table life_notes créée")
	return nil
}

// createRemindersTable crée la table reminders
func createRemindersTable(db *pg.DB, loggerInstance logger.Logger) error {
	// Créer la table si elle n'existe pas
	createQuery := `
	CREATE TABLE IF NOT EXISTS reminders (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
		transac_id UUID REFERENCES transactions(id) ON DELETE CASCADE,
		message TEXT NOT NULL,
		trigger_at TIMESTAMP WITH TIME ZONE NOT NULL,
		is_sent BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`

	_, err := db.Exec(createQuery)
	if err != nil {
		loggerInstance.Error("Erreur création table reminders", logger.Error(err))
		return err
	}

	// Migrer la structure existante si nécessaire
	migrationQueries := []string{
		// Ajouter la colonne task_id si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'reminders' AND column_name = 'task_id') THEN
				ALTER TABLE reminders ADD COLUMN task_id UUID REFERENCES tasks(id) ON DELETE CASCADE;
			END IF;
		END $$;`,

		// Ajouter la colonne transac_id si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'reminders' AND column_name = 'transac_id') THEN
				ALTER TABLE reminders ADD COLUMN transac_id UUID REFERENCES transactions(id) ON DELETE CASCADE;
			END IF;
		END $$;`,

		// Ajouter la colonne is_sent si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'reminders' AND column_name = 'is_sent') THEN
				ALTER TABLE reminders ADD COLUMN is_sent BOOLEAN DEFAULT FALSE;
			END IF;
		END $$;`,

		// Migrer scheduled_at vers trigger_at
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'reminders' AND column_name = 'scheduled_at') 
			   AND NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'reminders' AND column_name = 'trigger_at') THEN
				-- Renommer scheduled_at vers trigger_at
				ALTER TABLE reminders RENAME COLUMN scheduled_at TO trigger_at;
			END IF;
		END $$;`,

		// Supprimer les anciennes colonnes si elles existent
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'reminders' AND column_name = 'target_type') THEN
				ALTER TABLE reminders DROP COLUMN target_type;
			END IF;
		END $$;`,

		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'reminders' AND column_name = 'target_id') THEN
				ALTER TABLE reminders DROP COLUMN target_id;
			END IF;
		END $$;`,
	}

	// Exécuter les migrations
	for _, query := range migrationQueries {
		_, err := db.Exec(query)
		if err != nil {
			loggerInstance.Error("Erreur migration table reminders", logger.Error(err))
			return err
		}
	}

	// Créer les index s'ils n'existent pas
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_reminders_user_id ON reminders(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_trigger_at ON reminders(trigger_at);`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_task_id ON reminders(task_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_transac_id ON reminders(transac_id);`,
		`CREATE INDEX IF NOT EXISTS idx_reminders_is_sent ON reminders(is_sent);`,
	}

	for _, query := range indexQueries {
		_, err := db.Exec(query)
		if err != nil {
			loggerInstance.Error("Erreur création index reminders", logger.Error(err))
			return err
		}
	}

	loggerInstance.Info("Table reminders créée/mise à jour")
	return nil
}

// createNotificationsTable crée la table notifications
func createNotificationsTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		is_read BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
	CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table notifications", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table notifications créée")
	return nil
}

// createHabitsTable crée la table habits
func createHabitsTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS habits (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		frequency VARCHAR(50) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly')),
		target_time TIME,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_habits_user_id ON habits(user_id);
	CREATE INDEX IF NOT EXISTS idx_habits_is_active ON habits(is_active);
	CREATE INDEX IF NOT EXISTS idx_habits_frequency ON habits(frequency);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table habits", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table habits créée")
	return nil
}

// createMoodsTable crée la table moods
func createMoodsTable(db *pg.DB, loggerInstance logger.Logger) error {
	// Créer la table si elle n'existe pas
	createQuery := `
	CREATE TABLE IF NOT EXISTS moods (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		date DATE NOT NULL,
		feeling VARCHAR(50) NOT NULL,
		note TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(user_id, date)
	);
	`

	_, err := db.Exec(createQuery)
	if err != nil {
		loggerInstance.Error("Erreur création table moods", logger.Error(err))
		return err
	}

	// Migrer la structure existante si nécessaire
	migrationQueries := []string{
		// Ajouter la colonne feeling si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'moods' AND column_name = 'feeling') THEN
				ALTER TABLE moods ADD COLUMN feeling VARCHAR(50);
			END IF;
		END $$;`,

		// Ajouter la colonne note si elle n'existe pas
		`DO $$ 
		BEGIN 
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'moods' AND column_name = 'note') THEN
				ALTER TABLE moods ADD COLUMN note TEXT;
			END IF;
		END $$;`,

		// Migrer level vers feeling
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'moods' AND column_name = 'level') 
			   AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'moods' AND column_name = 'feeling') THEN
				-- Migrer les données de level vers feeling
				UPDATE moods SET feeling = CASE 
					WHEN level = 1 THEN 'very_bad'
					WHEN level = 2 THEN 'bad'
					WHEN level = 3 THEN 'neutral'
					WHEN level = 4 THEN 'good'
					WHEN level = 5 THEN 'very_good'
					ELSE 'neutral'
				END WHERE feeling IS NULL;
				
				-- Rendre feeling NOT NULL après la migration
				ALTER TABLE moods ALTER COLUMN feeling SET NOT NULL;
			END IF;
		END $$;`,

		// Migrer description vers note
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'moods' AND column_name = 'description') 
			   AND EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'moods' AND column_name = 'note') THEN
				-- Migrer les données de description vers note
				UPDATE moods SET note = description WHERE note IS NULL;
			END IF;
		END $$;`,

		// Supprimer les anciennes colonnes si elles existent
		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'moods' AND column_name = 'level') THEN
				ALTER TABLE moods DROP COLUMN level;
			END IF;
		END $$;`,

		`DO $$ 
		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'moods' AND column_name = 'description') THEN
				ALTER TABLE moods DROP COLUMN description;
			END IF;
		END $$;`,
	}

	// Exécuter les migrations
	for _, query := range migrationQueries {
		_, err := db.Exec(query)
		if err != nil {
			loggerInstance.Error("Erreur migration table moods", logger.Error(err))
			return err
		}
	}

	// Créer les index s'ils n'existent pas
	indexQueries := []string{
		`CREATE INDEX IF NOT EXISTS idx_moods_user_id ON moods(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_moods_date ON moods(date);`,
		`CREATE INDEX IF NOT EXISTS idx_moods_feeling ON moods(feeling);`,
	}

	for _, query := range indexQueries {
		_, err := db.Exec(query)
		if err != nil {
			loggerInstance.Error("Erreur création index moods", logger.Error(err))
			return err
		}
	}

	loggerInstance.Info("Table moods créée/mise à jour")
	return nil
}

// createWeeklySummariesTable crée la table weekly_summaries
func createWeeklySummariesTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS weekly_summaries (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		week_start_date DATE NOT NULL,
		tasks_completed INTEGER DEFAULT 0,
		total_revenue DECIMAL(10,2) DEFAULT 0,
		total_expense DECIMAL(10,2) DEFAULT 0,
		savings DECIMAL(10,2) DEFAULT 0,
		goal_progress_summary TEXT,
		mood_average DECIMAL(3,2) DEFAULT 0,
		notes TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(user_id, week_start_date)
	);
	
	CREATE INDEX IF NOT EXISTS idx_weekly_summaries_user_id ON weekly_summaries(user_id);
	CREATE INDEX IF NOT EXISTS idx_weekly_summaries_week_start_date ON weekly_summaries(week_start_date);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table weekly_summaries", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table weekly_summaries créée")
	return nil
}

// addTypeColumnToExoticTasks ajoute la colonne type à la table exotic_tasks si elle n'existe pas
func addTypeColumnToExoticTasks(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'exotic_tasks' AND column_name = 'type'
		) THEN
			ALTER TABLE exotic_tasks ADD COLUMN type VARCHAR(50) NOT NULL DEFAULT 'exotic';
		END IF;
	END $$;
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur ajout colonne type à exotic_tasks", logger.Error(err))
		return err
	}

	loggerInstance.Info("Colonne type ajoutée à exotic_tasks (si nécessaire)")
	return nil
}

// createAccountsTable crée la table accounts
func createAccountsTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL CHECK (type IN ('checking', 'savings', 'mobile_money', 'cash', 'bank')),
		balance DECIMAL(10,2) NOT NULL DEFAULT 0,
		currency VARCHAR(3) NOT NULL DEFAULT 'USD',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
	CREATE INDEX IF NOT EXISTS idx_accounts_type ON accounts(type);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table accounts", logger.Error(err))
		return err
	}

	// Ajouter la colonne icon si elle n'existe pas
	iconMigrationQuery := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'accounts' AND column_name = 'icon') THEN
			ALTER TABLE accounts ADD COLUMN icon VARCHAR(50);
		END IF;
	END $$;
	`

	_, err = db.Exec(iconMigrationQuery)
	if err != nil {
		loggerInstance.Error("Erreur ajout colonne icon", logger.Error(err))
		return err
	}

	// Ajouter la colonne color si elle n'existe pas
	colorMigrationQuery := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'accounts' AND column_name = 'color') THEN
			ALTER TABLE accounts ADD COLUMN color VARCHAR(20);
		END IF;
	END $$;
	`

	_, err = db.Exec(colorMigrationQuery)
	if err != nil {
		loggerInstance.Error("Erreur ajout colonne color", logger.Error(err))
		return err
	}

	// Ajouter la colonne account_number si elle n'existe pas
	accountNumberMigrationQuery := `
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'accounts' AND column_name = 'account_number') THEN
			ALTER TABLE accounts ADD COLUMN account_number VARCHAR(50);
		END IF;
	END $$;
	`

	_, err = db.Exec(accountNumberMigrationQuery)
	if err != nil {
		loggerInstance.Error("Erreur ajout colonne account_number", logger.Error(err))
		return err
	}

	// Mettre à jour la contrainte CHECK pour inclure les nouveaux types
	updateCheckConstraintQuery := `
	DO $$ 
	BEGIN 
		-- Supprimer l'ancienne contrainte si elle existe
		IF EXISTS (SELECT 1 FROM information_schema.check_constraints WHERE constraint_name = 'accounts_type_check') THEN
			ALTER TABLE accounts DROP CONSTRAINT accounts_type_check;
		END IF;
		
		-- Ajouter la nouvelle contrainte avec tous les types
		ALTER TABLE accounts ADD CONSTRAINT accounts_type_check CHECK (type IN ('checking', 'savings', 'mobile_money', 'cash', 'bank'));
	END $$;
	`

	_, err = db.Exec(updateCheckConstraintQuery)
	if err != nil {
		loggerInstance.Error("Erreur mise à jour contrainte CHECK type", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table accounts créée")
	return nil
}

// createTransactionsTable crée la table transactions
func createTransactionsTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
		category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
		type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
		amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
		description VARCHAR(255) NOT NULL,
		date DATE NOT NULL,
		recurring BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
	CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
	CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table transactions", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table transactions créée")
	return nil
}

// createBudgetsTable crée la table budgets
func createBudgetsTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS budgets (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
		amount_planned DECIMAL(10,2) NOT NULL CHECK (amount_planned > 0),
		amount_spent DECIMAL(10,2) NOT NULL DEFAULT 0,
		month INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
		year INTEGER NOT NULL CHECK (year >= 2020),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(user_id, category_id, month, year)
	);
	
	CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets(user_id);
	CREATE INDEX IF NOT EXISTS idx_budgets_category_id ON budgets(category_id);
	CREATE INDEX IF NOT EXISTS idx_budgets_month_year ON budgets(month, year);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table budgets", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table budgets créée")
	return nil
}

// createSavingGoalsTable crée la table saving_goals
func createSavingGoalsTable(db *pg.DB, loggerInstance logger.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS saving_goals (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		target_amount DECIMAL(10,2) NOT NULL CHECK (target_amount > 0),
		current_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
		deadline DATE,
		is_achieved BOOLEAN DEFAULT FALSE,
		frequency VARCHAR(20) CHECK (frequency IN ('weekly', 'monthly', 'yearly')),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_saving_goals_user_id ON saving_goals(user_id);
	CREATE INDEX IF NOT EXISTS idx_saving_goals_is_achieved ON saving_goals(is_achieved);
	CREATE INDEX IF NOT EXISTS idx_saving_goals_deadline ON saving_goals(deadline);
	`

	_, err := db.Exec(query)
	if err != nil {
		loggerInstance.Error("Erreur création table saving_goals", logger.Error(err))
		return err
	}

	loggerInstance.Info("Table saving_goals créée")
	return nil
}
