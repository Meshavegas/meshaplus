import React from 'react';
import { View, Text, TouchableOpacity, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';

interface AccountActionsProps {
  onAddTransaction: () => void;
  onTransfer: () => void;
  onExport: () => void;
  onShare: () => void;
}

export default function AccountActions({
  onAddTransaction,
  onTransfer,
  onExport,
  onShare
}: AccountActionsProps) {
  const actions = [
    {
      id: 'add-transaction',
      title: 'Nouvelle Transaction',
      icon: 'add',
      color: '#3b82f6',
      onPress: onAddTransaction
    },
    {
      id: 'transfer',
      title: 'Transfert',
      icon: 'swap-horizontal',
      color: '#10b981',
      onPress: onTransfer
    },
    {
      id: 'export',
      title: 'Exporter',
      icon: 'download',
      color: '#f59e0b',
      onPress: onExport
    },
    {
      id: 'share',
      title: 'Partager',
      icon: 'share-social',
      color: '#8b5cf6',
      onPress: onShare
    }
  ];

  return (
    <View className="bg-white p-4 rounded-lg shadow-sm mb-6">
      <Text className="text-lg font-bold text-gray-900 mb-4">Actions rapides</Text>
      <View className="flex flex-row flex-wrap gap-3 justify-between" style={{ width: '100%' }}>
        {actions.map((action) => (
          <TouchableOpacity
            key={action.id}
            className=" min-w-[100px] bg-gray-50 p-3 rounded-lg items-center"
            style={{ width: '30%', flex: 1 }}
            onPress={action.onPress}
          >
            <View 
              className="w-10 h-10 rounded-full items-center justify-center mb-2"
              style={{ backgroundColor: `${action.color}20` }}
            >
              <Ionicons name={action.icon as any} size={20} color={action.color} />
            </View>
            <Text className="text-xs text-gray-700 text-center font-medium">
              {action.title}
            </Text>
          </TouchableOpacity>
        ))}
      </View>
    </View>
  );
} 