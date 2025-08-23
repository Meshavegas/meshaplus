import { useEffect } from 'react'
import { View, Text } from 'react-native'
import { Redirect, router } from 'expo-router'
import { useIsAuthenticated } from '@/src/stores/authStore'

export default function Index() {
  const isAuthenticated = useIsAuthenticated()
  // const isAuthenticated = true


  // Rediriger vers le dashboard si l'utilisateur est connecté
  if (isAuthenticated) {
    return <Redirect href="/dashboard" />
  } else {
    // Rediriger vers la page de connexion si non connecté
    return <Redirect href="/auth/login" />
  }


  return (
    <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
      <Text>Redirection...</Text>
    </View>
  )
} 