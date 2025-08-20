import { View, Text, ScrollView, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';

export default function Revenue() {
  return (
    <ScrollView className="flex-1 bg-gray-50">
      <View className="p-4">
        <Text className="text-2xl font-bold text-gray-900 mb-6">Revenue</Text>
        
        {/* Add Revenue Button */}
        <TouchableOpacity className="bg-green-600 p-4 rounded-lg mb-6 flex-row items-center justify-center">
          <Ionicons name="add" size={20} color="white" />
          <Text className="text-white font-semibold ml-2">Add New Income</Text>
        </TouchableOpacity>
        
        {/* Revenue Summary */}
        <View className="bg-white p-4 rounded-lg shadow-sm mb-6">
          <Text className="text-lg font-semibold text-gray-900 mb-4">This Month's Revenue</Text>
          <Text className="text-3xl font-bold text-green-600 mb-2">$3,650.00</Text>
          <Text className="text-sm text-gray-500">+12.5% from last month</Text>
        </View>
        
        {/* Revenue Sources */}
        <View className="space-y-4 mb-6">
          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-lg font-semibold text-gray-900">Salary</Text>
              <Text className="text-green-600 font-semibold">$2,500.00</Text>
            </View>
            <View className="flex-row justify-between items-center">
              <Text className="text-sm text-gray-500">Monthly salary</Text>
              <Text className="text-sm text-gray-500">1 transaction</Text>
            </View>
          </View>
          
          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-lg font-semibold text-gray-900">Freelance</Text>
              <Text className="text-green-600 font-semibold">$850.00</Text>
            </View>
            <View className="flex-row justify-between items-center">
              <Text className="text-sm text-gray-500">Web development</Text>
              <Text className="text-sm text-gray-500">3 transactions</Text>
            </View>
          </View>
          
          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-lg font-semibold text-gray-900">Investments</Text>
              <Text className="text-green-600 font-semibold">$300.00</Text>
            </View>
            <View className="flex-row justify-between items-center">
              <Text className="text-sm text-gray-500">Dividends</Text>
              <Text className="text-sm text-gray-500">2 transactions</Text>
            </View>
          </View>
        </View>
        
        {/* Recent Income */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold text-gray-900 mb-4">Recent Income</Text>
          <View className="space-y-3">
            <View className="flex-row justify-between items-center">
              <View className="flex-row items-center">
                <View className="w-8 h-8 bg-green-100 rounded-full items-center justify-center mr-3">
                  <Ionicons name="briefcase" size={16} color="#16a34a" />
                </View>
                <View>
                  <Text className="text-gray-900 font-medium">Salary Deposit</Text>
                  <Text className="text-sm text-gray-500">Today</Text>
                </View>
              </View>
              <Text className="text-green-600 font-semibold">+$2,500.00</Text>
            </View>
            
            <View className="flex-row justify-between items-center">
              <View className="flex-row items-center">
                <View className="w-8 h-8 bg-blue-100 rounded-full items-center justify-center mr-3">
                  <Ionicons name="laptop" size={16} color="#2563eb" />
                </View>
                <View>
                  <Text className="text-gray-900 font-medium">Freelance Project</Text>
                  <Text className="text-sm text-gray-500">Yesterday</Text>
                </View>
              </View>
              <Text className="text-green-600 font-semibold">+$300.00</Text>
            </View>
            
            <View className="flex-row justify-between items-center">
              <View className="flex-row items-center">
                <View className="w-8 h-8 bg-yellow-100 rounded-full items-center justify-center mr-3">
                  <Ionicons name="trending-up" size={16} color="#ca8a04" />
                </View>
                <View>
                  <Text className="text-gray-900 font-medium">Investment Dividend</Text>
                  <Text className="text-sm text-gray-500">3 days ago</Text>
                </View>
              </View>
              <Text className="text-green-600 font-semibold">+$150.00</Text>
            </View>
          </View>
        </View>
      </View>
    </ScrollView>
  );
} 