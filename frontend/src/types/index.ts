// Types de base pour l'application

export interface User {
  id: string
  email: string
  name: string
  avatar?: string
  createdAt: Date
  updatedAt: Date
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