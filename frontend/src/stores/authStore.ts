import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { User } from '@/src/types'
import { storage } from '@/src/services/storage/mmkv'
import authApi from '../services/authService/authApi'
import { router } from 'expo-router'

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
  initializeAuth: () => void
}

type AuthStore = AuthState & AuthActions

// Storage personnalisé pour MMKV


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
        set({  error: null })
        try {
          const response = await authApi.login(email, password)

          // Vérifier que la réponse contient les données nécessaires
          if (response && response.data) {
            const user = response.data.user
            const token = response.data.access_token

            if (user && token) {
              set({
                user,
                token,
                isAuthenticated: true,
                isLoading: false,
                error: null
              })
            } else {
              throw new Error('Données de connexion invalides')
            }
          } else {
            throw new Error('Réponse invalide du serveur')
          }
        } catch (error) {
          console.error('Login error:', error)
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
            created_at: new Date(),
            updated_at: new Date(),
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
        router.push('/auth/login')
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

      initializeAuth: () => {
        // Cette fonction peut être appelée au démarrage de l'app
        // pour vérifier si l'utilisateur est déjà connecté
        const state = get()
        if (state.token && state.user) {
          console.log('User already authenticated:', state.user.email)
          // Optionnel: Vérifier la validité du token ici
          // await validateToken(state.token)
        } else {
          console.log('No authenticated user found')
          // Nettoyer l'état si nécessaire
          if (state.isAuthenticated) {
            set({
              user: null,
              token: null,
              isAuthenticated: false,
              error: null
            })
          }
        }
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => (
        {
          getItem: (name: string) => storage.getString(name) ?? null,
          setItem: (name: string, value: string) => storage.set(name, value),
          removeItem: (name: string) => storage.delete(name),
        }
      )),
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