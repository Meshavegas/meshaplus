import { View, Text, ScrollView, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';

export default function Goals() {
  return (
    <ScrollView className="flex-1 bg-gray-50">
      <View className="p-4">
        <Text className="text-2xl font-bold text-gray-900 mb-6">Financial Goals</Text>
        
        {/* Add Goal Button */}
        <TouchableOpacity className="bg-purple-600 p-4 rounded-lg mb-6 flex-row items-center justify-center">
          <Ionicons name="add" size={20} color="white" />
          <Text className="text-white font-semibold ml-2">Add New Goal</Text>
        </TouchableOpacity>
        
        {/* Active Goals */}
        <View className="space-y-4 mb-6">
          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-3">
              <Text className="text-lg font-semibold text-gray-900">Emergency Fund</Text>
              <View className="bg-green-100 px-2 py-1 rounded">
                <Text className="text-green-700 text-xs font-medium">Active</Text>
              </View>
            </View>
            <View className="mb-3">
              <View className="flex-row justify-between items-center mb-1">
                <Text className="text-sm text-gray-500">Progress</Text>
                <Text className="text-sm text-gray-500">$8,000 / $10,000</Text>
              </View>
              <View className="w-full bg-gray-200 rounded-full h-2">
                <View className="bg-green-500 h-2 rounded-full" style={{ width: '80%' }} />
              </View>
            </View>
            <Text className="text-sm text-gray-600">Target: $10,000 • Due: Dec 2024</Text>
          </View>
          
          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-3">
              <Text className="text-lg font-semibold text-gray-900">Vacation Fund</Text>
              <View className="bg-blue-100 px-2 py-1 rounded">
                <Text className="text-blue-700 text-xs font-medium">Active</Text>
              </View>
            </View>
            <View className="mb-3">
              <View className="flex-row justify-between items-center mb-1">
                <Text className="text-sm text-gray-500">Progress</Text>
                <Text className="text-sm text-gray-500">$1,200 / $3,000</Text>
              </View>
              <View className="w-full bg-gray-200 rounded-full h-2">
                <View className="bg-blue-500 h-2 rounded-full" style={{ width: '40%' }} />
              </View>
            </View>
            <Text className="text-sm text-gray-600">Target: $3,000 • Due: Mar 2025</Text>
          </View>
          
          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-3">
              <Text className="text-lg font-semibold text-gray-900">New Car</Text>
              <View className="bg-yellow-100 px-2 py-1 rounded">
                <Text className="text-yellow-700 text-xs font-medium">Planning</Text>
              </View>
            </View>
            <View className="mb-3">
              <View className="flex-row justify-between items-center mb-1">
                <Text className="text-sm text-gray-500">Progress</Text>
                <Text className="text-sm text-gray-500">$0 / $25,000</Text>
              </View>
              <View className="w-full bg-gray-200 rounded-full h-2">
                <View className="bg-yellow-500 h-2 rounded-full" style={{ width: '0%' }} />
              </View>
            </View>
            <Text className="text-sm text-gray-600">Target: $25,000 • Due: Dec 2025</Text>
          </View>
        </View>
        
        {/* Completed Goals */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold text-gray-900 mb-4">Completed Goals</Text>
          <View className="space-y-3">
            <View className="flex-row justify-between items-center">
              <View className="flex-row items-center">
                <View className="w-8 h-8 bg-green-100 rounded-full items-center justify-center mr-3">
                  <Ionicons name="checkmark" size={16} color="#16a34a" />
                </View>
                <View>
                  <Text className="text-gray-900 font-medium">Laptop Purchase</Text>
                  <Text className="text-sm text-gray-500">Completed Dec 2023</Text>
                </View>
              </View>
              <Text className="text-green-600 font-semibold">$1,200</Text>
            </View>
            
            <View className="flex-row justify-between items-center">
              <View className="flex-row items-center">
                <View className="w-8 h-8 bg-green-100 rounded-full items-center justify-center mr-3">
                  <Ionicons name="checkmark" size={16} color="#16a34a" />
                </View>
                <View>
                  <Text className="text-gray-900 font-medium">Phone Upgrade</Text>
                  <Text className="text-sm text-gray-500">Completed Nov 2023</Text>
                </View>
              </View>
              <Text className="text-green-600 font-semibold">$800</Text>
            </View>
          </View>
        </View>
      </View>
    </ScrollView>
  );
} 