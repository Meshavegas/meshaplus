import React from 'react'
import { View, Text, TouchableOpacity, ActivityIndicator } from 'react-native'
import { Ionicons } from '@expo/vector-icons'
import { useOfflineSync } from '@/src/services/sync/offlineSync'

interface SyncStatusProps {
  showDetails?: boolean
  onSyncPress?: () => void
}

export const SyncStatus: React.FC<SyncStatusProps> = ({ 
  showDetails = false, 
  onSyncPress 
}) => {
  const { 
    isOnline, 
    lastSync, 
    pendingChanges, 
    isSyncing, 
    syncProgress,
    forceSync 
  } = useOfflineSync()

  const handleSyncPress = () => {
    if (onSyncPress) {
      onSyncPress()
    } else {
      forceSync()
    }
  }

  const getStatusColor = () => {
    if (isSyncing) return 'blue'
    if (!isOnline) return 'orange'
    if (pendingChanges > 0) return 'yellow'
    return 'green'
  }

  const getStatusText = () => {
    if (isSyncing) return 'Synchronisation...'
    if (!isOnline) return 'Hors ligne'
    if (pendingChanges > 0) return `${pendingChanges} modification(s) en attente`
    return 'Synchronisé'
  }

  const getStatusIcon = () => {
    if (isSyncing) return 'sync'
    if (!isOnline) return 'cloud-offline'
    if (pendingChanges > 0) return 'cloud-upload'
    return 'cloud-done'
  }

  return (
    <View className="bg-white rounded-lg p-4 shadow-sm">
      <View className="flex-row items-center justify-between">
        <View className="flex-row items-center">
          <View className={`w-3 h-3 rounded-full mr-2 bg-${getStatusColor()}-500`} />
          {/* <Text className="text-sm font-medium text-gray-700">
            {getStatusText()}
          </Text> */}
        </View>
        
        <TouchableOpacity 
          onPress={handleSyncPress}
          disabled={isSyncing}
          className={`p-2 rounded-full ${isSyncing ? 'opacity-50' : ''}`}
        >
          {isSyncing ? (
            <ActivityIndicator size="small" color="#3b82f6" />
          ) : (
            <Ionicons 
              name={getStatusIcon() as any} 
              size={20} 
              color={getStatusColor() === 'green' ? '#10b981' : '#f59e0b'} 
            />
          )}
        </TouchableOpacity>
      </View>

      {showDetails && (
        <View className="mt-3 pt-3 border-t border-gray-100">
          <View className="flex-row justify-between items-center mb-2">
            <Text className="text-xs text-gray-500">Connexion</Text>
            <Text className={`text-xs font-medium ${isOnline ? 'text-green-600' : 'text-orange-600'}`}>
              {isOnline ? 'En ligne' : 'Hors ligne'}
            </Text>
          </View>
          
          {lastSync && (
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-xs text-gray-500">Dernière sync</Text>
              <Text className="text-xs text-gray-600">
                {new Date(lastSync).toLocaleString()}
              </Text>
            </View>
          )}
          
          {pendingChanges > 0 && (
            <View className="flex-row justify-between items-center mb-2">
              <Text className="text-xs text-gray-500">En attente</Text>
              <Text className="text-xs font-medium text-yellow-600">
                {pendingChanges} modification(s)
              </Text>
            </View>
          )}
          
          {isSyncing && syncProgress > 0 && (
            <View className="mt-2">
              <View className="flex-row justify-between items-center mb-1">
                <Text className="text-xs text-gray-500">Progression</Text>
                <Text className="text-xs text-gray-600">{syncProgress}%</Text>
              </View>
              <View className="w-full bg-gray-200 rounded-full h-1">
                <View 
                  className="bg-blue-500 h-1 rounded-full" 
                  style={{ width: `${syncProgress}%` }} 
                />
              </View>
            </View>
          )}
        </View>
      )}
    </View>
  )
}

export default SyncStatus 