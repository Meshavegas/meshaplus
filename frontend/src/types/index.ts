// Types de base pour l'application

export interface User {
  id: string
  email: string
  name: string
  avatar?: string
  created_at: Date
  updated_at: Date
}

export interface Expense {
  id: string
  userId: string
  title: string
  amount: number
  category: string
  date: Date
  description?: string
  tags?: string[]
  createdAt: Date
  updatedAt: Date
}

export interface Revenue {
  id: string
  userId: string
  title: string
  amount: number
  category: string
  date: Date
  description?: string
  createdAt: Date
  updatedAt: Date
}

export interface Habit {
  id: string
  userId: string
  title: string
  description?: string
  frequency: 'daily' | 'weekly' | 'monthly'
  target: number
  currentStreak: number
  longestStreak: number
  color: string
  icon: string
  createdAt: Date
  updatedAt: Date
}

export interface HabitLog {
  id: string
  habitId: string
  userId: string
  date: Date
  completed: boolean
  notes?: string
  createdAt: Date
}

export interface Goal {
  id: string
  userId: string
  title: string
  description?: string
  targetAmount: number
  currentAmount: number
  deadline: Date
  category: string
  color: string
  createdAt: Date
  updatedAt: Date
}

export interface Task {
  id: string
  userId: string
  title: string
  description?: string
  completed: boolean
  priority: 'low' | 'medium' | 'high'
  dueDate?: Date
  category: string
  createdAt: Date
  updatedAt: Date
}

// Types pour l'API
export interface ApiResponse<T> {
  data: T
  message: string
  success: boolean
}


export interface CreateTransactionRequest {
  account_id: string
  categoryId?: string
  type: 'income' | 'expense'
  amount: number
  description?: string
  date: Date
  recurring?: boolean
}


// Types pour les statistiques
export interface TaskStats {
  totalTasks: number
  completedTasks: number
  pendingTasks: number
  completionRate: number
}

export interface TransactionStats {
  period: string
  totalIncome: number
  totalExpense: number
  netAmount: number
  count: number
}

export interface BudgetStats {
  period: string
  totalPlanned: number
  totalSpent: number
  remaining: number
  budgetCount: number
  overBudgetCount: number
  utilizationRate: number
}

export interface FinanceOverview {
  totalRevenue: number
  totalExpense: number
  balance: number
  monthlyBudget: number
  savings: number
}

export interface TasksOverview {
  totalTasks: number
  completed: number
  pending: number
  incoming: number
  running: number
}

export interface DashboardData {
  userId: string
  tasksOverview: TasksOverview
  financeOverview: FinanceOverview
  recentTransactions: Transaction[]
  recentTasks: Task[]
}

// Types pour les tâches
export interface Task {
  id: string
  userId: string
  title: string
  description?: string
  priority: 'low' | 'medium' | 'high'
  status: 'pending' | 'running' | 'completed' | 'expired'
  dueDate?: Date
  categoryId?: string
  durationPlanned?: number
  durationSpent?: number
  createdAt: Date
  updatedAt: Date
}

// Types pour les budgets
export interface Budget {
  id: string
  userId: string
  categoryId: string
  amountPlanned: number
  amountSpent: number
  month: number
  year: number
  createdAt: Date
  updatedAt: Date
}

// Types pour les objectifs d'épargne
export interface SavingGoal {
  id: string
  userId: string
  title: string
  targetAmount: number
  currentAmount: number
  deadline?: Date
  isAchieved: boolean
  frequency: 'weekly' | 'monthly' | 'yearly'
  createdAt: Date
  updatedAt: Date
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  hasMore: boolean
}

// Types pour la navigation
export type RootStackParamList = {
  '(auth)': undefined
  '(tabs)': undefined
  modal: { id: string }
}

export type TabParamList = {
  index: undefined
  finances: undefined
  goals: undefined
  tasks: undefined
  profile: undefined
}

// Types pour les stores
export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
}

export interface SyncState {
  isOnline: boolean
  lastSync: Date | null
  pendingChanges: number
  isSyncing: boolean
} 