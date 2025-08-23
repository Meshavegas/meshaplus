import { useEffect } from 'react'
import { useAuthStore } from '@/src/stores/authStore'

export const useAuthInit = () => {
  const { initializeAuth, isAuthenticated, user, token } = useAuthStore()

  useEffect(() => {
    // Initialiser l'authentification au démarrage
    initializeAuth()
    
    // Log pour debug
    console.log('Auth state initialized:', {
      isAuthenticated,
      hasUser: !!user,
      hasToken: !!token,
      userEmail: user?.email
    })
  }, [])

  return {
    isAuthenticated,
    user,
    token
  }
} 