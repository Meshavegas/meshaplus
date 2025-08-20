import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';

import '@/global.css';

export default function RootLayout() {
  return (
    <>
      <Stack>
        <Stack.Screen name="index" options={{ title: 'Home' }} />
        <Stack.Screen name="dashboard" options={{ title: 'Dashboard', headerShown: false }} />
        <Stack.Screen name="settings" options={{ title: 'Settings' }} />
      </Stack>
      <StatusBar style="auto" />
    </>
  );
} 