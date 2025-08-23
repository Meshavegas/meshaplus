import { apiHelpers } from "../api/client"
import { SavingGoal } from "@/src/types"

const savingGoalApi = {
  getSavingGoals: async (): Promise<SavingGoal[]> => {
    try {
      const response = await apiHelpers.get<{data: SavingGoal[]}>('/saving-goals')
      return response.data
    } catch (error) {
      console.error('Get saving goals error:', error)
      throw error
    }
  },

  getSavingGoal: async (id: string): Promise<SavingGoal> => {
    try {
      const response = await apiHelpers.get<{data: SavingGoal}>(`/saving-goals/${id}`)
      return response.data
    } catch (error) {
      console.error('Get saving goal error:', error)
      throw error
    }
  },

  createSavingGoal: async (goal: {
    title: string
    targetAmount: number
    deadline?: Date
    frequency: 'weekly' | 'monthly' | 'yearly'
  }): Promise<SavingGoal> => {
    try {
      const response = await apiHelpers.post<{data: SavingGoal}>('/saving-goals', goal)
      return response.data
    } catch (error) {
      console.error('Create saving goal error:', error)
      throw error
    }
  },

  updateSavingGoal: async (id: string, goal: Partial<SavingGoal>): Promise<SavingGoal> => {
    try {
      const response = await apiHelpers.put<{data: SavingGoal}>(`/saving-goals/${id}`, goal)
      return response.data
    } catch (error) {
      console.error('Update saving goal error:', error)
      throw error
    }
  },

  deleteSavingGoal: async (id: string): Promise<void> => {
    try {
      await apiHelpers.delete(`/saving-goals/${id}`)
    } catch (error) {
      console.error('Delete saving goal error:', error)
      throw error
    }
  },

  addContribution: async (id: string, amount: number): Promise<SavingGoal> => {
    try {
      const response = await apiHelpers.post<{data: SavingGoal}>(`/saving-goals/${id}/contribute`, { amount })
      return response.data
    } catch (error) {
      console.error('Add contribution error:', error)
      throw error
    }
  }
}

export default savingGoalApi 