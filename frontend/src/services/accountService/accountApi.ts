import apiClient, { apiHelpers } from "../api/client"
import { Account } from "@/src/types"

const accountApi = {
  getAccounts: async (): Promise<Account[]> => {
    try {
      const response = await apiHelpers.get<{data: {accounts: Account[]}}>('/accounts')
      return response.data.accounts
    } catch (error) {
      console.error('Get accounts error:', error)
      throw error
    }
  },

  getAccount: async (id: string): Promise<Account> => {
    try {
      const response = await apiClient.get(`/accounts/${id}`)
      return response.data.data
    } catch (error) {
      console.error('Get account error:', error)
      throw error
    }
  },

  createAccount: async (account: { name: string; type: string; balance: number; currency: string }): Promise<Account> => {
    try {
      const response = await apiClient.post('/accounts', account)
      return response.data.data
    } catch (error) {
      console.error('Create account error:', error)
      throw error
    }
  },

  updateAccount: async (id: string, account: Partial<Account>): Promise<Account> => {
    try {
      const response = await apiClient.put(`/accounts/${id}`, account)
      return response.data.data
    } catch (error) {
      console.error('Update account error:', error)
      throw error
    }
  }
}

export default accountApi 