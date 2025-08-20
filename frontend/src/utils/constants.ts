// Constantes de l'application

export const APP_CONFIG = {
  name: 'MeshaPlus',
  version: '1.0.0',
  description: 'Application de suivi des finances et des habitudes',
} as const

export const API_CONFIG = {
  baseURL: process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api',
  timeout: 10000,
} as const

export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  USER_DATA: 'user_data',
  EXPENSES: 'expenses',
  REVENUES: 'revenues',
  HABITS: 'habits',
  HABIT_LOGS: 'habit_logs',
  GOALS: 'goals',
  TASKS: 'tasks',
  SETTINGS: 'settings',
  OFFLINE_QUEUE: 'offline_queue',
  LAST_SYNC: 'last_sync',
} as const

export const ANIMATION_CONFIG = {
  duration: {
    fast: 150,
    normal: 300,
    slow: 500,
  },
  easing: {
    ease: 'ease',
    easeIn: 'ease-in',
    easeOut: 'ease-out',
    easeInOut: 'ease-in-out',
  },
} as const

export const CHART_COLORS = {
  primary: '#0ea5e9',
  secondary: '#d946ef',
  success: '#22c55e',
  warning: '#f59e0b',
  error: '#ef4444',
  habit: {
    blue: '#3b82f6',
    green: '#10b981',
    purple: '#8b5cf6',
    orange: '#f59e0b',
    pink: '#ec4899',
    teal: '#14b8a6',
    indigo: '#6366f1',
    red: '#ef4444',
  },
} as const

export const DATE_FORMATS = {
  short: 'dd/MM/yyyy',
  long: 'dd MMMM yyyy',
  time: 'HH:mm',
  datetime: 'dd/MM/yyyy HH:mm',
} as const

export const CURRENCY_CONFIG = {
  symbol: '€',
  code: 'EUR',
  locale: 'fr-FR',
} as const

export const NOTIFICATION_CONFIG = {
  habitReminder: {
    title: 'Rappel Habitude',
    body: 'Il est temps de compléter votre habitude !',
  },
  goalReminder: {
    title: 'Rappel Objectif',
    body: 'N\'oubliez pas votre objectif !',
  },
  taskReminder: {
    title: 'Rappel Tâche',
    body: 'Vous avez une tâche à accomplir !',
  },
} as const

export const VALIDATION_RULES = {
  email: {
    required: 'L\'email est requis',
    invalid: 'L\'email n\'est pas valide',
  },
  password: {
    required: 'Le mot de passe est requis',
    minLength: 'Le mot de passe doit contenir au moins 8 caractères',
    strong: 'Le mot de passe doit contenir au moins une majuscule, une minuscule, un chiffre et un caractère spécial',
  },
  name: {
    required: 'Le nom est requis',
    minLength: 'Le nom doit contenir au moins 2 caractères',
  },
  amount: {
    required: 'Le montant est requis',
    positive: 'Le montant doit être positif',
  },
} as const

export default {
  APP_CONFIG,
  API_CONFIG,
  STORAGE_KEYS,
  ANIMATION_CONFIG,
  CHART_COLORS,
  DATE_FORMATS,
  CURRENCY_CONFIG,
  NOTIFICATION_CONFIG,
  VALIDATION_RULES,
} 