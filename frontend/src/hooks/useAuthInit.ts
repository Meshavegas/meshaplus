import { useEffect } from 'react'
import { useAuthStore } from '@/src/stores/authStore'
import { apiHelpers } from '../services/api/client'

export const useAuthInit = () => {
  const { initializeAuth, isAuthenticated, user, token } = useAuthStore()

  useEffect(() => {
    // Initialiser l'authentification au démarrage
    initializeAuth()
    apiHelpers.setToken(token ?? '')
    
    // Log pour debug
    console.log('Auth state initialized:', {
      isAuthenticated,
      hasUser: !!user,
      hasToken: !!token,
      userEmail: user?.email
    })
  }, [isAuthenticated])

  return {
    isAuthenticated,
    user,
    token
  }
} 