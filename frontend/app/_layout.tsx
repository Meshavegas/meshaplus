import { Redirect, router, Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { View, ActivityIndicator } from 'react-native';
import { PortalHost } from '@rn-primitives/portal';

import '@/global.css';
import { useIsAuthenticated, useAuthLoading } from '@/src';
import { useAuthInit } from '@/src/hooks/useAuthInit';
import { colors } from '@/src/theme';
import { useEffect } from 'react';

export default function RootLayout() {
  // Initialiser l'authentification au démarrage
  const { isAuthenticated } = useAuthInit()
  const isLoading = useAuthLoading()
  console.log('isAuthenticated', isAuthenticated, 'isLoading:', isLoading);

  // Afficher un loader pendant l'initialisation
  if (isLoading) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: colors.background.primary }}>
        <ActivityIndicator size="large" color={colors.primary[500]} />
      </View>
    )
  }

  return (
    <>
      <Stack screenOptions={{ headerShown: false }}>
        {isAuthenticated ? (
          <>
            <Stack.Screen name="index" options={{ title: 'Home' }} />
            <Stack.Screen name="dashboard" options={{ title: 'Dashboard', headerShown: false }} />
            <Stack.Screen name="settings" options={{ title: 'Settings', headerShown: false }} />
          </>
        ) : (
          <>
            <Stack.Screen name="auth/login" options={{ title: 'Login', headerShown: false }} />
            <Stack.Screen name="auth/register" options={{ title: 'Register', headerShown: false }} />
          </>
        )}
      </Stack>
      <PortalHost />
      <StatusBar style="auto" />
    </>
  );
} 