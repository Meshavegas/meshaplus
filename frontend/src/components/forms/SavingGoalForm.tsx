import React, { useState, useEffect } from 'react';
import { View, Text, Modal, TouchableOpacity, TextInput, ScrollView, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';

interface SavingGoalFormProps {
  visible: boolean;
  onClose: () => void;
  onSubmit: (goalData: SavingGoalFormData) => void;
  initialData?: Partial<SavingGoalFormData>;
}

interface SavingGoalFormData {
  title: string;
  target_amount: number;
  frequency: 'weekly' | 'monthly' | 'yearly';
  deadline: Date;
}

export default function SavingGoalForm({ 
  visible, 
  onClose, 
  onSubmit, 
  initialData 
}: SavingGoalFormProps) {
  const [formData, setFormData] = useState<SavingGoalFormData>({
    title: '',
    target_amount: 0,
    frequency: 'monthly',
    deadline: new Date()
  });

  const [errors, setErrors] = useState<Partial<SavingGoalFormData>>({});

  useEffect(() => {
    if (initialData) {
      const deadline = initialData.deadline ? new Date(initialData.deadline) : new Date();
      setFormData(prev => ({ 
        ...prev, 
        ...initialData,
        deadline: deadline
      }));
    }
  }, [initialData]);

  const validateForm = () => {
    const newErrors: Partial<SavingGoalFormData> = {};

    if (!formData.title.trim()) {
      newErrors.title = 'Le titre de l\'objectif est requis';
    }

    if (formData.target_amount <= 0) {
      newErrors.target_amount = 'Le montant cible doit être supérieur à 0';
    }

    if (formData.deadline <= new Date()) {
      newErrors.deadline = 'La date limite doit être dans le futur';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = () => {
    if (validateForm()) {
      onSubmit({
        ...formData,
        deadline: formData.deadline.toISOString()
      });
      handleClose();
    }
  };

  const handleClose = () => {
    setFormData({
      title: '',
      target_amount: 0,
      frequency: 'monthly',
      deadline: new Date()
    });
    setErrors({});
    onClose();
  };

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('fr-FR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('fr-FR', {
      style: 'currency',
      currency: 'EUR'
    }).format(amount);
  };

  const calculateMonthlyContribution = () => {
    const now = new Date();
    const monthsDiff = (formData.deadline.getFullYear() - now.getFullYear()) * 12 + 
                      (formData.deadline.getMonth() - now.getMonth());
    return monthsDiff > 0 ? formData.target_amount / monthsDiff : formData.target_amount;
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
              {initialData ? 'Modifier l\'objectif' : 'Nouvel objectif d\'épargne'}
            </Text>
            <TouchableOpacity onPress={handleSubmit}>
              <Text className="text-blue-600 font-medium">Enregistrer</Text>
            </TouchableOpacity>
          </View>
        </View>

        <ScrollView className="flex-1 p-4">
          {/* Titre de l'objectif */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Titre de l'objectif</Text>
            <TextInput
              className={`border rounded-lg px-3 py-2 text-gray-900 ${
                errors.title ? 'border-red-500' : 'border-gray-300'
              }`}
              placeholder="Ex: Vacances d'été"
              value={formData.title}
              onChangeText={(text) => setFormData(prev => ({ ...prev, title: text }))}
            />
            {errors.title && (
              <Text className="text-red-500 text-sm mt-1">{errors.title}</Text>
            )}
          </View>

          {/* Montant cible */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Montant cible</Text>
            <TextInput
              className={`border rounded-lg px-3 py-2 text-gray-900 ${
                errors.target_amount ? 'border-red-500' : 'border-gray-300'
              }`}
              placeholder="0.00"
              keyboardType="numeric"
              value={formData.target_amount.toString()}
              onChangeText={(text) => {
                const amount = parseFloat(text) || 0;
                setFormData(prev => ({ ...prev, target_amount: amount }));
              }}
            />
            {errors.target_amount && (
              <Text className="text-red-500 text-sm mt-1">{errors.target_amount}</Text>
            )}
          </View>

          {/* Fréquence */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Fréquence de suivi</Text>
            <View className="flex-row gap-2">
              {[
                { value: 'weekly', label: 'Hebdomadaire' },
                { value: 'monthly', label: 'Mensuel' },
                { value: 'yearly', label: 'Annuel' }
              ].map((freq) => (
                <TouchableOpacity
                  key={freq.value}
                  className={`flex-1 py-2 px-3 rounded-lg border ${
                    formData.frequency === freq.value
                      ? 'bg-blue-600 border-blue-600'
                      : 'bg-white border-gray-300'
                  }`}
                  onPress={() => setFormData(prev => ({ ...prev, frequency: freq.value as any }))}
                >
                  <Text
                    className={`text-center font-medium ${
                      formData.frequency === freq.value ? 'text-white' : 'text-gray-700'
                    }`}
                  >
                    {freq.label}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          {/* Date limite */}
          <View className="bg-white p-4 rounded-lg shadow-sm mb-4">
            <Text className="text-sm font-medium text-gray-700 mb-2">Date limite</Text>
            <View className="flex-row gap-2">
              <View className="flex-1">
                <Text className="text-xs text-gray-600 mb-1">Jour</Text>
                <TextInput
                  className={`border rounded-lg px-3 py-2 text-gray-900 ${
                    errors.deadline ? 'border-red-500' : 'border-gray-300'
                  }`}
                  placeholder="JJ"
                  keyboardType="numeric"
                  value={formData.deadline.getDate().toString()}
                  onChangeText={(text) => {
                    const day = parseInt(text) || 1;
                    const newDate = new Date(formData.deadline);
                    newDate.setDate(day);
                    setFormData(prev => ({ ...prev, deadline: newDate }));
                  }}
                />
              </View>
              <View className="flex-1">
                <Text className="text-xs text-gray-600 mb-1">Mois</Text>
                <TextInput
                  className={`border rounded-lg px-3 py-2 text-gray-900 ${
                    errors.deadline ? 'border-red-500' : 'border-gray-300'
                  }`}
                  placeholder="MM"
                  keyboardType="numeric"
                  value={(formData.deadline.getMonth() + 1).toString()}
                  onChangeText={(text) => {
                    const month = parseInt(text) || 1;
                    const newDate = new Date(formData.deadline);
                    newDate.setMonth(month - 1);
                    setFormData(prev => ({ ...prev, deadline: newDate }));
                  }}
                />
              </View>
              <View className="flex-1">
                <Text className="text-xs text-gray-600 mb-1">Année</Text>
                <TextInput
                  className={`border rounded-lg px-3 py-2 text-gray-900 ${
                    errors.deadline ? 'border-red-500' : 'border-gray-300'
                  }`}
                  placeholder="AAAA"
                  keyboardType="numeric"
                  value={formData.deadline.getFullYear().toString()}
                  onChangeText={(text) => {
                    const year = parseInt(text) || new Date().getFullYear();
                    const newDate = new Date(formData.deadline);
                    newDate.setFullYear(year);
                    setFormData(prev => ({ ...prev, deadline: newDate }));
                  }}
                />
              </View>
            </View>
            {errors.deadline && (
              <Text className="text-red-500 text-sm mt-1">{errors.deadline as string}</Text>
            )}
          </View>

          {/* Résumé */}
          <View className="bg-blue-50 p-4 rounded-lg mb-4">
            <Text className="text-sm font-medium text-blue-900 mb-2">Résumé</Text>
            <View className="space-y-2">
              <Text className="text-blue-800">
                Objectif : {formData.title || 'sans nom'}
              </Text>
              <Text className="text-blue-800">
                Montant cible : {formatCurrency(formData.target_amount)}
              </Text>
              <Text className="text-blue-800">
                Date limite : {formatDate(formData.deadline)}
              </Text>
              <Text className="text-blue-800 font-medium">
                Contribution mensuelle suggérée : {formatCurrency(calculateMonthlyContribution())}
              </Text>
            </View>
          </View>

          {/* Conseils */}
          <View className="bg-green-50 p-4 rounded-lg mb-4">
            <View className="flex-row items-start">
              <Ionicons name="bulb-outline" size={20} color="#16a34a" />
              <View className="ml-2 flex-1">
                <Text className="text-sm font-medium text-green-900 mb-1">Conseils</Text>
                <Text className="text-sm text-green-800">
                  • Définissez un montant réaliste selon vos revenus{'\n'}
                  • Choisissez une date limite motivante{'\n'}
                  • Suivez régulièrement vos progrès{'\n'}
                  • Célébrez les étapes importantes
                </Text>
              </View>
            </View>
          </View>
        </ScrollView>
      </View>
    </Modal>
  );
} 