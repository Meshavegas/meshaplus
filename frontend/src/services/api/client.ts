// Temporairement commenté jusqu'à l'installation d'axios
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { QueryClient } from '@tanstack/react-query'
import { useAuthStore } from '@/src/stores/authStore'

// Configuration de l'API
const API_BASE_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api'

// Client Axios avec intercepteurs - temporairement commenté
export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Intercepteur pour ajouter le token d'authentification
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Intercepteur pour gérer les erreurs d'authentification
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  async (error) => {
    if (error.response?.status === 401) {
      // Token expiré, déconnexion automatique
      useAuthStore.getState().logout()
    }
    return Promise.reject(error)
  }
)

// Configuration TanStack Query
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      gcTime: 10 * 60 * 1000, // 10 minutes (anciennement cacheTime)
      retry: (failureCount, error: any) => {
        // Ne pas retry sur les erreurs 4xx
        if (error?.response?.status >= 400 && error?.response?.status < 500) {
          return false
        }
        return failureCount < 3
      },
      networkMode: 'online',
    },
    mutations: {
      retry: false,
    },
  },
})

// Types pour les requêtes API
export interface ApiRequestConfig extends AxiosRequestConfig {
  skipAuth?: boolean
}

export interface ApiResponse<T = any> {
  data: T
  message: string
  success: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  hasMore: boolean
}

// Helpers pour les requêtes API
export const apiHelpers = {
  // GET request
  get: async <T>(url: string, config?: ApiRequestConfig): Promise<T> => {
    const response = await apiClient.get<T>(url, config)
    return response.data
  },

  // POST request
  post: async <T>(url: string, data?: any, config?: ApiRequestConfig): Promise<T> => {
    const response = await apiClient.post<T>(url, data, config)
    return response.data
  },

  // PUT request
  put: async <T>(url: string, data?: any, config?: ApiRequestConfig): Promise<T> => {
    const response = await apiClient.put<T>(url, data, config)
    return response.data
  },

  // PATCH request
  patch: async <T>(url: string, data?: any, config?: ApiRequestConfig): Promise<T> => {
    const response = await apiClient.patch<T>(url, data, config)
    return response.data
  },

  // DELETE request
  delete: async <T>(url: string, config?: ApiRequestConfig): Promise<T> => {
    const response = await apiClient.delete<T>(url, config)
    return response.data
  },
}

// Clés de cache pour TanStack Query
export const queryKeys = {
  auth: {
    user: ['auth', 'user'] as const,
    profile: ['auth', 'profile'] as const,
  },
  expenses: {
    all: ['expenses'] as const,
    lists: () => [...queryKeys.expenses.all, 'list'] as const,
    list: (filters: any) => [...queryKeys.expenses.lists(), filters] as const,
    details: () => [...queryKeys.expenses.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.expenses.details(), id] as const,
  },
  revenues: {
    all: ['revenues'] as const,
    lists: () => [...queryKeys.revenues.all, 'list'] as const,
    list: (filters: any) => [...queryKeys.revenues.lists(), filters] as const,
    details: () => [...queryKeys.revenues.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.revenues.details(), id] as const,
  },
  habits: {
    all: ['habits'] as const,
    lists: () => [...queryKeys.habits.all, 'list'] as const,
    list: (filters: any) => [...queryKeys.habits.lists(), filters] as const,
    details: () => [...queryKeys.habits.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.habits.details(), id] as const,
    logs: (habitId: string) => [...queryKeys.habits.detail(habitId), 'logs'] as const,
  },
  goals: {
    all: ['goals'] as const,
    lists: () => [...queryKeys.goals.all, 'list'] as const,
    list: (filters: any) => [...queryKeys.goals.lists(), filters] as const,
    details: () => [...queryKeys.goals.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.goals.details(), id] as const,
  },
  tasks: {
    all: ['tasks'] as const,
    lists: () => [...queryKeys.tasks.all, 'list'] as const,
    list: (filters: any) => [...queryKeys.tasks.lists(), filters] as const,
    details: () => [...queryKeys.tasks.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.tasks.details(), id] as const,
  },
} as const

export default apiClient 