import apiClient, { apiHelpers } from "../api/client"


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

  getAccount: async (id: string): Promise<AccountDetails> => {
    try {
      const response = await apiHelpers.get<{data: AccountDetails}>(`/accounts/${id}/details`)
      // console.log(JSON.stringify(response.data, null, 2), "response.data account")
      return response.data
    } catch (error) {
      console.error('Get account error:', error)
      throw error
    }
  },

  createAccount: async (account: IAccountRequest): Promise<Account> => {
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
  },

  deleteAccount: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/accounts/${id}`)
    } catch (error) {
      console.error('Delete account error:', error)
      throw error
    }
  }
}

export default accountApi 