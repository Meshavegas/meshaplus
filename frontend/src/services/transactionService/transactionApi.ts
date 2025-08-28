import apiClient, { apiHelpers } from "../api/client"
import { CreateTransactionRequest } from "@/src/types"

const transactionApi = {
  createTransaction: async (transaction: CreateTransactionRequest): Promise<Transaction> => {
    try {
      // Remove keys with undefined values from the transaction object
      const cleanedTransaction = Object.fromEntries(
        Object.entries(transaction).filter(([_, v]) => v !== undefined)
      )
      
      const response = await apiHelpers.post<{ data: Transaction }>('/transactions', cleanedTransaction)
      return response.data
    } catch (error) {
      console.error('Create transaction error:', error)
      throw error
    }
  },

  getTransactions: async (): Promise<Transaction[]> => {
    try {
      const response = await apiHelpers.get<{data: {transactions: Transaction[]}}>('/transactions')
      return response.data.transactions
    } catch (error) {
      console.error('Get transactions error:', error)
      throw error
    }
  },

  getTransaction: async (id: string): Promise<Transaction> => {
    try {
      const response = await apiClient.get(`/transactions/${id}`)
      console.log(response.data.data,"response.data.data transaction")
      return response.data.data
    } catch (error) {
      console.error('Get transaction error:', error)
      throw error
    }
  },

  updateTransaction: async (id: string, transaction: Partial<CreateTransactionRequest>): Promise<Transaction> => {
    try {
      const response = await apiClient.put(`/transactions/${id}`, transaction)
      return response.data.data
    } catch (error) {
      console.error('Update transaction error:', error)
      throw error
    }
  },

  deleteTransaction: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/transactions/${id}`)
    } catch (error) {
      console.error('Delete transaction error:', error)
      throw error
    }
  },

  getTransactionsByAccount: async (accountId: string): Promise<Transaction[]> => {
    try {
      const response = await apiClient.get(`/transactions?accountId=${accountId}`)
      return response.data.data
    } catch (error) {
      console.error('Get transactions by account error:', error)
      throw error
    }
  }
}

export default transactionApi 