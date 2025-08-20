import { View, Text } from 'react-native';
import { Slot } from 'expo-router';

export default function DashboardLayout() {
  return (
    <View className="flex-1 bg-gray-50">
      <Slot />
    </View>
  );
}