import React from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import Icon from './Icons';

interface BudgetCardProps {
    budget: Budget;
    onPress?: () => void;
    onEdit?: () => void;
    onDelete?: () => void;
}

export default function BudgetCard({ budget, onPress, onEdit, onDelete }: BudgetCardProps) {
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
            case 'good':
                return 'Actif';
            case 'overdue':
                return 'Dépassé';
            case 'completed':
                return 'Terminé';
            default:
                return status;
        }
    };

    const getProgressColor = (percentage: number) => {
        if (percentage >= 90) return '#ef4444'; // Rouge pour > 90%
        if (percentage >= 75) return '#f59e0b'; // Orange pour > 75%
        return '#10b981'; // Vert pour < 75%
    };

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('fr-FR', {
            style: 'currency',
            currency: 'EUR'
        }).format(amount);
    };

    const formatPeriod = (period: string) => {
        switch (period) {
            case 'weekly':
                return 'semaine';
            case 'monthly':
                return 'mois';
            case 'yearly':
                return 'année';
            default:
                return "mois";
        }
    };

    return (
        <TouchableOpacity
            className="bg-white p-4 rounded-lg shadow-sm mb-3 w-full"
            onPress={onPress}
            activeOpacity={0.7}
        >
            {/* Header */}
            <View className="flex-row justify-between items-start mb-3">
                <View className="flex-1">
                    <View className="flex-row items-center">
                        <View
                            className="w-10 h-10 rounded-full items-center justify-center mr-2"
                            style={{ backgroundColor: `${budget.category?.color ?? '#dc2626'}30` }}
                        >
                            <Icon name={budget.category?.icon ?? 'restaurant'} size={16} color={budget.category?.color ?? '#dc2626'} />
                        </View>
                        <View className="flex-1">

                            <Text className="text-lg font-bold text-gray-900">{budget.category?.name}</Text>
                            <Text className="text-sm text-gray-500">
                                {formatCurrency(budget.amount_planned)} par {formatPeriod(budget.period)}
                            </Text>
                        </View>
                    </View>

                </View>

                {/* Status Badge */}
                {/* <View 
          className="rounded-full w-fit h-fit"
          style={{ backgroundColor: `${getStatusColor(budget.status)}20` }}
        >
          <Text 
            className="text-xs font-medium"
            style={{ color: getStatusColor(budget.status) }}
          >
            {getStatusText(budget.status)}
          </Text>
        </View> */}
            </View>

            {/* Progress Bar */}
            <View className="mb-3">
                <View className="flex-row justify-between items-center mb-1">
                    <Text className="text-sm text-gray-600">
                        {formatCurrency(budget.amount_spent)} / {formatCurrency(budget.amount_planned)}
                    </Text>
                    <Text className="text-sm font-medium text-gray-900">
                        {budget.percentage_used.toFixed(0)}%
                    </Text>
                </View>

                <View className="w-full bg-gray-200 rounded-full h-2">
                    <View
                        className="h-2 rounded-full"
                        style={{
                            width: `${Math.min(budget.percentage_used, 100)}%`,
                            backgroundColor: getProgressColor(budget.percentage_used)
                        }}
                    />
                </View>
            </View>

            {/* Details */}
            <View className="flex-row justify-between items-center">
                <View className="flex-row items-center">
                    <Ionicons name="time-outline" size={16} color="#6b7280" />
                    <Text className="text-sm text-gray-600 ml-1">
                        {budget.days_remaining} jours restants
                    </Text>
                </View>

                {/* Actions */}
                <View className="flex-row items-center">
                    {onEdit && (
                        <TouchableOpacity
                            onPress={onEdit}
                            className="p-1 mr-2"
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

            {/* Remaining Amount */}
            {budget.remaining_amount > 0 && (
                <View className="mt-2 pt-2 border-t border-gray-100">
                    <Text className="text-sm text-gray-600">
                        Reste à dépenser : <Text className="font-medium text-green-600">
                            {formatCurrency(budget.remaining_amount)}
                        </Text>
                    </Text>
                </View>
            )}

            {/* Warning for overdue budgets */}
            {budget.status === 'overdue' && (
                <View className="mt-2 p-2 bg-red-50 rounded-lg">
                    <View className="flex-row items-center">
                        <Ionicons name="warning-outline" size={16} color="#ef4444" />
                        <Text className="text-sm text-red-600 ml-1">
                            Budget dépassé de {formatCurrency(budget.amount_spent - budget.amount_planned)}
                        </Text>
                    </View>
                </View>
            )}
        </TouchableOpacity>
    );
} 