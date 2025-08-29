import { apiHelpers } from "../api/client"

const financeApi = {
  getFinancialSummary: async (): Promise<FinanceDashboard> => {
    const response = await apiHelpers.get<{data:FinanceDashboard}>('/finance/dashboard')
    console.log(JSON.stringify(response, null, 2), "response finance dashboard");
    return response.data
  }
}

export default financeApi