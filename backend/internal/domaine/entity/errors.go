package entity

import "errors"

// Erreurs du domaine User
var (
	ErrUserNotFound      = errors.New("utilisateur non trouvé")
	ErrUserAlreadyExists = errors.New("utilisateur déjà existant")
	ErrInvalidUserData   = errors.New("données utilisateur invalides")
	ErrUserInactive      = errors.New("utilisateur inactif")
)

// Erreurs du domaine File
var (
	ErrFileNotFound     = errors.New("fichier non trouvé")
	ErrFileUploadFailed = errors.New("échec de l'upload du fichier")
	ErrInvalidFileType  = errors.New("type de fichier non autorisé")
	ErrFileTooLarge     = errors.New("fichier trop volumineux")
)

// Erreurs génériques
var (
	ErrInternalServer = errors.New("erreur interne du serveur")
	ErrValidation     = errors.New("erreur de validation")
	ErrUnauthorized   = errors.New("non autorisé")
	ErrForbidden      = errors.New("accès interdit")
	ErrBadRequest     = errors.New("requête invalide")
	ErrAccessDenied   = errors.New("accès refusé")
)

// Erreurs du domaine Task
var (
	ErrTaskNotFound    = errors.New("tâche non trouvée")
	ErrInvalidTaskData = errors.New("données de tâche invalides")
)

// Erreurs du domaine DailyRoutine
var (
	ErrDailyRoutineNotFound    = errors.New("routine quotidienne non trouvée")
	ErrInvalidDailyRoutineData = errors.New("données de routine quotidienne invalides")
)

// Erreurs du domaine Goal
var (
	ErrGoalNotFound    = errors.New("objectif non trouvé")
	ErrInvalidGoalData = errors.New("données d'objectif invalides")
)

// Erreurs du domaine ExoticTask
var (
	ErrExoticTaskNotFound    = errors.New("tâche exotique non trouvée")
	ErrInvalidExoticTaskData = errors.New("données de tâche exotique invalides")
)

// Erreurs du domaine Expense
var (
	ErrExpenseNotFound    = errors.New("dépense non trouvée")
	ErrInvalidExpenseData = errors.New("données de dépense invalides")
)

// Erreurs du domaine Revenue
var (
	ErrRevenueNotFound    = errors.New("revenu non trouvé")
	ErrInvalidRevenueData = errors.New("données de revenu invalides")
)

// Erreurs du domaine SavingStrategy
var (
	ErrSavingStrategyNotFound    = errors.New("stratégie d'épargne non trouvée")
	ErrInvalidSavingStrategyData = errors.New("données de stratégie d'épargne invalides")
)
