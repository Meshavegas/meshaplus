import apiClient, { apiHelpers } from "../api/client";

interface SavingGoal {
  id: string;
  userId: string;
  title: string;
  target_amount: number;
  current_amount: number;
  frequency: 'weekly' | 'monthly' | 'yearly';
  deadline: string;
  status: 'active' | 'completed' | 'overdue';
  percentage_achieved: number;
  remaining_amount: number;
  days_remaining: number;
  created_at: string;
  updated_at: string;
}

interface CreateSavingGoalRequest {
  title: string;
  target_amount: number;
  frequency: 'weekly' | 'monthly' | 'yearly';
  deadline: string;
}

interface UpdateSavingGoalRequest {
  title?: string;
  target_amount?: number;
  frequency?: 'weekly' | 'monthly' | 'yearly';
  deadline?: string;
}

const savingGoalApi = {
  getSavingGoals: async (): Promise<SavingGoal[]> => {
    try {
      const response = await apiHelpers.get<{data: SavingGoal[]}>('/saving-goals');
      return response.data;
    } catch (error) {
      console.error('Get saving goals error:', error);
      throw error;
    }
  },

  getSavingGoal: async (id: string): Promise<SavingGoal> => {
    try {
      const response = await apiClient.get(`/saving-goals/${id}`);
      return response.data.data;
    } catch (error) {
      console.error('Get saving goal error:', error);
      throw error;
    }
  },

  createSavingGoal: async (goal: CreateSavingGoalRequest): Promise<SavingGoal> => {
    try {
      const response = await apiClient.post('/saving-goals', goal);
      return response.data.data;
    } catch (error) {
      console.error('Create saving goal error:', error);
      throw error;
    }
  },

  updateSavingGoal: async (id: string, goal: UpdateSavingGoalRequest): Promise<SavingGoal> => {
    try {
      const response = await apiClient.put(`/saving-goals/${id}`, goal);
      return response.data.data;
    } catch (error) {
      console.error('Update saving goal error:', error);
      throw error;
    }
  },

  deleteSavingGoal: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/saving-goals/${id}`);
    } catch (error) {
      console.error('Delete saving goal error:', error);
      throw error;
    }
  },

  addContribution: async (id: string, amount: number): Promise<SavingGoal> => {
    try {
      const response = await apiClient.post(`/saving-goals/${id}/contributions`, { amount });
      return response.data.data;
    } catch (error) {
      console.error('Add contribution error:', error);
      throw error;
    }
  },

  getSavingGoalsWithStatus: async (): Promise<SavingGoal[]> => {
    try {
      const response = await apiHelpers.get<{data: {saving_goals: SavingGoal[]}}>('/saving-goals/with-status');
      return response.data.saving_goals;
    } catch (error) {
      console.error('Get saving goals with status error:', error);
      throw error;
    }
  }
};

export default savingGoalApi;
export type { SavingGoal, CreateSavingGoalRequest, UpdateSavingGoalRequest }; 