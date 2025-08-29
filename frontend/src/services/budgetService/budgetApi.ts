import apiClient, { apiHelpers } from "../api/client";

interface CreateBudgetRequest {
  name: string;
  categoryId: string;
  amountPlanned: number;
  period: 'weekly' | 'monthly' | 'yearly';
  description?: string;
}

interface UpdateBudgetRequest {
  name?: string;
  categoryId?: string;
  amountPlanned?: number;
  period?: 'weekly' | 'monthly' | 'yearly';
  description?: string;
}

const budgetApi = {
  getBudgets: async (): Promise<Budget[]> => {
    try {
      const response = await apiHelpers.get<{data: Budget[]}>('/budgets');
      return response.data;
    } catch (error) {
      console.error('Get budgets error:', error);
      throw error;
    }
  },

  getBudget: async (id: string): Promise<Budget> => {
    try {
      const response = await apiClient.get(`/budgets/${id}`);
      return response.data.data;
    } catch (error) {
      console.error('Get budget error:', error);
      throw error;
    }
  },

  createBudget: async (budget: CreateBudgetRequest): Promise<Budget> => {
    try {
      const response = await apiClient.post('/budgets', budget);
      return response.data.data;
    } catch (error) {
      console.error('Create budget error:', error);
      throw error;
    }
  },

  updateBudget: async (id: string, budget: UpdateBudgetRequest): Promise<Budget> => {
    try {
      const response = await apiClient.put(`/budgets/${id}`, budget);
      return response.data.data;
    } catch (error) {
      console.error('Update budget error:', error);
      throw error;
    }
  },

  deleteBudget: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/budgets/${id}`);
    } catch (error) {
      console.error('Delete budget error:', error);
      throw error;
    }
  },

  getBudgetsByCategory: async (categoryId: string): Promise<Budget[]> => {
    try {
      const response = await apiClient.get(`/budgets?categoryId=${categoryId}`);
      return response.data.data;
    } catch (error) {
      console.error('Get budgets by category error:', error);
      throw error;
    }
  },

  getBudgetsWithStatus: async (): Promise<Budget[]> => {
    try {
      const response = await apiHelpers.get<{data: {budgets: Budget[]}}>('/budgets/with-status');
      return response.data.budgets;
    } catch (error) {
      console.error('Get budgets with status error:', error);
      throw error;
    }
  }
};

export default budgetApi;
export type { CreateBudgetRequest, UpdateBudgetRequest }; 