import { apiHelpers } from "../api/client"
import { Budget, BudgetStats } from "@/src/types"

const budgetApi = {
  getBudgets: async (): Promise<Budget[]> => {
    try {
      const response = await apiHelpers.get<{data: Budget[]}>('/budgets')
      return response.data
    } catch (error) {
      console.error('Get budgets error:', error)
      throw error
    }
  },

  getBudget: async (id: string): Promise<Budget> => {
    try {
      const response = await apiHelpers.get<{data: Budget}>(`/budgets/${id}`)
      return response.data
    } catch (error) {
      console.error('Get budget error:', error)
      throw error
    }
  },

  createBudget: async (budget: {
    categoryId: string
    amountPlanned: number
    month: number
    year: number
  }): Promise<Budget> => {
    try {
      const response = await apiHelpers.post<{data: Budget}>('/budgets', budget)
      return response.data
    } catch (error) {
      console.error('Create budget error:', error)
      throw error
    }
  },

  updateBudget: async (id: string, budget: Partial<Budget>): Promise<Budget> => {
    try {
      const response = await apiHelpers.put<{data: Budget}>(`/budgets/${id}`, budget)
      return response.data
    } catch (error) {
      console.error('Update budget error:', error)
      throw error
    }
  },

  deleteBudget: async (id: string): Promise<void> => {
    try {
      await apiHelpers.delete(`/budgets/${id}`)
    } catch (error) {
      console.error('Delete budget error:', error)
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

export default budgetApi 