interface FinanceDashboard {
    accounts: Account[],
    current_month_transactions: Transaction[],
    budgets_with_status: Budget[],
    saving_goals: SavingGoal[],
    summary: FinanceSummary
}

interface FinanceSummary {
    total_balance: number,
    monthly_income: number,
    monthly_expenses: number,
    monthly_savings: number,
    total_debts: number,
    budgets_overspent: number,
    saving_goals_achieved: number
}

interface IAccountRequest {
    id?: string
    name: string
    type: string
    balance: float64
    currency: string,
    color: string,
    icon: string,
    accountNumber: string
}