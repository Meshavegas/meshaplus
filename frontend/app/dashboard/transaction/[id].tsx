import React, { useEffect, useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import transactionApi from '@/src/services/transactionService/transactionApi';
import categoryApi from '@/src/services/categoryService/categoryApi';
import accountApi from '@/src/services/accountService/accountApi';

 

interface Category {
  id: string;
  name: string;
  type: 'income' | 'expense';
  icon: string;
  color: string;
}

interface Account {
  id: string;
  name: string;
  type: string;
  balance: number;
  currency: string;
}

export default function TransactionDetail() {
  const res= useLocalSearchParams();
  const id = res.id as string;
  const router = useRouter();
  const [transaction, setTransaction] = useState<Transaction | null>(null);
  const [category, setCategory] = useState<Category | null>(null);
  const [account, setAccount] = useState<Account | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadTransactionData();
  }, [id]);

  const loadTransactionData = async () => {
    try {
      setLoading(true);
      const transactionData = await transactionApi.getTransaction(id as string);
      setTransaction(transactionData);

      // Charger les données de la catégorie et du compte
      const [categoryData, accountData] = await Promise.all([
        categoryApi.getCategory(transactionData.category_id),
        accountApi.getAccount(transactionData.account_id)
      ]);

      setCategory(categoryData);
      setAccount(accountData);
    } catch (error) {
      console.error('Error loading transaction data:', error);
      Alert.alert('Erreur', 'Impossible de charger les données de la transaction');
    } finally {
      setLoading(false);
    }
  };

  const handleEditTransaction = () => {
    Alert.alert('Modifier la transaction', 'Fonctionnalité à implémenter');
  };

  const handleDeleteTransaction = () => {
    Alert.alert(
      'Supprimer la transaction',
      'Êtes-vous sûr de vouloir supprimer cette transaction ?',
      [
        { text: 'Annuler', style: 'cancel' },
        {
          text: 'Supprimer',
          style: 'destructive',
          onPress: async () => {
            try {
              await transactionApi.deleteTransaction(id as string);
              router.back();
            } catch (error) {
              Alert.alert('Erreur', 'Impossible de supprimer la transaction');
            }
          }
        }
      ]
    );
  };

  const formatCurrency = (amount: number) => {
    if (!amount) return '';
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: account?.currency || 'EUR'
    }).format(amount);
  };

  const formatDate = (date: Date | string | undefined) => {
    console.log(date, "date")
    if (!date) return '';
    
    const dateObj = typeof date === 'string' ? new Date(date) : date;
    
    if (isNaN(dateObj.getTime())) {
      return 'Date invalide';
    }
    
    return dateObj.toLocaleDateString('fr-FR', {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getTransactionIcon = (categoryName: string) => {
    const categoryMap: { [key: string]: string } = {
      'alimentation': 'restaurant',
      'transport': 'car',
      'logement': 'home',
      'sante': 'medical',
      'loisirs': 'game-controller',
      'shopping': 'bag',
      'default': 'card'
    };
    return categoryMap[categoryName.toLowerCase()] || categoryMap.default;
  };

  const getTransactionColor = (type: string) => {
    return type === 'income' ? '#16a34a' : '#dc2626';
  };

  const getTransactionBgColor = (type: string) => {
    return type === 'income' ? '#dcfce7' : '#fef2f2';
  };

  if (loading) {
    return (
      <View className="flex-1 bg-gray-50 justify-center items-center">
        <Text className="text-gray-500">Chargement...</Text>
      </View>
    );
  }

  if (!transaction) {
    return (
      <View className="flex-1 bg-gray-50 justify-center items-center">
        <Text className="text-gray-500">Transaction non trouvée</Text>
      </View>
    );
  }

  return (
    <View className="flex-1 bg-gray-50">
      {/* Header */}
      <View className="bg-white shadow-sm">
        <View className="flex-row items-center p-4 pt-12">
          <TouchableOpacity onPress={() => router.back()} className="mr-4">
            <Ionicons name="arrow-back" size={24} color="#374151" />
          </TouchableOpacity>
          <View className="flex-1">
            <Text className="text-xl font-bold text-gray-900">Détails de la transaction</Text>
            <Text className="text-sm text-gray-500">{formatDate(transaction.created_at)}</Text>
          </View>
          <TouchableOpacity onPress={handleEditTransaction} className="mr-2">
            <Ionicons name="create-outline" size={24} color="#374151" />
          </TouchableOpacity>
          <TouchableOpacity onPress={handleDeleteTransaction}>
            <Ionicons name="trash-outline" size={24} color="#dc2626" />
          </TouchableOpacity>
        </View>
      </View>

      <ScrollView className="p-4">
        {/* Transaction Summary Card */}
        <View className="bg-white p-6 rounded-lg shadow-sm mb-6">
          <View className="flex-row items-center mb-4">
            <View 
              className="w-16 h-16 rounded-full items-center justify-center mr-4"
              style={{ backgroundColor: getTransactionBgColor(transaction.type) }}
            >
              <Ionicons 
                name={getTransactionIcon(category?.name || 'default') as any} 
                size={32} 
                color={getTransactionColor(transaction.type)} 
              />
            </View>
            <View className="flex-1">
              <Text className="text-2xl font-bold text-gray-900">
                {transaction.type === 'income' ? '+' : '-'}{formatCurrency(transaction.amount)}
              </Text>
              <Text className="text-lg text-gray-700">{transaction.description}</Text>
              <Text className="text-sm text-gray-500">{category?.name}</Text>
            </View>
          </View>

          <View className="flex-row justify-between py-3 border-t border-gray-100">
            <View className="items-center flex-1">
              <Text className="text-sm text-gray-500">Type</Text>
              <Text className={`font-semibold ${
                transaction.type === 'income' ? 'text-green-600' : 'text-red-600'
              }`}>
                {transaction.type === 'income' ? 'Revenu' : 'Dépense'}
              </Text>
            </View>
            <View className="items-center flex-1">
              <Text className="text-sm text-gray-500">Récurrent</Text>
              <Text className="font-semibold text-gray-900">
                {transaction.recurring ? 'Oui' : 'Non'}
              </Text>
            </View>
            <View className="items-center flex-1">
              <Text className="text-sm text-gray-500">Compte</Text>
              <Text className="font-semibold text-gray-900">{account?.name}</Text>
            </View>
          </View>
        </View>

        {/* Transaction Details */}
        <View className="bg-white p-6 rounded-lg shadow-sm mb-6">
          <Text className="text-lg font-bold text-gray-900 mb-4">Informations détaillées</Text>
          <View className="space-y-4">
            <View className="flex-row justify-between">
              <Text className="text-gray-600">Description</Text>
              <Text className="text-gray-900 font-medium text-right flex-1 ml-4">
                {transaction.description}
              </Text>
            </View>
            <View className="flex-row justify-between">
              <Text className="text-gray-600">Montant</Text>
              <Text className={`font-semibold ${
                transaction.type === 'income' ? 'text-green-600' : 'text-red-600'
              }`}>
                {formatCurrency(transaction.amount)}
              </Text>
            </View>
            <View className="flex-row justify-between">
              <Text className="text-gray-600">Catégorie</Text>
              <Text className="text-gray-900 font-medium">{category?.name}</Text>
            </View>
            <View className="flex-row justify-between">
              <Text className="text-gray-600">Compte</Text>
              <Text className="text-gray-900 font-medium">{account?.name}</Text>
            </View>
            <View className="flex-row justify-between">
              <Text className="text-gray-600">Date</Text>
              <Text className="text-gray-900 font-medium">{formatDate(transaction.created_at)}</Text>
            </View>
            <View className="flex-row justify-between">
              <Text className="text-gray-600">Récurrent</Text>
              <Text className="text-gray-900 font-medium">
                {transaction.recurring ? 'Oui' : 'Non'}
              </Text>
            </View>
            <View className="flex-row justify-between">
              <Text className="text-gray-600">Créé le</Text>
              <Text className="text-gray-900 font-medium">{formatDate(transaction.created_at)}</Text>
            </View>
            <View className="flex-row justify-between">
              <Text className="text-gray-600">Modifié le</Text>
              <Text className="text-gray-900 font-medium">{formatDate(transaction.updated_at)}</Text>
            </View>
          </View>
        </View>

        {/* Quick Actions */}
        <View className="bg-white p-6 rounded-lg shadow-sm mb-6">
          <Text className="text-lg font-bold text-gray-900 mb-4">Actions</Text>
          <View className="gap-3">
            <TouchableOpacity
              className="flex-row items-center p-3 bg-blue-50 rounded-lg"
              onPress={handleEditTransaction}
            >
              <Ionicons name="create-outline" size={20} color="#3b82f6" />
              <Text className="text-blue-600 font-medium ml-3">Modifier la transaction</Text>
            </TouchableOpacity>
            <TouchableOpacity
              className="flex-row items-center p-3 bg-green-50 rounded-lg"
              onPress={() => Alert.alert('Dupliquer', 'Fonctionnalité à implémenter')}
            >
              <Ionicons name="copy-outline" size={20} color="#10b981" />
              <Text className="text-green-600 font-medium ml-3">Dupliquer la transaction</Text>
            </TouchableOpacity>
            <TouchableOpacity
              className="flex-row items-center p-3 bg-orange-50 rounded-lg"
              onPress={() => Alert.alert('Marquer comme récurrent', 'Fonctionnalité à implémenter')}
            >
              <Ionicons name="repeat-outline" size={20} color="#f59e0b" />
              <Text className="text-orange-600 font-medium ml-3">
                {transaction.recurring ? 'Modifier la récurrence' : 'Marquer comme récurrent'}
              </Text>
            </TouchableOpacity>
            <TouchableOpacity
              className="flex-row items-center p-3 bg-red-50 rounded-lg"
              onPress={handleDeleteTransaction}
            >
              <Ionicons name="trash-outline" size={20} color="#dc2626" />
              <Text className="text-red-600 font-medium ml-3">Supprimer la transaction</Text>
            </TouchableOpacity>
          </View>
        </View>
      </ScrollView>
    </View>
  );
} 