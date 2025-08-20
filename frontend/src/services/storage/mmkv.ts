import { MMKV } from 'react-native-mmkv'

// Configuration MMKV pour le stockage local rapide
export const storage = new MMKV({
  id: 'meshaplus-storage',
  encryptionKey: 'meshaplus-encryption-key',
})

// Clés de stockage
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

// Helpers pour MMKV
export const mmkvHelpers = {
  // String
  setString: (key: string, value: string) => {
    storage.set(key, value)
  },
  getString: (key: string) => {
    return storage.getString(key)
  },
  
  // Number
  setNumber: (key: string, value: number) => {
    storage.set(key, value)
  },
  getNumber: (key: string) => {
    return storage.getNumber(key)
  },
  
  // Boolean
  setBoolean: (key: string, value: boolean) => {
    storage.set(key, value)
  },
  getBoolean: (key: string) => {
    return storage.getBoolean(key)
  },
  
  // Object
  setObject: <T>(key: string, value: T) => {
    storage.set(key, JSON.stringify(value))
  },
  getObject: <T>(key: string): T | null => {
    const value = storage.getString(key)
    if (!value) return null
    try {
      return JSON.parse(value) as T
    } catch {
      return null
    }
  },
  
  // Array
  setArray: <T>(key: string, value: T[]) => {
    storage.set(key, JSON.stringify(value))
  },
  getArray: <T>(key: string): T[] => {
    const value = storage.getString(key)
    if (!value) return []
    try {
      return JSON.parse(value) as T[]
    } catch {
      return []
    }
  },
  
  // Delete
  delete: (key: string) => {
    storage.delete(key)
  },
  
  // Clear all
  clearAll: () => {
    storage.clearAll()
  },
  
  // Check if key exists
  contains: (key: string) => {
    return storage.contains(key)
  },
  
  // Get all keys
  getAllKeys: () => {
    return storage.getAllKeys()
  },
} as const

// Types pour les données stockées
export interface StoredData {
  expenses: any[]
  revenues: any[]
  habits: any[]
  habitLogs: any[]
  goals: any[]
  tasks: any[]
  settings: any
  offlineQueue: any[]
  lastSync: string | null
}

export default storage 