import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, TouchableOpacity, StatusBar, Alert, ActivityIndicator } from 'react-native';
import { Icon, FONTAWESOME_6 } from '@/src/components/Icons';
import { SyncStatus } from '@/src/components/ui/SyncStatus';
import CircularProgress from '@/src/components/ui/CircularProgress';
import { useAuthStore } from '@/src';
import { TransactionForm } from '@/src/components/forms/TransactionForm';
import transactionApi from '@/src/services/transactionService/transactionApi';
import categoryApi from '@/src/services/categoryService/categoryApi';
import accountApi from '@/src/services/accountService/accountApi';
import taskApi from '@/src/services/taskService/taskApi';
import budgetApi from '@/src/services/budgetService/budgetApi';
import savingGoalApi from '@/src/services/savingGoalService/savingGoalApi';
import dashboardApi from '@/src/services/dashboardService/dashboardApi';
import {
  CreateTransactionRequest,
  Category,
  Account,
  Task,
  Budget,
  SavingGoal,
  DashboardData,
  TaskStats,
  TransactionStats,
  BudgetStats
} from '@/src/types';
import { TaskForm } from '@/src/components/forms/TaskForm';
import { Link, Redirect } from 'expo-router';

export default function Overview() {
  const { user, require_preferences } = useAuthStore()

  const [showTransactionForm, setShowTransactionForm] = useState(false)
  const [showTaskForm, setShowTaskForm] = useState(false)
  const [transactionType, setTransactionType] = useState<'income' | 'expense'>('expense')
  const [categories, setCategories] = useState<Category[]>([])
  const [accounts, setAccounts] = useState<Account[]>([])
  const [tasks, setTasks] = useState<Task[]>([])
  const [recentTransactions, setRecentTransactions] = useState<any[]>([])
  const [taskStats, setTaskStats] = useState<TaskStats | null>(null)
  const [transactionStats, setTransactionStats] = useState<TransactionStats | null>(null)
  const [budgetStats, setBudgetStats] = useState<BudgetStats | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  // Charger les données au montage du composant
  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setIsLoading(true)
      const [
        categoriesData,
        accountsData,
        tasksData,
        transactionsData,
        // taskStatsData,
        // transactionStatsData,
        budgetStatsData
      ] = await Promise.all([
        categoryApi.getCategories(),
        accountApi.getAccounts(),
        taskApi.getTasks(),
        transactionApi.getTransactions(),
        // taskApi.getTaskStats(),
        // dashboardApi.getTransactionStats(),
        budgetApi.getBudgets()
      ])

      console.log("categoriesData", JSON.stringify(budgetStatsData, null, 2))
      setCategories(categoriesData)
      setAccounts(accountsData)
      setTasks(tasksData)
      setRecentTransactions((transactionsData ?? []).slice(0, 5)) // Limiter à 5 transactions récentes
      // setTaskStats(taskStatsData)
      // setTransactionStats(transactionStatsData)
      // setBudgetStats(budgetStatsData)
    } catch (error) {
      console.error('Error loading data:', error)
      Alert.alert('Erreur', 'Impossible de charger les données')
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreateTransaction = async (transaction: CreateTransactionRequest) => {
    try {
      await transactionApi.createTransaction(transaction)
      Alert.alert('Succès', 'Transaction créée avec succès')
      // Rafraîchir les données après création
      loadData()
    } catch (error) {
      console.error('Error creating transaction:', error)
      throw error
    }
  }

  const handleOpenTransactionForm = (type: 'income' | 'expense') => {
    setTransactionType(type)
    setShowTransactionForm(true)
  }

  const handleCreateTask = async (task: any) => {
    try {
      await taskApi.createTask(task)
      Alert.alert('Succès', 'Tâche créée avec succès')
      // Rafraîchir les données après création
      loadData()
    } catch (error) {
      console.error('Error creating task:', error)
      throw error
    }
  }

  // Catégories de tâches
  const taskCategories = [
    { id: 'work', name: 'Travail', color: '#3b82f6', icon: 'fa6:briefcase' },
    { id: 'personal', name: 'Personnel', color: '#10b981', icon: 'fa6:user' },
    { id: 'health', name: 'Santé', color: '#ef4444', icon: 'fa6:heart-pulse' },
    { id: 'finance', name: 'Finance', color: '#f59e0b', icon: 'fa6:wallet' },
    { id: 'education', name: 'Éducation', color: '#8b5cf6', icon: 'fa6:graduation-cap' },
    { id: 'home', name: 'Maison', color: '#06b6d4', icon: 'fa6:house' },
  ]

  if (isLoading) {
    return (
      <View className="flex-1 bg-gray-50 mt-safe justify-center items-center">
        <ActivityIndicator size="large" color="#0ea5e9" />
        <Text className="mt-4 text-gray-600">Chargement des données...</Text>
      </View>
    )
  }
   

  if (require_preferences) {
    return <Redirect href="/dashboard/wizard" />
  }

  return (
    <ScrollView className="flex-1 bg-gray-50 mt-safe">
      <StatusBar barStyle="dark-content" backgroundColor="#ffffff" />
      <View className="p-4">
        <View className='flex-row items-center justify-between gap-2'>
          <View className='flex-row items-center gap-2'>
            <View className='h-12 w-12 bg-primary items-center justify-center rounded-full'>
              <Text className="text-2xl font-bold text-white">MV</Text>
            </View>
            <View className='flex-col'>
              <Text className="text-md font-bold text-gray-900">{user?.name}</Text>
              <Text className="text-xs text-gray-500">Bienvenue sur votre tableau de bord</Text>
            </View>
          </View>
          <SyncStatus />
        </View>

        {/* Daily Budget Status */}
        <View className="bg-white  rounded-lg shadow-sm mb-6 mt-3  overflow-hidden">


          <View className='flex-row items-center bg-primary/30 '>
            <CircularProgress
              progressPercent={taskStats ? (taskStats.completedTasks / taskStats.totalTasks) * 100 : 0}
              size={60}
              strokeWidth={10}
              text={`${taskStats?.completedTasks || 0}/${taskStats?.totalTasks || 0}`}
              textSize={20}
            />
            <View className=' flex flex-col gap-2'>
              <Text className="text-md font-bold text-gray-900 ">Taches</Text>
              <Text>{taskStats?.pendingTasks || 0} tâches en attente</Text>
            </View>
          </View>


          <View className="bg-teal/70 px-4 py-4">
            <View className="flex-row justify-between items-center">
              <Text className="text-sm text-white">Budget mensuel</Text>
              <Text className="text-white font-semibold">
                €{budgetStats?.totalSpent?.toFixed(0) || 0} / €{budgetStats?.totalPlanned?.toFixed(0) || 0}
              </Text>
            </View>
            <View className="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
              <View
                className="bg-teal-500 h-3 rounded-full"
                style={{
                  width: `${budgetStats && budgetStats.totalPlanned > 0
                    ? Math.min((budgetStats.totalSpent / budgetStats.totalPlanned) * 100, 100)
                    : 0}%`
                }}
              />
            </View>
            <Text className="text-xs text-white mt-1">
              €{budgetStats ? (budgetStats.totalPlanned - budgetStats.totalSpent)?.toFixed(0) : 0} restant ce mois
            </Text>
          </View>
        </View>


        <View>
          <Text className="text-md text-black">Mois en cours</Text>
          <View className="space-y-4 mb-2 mt-2 flex flex-row gap-2">

            <View className="bg-white px-4 py-2 flex flex-col justify-between rounded-lg shadow-sm flex-1">
              <Text className="text-[12px] text-gray-500">Solde Total</Text>
              <Text className="text-md font-bold text-green-600">€{transactionStats?.netAmount?.toFixed(0) || 0}</Text>
              <Text className="text-xs text-green-600">
                {transactionStats?.netAmount && transactionStats.netAmount > 0 ? '+' : ''}€{transactionStats?.netAmount?.toFixed(0) || 0} ce mois
              </Text>
            </View>

            <View className="bg-white px-4 py-2 flex flex-col justify-between rounded-lg shadow-sm flex-1">
              <Text className="text-[12px] text-gray-500">Dépenses</Text>
              <Text className="text-md font-bold text-red-600">€{transactionStats?.totalExpense?.toFixed(0) || 0}</Text>
              <Text className="text-xs text-red-600">
                {budgetStats && budgetStats.totalPlanned > 0
                  ? `${((budgetStats.totalSpent / budgetStats.totalPlanned) * 100).toFixed(0)}% du budget`
                  : '0% du budget'
                }
              </Text>
            </View>

            <View className="bg-white px-4 py-2 flex flex-col justify-between rounded-lg shadow-sm flex-1">
              <Text className="text-[12px] text-gray-500">Revenus</Text>
              <Text className="text-md font-bold text-blue-600">€{transactionStats?.totalIncome?.toFixed(0) || 0}</Text>
              <Text className="text-xs text-blue-600">
                {transactionStats?.totalIncome && transactionStats.totalIncome > 0 ? '+' : ''}€{transactionStats?.totalIncome?.toFixed(0) || 0}
              </Text>
            </View>
          </View>
        </View>

        <View className='w-full h-[200px] bg-white rounded-lg shadow-sm mb-6 items-center justify-center'>
          <Text className="text-md text-black">Statistique</Text>
        </View>

        {/* Daily Tasks */}
        <View className="bg-white p-4 rounded-lg shadow-sm mb-6">
          <Text className="text-lg font-semibold text-gray-900 mb-4">Tâches Récentes</Text>
          <View className="space-y-3">
            {tasks?.slice(0, 3).map((task, index) => (
              <View key={task.id} className="flex-row items-center">
                <View className={`w-5 h-5 rounded-full items-center justify-center mr-3 ${task.status === 'completed' ? 'bg-green-100' : 'bg-gray-200'
                  }`}>
                  {task.status === 'completed' ? (
                    <Icon name="checkmark" size={12} color="#10b981" />
                  ) : (
                    <Icon name="time" size={12} color="#9ca3af" />
                  )}
                </View>
                <Text className="text-gray-700 flex-1">{task.title}</Text>
              </View>
            ))}
            {tasks?.length === 0 && (
              <Text className="text-gray-500 text-center py-4">Aucune tâche pour le moment</Text>
            )}
          </View>
        </View>

        {/* Quick Actions */}
        <View className="space-y-3 mb-6">
          <Text className="text-lg font-semibold text-gray-900">Actions Rapides</Text>
          <View className="flex-row gap-2">
            <TouchableOpacity
              className="flex-1 bg-[#dc2626] p-4 rounded-lg items-center"
              onPress={() => handleOpenTransactionForm('expense')}
            >
              <Icon name={`fa6:${FONTAWESOME_6.CHART_LINE}`} size={20} color="white" />
              <Text className="text-white font-medium mt-1">Dépense</Text>
            </TouchableOpacity>
            <TouchableOpacity
              className="flex-1 bg-green-600 p-4 rounded-lg items-center"
              onPress={() => handleOpenTransactionForm('income')}
            >
              <Icon name="add" size={20} color="white" />
              <Text className="text-white font-medium mt-1">Revenu</Text>
            </TouchableOpacity>
            <TouchableOpacity
              className="flex-1 bg-blue-600 p-4 rounded-lg items-center"
              onPress={() => setShowTaskForm(true)}
            >
              <Icon name="flag" size={20} color="white" />
              <Text className="text-white font-medium mt-1">Taches</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Recent Activity */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold text-gray-900 mb-4">Activité Récente</Text>
          <View className="space-y-3">
            {recentTransactions?.slice(0, 3).map((transaction) => (
              <View key={transaction.id} className="flex-row justify-between items-center">
                <View className="flex-row items-center">
                  <View className={`w-8 h-8 rounded-full items-center justify-center mr-3 ${transaction.type === 'income' ? 'bg-green-100' : 'bg-red-100'
                    }`}>
                    <Icon
                      name={transaction.type === 'income' ? 'arrow-up' : 'arrow-down'}
                      size={16}
                      color={transaction.type === 'income' ? '#16a34a' : '#dc2626'}
                    />
                  </View>
                  <View>
                    <Text className="text-gray-900 font-medium">{transaction.description || 'Transaction'}</Text>
                    <Text className="text-xs text-gray-500">
                      {new Date(transaction.date).toLocaleDateString('fr-FR')}
                    </Text>
                  </View>
                </View>
                <Text className={`font-semibold ${transaction.type === 'income' ? 'text-green-600' : 'text-red-600'
                  }`}>
                  {transaction.type === 'income' ? '+' : '-'}€{transaction.amount.toFixed(2)}
                </Text>
              </View>
            ))}
            {recentTransactions.length === 0 && (
              <Text className="text-gray-500 text-center py-4">Aucune transaction récente</Text>
            )}
          </View>
        </View>
      </View>

      {/* Transaction Form Modal */}
      <TransactionForm
        visible={showTransactionForm}
        onClose={() => setShowTransactionForm(false)}
        onSubmit={handleCreateTransaction}
        categories={categories}
        accounts={accounts}
        defaultType={transactionType}
      />
      <TaskForm
        visible={showTaskForm}
        onClose={() => setShowTaskForm(false)}
        onSubmit={handleCreateTask}
        categories={taskCategories}
        defaultStatus="todo"
      />
    </ScrollView>
  );
} 