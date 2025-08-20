import React from 'react'
import NetInfo from '@react-native-community/netinfo'
import { storage, STORAGE_KEYS, mmkvHelpers } from '@/src/services/storage/mmkv'
import { apiClient } from '@/src/services/api/client'
import { queryClient } from '@/src/services/api/client'

export interface OfflineAction {
  id: string
  type: 'CREATE' | 'UPDATE' | 'DELETE'
  entity: 'expense' | 'revenue' | 'habit' | 'goal' | 'task'
  data: any
  timestamp: number
  retryCount: number
  maxRetries: number
}

export interface SyncState {
  isOnline: boolean
  lastSync: Date | null
  pendingChanges: number
  isSyncing: boolean
  syncProgress: number
}

class OfflineSyncService {
  private syncState: SyncState = {
    isOnline: true,
    lastSync: null,
    pendingChanges: 0,
    isSyncing: false,
    syncProgress: 0,
  }

  private listeners: ((state: SyncState) => void)[] = []

  constructor() {
    this.initializeSync()
  }

  // Initialiser la synchronisation
  private async initializeSync() {
    // Charger l'état initial
    const lastSyncStr = storage.getString(STORAGE_KEYS.LAST_SYNC)
    const pendingActionsStr = storage.getString(STORAGE_KEYS.OFFLINE_QUEUE)
    
    this.syncState.lastSync = lastSyncStr ? new Date(lastSyncStr) : null
    this.syncState.pendingChanges = pendingActionsStr ? JSON.parse(pendingActionsStr).length : 0

    // Écouter les changements de connectivité
    NetInfo.addEventListener((state) => {
      const isOnline = state.isConnected && state.isInternetReachable
      this.updateSyncState({ isOnline: !!isOnline })
      
      // Synchroniser automatiquement quand on revient en ligne
      if (isOnline && this.syncState.pendingChanges > 0) {
        this.syncPendingActions()
      }
    })

    // Vérifier la connectivité initiale
    const netInfo = await NetInfo.fetch()
    this.updateSyncState({ isOnline: !!(netInfo.isConnected && netInfo.isInternetReachable) })
  }

  // Mettre à jour l'état de synchronisation
  private updateSyncState(updates: Partial<SyncState>) {
    this.syncState = { ...this.syncState, ...updates }
    this.notifyListeners()
  }

  // Notifier les listeners
  private notifyListeners() {
    this.listeners.forEach(listener => listener(this.syncState))
  }

  // Ajouter un listener
  public addListener(listener: (state: SyncState) => void) {
    this.listeners.push(listener)
    return () => {
      this.listeners = this.listeners.filter(l => l !== listener)
    }
  }

  // Ajouter une action en attente
  public addPendingAction(action: Omit<OfflineAction, 'id' | 'timestamp' | 'retryCount' | 'maxRetries'>) {
    const pendingAction: OfflineAction = {
      ...action,
      id: `${action.entity}_${Date.now()}_${Math.random()}`,
      timestamp: Date.now(),
      retryCount: 0,
      maxRetries: 3,
    }

    const pendingActions = this.getPendingActions()
    pendingActions.push(pendingAction)
    this.savePendingActions(pendingActions)

    this.updateSyncState({ 
      pendingChanges: this.syncState.pendingChanges + 1 
    })

    // Synchroniser immédiatement si en ligne
    if (this.syncState.isOnline) {
      this.syncPendingActions()
    }
  }

  // Obtenir les actions en attente
  private getPendingActions(): OfflineAction[] {
    const pendingActionsStr = storage.getString(STORAGE_KEYS.OFFLINE_QUEUE)
    return pendingActionsStr ? JSON.parse(pendingActionsStr) : []
  }

  // Sauvegarder les actions en attente
  private savePendingActions(actions: OfflineAction[]) {
    storage.set(STORAGE_KEYS.OFFLINE_QUEUE, JSON.stringify(actions))
  }

  // Synchroniser les actions en attente
  public async syncPendingActions(): Promise<void> {
    if (this.syncState.isSyncing || !this.syncState.isOnline) return

    this.updateSyncState({ isSyncing: true, syncProgress: 0 })

    try {
      const pendingActions = this.getPendingActions()
      if (pendingActions.length === 0) {
        this.updateSyncState({ isSyncing: false, syncProgress: 100 })
        return
      }

      const successfulActions: string[] = []
      const failedActions: OfflineAction[] = []

      for (let i = 0; i < pendingActions.length; i++) {
        const action = pendingActions[i]
        
        try {
          await this.executeAction(action)
          successfulActions.push(action.id)
        } catch (error) {
          console.error(`Erreur lors de la synchronisation de l'action ${action.id}:`, error)
          
          // Incrémenter le compteur de tentatives
          action.retryCount++
          if (action.retryCount < action.maxRetries) {
            failedActions.push(action)
          }
        }

        // Mettre à jour le progrès
        this.updateSyncState({ 
          syncProgress: Math.round(((i + 1) / pendingActions.length) * 100) 
        })
      }

      // Mettre à jour la file d'attente
      const remainingActions = failedActions
      this.savePendingActions(remainingActions)

      // Mettre à jour l'état
      const lastSync = new Date()
      storage.set(STORAGE_KEYS.LAST_SYNC, lastSync.toISOString())

      this.updateSyncState({
        lastSync,
        pendingChanges: remainingActions.length,
        isSyncing: false,
        syncProgress: 100,
      })

      // Invalider le cache pour forcer le re-fetch
      queryClient.invalidateQueries()

    } catch (error) {
      console.error('Erreur lors de la synchronisation:', error)
      this.updateSyncState({ isSyncing: false, syncProgress: 0 })
    }
  }

  // Exécuter une action
  private async executeAction(action: OfflineAction): Promise<void> {
    const endpoint = `/${action.entity}s`
    
    switch (action.type) {
      case 'CREATE':
        await apiClient.post(endpoint, action.data)
        break
      case 'UPDATE':
        await apiClient.put(`${endpoint}/${action.data.id}`, action.data)
        break
      case 'DELETE':
        await apiClient.delete(`${endpoint}/${action.data.id}`)
        break
      default:
        throw new Error(`Type d'action non supporté: ${action.type}`)
    }
  }

  // Forcer la synchronisation
  public async forceSync(): Promise<void> {
    await this.syncPendingActions()
  }

  // Obtenir l'état actuel
  public getSyncState(): SyncState {
    return { ...this.syncState }
  }

  // Vérifier la connectivité
  public async checkConnectivity(): Promise<boolean> {
    try {
      const netInfo = await NetInfo.fetch()
      const isOnline = netInfo.isConnected && netInfo.isInternetReachable
      this.updateSyncState({ isOnline: !!isOnline })
      return !!isOnline
    } catch (error) {
      console.error('Erreur lors de la vérification de la connectivité:', error)
      return false
    }
  }

  // Nettoyer les actions anciennes (plus de 7 jours)
  public cleanupOldActions(): void {
    const pendingActions = this.getPendingActions()
    const sevenDaysAgo = Date.now() - (7 * 24 * 60 * 60 * 1000)
    
    const validActions = pendingActions.filter(action => 
      action.timestamp > sevenDaysAgo
    )

    if (validActions.length !== pendingActions.length) {
      this.savePendingActions(validActions)
      this.updateSyncState({ pendingChanges: validActions.length })
    }
  }
}

// Instance singleton
export const offlineSyncService = new OfflineSyncService()

// Hook pour utiliser le service
export const useOfflineSync = () => {
  const [syncState, setSyncState] = React.useState<SyncState>(
    offlineSyncService.getSyncState()
  )

  React.useEffect(() => {
    const unsubscribe = offlineSyncService.addListener(setSyncState)
    return unsubscribe
  }, [])

  return {
    ...syncState,
    addPendingAction: offlineSyncService.addPendingAction.bind(offlineSyncService),
    syncPendingActions: offlineSyncService.syncPendingActions.bind(offlineSyncService),
    forceSync: offlineSyncService.forceSync.bind(offlineSyncService),
    checkConnectivity: offlineSyncService.checkConnectivity.bind(offlineSyncService),
  }
}

export default offlineSyncService 