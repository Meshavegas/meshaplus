import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { SyncStatus } from '@/src/components/ui/SyncStatus';

export default function Finances() {
  const [activeTab, setActiveTab] = useState<'expenses' | 'revenues' | 'budgets'>('expenses');

  const renderExpensesTab = () => (
    <View className="space-y-4">
      {/* Quick Add Button */}
      <TouchableOpacity className="bg-red-600 p-4 rounded-lg flex-row items-center justify-center">
        <Ionicons name="add" size={20} color="white" />
        <Text className="text-white font-semibold ml-2">Ajouter une Dépense</Text>
      </TouchableOpacity>

      {/* Budget Overview */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <Text className="text-lg font-semibold text-gray-900 mb-3">Budget du Mois</Text>
        <View className="space-y-2">
          <View className="flex-row justify-between items-center">
            <Text className="text-sm text-gray-600">Dépenses</Text>
            <Text className="text-red-600 font-semibold">$1,200 / $2,000</Text>
          </View>
          <View className="w-full bg-gray-200 rounded-full h-2">
            <View className="bg-red-500 h-2 rounded-full" style={{ width: '60%' }} />
          </View>
          <Text className="text-xs text-gray-500">40% restant</Text>
        </View>
      </View>

      {/* Categories */}
      <View className="space-y-3">
        <Text className="text-lg font-semibold text-gray-900">Catégories</Text>
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <View className="flex-row justify-between items-center mb-3">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-red-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="restaurant" size={16} color="#dc2626" />
              </View>
              <Text className="text-gray-900 font-medium">Alimentation</Text>
            </View>
            <Text className="text-red-600 font-semibold">$450</Text>
          </View>
          <View className="flex-row justify-between items-center mb-3">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-blue-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="car" size={16} color="#2563eb" />
              </View>
              <Text className="text-gray-900 font-medium">Transport</Text>
            </View>
            <Text className="text-red-600 font-semibold">$280</Text>
          </View>
          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-purple-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="shirt" size={16} color="#9333ea" />
              </View>
              <Text className="text-gray-900 font-medium">Shopping</Text>
            </View>
            <Text className="text-red-600 font-semibold">$320</Text>
          </View>
        </View>
      </View>

      {/* Recurring Expenses */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <Text className="text-lg font-semibold text-gray-900 mb-3">Dépenses Récurrentes</Text>
        <View className="space-y-2">
          <View className="flex-row justify-between items-center">
            <Text className="text-gray-700">Loyer</Text>
            <Text className="text-red-600">-$800</Text>
          </View>
          <View className="flex-row justify-between items-center">
            <Text className="text-gray-700">Électricité</Text>
            <Text className="text-red-600">-$120</Text>
          </View>
          <View className="flex-row justify-between items-center">
            <Text className="text-gray-700">Internet</Text>
            <Text className="text-red-600">-$60</Text>
          </View>
        </View>
      </View>
    </View>
  );

  const renderRevenuesTab = () => (
    <View className="space-y-4">
      {/* Quick Add Button */}
      <TouchableOpacity className="bg-green-600 p-4 rounded-lg flex-row items-center justify-center">
        <Ionicons name="add" size={20} color="white" />
        <Text className="text-white font-semibold ml-2">Ajouter un Revenu</Text>
      </TouchableOpacity>

      {/* Revenue Summary */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <Text className="text-lg font-semibold text-gray-900 mb-3">Revenus du Mois</Text>
        <Text className="text-3xl font-bold text-green-600 mb-2">$3,650</Text>
        <Text className="text-sm text-gray-500">+12.5% vs mois dernier</Text>
      </View>

      {/* Revenue Sources */}
      <View className="space-y-3">
        <Text className="text-lg font-semibold text-gray-900">Sources de Revenus</Text>
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <View className="flex-row justify-between items-center mb-3">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-green-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="briefcase" size={16} color="#16a34a" />
              </View>
              <Text className="text-gray-900 font-medium">Salaire</Text>
            </View>
            <Text className="text-green-600 font-semibold">$2,500</Text>
          </View>
          <View className="flex-row justify-between items-center mb-3">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-blue-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="laptop" size={16} color="#2563eb" />
              </View>
              <Text className="text-gray-900 font-medium">Freelance</Text>
            </View>
            <Text className="text-green-600 font-semibold">$850</Text>
          </View>
          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-yellow-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="trending-up" size={16} color="#ca8a04" />
              </View>
              <Text className="text-gray-900 font-medium">Investissements</Text>
            </View>
            <Text className="text-green-600 font-semibold">$300</Text>
          </View>
        </View>
      </View>
    </View>
  );

  const renderBudgetsTab = () => (
    <View className="space-y-4">
      {/* Budget Alerts */}
      <View className="bg-orange-50 border border-orange-200 rounded-lg p-4">
        <View className="flex-row items-center mb-2">
          <Ionicons name="warning" size={20} color="#f59e0b" />
          <Text className="text-orange-800 font-semibold ml-2">Alertes Budget</Text>
        </View>
        <Text className="text-orange-700 text-sm">
          Vous avez dépensé 80% de votre budget alimentation
        </Text>
      </View>

      {/* Budget Categories */}
      <View className="space-y-3">
        <Text className="text-lg font-semibold text-gray-900">Budgets par Catégorie</Text>
        <View className="space-y-3">
          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-gray-900 font-medium">Alimentation</Text>
              <Text className="text-gray-600">$450 / $500</Text>
            </View>
            <View className="w-full bg-gray-200 rounded-full h-2">
              <View className="bg-orange-500 h-2 rounded-full" style={{ width: '90%' }} />
            </View>
          </View>

          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-gray-900 font-medium">Transport</Text>
              <Text className="text-gray-600">$280 / $400</Text>
            </View>
            <View className="w-full bg-gray-200 rounded-full h-2">
              <View className="bg-blue-500 h-2 rounded-full" style={{ width: '70%' }} />
            </View>
          </View>

          <View className="bg-white p-4 rounded-lg shadow-sm">
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-gray-900 font-medium">Loisirs</Text>
              <Text className="text-gray-600">$150 / $300</Text>
            </View>
            <View className="w-full bg-gray-200 rounded-full h-2">
              <View className="bg-green-500 h-2 rounded-full" style={{ width: '50%' }} />
            </View>
          </View>
        </View>
      </View>
    </View>
  );

  return (
    <ScrollView className="flex-1 bg-gray-50">
      <View className="p-4 mt-safe">
        {/* Tab Navigation */}
        <View className='flex-row justify-end gap-4'>
          <Text>
            View Stats
          </Text>
          <Text>
            View Goals
          </Text>
        </View>
        <View className="flex-row bg-white rounded-lg p-1 mb-6">
          <TouchableOpacity
            className={`flex-1 py-2 px-4 rounded-md ${activeTab === 'expenses' ? 'bg-red-100' : ''}`}
            onPress={() => setActiveTab('expenses')}
          >
            <Text className={`text-center font-medium ${activeTab === 'expenses' ? 'text-red-700' : 'text-gray-600'}`}>
              Dépenses
            </Text>
          </TouchableOpacity>
          <TouchableOpacity
            className={`flex-1 py-2 px-4 rounded-md ${activeTab === 'revenues' ? 'bg-green-100' : ''}`}
            onPress={() => setActiveTab('revenues')}
          >
            <Text className={`text-center font-medium ${activeTab === 'revenues' ? 'text-green-700' : 'text-gray-600'}`}>
              Revenus
            </Text>
          </TouchableOpacity>
          <TouchableOpacity
            className={`flex-1 py-2 px-4 rounded-md ${activeTab === 'budgets' ? 'bg-blue-100' : ''}`}
            onPress={() => setActiveTab('budgets')}
          >
            <Text className={`text-center font-medium ${activeTab === 'budgets' ? 'text-blue-700' : 'text-gray-600'}`}>
              Budgets
            </Text>
          </TouchableOpacity>
        </View>

        {/* Tab Content */}
        {activeTab === 'expenses' && renderExpensesTab()}
        {activeTab === 'revenues' && renderRevenuesTab()}
        {activeTab === 'budgets' && renderBudgetsTab()}
      </View>
    </ScrollView>
  );
} 