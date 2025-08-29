package handler

import (
	"backend/internal/domaine/entity"
	"backend/internal/service"
	"backend/pkg/logger"
	"backend/pkg/response"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// FinanceDashboardHandler gère les requêtes pour le tableau de bord financier
type FinanceDashboardHandler struct {
	accountService     *service.AccountService
	transactionService *service.TransactionService
	budgetService      *service.BudgetService
	savingGoalService  *service.SavingGoalService
	logger             logger.Logger
}

// NewFinanceDashboardHandler crée une nouvelle instance de FinanceDashboardHandler
func NewFinanceDashboardHandler(
	accountService *service.AccountService,
	transactionService *service.TransactionService,
	budgetService *service.BudgetService,
	savingGoalService *service.SavingGoalService,
	logger logger.Logger,
) *FinanceDashboardHandler {
	return &FinanceDashboardHandler{
		accountService:     accountService,
		transactionService: transactionService,
		budgetService:      budgetService,
		savingGoalService:  savingGoalService,
		logger:             logger,
	}
}

// BudgetWithStatus représente un budget avec son statut
type BudgetWithStatus struct {
	*entity.Budget
	Status          string  `json:"status"`           // "good", "warning", "danger"
	PercentageUsed  float64 `json:"percentage_used"`  // Pourcentage utilisé
	RemainingAmount float64 `json:"remaining_amount"` // Montant restant
	DaysRemaining   int     `json:"days_remaining"`   // Jours restants dans la période
}

// DebtInfo représente les informations de dette
type DebtInfo struct {
	AccountID   uuid.UUID `json:"account_id"`
	AccountName string    `json:"account_name"`
	DebtAmount  float64   `json:"debt_amount"` // Montant négatif du solde
	Currency    string    `json:"currency"`
}

// FinanceDashboardResponse représente la réponse du tableau de bord
type FinanceDashboardResponse struct {
	// Comptes
	Accounts []*entity.Account `json:"accounts"`

	// Transactions du mois
	CurrentMonthTransactions []*entity.Transaction `json:"current_month_transactions"`

	// Budgets avec leur état
	BudgetsWithStatus []*BudgetWithStatus `json:"budgets_with_status"`

	// Objectifs d'épargne
	SavingGoals []*entity.SavingGoal `json:"saving_goals"`

	// Dettes (comptes avec solde négatif)
	Debts []*DebtInfo `json:"debts"`

	// Transactions récurrentes
	RecurringTransactions []*entity.Transaction `json:"recurring_transactions"`

	// Statistiques résumées
	Summary struct {
		TotalBalance        float64 `json:"total_balance"`
		MonthlyIncome       float64 `json:"monthly_income"`
		MonthlyExpenses     float64 `json:"monthly_expenses"`
		MonthlySavings      float64 `json:"monthly_savings"`
		TotalDebts          float64 `json:"total_debts"`
		BudgetsOverspent    int     `json:"budgets_overspent"`
		SavingGoalsAchieved int     `json:"saving_goals_achieved"`
	} `json:"summary"`
}

// GetFinanceDashboard godoc
// @Summary Récupérer le tableau de bord financier
// @Description Récupère toutes les informations financières pour le tableau de bord
// @Tags Finance
// @Produce json
// @Success 200 {object} response.Response{data=FinanceDashboardResponse} "Données du tableau de bord"
// @Failure 401 {object} response.ErrorResponse "Non autorisé"
// @Failure 500 {object} response.ErrorResponse "Erreur serveur"
// @Security BearerAuth
// @Router /finance/dashboard [get]
func (h *FinanceDashboardHandler) GetFinanceDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Utilisateur non authentifié", nil)
		return
	}

	var dashboardData FinanceDashboardResponse

	// 1. Récupérer tous les comptes
	accounts, err := h.accountService.GetAccountsByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Erreur récupération comptes", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la récupération des comptes", err)
		return
	}
	dashboardData.Accounts = accounts

	// 2. Récupérer les transactions du mois actuel
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	currentMonthTransactions, err := h.transactionService.GetTransactionsByDateRange(
		r.Context(),
		userID,
		startOfMonth.Format("2006-01-02"),
		endOfMonth.Format("2006-01-02"),
	)
	if err != nil {
		h.logger.Error("Erreur récupération transactions du mois", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la récupération des transactions", err)
		return
	}
	dashboardData.CurrentMonthTransactions = currentMonthTransactions

	// 3. Récupérer les budgets avec leur état
	budgets, err := h.budgetService.GetBudgetsByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Erreur récupération budgets", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la récupération des budgets", err)
		return
	}

	var budgetsWithStatus []*BudgetWithStatus

	//calculer amount_spent pour chaque budget
	for _, budget := range budgets {
		var budgetTransactions []*entity.Transaction
		for _, tx := range currentMonthTransactions {
			if tx.CategoryID != nil && budget.Category != nil && *tx.CategoryID == budget.CategoryID {
				budgetTransactions = append(budgetTransactions, tx)
			}
		}
		h.logger.Info("budgetTransactions", logger.Any("budgetTransactions", budgetTransactions))
		var amountSpent float64
		for _, tx := range budgetTransactions {
			amountSpent += tx.Amount
		}
		budget.AmountSpent = amountSpent
	}

	for _, budget := range budgets {
		budgetStatus := h.calculateBudgetStatus(budget)
		budgetsWithStatus = append(budgetsWithStatus, budgetStatus)
	}
	dashboardData.BudgetsWithStatus = budgetsWithStatus

	// 4. Récupérer les objectifs d'épargne
	savingGoals, err := h.savingGoalService.GetSavingGoalsByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Erreur récupération objectifs d'épargne", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la récupération des objectifs d'épargne", err)
		return
	}
	dashboardData.SavingGoals = savingGoals

	// 5. Identifier les dettes (comptes avec solde négatif)
	var debts []*DebtInfo
	for _, account := range accounts {
		if account.Balance < 0 {
			debt := &DebtInfo{
				AccountID:   account.ID,
				AccountName: account.Name,
				DebtAmount:  -account.Balance, // Convertir en montant positif
				Currency:    account.Currency,
			}
			debts = append(debts, debt)
		}
	}
	dashboardData.Debts = debts

	// 6. Récupérer les transactions récurrentes
	allTransactions, err := h.transactionService.GetTransactionsByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("Erreur récupération transactions", logger.Error(err))
		response.Error(w, http.StatusInternalServerError, "Erreur lors de la récupération des transactions", err)
		return
	}

	var recurringTransactions []*entity.Transaction
	for _, transaction := range allTransactions {
		if transaction.Recurring {
			recurringTransactions = append(recurringTransactions, transaction)
		}
	}
	dashboardData.RecurringTransactions = recurringTransactions

	// 7. Calculer les statistiques résumées
	dashboardData.Summary = h.calculateSummary(accounts, currentMonthTransactions, budgetsWithStatus, savingGoals, debts)

	response.Success(w, http.StatusOK, "Tableau de bord récupéré avec succès", dashboardData)
}

// calculateBudgetStatus calcule le statut d'un budget
func (h *FinanceDashboardHandler) calculateBudgetStatus(budget *entity.Budget) *BudgetWithStatus {
	var percentageUsed float64
	if budget.AmountPlanned > 0 {
		percentageUsed = (budget.AmountSpent / budget.AmountPlanned) * 100
	}

	remainingAmount := budget.AmountPlanned - budget.AmountSpent

	// Déterminer le statut
	var status string
	switch {
	case percentageUsed >= 100:
		status = "danger" // Dépassé
	case percentageUsed >= 80:
		status = "warning" // Presque dépassé
	default:
		status = "good" // Bon état
	}

	// Calculer les jours restants (simplifié pour mensuel)
	now := time.Now()
	var daysRemaining int
	switch budget.Period {
	case "monthly":
		endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())
		daysRemaining = int(endOfMonth.Sub(now).Hours() / 24)
	case "yearly":
		endOfYear := time.Date(now.Year()+1, 1, 0, 23, 59, 59, 0, now.Location())
		daysRemaining = int(endOfYear.Sub(now).Hours() / 24)
	case "weekly":
		// Calculer le dimanche de la semaine
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Dimanche = 7
		}
		daysUntilSunday := 7 - weekday
		endOfWeek := now.AddDate(0, 0, daysUntilSunday)
		daysRemaining = int(endOfWeek.Sub(now).Hours() / 24)
	default:
		daysRemaining = 30 // Valeur par défaut
	}

	return &BudgetWithStatus{
		Budget:          budget,
		Status:          status,
		PercentageUsed:  percentageUsed,
		RemainingAmount: remainingAmount,
		DaysRemaining:   daysRemaining,
	}
}

// calculateSummary calcule les statistiques résumées
func (h *FinanceDashboardHandler) calculateSummary(
	accounts []*entity.Account,
	transactions []*entity.Transaction,
	budgets []*BudgetWithStatus,
	savingGoals []*entity.SavingGoal,
	debts []*DebtInfo,
) struct {
	TotalBalance        float64 `json:"total_balance"`
	MonthlyIncome       float64 `json:"monthly_income"`
	MonthlyExpenses     float64 `json:"monthly_expenses"`
	MonthlySavings      float64 `json:"monthly_savings"`
	TotalDebts          float64 `json:"total_debts"`
	BudgetsOverspent    int     `json:"budgets_overspent"`
	SavingGoalsAchieved int     `json:"saving_goals_achieved"`
} {
	var summary struct {
		TotalBalance        float64 `json:"total_balance"`
		MonthlyIncome       float64 `json:"monthly_income"`
		MonthlyExpenses     float64 `json:"monthly_expenses"`
		MonthlySavings      float64 `json:"monthly_savings"`
		TotalDebts          float64 `json:"total_debts"`
		BudgetsOverspent    int     `json:"budgets_overspent"`
		SavingGoalsAchieved int     `json:"saving_goals_achieved"`
	}

	// Calculer le solde total
	for _, account := range accounts {
		if account.Balance > 0 { // Ne compter que les soldes positifs pour le total
			summary.TotalBalance += account.Balance
		}
	}

	// Calculer les revenus et dépenses du mois
	for _, transaction := range transactions {
		switch transaction.Type {
		case "income":
			summary.MonthlyIncome += transaction.Amount
		case "expense":
			summary.MonthlyExpenses += transaction.Amount
		case "saving":
			summary.MonthlySavings += transaction.Amount
		}
	}

	// Calculer le total des dettes
	for _, debt := range debts {
		summary.TotalDebts += debt.DebtAmount
	}

	// Compter les budgets dépassés
	for _, budget := range budgets {
		if budget.Status == "danger" {
			summary.BudgetsOverspent++
		}
	}

	// Compter les objectifs d'épargne atteints
	for _, goal := range savingGoals {
		if goal.IsAchieved {
			summary.SavingGoalsAchieved++
		}
	}

	return summary
}
