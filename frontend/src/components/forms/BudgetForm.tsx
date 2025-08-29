import React, { useState, useEffect } from 'react';
import { View, Text, Modal, TouchableOpacity, TextInput, ScrollView, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import Icon from '../Icons';

interface BudgetFormProps {
  visible: boolean;
  onClose: () => void;
  onSubmit: (budgetData: BudgetFormData) => void;
  categories: Category[];
  initialData?: Partial<BudgetFormData>;
}

interface BudgetFormData {
  name: string;
  categoryId: string;
  amountPlanned: number;
  period: 'monthly' | 'weekly' | 'yearly';
  description?: string;
}

export default function BudgetForm({ 
  visible, 
  onClose, 
  onSubmit, 
  categories,
  initialData 
}: BudgetFormProps) {
  const [formData, setFormData] = useState<BudgetFormData>({
    name: '',
    categoryId: '',
    amountPlanned: 0,
    period: 'monthly',
    description: ''
  });

  const [errors, setErrors] = useState<Partial<BudgetFormData>>({});

  useEffect(() => {
    if (initialData) {
      setFormData(prev => ({ ...prev, ...initialData }));
    }
  }, [initialData]);

  const validateForm = () => {
    const newErrors: Partial<BudgetFormData> = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Le nom du budget est requis';
    }

    if (!formData.categoryId) {
      newErrors.categoryId = 'La catégorie est requise';
    }

    if (formData.amountPlanned <= 0) {
      newErrors.amountPlanned = 'Le montant doit être supérieur à 0';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = () => {
    if (validateForm()) {
      onSubmit(formData);
      handleClose();
    }
  };

  const handleClose = () => {
    setFormData({
      name: '',
      categoryId: '',
      amountPlanned: 0,
      period: 'monthly',
      description: ''
    });
    setErrors({});
    onClose();
  };

  const expenseCategories = categories.filter(cat => cat.type === 'expense');

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
              {initialData ? 'Modifier le budget' : 'Nouveau budget'}
            </Text>
            <TouchableOpacity onPress={handleSubmit}>
              <Text className="text-blue-600 font-medium">Enregistrer</Text>
            </TouchableOpacity>
          </View>
        </View>

        <ScrollView className="flex-1 p-4">
          {/* Nom du budget */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Nom du budget</Text>
            <TextInput
              className={`border rounded-lg px-3 py-2 text-gray-900 ${
                errors.name ? 'border-red-500' : 'border-gray-300'
              }`}
              placeholder="Ex: Budget alimentation"
              value={formData.name}
              onChangeText={(text) => setFormData(prev => ({ ...prev, name: text }))}
            />
            {errors.name && (
              <Text className="text-red-500 text-sm mt-1">{errors.name}</Text>
            )}
          </View>

          {/* Catégorie */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Catégorie</Text>
            <View className="border border-gray-300 rounded-lg">
              {expenseCategories.map((category) => (
                <TouchableOpacity
                  key={category.id}
                  className={`flex-row items-center p-3 border-b border-gray-100 ${
                    formData.categoryId === category.id ? 'bg-blue-50' : ''
                  }`}
                  onPress={() => setFormData(prev => ({ ...prev, categoryId: category.id }))}
                >
                  <View 
                    className="w-8 h-8 rounded-full items-center justify-center mr-3"
                    style={{ backgroundColor: `${category.color}20` }}
                  >
                    <Icon name={category.icon as any} size={16} color={category.color} />
                  </View>
                  <Text className="flex-1 text-gray-900">{category.name}</Text>
                  {formData.categoryId === category.id && (
                    <Ionicons name="checkmark" size={20} color="#3b82f6" />
                  )}
                </TouchableOpacity>
              ))}
            </View>
            {errors.categoryId && (
              <Text className="text-red-500 text-sm mt-1">{errors.categoryId}</Text>
            )}
          </View>

          {/* Montant prévu */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Montant prévu</Text>
            <TextInput
              className={`border rounded-lg px-3 py-2 text-gray-900 ${
                errors.amountPlanned ? 'border-red-500' : 'border-gray-300'
              }`}
              placeholder="0.00"
              keyboardType="numeric"
              value={formData.amountPlanned.toString()}
              onChangeText={(text) => {
                const amount = parseFloat(text) || 0;
                setFormData(prev => ({ ...prev, amountPlanned: amount }));
              }}
            />
            {errors.amountPlanned && (
              <Text className="text-red-500 text-sm mt-1">{errors.amountPlanned}</Text>
            )}
          </View>

          {/* Période */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Période</Text>
            <View className="flex-row gap-2">
              {[
                { value: 'weekly', label: 'Hebdomadaire' },
                { value: 'monthly', label: 'Mensuel' },
                { value: 'yearly', label: 'Annuel' }
              ].map((period) => (
                <TouchableOpacity
                  key={period.value}
                  className={`flex-1 py-2 px-3 rounded-lg border ${
                    formData.period === period.value
                      ? 'bg-blue-600 border-blue-600'
                      : 'bg-white border-gray-300'
                  }`}
                  onPress={() => setFormData(prev => ({ ...prev, period: period.value as any }))}
                >
                  <Text
                    className={`text-center font-medium ${
                      formData.period === period.value ? 'text-white' : 'text-gray-700'
                    }`}
                  >
                    {period.label}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          {/* Description */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Description (optionnel)</Text>
            <TextInput
              className="border border-gray-300 rounded-lg px-3 py-2 text-gray-900"
              placeholder="Description du budget..."
              multiline
              numberOfLines={3}
              value={formData.description}
              onChangeText={(text) => setFormData(prev => ({ ...prev, description: text }))}
            />
          </View>

          {/* Résumé */}
          <View className="bg-blue-50 p-4 rounded-lg mb-4">
            <Text className="text-sm font-medium text-blue-900 mb-2">Résumé</Text>
            <Text className="text-blue-800">
              Budget {formData.name || 'sans nom'} - {formData.amountPlanned}€ par {formData.period === 'weekly' ? 'semaine' : formData.period === 'monthly' ? 'mois' : 'année'}
            </Text>
          </View>
        </ScrollView>
      </View>
    </Modal>
  );
} 