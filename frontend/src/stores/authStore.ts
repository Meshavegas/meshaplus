import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { User } from '@/src/types'
import { storage } from '@/src/services/storage/mmkv'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
}

interface AuthActions {
  login: (email: string, password: string) => Promise<void>
  register: (userData: { name: string; email: string; password: string }) => Promise<void>
  logout: () => void
  setUser: (user: User) => void
  setToken: (token: string) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  clearError: () => void
}

type AuthStore = AuthState & AuthActions

// Storage personnalisé pour MMKV
const mmkvStorage = {
  getItem: (name: string) => {
    const value = storage.getString(name)
    return Promise.resolve(value ? JSON.parse(value) : null)
  },
  setItem: (name: string, value: string) => {
    storage.set(name, value)
    return Promise.resolve()
  },
  removeItem: (name: string) => {
    storage.delete(name)
    return Promise.resolve()
  },
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // État initial
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      // Actions
      login: async (email: string, password: string) => {
        set({ isLoading: true, error: null })
        try {
          // TODO: Appel API de connexion
          // const response = await authApi.login(email, password)
          // set({ user: response.user, token: response.token, isAuthenticated: true })
          
          // Simulation pour l'instant
          const mockUser: User = {
            id: '1',
            email,
            name: 'John Doe',
            createdAt: new Date(),
            updatedAt: new Date(),
          }
          const mockToken = 'mock-jwt-token'
          
          set({ 
            user: mockUser, 
            token: mockToken, 
            isAuthenticated: true,
            isLoading: false 
          })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Erreur de connexion',
            isLoading: false 
          })
        }
      },

      register: async (userData) => {
        set({ isLoading: true, error: null })
        try {
          // TODO: Appel API d'inscription
          // const response = await authApi.register(userData)
          // set({ user: response.user, token: response.token, isAuthenticated: true })
          
          // Simulation pour l'instant
          const mockUser: User = {
            id: '1',
            email: userData.email,
            name: userData.name,
            createdAt: new Date(),
            updatedAt: new Date(),
          }
          const mockToken = 'mock-jwt-token'
          
          set({ 
            user: mockUser, 
            token: mockToken, 
            isAuthenticated: true,
            isLoading: false 
          })
        } catch (error) {
          set({ 
            error: error instanceof Error ? error.message : 'Erreur d\'inscription',
            isLoading: false 
          })
        }
      },

      logout: () => {
        set({ 
          user: null, 
          token: null, 
          isAuthenticated: false,
          error: null 
        })
      },

      setUser: (user: User) => {
        set({ user })
      },

      setToken: (token: string) => {
        set({ token, isAuthenticated: true })
      },

      setLoading: (loading: boolean) => {
        set({ isLoading: loading })
      },

      setError: (error: string | null) => {
        set({ error })
      },

      clearError: () => {
        set({ error: null })
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => mmkvStorage),
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)

// Sélecteurs pour optimiser les re-renders
export const useAuthUser = () => useAuthStore((state) => state.user)
export const useAuthToken = () => useAuthStore((state) => state.token)
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated)
export const useAuthLoading = () => useAuthStore((state) => state.isLoading)
export const useAuthError = () => useAuthStore((state) => state.error) 