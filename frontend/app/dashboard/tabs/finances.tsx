import React, { useEffect, useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';
import { SyncStatus } from '@/src/components/ui/SyncStatus';
import financeApi from '@/src/services/finance';
import Icon from '@/src/components/Icons';
import { AccountForm } from '@/src/components/forms/AccountFrom';
import { TransactionForm } from '@/src/components/forms/TransactionForm';
import accountApi from '@/src/services/accountService/accountApi';
import transactionApi from '@/src/services/transactionService/transactionApi';
import categoryApi from '@/src/services/categoryService/categoryApi';
import { Category } from '@/src/types';

export default function Finances() {
  const router = useRouter();
  const [expandedTransaction, setExpandedTransaction] = useState<string | null>(null);
  const [financeSummary, setFinanceSummary] = useState<FinanceDashboard | null>(null);
  const [showAccountForm, setShowAccountForm] = useState(false);
  const [showTransactionForm, setShowTransactionForm] = useState(false);
  const [accounts, setAccounts] = useState<any[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);

  useEffect(() => {
    const loadData = async () => {
      try {
        const [financeData, accountsData, categoriesData] = await Promise.all([
          financeApi.getFinancialSummary(),
          accountApi.getAccounts(),
          categoryApi.getCategories()
        ]);

        setFinanceSummary(financeData);
        setAccounts(accountsData);
        setCategories(categoriesData);
      } catch (error) {
        console.error('Error loading data:', error);
      }
    };

    loadData();
  }, [])

  // Section Header Component
  const SectionHeader = ({ title, onAddPress, icon = "add" }: { title: string; onAddPress: () => void; icon?: string }) => (
    <View className="flex-row justify-between items-center mb-4">
      <Text className="text-xl font-bold text-gray-900">{title}</Text>
      <TouchableOpacity
        className="bg-blue-600 px-4 py-2 rounded-lg flex-row items-center"
        onPress={onAddPress}
      >
        <Ionicons name={icon as any} size={16} color="white" />
        <Text className="text-white font-medium ml-1">Ajouter</Text>
      </TouchableOpacity>
    </View>
  );

  // Handle account creation
  const handleCreateAccount = async (accountData: any) => {
    try {
      await accountApi.createAccount(accountData);

      // Refresh the finance summary
      const updatedData = await financeApi.getFinancialSummary();
      setFinanceSummary(updatedData);

      Alert.alert('Succès', 'Compte créé avec succès');
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de la création du compte');
    }
  };

  // Handle transaction creation
  const handleCreateTransaction = async (transactionData: any) => {
    try {
      await transactionApi.createTransaction(transactionData);

      // Refresh the finance summary
      const updatedData = await financeApi.getFinancialSummary();
      setFinanceSummary(updatedData);

      Alert.alert('Succès', 'Transaction créée avec succès');
    } catch (error) {
      Alert.alert('Erreur', 'Erreur lors de la création de la transaction');
    }
  };

  // Accounts Section
  const renderAccountsSection = () => (
    <View className="mb-8">
      <SectionHeader
        title="Comptes"
        onAddPress={() => setShowAccountForm(true)}
      />
      <View className="flex-row flex-wrap gap-3">
        {financeSummary?.accounts.map((account) => (
          <TouchableOpacity
            key={account.id}
            className="bg-white p-4 rounded-lg shadow-sm flex-1 min-w-[140px]"
            onPress={() => router.push(`/dashboard/account/${account.id}`)}
          >
            <View className="flex-row items-center mb-2">
              <View className={`w-8 h-8 rounded-full items-center justify-center mr-2`}
                style={{ backgroundColor: `#${account.color?.replace('#', '')}30` }}
              >
                <Icon name={account.icon} size={16} color={account.color} />
              </View>
              <Text className="text-gray-900 font-medium">{account.name}</Text>
            </View>
            <Text className="text-2xl font-bold text-gray-900">{account.balance}</Text>
            <Text className="text-sm text-gray-500">{account.currency}</Text>
          </TouchableOpacity>
        ))}
      </View>
    </View>
  );

  // Transactions Section
  const renderTransactionsSection = () => (
    <View className="mb-8">
      <SectionHeader
        title="Transactions Récentes"
        onAddPress={() => setShowTransactionForm(true)}
      />
      <View className="flex gap-2">
        {/* Transaction 1 */}
        <TouchableOpacity
          className="bg-white p-4 rounded-lg shadow-sm"
          onPress={() => setExpandedTransaction(expandedTransaction === '1' ? null : '1')}
        >
          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center flex-1">
              <View className="w-10 h-10 bg-red-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="restaurant" size={18} color="#dc2626" />
              </View>
              <View className="flex-1">
                <Text className="text-gray-900 font-medium">Supermarché</Text>
                <Text className="text-sm text-gray-500">Alimentation • Aujourd'hui</Text>
              </View>
            </View>
            <View className="items-end">
              <Text className="text-red-600 font-semibold">-$45.20</Text>
              <Ionicons
                name={expandedTransaction === '1' ? "chevron-up" : "chevron-down"}
                size={16}
                color="#9ca3af"
              />
            </View>
          </View>
          {expandedTransaction === '1' && (
            <View className="mt-4 pt-4 border-t border-gray-100">
              <View className="space-y-2">
                <View className="flex-row justify-between">
                  <Text className="text-gray-600">Compte:</Text>
                  <Text className="text-gray-900">Compte Courant</Text>
                </View>
                <View className="flex-row justify-between">
                  <Text className="text-gray-600">Date:</Text>
                  <Text className="text-gray-900">15 Nov 2023, 14:30</Text>
                </View>
                <View className="flex-row justify-between">
                  <Text className="text-gray-600">Récurrent:</Text>
                  <Text className="text-gray-900">Non</Text>
                </View>
                <View className="flex-row justify-between">
                  <Text className="text-gray-600">Description:</Text>
                  <Text className="text-gray-900">Courses hebdomadaires</Text>
                </View>
              </View>
            </View>
          )}
        </TouchableOpacity>

        {/* Transaction 2 */}
        <TouchableOpacity
          className="bg-white p-4 rounded-lg shadow-sm"
          onPress={() => setExpandedTransaction(expandedTransaction === '2' ? null : '2')}
        >
          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center flex-1">
              <View className="w-10 h-10 bg-green-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="briefcase" size={18} color="#16a34a" />
              </View>
              <View className="flex-1">
                <Text className="text-gray-900 font-medium">Salaire</Text>
                <Text className="text-sm text-gray-500">Revenus • 1er Nov</Text>
              </View>
            </View>
            <View className="items-end">
              <Text className="text-green-600 font-semibold">+$2,500</Text>
              <Ionicons
                name={expandedTransaction === '2' ? "chevron-up" : "chevron-down"}
                size={16}
                color="#9ca3af"
              />
            </View>
          </View>
          {expandedTransaction === '2' && (
            <View className="mt-4 pt-4 border-t border-gray-100">
              <View className="space-y-2">
                <View className="flex-row justify-between">
                  <Text className="text-gray-600">Compte:</Text>
                  <Text className="text-gray-900">Compte Courant</Text>
                </View>
                <View className="flex-row justify-between">
                  <Text className="text-gray-600">Date:</Text>
                  <Text className="text-gray-900">01 Nov 2023, 09:00</Text>
                </View>
                <View className="flex-row justify-between">
                  <Text className="text-gray-600">Récurrent:</Text>
                  <Text className="text-gray-900">Mensuel</Text>
                </View>
              </View>
            </View>
          )}
        </TouchableOpacity>
      </View>
    </View>
  );

  // Budgets Section
  const renderBudgetsSection = () => (
    <View className="mb-8">
      <SectionHeader
        title="Budgets"
        onAddPress={() => Alert.alert('Ajouter un budget', 'Fonctionnalité à implémenter')}
      />
      <View className="flex gap-2">
        {/* Budget Item 1 */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <View className="flex-row justify-between items-center mb-2">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-red-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="restaurant" size={16} color="#dc2626" />
              </View>
              <Text className="text-gray-900 font-medium">Alimentation</Text>
            </View>
            <Text className="text-gray-600">$450 / $500</Text>
          </View>
          <View className="w-full bg-gray-200 rounded-full h-2">
            <View className="bg-orange-500 h-2 rounded-full" style={{ width: '90%' }} />
          </View>
          <Text className="text-xs text-gray-500 mt-1">90% utilisé</Text>
        </View>

        {/* Budget Item 2 */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <View className="flex-row justify-between items-center mb-2">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-blue-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="car" size={16} color="#2563eb" />
              </View>
              <Text className="text-gray-900 font-medium">Transport</Text>
            </View>
            <Text className="text-gray-600">$280 / $400</Text>
          </View>
          <View className="w-full bg-gray-200 rounded-full h-2">
            <View className="bg-blue-500 h-2 rounded-full" style={{ width: '70%' }} />
          </View>
          <Text className="text-xs text-gray-500 mt-1">70% utilisé</Text>
        </View>
      </View>
    </View>
  );

  // Saving Goals Section
  const renderSavingGoalsSection = () => (
    <View className="mb-8">
      <SectionHeader
        title="Objectifs d'Épargne"
        onAddPress={() => Alert.alert('Ajouter un objectif', 'Fonctionnalité à implémenter')}
        icon="flag"
      />
      <View className="flex gap-2">
        {/* Saving Goal 1 */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <View className="flex-row justify-between items-center mb-2">
            <Text className="text-gray-900 font-medium">Vacances d'été</Text>
            <Text className="text-gray-600">$1,200 / $2,000</Text>
          </View>
          <View className="w-full bg-gray-200 rounded-full h-2">
            <View className="bg-green-500 h-2 rounded-full" style={{ width: '60%' }} />
          </View>
          <View className="flex-row justify-between mt-2">
            <Text className="text-xs text-gray-500">60% atteint</Text>
            <Text className="text-xs text-gray-500">Échéance: Juin 2024</Text>
          </View>
        </View>

        {/* Saving Goal 2 */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <View className="flex-row justify-between items-center mb-2">
            <Text className="text-gray-900 font-medium">Fonds d'urgence</Text>
            <Text className="text-gray-600">$3,500 / $5,000</Text>
          </View>
          <View className="w-full bg-gray-200 rounded-full h-2">
            <View className="bg-blue-500 h-2 rounded-full" style={{ width: '70%' }} />
          </View>
          <View className="flex-row justify-between mt-2">
            <Text className="text-xs text-gray-500">70% atteint</Text>
            <Text className="text-xs text-gray-500">Échéance: Déc 2024</Text>
          </View>
        </View>
      </View>
    </View>
  );

  // Debts Section
  const renderDebtsSection = () => (
    <View className="mb-8">
      <SectionHeader
        title="Dettes"
        onAddPress={() => Alert.alert('Ajouter une dette', 'Fonctionnalité à implémenter')}
        icon="warning"
      />
      <View className="flex gap-2">
        {/* Debt 1 */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <View className="flex-row justify-between items-center mb-2">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-red-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="card" size={16} color="#dc2626" />
              </View>
              <Text className="text-gray-900 font-medium">Carte de Crédit</Text>
            </View>
            <Text className="text-red-600 font-semibold">$1,250</Text>
          </View>
          <View className="flex-row justify-between">
            <Text className="text-xs text-gray-500">Taux: 18.9%</Text>
            <Text className="text-xs text-gray-500">Paiement min: $50/mois</Text>
          </View>
        </View>

        {/* Debt 2 */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <View className="flex-row justify-between items-center mb-2">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-orange-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="home" size={16} color="#f59e0b" />
              </View>
              <Text className="text-gray-900 font-medium">Prêt Étudiant</Text>
            </View>
            <Text className="text-red-600 font-semibold">$15,000</Text>
          </View>
          <View className="flex-row justify-between">
            <Text className="text-xs text-gray-500">Taux: 5.2%</Text>
            <Text className="text-xs text-gray-500">Paiement: $200/mois</Text>
          </View>
        </View>
      </View>
    </View>
  );

  // Recurring Transactions Section
  const renderRecurringSection = () => (
    <View className="mb-8">
      <SectionHeader
        title="Transactions Récurrentes"
        onAddPress={() => Alert.alert('Ajouter une récurrence', 'Fonctionnalité à implémenter')}
        icon="repeat"
      />

      {/* Recurring Expenses */}
      <View className="mb-4">
        <Text className="text-lg font-semibold text-gray-900 mb-3">Dépenses Récurrentes</Text>
        <View className="bg-white p-4 rounded-lg shadow-sm flex gap-2">
          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-red-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="home" size={16} color="#dc2626" />
              </View>
              <View>
                <Text className="text-gray-900 font-medium">Loyer</Text>
                <Text className="text-xs text-gray-500">Mensuel</Text>
              </View>
            </View>
            <Text className="text-red-600 font-semibold">-$800</Text>
          </View>

          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-yellow-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="flash" size={16} color="#f59e0b" />
              </View>
              <View>
                <Text className="text-gray-900 font-medium">Électricité</Text>
                <Text className="text-xs text-gray-500">Mensuel</Text>
              </View>
            </View>
            <Text className="text-red-600 font-semibold">-$120</Text>
          </View>
        </View>
      </View>

      {/* Recurring Income */}
      <View>
        <Text className="text-lg font-semibold text-gray-900 mb-3">Revenus Récurrents</Text>
        <View className="bg-white p-4 rounded-lg shadow-sm flex gap-2">
          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-green-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="briefcase" size={16} color="#16a34a" />
              </View>
              <View>
                <Text className="text-gray-900 font-medium">Salaire</Text>
                <Text className="text-xs text-gray-500">Mensuel</Text>
              </View>
            </View>
            <Text className="text-green-600 font-semibold">+$2,500</Text>
          </View>

          <View className="flex-row justify-between items-center">
            <View className="flex-row items-center">
              <View className="w-8 h-8 bg-blue-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="laptop" size={16} color="#2563eb" />
              </View>
              <View>
                <Text className="text-gray-900 font-medium">Freelance</Text>
                <Text className="text-xs text-gray-500">Hebdomadaire</Text>
              </View>
            </View>
            <Text className="text-green-600 font-semibold">+$200</Text>
          </View>
        </View>
      </View>
    </View>
  );

  return (
    <ScrollView className="flex-1 bg-gray-50">
      <View className="p-4 pt-12">
        {/* Header */}
        <View className="flex-row justify-between items-center mb-6">
          <Text className="text-2xl font-bold text-gray-900">Finances</Text>
          <View className="flex-row gap-4">
            <TouchableOpacity>
              <Text className="text-blue-600 font-medium">Statistiques</Text>
            </TouchableOpacity>
            <TouchableOpacity>
              <Text className="text-blue-600 font-medium">Objectifs</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* All Sections */}
        {renderAccountsSection()}
        {renderTransactionsSection()}
        {renderBudgetsSection()}
        {renderSavingGoalsSection()}
        {renderDebtsSection()}
        {renderRecurringSection()}
      </View>

      {/* Account Form Modal */}
      <AccountForm
        visible={showAccountForm}
        onClose={() => setShowAccountForm(false)}
        onSubmit={(account) => handleCreateAccount(account)}
      />

      {/* Transaction Form Modal */}
      <TransactionForm
        visible={showTransactionForm}
        onClose={() => setShowTransactionForm(false)}
        onSubmit={handleCreateTransaction}
        categories={categories}
        accounts={accounts}
      />
    </ScrollView>
  );
} 