import apiClient, { apiHelpers } from "../api/client"

const authApi = {
    login: async (email: string, password: string): Promise<LoginResponse> => {
        try {
            const response = await apiHelpers.post<LoginResponse>(`/auth/login`, { email, password })
            console.log('login', response.data, {
                email,
                password
            }, JSON.stringify(response.data, null, 2));

            apiHelpers.setToken(response.data.access_token)
            return response
        } catch (error) {
            throw error
        }
    }
}

export default authApi