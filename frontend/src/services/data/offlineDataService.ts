import { storage, STORAGE_KEYS, mmkvHelpers } from '@/src/services/storage/mmkv'
import { offlineSyncService } from '@/src/services/sync/offlineSync'

export interface OfflineData<T> {
  data: T[]
  lastSync: string | null
  version: number
}

export interface DataEntity {
  id: string
  createdAt: string
  updatedAt: string
}

class OfflineDataService {
  // Obtenir les données depuis le cache local
  public getLocalData<T extends DataEntity>(entity: string): T[] {
    const data = mmkvHelpers.getObject<OfflineData<T>>(`${entity}_data`)
    return data?.data || []
  }

  // Sauvegarder les données en local
  public saveLocalData<T extends DataEntity>(entity: string, data: T[]): void {
    const offlineData: OfflineData<T> = {
      data,
      lastSync: new Date().toISOString(),
      version: Date.now(),
    }
    mmkvHelpers.setObject(`${entity}_data`, offlineData)
  }

  // Ajouter un élément en local et en attente de sync
  public addItem<T extends DataEntity>(
    entity: string, 
    item: Omit<T, 'id' | 'createdAt' | 'updatedAt'>
  ): T {
    const localData = this.getLocalData<T>(entity)
    
    const newItem: T = {
      ...item,
      id: `local_${Date.now()}_${Math.random()}`,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    } as T

    // Ajouter en local
    localData.push(newItem)
    this.saveLocalData(entity, localData)

    // Ajouter à la file d'attente de sync
    offlineSyncService.addPendingAction({
      type: 'CREATE',
      entity: entity as any,
      data: newItem,
    })

    return newItem
  }

  // Mettre à jour un élément en local et en attente de sync
  public updateItem<T extends DataEntity>(
    entity: string, 
    id: string, 
    updates: Partial<T>
  ): T | null {
    const localData = this.getLocalData<T>(entity)
    const index = localData.findIndex(item => item.id === id)
    
    if (index === -1) return null

    const updatedItem: T = {
      ...localData[index],
      ...updates,
      updatedAt: new Date().toISOString(),
    }

    // Mettre à jour en local
    localData[index] = updatedItem
    this.saveLocalData(entity, localData)

    // Ajouter à la file d'attente de sync
    offlineSyncService.addPendingAction({
      type: 'UPDATE',
      entity: entity as any,
      data: updatedItem,
    })

    return updatedItem
  }

  // Supprimer un élément en local et en attente de sync
  public deleteItem<T extends DataEntity>(entity: string, id: string): boolean {
    const localData = this.getLocalData<T>(entity)
    const index = localData.findIndex(item => item.id === id)
    
    if (index === -1) return false

    const deletedItem = localData[index]

    // Supprimer en local
    localData.splice(index, 1)
    this.saveLocalData(entity, localData)

    // Ajouter à la file d'attente de sync
    offlineSyncService.addPendingAction({
      type: 'DELETE',
      entity: entity as any,
      data: deletedItem,
    })

    return true
  }

  // Obtenir un élément par ID
  public getItem<T extends DataEntity>(entity: string, id: string): T | null {
    const localData = this.getLocalData<T>(entity)
    return localData.find(item => item.id === id) || null
  }

  // Rechercher des éléments avec filtres
  public searchItems<T extends DataEntity>(
    entity: string, 
    filters: Partial<T> = {}
  ): T[] {
    const localData = this.getLocalData<T>(entity)
    
    return localData.filter(item => {
      return Object.entries(filters).every(([key, value]) => {
        if (value === undefined || value === null) return true
        return item[key as keyof T] === value
      })
    })
  }

  // Synchroniser les données depuis le serveur
  public async syncFromServer<T extends DataEntity>(
    entity: string, 
    fetchFunction: () => Promise<T[]>
  ): Promise<void> {
    try {
      const serverData = await fetchFunction()
      this.saveLocalData(entity, serverData)
    } catch (error) {
      console.error(`Erreur lors de la synchronisation de ${entity}:`, error)
      throw error
    }
  }

  // Nettoyer les données anciennes
  public cleanupOldData(entity: string, daysToKeep: number = 30): void {
    const localData = this.getLocalData(entity)
    const cutoffDate = new Date()
    cutoffDate.setDate(cutoffDate.getDate() - daysToKeep)

    const filteredData = localData.filter(item => {
      const itemDate = new Date(item.updatedAt)
      return itemDate > cutoffDate
    })

    if (filteredData.length !== localData.length) {
      this.saveLocalData(entity, filteredData)
    }
  }

  // Obtenir les statistiques de cache
  public getCacheStats(): Record<string, { count: number; lastSync: string | null }> {
    const entities = ['expenses', 'revenues', 'habits', 'goals', 'tasks']
    const stats: Record<string, { count: number; lastSync: string | null }> = {}

    entities.forEach(entity => {
      const data = mmkvHelpers.getObject<OfflineData<any>>(`${entity}_data`)
      stats[entity] = {
        count: data?.data?.length || 0,
        lastSync: data?.lastSync || null,
      }
    })

    return stats
  }

  // Vider le cache
  public clearCache(entity?: string): void {
    if (entity) {
      storage.delete(`${entity}_data`)
    } else {
      const entities = ['expenses', 'revenues', 'habits', 'goals', 'tasks']
      entities.forEach(entity => {
        storage.delete(`${entity}_data`)
      })
    }
  }
}

export const offlineDataService = new OfflineDataService()
export default offlineDataService 