import { apiHelpers } from "../api/client"
import { DashboardData, TransactionStats, TaskStats, BudgetStats } from "@/src/types"

const dashboardApi = {
  getDashboardData: async (): Promise<DashboardData> => {
    try {
      // Charger toutes les données en parallèle
      const [
        transactions,
        tasks,
        taskStats,
        transactionStats,
        budgetStats
      ] = await Promise.all([
        apiHelpers.get<{data: any[]}>('/transactions?limit=5'),
        apiHelpers.get<{data: any[]}>('/tasks?limit=5'),
        apiHelpers.get<{data: TaskStats}>('/tasks/stats'),
        apiHelpers.get<{data: TransactionStats}>('/transactions/stats'),
        apiHelpers.get<{data: BudgetStats}>('/budgets/stats')
      ])

      // Calculer les données du dashboard
      const financeOverview = {
        totalRevenue: transactionStats.data.totalIncome,
        totalExpense: transactionStats.data.totalExpense,
        balance: transactionStats.data.netAmount,
        monthlyBudget: budgetStats.data.totalPlanned,
        savings: budgetStats.data.remaining
      }

      const tasksOverview = {
        totalTasks: taskStats.data.totalTasks,
        completed: taskStats.data.completedTasks,
        pending: taskStats.data.pendingTasks,
        incoming: 0, // À calculer selon la logique métier
        running: 0   // À calculer selon la logique métier
      }

      return {
        userId: '', // Sera rempli par le store d'auth
        tasksOverview,
        financeOverview,
        recentTransactions: transactions.data,
        recentTasks: tasks.data
      }
    } catch (error) {
      console.error('Get dashboard data error:', error)
      throw error
    }
  },

  getTransactionStats: async (period: string = 'month'): Promise<TransactionStats> => {
    try {
      const response = await apiHelpers.get<{data: TransactionStats}>(`/transactions/stats?period=${period}`)
      return response.data
    } catch (error) {
      console.error('Get transaction stats error:', error)
      throw error
    }
  },

  getTaskStats: async (): Promise<TaskStats> => {
    try {
      const response = await apiHelpers.get<{data: TaskStats}>('/tasks/stats')
      return response.data
    } catch (error) {
      console.error('Get task stats error:', error)
      throw error
    }
  },

  getBudgetStats: async (month?: number, year?: number): Promise<BudgetStats> => {
    try {
      const params = new URLSearchParams()
      if (month) params.append('month', month.toString())
      if (year) params.append('year', year.toString())
      
      const url = `/budgets/stats${params.toString() ? `?${params.toString()}` : ''}`
      const response = await apiHelpers.get<{data: BudgetStats}>(url)
      return response.data
    } catch (error) {
      console.error('Get budget stats error:', error)
      throw error
    }
  }
}

export default dashboardApi 