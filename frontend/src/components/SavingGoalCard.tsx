import React, { useMemo } from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { SavingGoal } from '@/src/services/savingGoalService/savingGoalApi';

interface SavingGoalCardProps {
  goal: SavingGoal;
  onPress?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  onContribute?: () => void;
}

export default function SavingGoalCard({ 
  goal, 
  onPress, 
  onEdit, 
  onDelete, 
  onContribute 
}: SavingGoalCardProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return '#10b981';
      case 'overdue':
        return '#ef4444';
      case 'completed':
        return '#3b82f6';
      default:
        return '#6b7280';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'active':
        return 'En cours';
      case 'overdue':
        return 'En retard';
      case 'completed':
        return 'Atteint';
      default:
        return 'Inconnu';
    }
  };

  const getProgressColor = (percentage: number) => {
    if (percentage >= 100) return '#3b82f6'; // Bleu pour atteint
    if (percentage >= 75) return '#10b981'; // Vert pour > 75%
    if (percentage >= 50) return '#f59e0b'; // Orange pour > 50%
    return '#ef4444'; // Rouge pour < 50%
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: 'EUR'
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('fr-FR', {
      year: 'numeric',
      month: 'short'
    });
  };

  const formatFrequency = (frequency: string) => {
    switch (frequency) {
      case 'weekly':
        return 'Hebdomadaire';
      case 'monthly':
        return 'Mensuel';
      case 'yearly':
        return 'Annuel';
      default:
        return frequency;
    }
  };
  const percentage = useMemo(() => {
    return (goal?.current_amount??0 / goal?.target_amount) * 100;
  }, [goal]);

  const getDaysRemaining = () => {
    const now = new Date();
    const deadline = new Date(goal.deadline);
    const diffTime = deadline.getTime() - now.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  };

  const getUrgencyIcon = () => {
    const daysRemaining = getDaysRemaining();
    if (daysRemaining <= 7) return 'warning';
    if (daysRemaining <= 30) return 'time';
    return 'calendar';
  };

  const getUrgencyColor = () => {
    const daysRemaining = getDaysRemaining();
    if (daysRemaining <= 7) return '#ef4444';
    if (daysRemaining <= 30) return '#f59e0b';
    return '#6b7280';
  };

  return (
    <TouchableOpacity
      className="bg-white p-4 rounded-lg shadow-sm mb-3"
      onPress={onPress}
      activeOpacity={0.7}
    >
      {/* Header */}
      <View className="flex-row justify-between items-start mb-3">
        <View className="flex-1">
          <Text className="text-lg font-bold text-gray-900">{goal.title}</Text>
          <Text className="text-sm text-gray-500">
            {formatCurrency(goal?.target_amount)} • {formatFrequency(goal?.frequency)}
          </Text>
        </View>
        
        {/* Status Badge */}
        <View 
          className="px-2 py-1 rounded-full w-20 h-8 flex items-center justify-center"
          style={{ backgroundColor: `${getStatusColor(goal?.status)}20` }}
        >
          <Text 
            className="text-xs font-medium"
            style={{ color: getStatusColor(goal?.status) }}
          >
            {getStatusText(goal?.status)}
          </Text>
        </View>
      </View>

      {/* Progress Bar */}
      <View className="mb-3">
        <View className="flex-row justify-between items-center mb-1">
          <Text className="text-sm text-gray-600">
            {formatCurrency(goal?.current_amount)} / {formatCurrency(goal?.target_amount)}
          </Text>
          <Text className="text-sm font-medium text-gray-900">
            {goal?.percentage_achieved?.toFixed(0)}%
          </Text>
        </View>
        
        <View className="w-full bg-gray-200 rounded-full h-2">
          <View 
            className="h-2 rounded-full"
            style={{ 
              width: `${Math.min(percentage, 100)}%`,
              backgroundColor: getProgressColor(percentage)
            }}
          />
        </View>
      </View>

      {/* Details */}
      <View className="flex-row justify-between items-center mb-3">
        <View className="flex-row items-center">
          <Ionicons 
            name={getUrgencyIcon() as any} 
            size={16} 
            color={getUrgencyColor()} 
          />
          <Text className="text-sm text-gray-600 ml-1">
            {getDaysRemaining()} jours restants
          </Text>
        </View>

        <Text className="text-sm text-gray-500">
          Échéance: {formatDate(goal.deadline)}
        </Text>
      </View>

      {/* Remaining Amount */}
      {goal?.target_amount - (goal?.current_amount??0) > 0 && (
        <View className="mb-3 p-3 bg-blue-50 rounded-lg">
          <Text className="text-sm text-blue-800">
            Reste à épargner : <Text className="font-medium">
              {formatCurrency(goal?.target_amount - (goal?.current_amount??0))}
            </Text>
          </Text>
        </View>
      )}

      {/* Actions */}
      <View className="flex-row justify-between items-center pt-2 border-t border-gray-100 pt-2">
        <View className="flex-row items-center pt-2">
          {onContribute && (
            <TouchableOpacity
              onPress={onContribute}
              className="flex-row items-center bg-green-100 px-2 py-1 rounded-full mr-2 pt-2"
            >
              <Ionicons name="add-circle-outline" size={16} color="#16a34a" />
              <Text className="text-sm font-medium text-green-700 ml-1">
                Contribuer
              </Text>
            </TouchableOpacity>
          )}
        </View>

        <View className="flex-row items-center">
          {onEdit && (
            <TouchableOpacity
              onPress={onEdit}
              className="p-1 mr-2 pt-2"
              hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
            >
              <Ionicons name="create-outline" size={16} color="#3b82f6" />
            </TouchableOpacity>
          )}
          
          {onDelete && (
            <TouchableOpacity
              onPress={onDelete}
              className="p-1"
              hitSlop={{ top: 10, bottom: 10, left: 10, right: 10 }}
            >
              <Ionicons name="trash-outline" size={16} color="#ef4444" />
            </TouchableOpacity>
          )}
        </View>
      </View>

      {/* Celebration for completed goals */}
      {goal.status === 'completed' && (
        <View className="mt-3 p-3 bg-green-50 rounded-lg">
          <View className="flex-row items-center">
            <Ionicons name="trophy" size={20} color="#16a34a" />
            <Text className="text-sm font-medium text-green-800 ml-2">
              Félicitations ! Objectif atteint ! 🎉
            </Text>
          </View>
        </View>
      )}

      {/* Warning for overdue goals */}
      {goal.status === 'overdue' && (
        <View className="mt-3 p-3 bg-red-50 rounded-lg">
          <View className="flex-row items-center">
            <Ionicons name="warning-outline" size={20} color="#ef4444" />
            <Text className="text-sm font-medium text-red-800 ml-2">
              Objectif en retard - {formatCurrency(goal.remaining_amount)} restants
            </Text>
          </View>
        </View>
      )}
    </TouchableOpacity>
  );
} 