import React, { useState } from 'react';
import { View, Text, Modal, TouchableOpacity, TextInput, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { SavingGoal } from '@/src/services/savingGoalService/savingGoalApi';

interface ContributionFormProps {
  visible: boolean;
  onClose: () => void;
  onSubmit: (amount: number) => void;
  goal: SavingGoal;
}

export default function ContributionForm({ 
  visible, 
  onClose, 
  onSubmit, 
  goal 
}: ContributionFormProps) {
  const [amount, setAmount] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = () => {
    const numAmount = parseFloat(amount);
    
    if (!amount || numAmount <= 0) {
      setError('Veuillez entrer un montant valide');
      return;
    }

    if (numAmount > goal.remaining_amount) {
      setError(`Le montant ne peut pas dépasser ${formatCurrency(goal.remaining_amount)}`);
      return;
    }

    onSubmit(numAmount);
    handleClose();
  };

  const handleClose = () => {
    setAmount('');
    setError('');
    onClose();
  };

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: 'EUR'
    }).format(value);
  };

  const getSuggestedAmounts = () => {
    const remaining = goal.remaining_amount;
    const suggestions = [];
    
    // 10% du montant restant
    if (remaining * 0.1 > 0) {
      suggestions.push(Math.round(remaining * 0.1));
    }
    
    // 25% du montant restant
    if (remaining * 0.25 > 0) {
      suggestions.push(Math.round(remaining * 0.25));
    }
    
    // 50% du montant restant
    if (remaining * 0.5 > 0) {
      suggestions.push(Math.round(remaining * 0.5));
    }
    
    // Montant restant complet
    suggestions.push(remaining);
    
    return suggestions.slice(0, 3); // Limiter à 3 suggestions
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      presentationStyle="pageSheet"
    >
      <View className="flex-1 bg-gray-50">
        {/* Header */}
        <View className="bg-white shadow-sm">
          <View className="flex-row items-center justify-between p-4 pt-12">
            <TouchableOpacity onPress={handleClose}>
              <Text className="text-blue-600 font-medium">Annuler</Text>
            </TouchableOpacity>
            <Text className="text-lg font-bold text-gray-900">
              Ajouter une contribution
            </Text>
            <TouchableOpacity onPress={handleSubmit}>
              <Text className="text-blue-600 font-medium">Ajouter</Text>
            </TouchableOpacity>
          </View>
        </View>

        <View className="p-4">
          {/* Goal Info */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-lg font-bold text-gray-900 mb-2">{goal.title}</Text>
            <View className="space-y-1">
              <Text className="text-sm text-gray-600">
                Objectif : {formatCurrency(goal.target_amount)}
              </Text>
              <Text className="text-sm text-gray-600">
                Épargné : {formatCurrency(goal.current_amount)}
              </Text>
              <Text className="text-sm text-gray-600">
                Reste : {formatCurrency(goal.remaining_amount)}
              </Text>
              <Text className="text-sm text-gray-600">
                Progression : {goal.percentage_achieved.toFixed(0)}%
              </Text>
            </View>
          </View>

          {/* Amount Input */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Montant de la contribution</Text>
            <TextInput
              className={`border rounded-lg px-3 py-3 text-gray-900 text-lg ${
                error ? 'border-red-500' : 'border-gray-300'
              }`}
              placeholder="0.00"
              keyboardType="numeric"
              value={amount}
              onChangeText={(text) => {
                setAmount(text);
                setError('');
              }}
            />
            {error && (
              <Text className="text-red-500 text-sm mt-1">{error}</Text>
            )}
          </View>

          {/* Quick Amount Suggestions */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Suggestions rapides</Text>
            <View className="flex-row flex-wrap gap-2">
              {getSuggestedAmounts().map((suggestedAmount) => (
                <TouchableOpacity
                  key={suggestedAmount}
                  className="bg-blue-100 px-3 py-2 rounded-lg"
                  onPress={() => setAmount(suggestedAmount.toString())}
                >
                  <Text className="text-blue-700 font-medium">
                    {formatCurrency(suggestedAmount)}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          {/* Impact Preview */}
          {amount && parseFloat(amount) > 0 && (
            <View className="bg-green-50 p-4 rounded-lg mb-4">
              <Text className="text-sm font-medium text-green-900 mb-2">Impact de cette contribution</Text>
              <View className="space-y-1">
                <Text className="text-sm text-green-800">
                  Nouveau total : {formatCurrency(goal.current_amount + parseFloat(amount))}
                </Text>
                <Text className="text-sm text-green-800">
                  Nouvelle progression : {((goal.current_amount + parseFloat(amount)) / goal.target_amount * 100).toFixed(0)}%
                </Text>
                <Text className="text-sm text-green-800">
                  Reste après contribution : {formatCurrency(goal.remaining_amount - parseFloat(amount))}
                </Text>
              </View>
            </View>
          )}

          {/* Motivation */}
          <View className="bg-blue-50 p-4 rounded-lg">
            <View className="flex-row items-start">
              <Ionicons name="heart-outline" size={20} color="#3b82f6" />
              <View className="ml-2 flex-1">
                <Text className="text-sm font-medium text-blue-900 mb-1">Motivation</Text>
                <Text className="text-sm text-blue-800">
                  Chaque contribution vous rapproche de votre objectif ! 
                  Même les petits montants comptent. 💪
                </Text>
              </View>
            </View>
          </View>
        </View>
      </View>
    </Modal>
  );
} 