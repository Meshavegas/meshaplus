import { View, Text, TouchableOpacity } from 'react-native';
import { Link } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { ScreenContent } from '@/components/ScreenContent';

export default function Home() {
  return (
    <ScreenContent title="Home" path="app/index.tsx">
      <View className="space-y-4 mt-6">
        <Link href="/dashboard/tabs/overview" asChild>
          <TouchableOpacity className="bg-blue-600 p-4 rounded-lg flex-row items-center justify-center">
            <Ionicons name="analytics" size={20} color="white" />
            <Text className="text-white font-semibold ml-2">Go to Dashboard</Text>
          </TouchableOpacity>
        </Link>
        
        <Link href="/settings" asChild>
          <TouchableOpacity className="bg-gray-600 p-4 rounded-lg flex-row items-center justify-center">
            <Ionicons name="settings" size={20} color="white" />
            <Text className="text-white font-semibold ml-2">Settings</Text>
          </TouchableOpacity>
        </Link>
      </View>
    </ScreenContent>
  );
} 