import React, { useState } from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  TextInput,
  Alert,
  Modal,
  KeyboardAvoidingView,
  Platform,
} from 'react-native';
import { Icon } from '@/src/components/Icons';
import { Button } from '@/src/components/ui/Button';

// Types pour les données du wizard
export interface WizardData {
  income: {
    sources: string[];
    monthly_total: number;
    accounts: string[];
    has_debt: boolean;
    debt_amount?: number;
  };
  expenses: {
    top_categories: string[];
    food: number;
    transport: number;
    housing: number;
    subscriptions: number;
    alerts_enabled: boolean;
    auto_budget: boolean;
  };
  goals: {
    main_goal: string;
    secondary_goal?: string;
    savings_target: number;
    deadline: string;
    advice_enabled: boolean;
  };
  habits: {
    planning_time: string;
    daily_focus_time: string;
    custom_habit?: string;
    summary_type: string;
  };
}

interface SetupWizardProps {
  visible: boolean;
  onClose: () => void;
  onComplete: (data: WizardData) => Promise<void>;
}

export const SetupWizard: React.FC<SetupWizardProps> = ({
  visible,
  onClose,
  onComplete,
}) => {
  const [currentStep, setCurrentStep] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [wizardData, setWizardData] = useState<WizardData>({
    income: {
      sources: [],
      monthly_total: 0,
      accounts: [],
      has_debt: false,
    },
    expenses: {
      top_categories: [],
      food: 0,
      transport: 0,
      housing: 0,
      subscriptions: 0,
      alerts_enabled: true,
      auto_budget: true,
    },
    goals: {
      main_goal: '',
      savings_target: 0,
      deadline: '3 mois',
      advice_enabled: true,
    },
    habits: {
      planning_time: 'Matin',
      daily_focus_time: '30min',
      summary_type: 'Hebdomadaire',
    },
  });

  const totalSteps = 5;

  const updateWizardData = (section: keyof WizardData, updates: Partial<WizardData[keyof WizardData]>) => {
    setWizardData(prev => ({
      ...prev,
      [section]: { ...prev[section], ...updates },
    }));
  };

  const nextStep = () => {
    if (currentStep < totalSteps - 1) {
      setCurrentStep(currentStep + 1);
    }
  };

  const prevStep = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleComplete = async () => {
    setIsLoading(true);
    try {
      await onComplete(wizardData);
      Alert.alert('Succès', 'Configuration terminée avec succès !');
      onClose();
    } catch (error) {
      Alert.alert('Erreur', 'Impossible de sauvegarder la configuration');
    } finally {
      setIsLoading(false);
    }
  };

  const renderWelcome = () => (
    <View className="flex-1 justify-center items-center px-6">
      <View className="bg-blue-100 w-24 h-24 rounded-full items-center justify-center mb-8">
        <Icon name="fa6:rocket" size={48} color="#3b82f6" />
      </View>
      <Text className="text-2xl font-bold text-gray-900 text-center mb-4">
        Bienvenue ! 🎉
      </Text>
      <Text className="text-base text-gray-600 text-center leading-6 mb-8">
        On va t'aider à mieux gérer ton temps et ton argent.{'\n'}
        Réponds à quelques questions pour personnaliser ton expérience.
      </Text>
      <Button
        title="Commencer"
        onPress={nextStep}
        variant="primary"
        size="large"
        fullWidth
      />
    </View>
  );

  const renderIncomeStep = () => (
    <ScrollView className="flex-1 px-6">
      <Text className="text-xl font-bold text-gray-900 mb-6">Revenus & Comptes</Text>
      
      <View className="mb-6">
        <Text className="text-base font-semibold text-gray-900 mb-3">
          Montant total moyen par mois (FCFA)
        </Text>
        <TextInput
          className="border border-gray-300 rounded-lg px-4 py-3 text-base"
          placeholder="250000"
          keyboardType="numeric"
          value={wizardData.income.monthly_total.toString()}
          onChangeText={(text) => updateWizardData('income', { monthly_total: parseInt(text) || 0 })}
        />
      </View>

      <View className="mb-6">
        <Text className="text-base font-semibold text-gray-900 mb-3">
          As-tu des dettes/crédits à rembourser ?
        </Text>
        <View className="flex-row gap-3">
          <TouchableOpacity
            className={`flex-1 py-3 rounded-lg border-2 ${
              wizardData.income.has_debt ? 'border-red-500 bg-red-50' : 'border-gray-300'
            }`}
            onPress={() => updateWizardData('income', { has_debt: true })}
          >
            <Text className={`text-center ${
              wizardData.income.has_debt ? 'text-red-600 font-medium' : 'text-gray-600'
            }`}>
              Oui
            </Text>
          </TouchableOpacity>
          <TouchableOpacity
            className={`flex-1 py-3 rounded-lg border-2 ${
              !wizardData.income.has_debt ? 'border-green-500 bg-green-50' : 'border-gray-300'
            }`}
            onPress={() => updateWizardData('income', { has_debt: false })}
          >
            <Text className={`text-center ${
              !wizardData.income.has_debt ? 'text-green-600 font-medium' : 'text-gray-600'
            }`}>
              Non
            </Text>
          </TouchableOpacity>
        </View>
      </View>
    </ScrollView>
  );

  const renderExpensesStep = () => (
    <ScrollView className="flex-1 px-6">
      <Text className="text-xl font-bold text-gray-900 mb-6">Dépenses & Budget</Text>
      
      <View className="mb-6">
        <Text className="text-base font-semibold text-gray-900 mb-3">
          Montants mensuels estimés (FCFA)
        </Text>
        <View className="space-y-3">
          <View>
            <Text className="text-sm text-gray-600 mb-1">Nourriture</Text>
            <TextInput
              className="border border-gray-300 rounded-lg px-4 py-3 text-base"
              placeholder="60000"
              keyboardType="numeric"
              value={wizardData.expenses.food.toString()}
              onChangeText={(text) => updateWizardData('expenses', { food: parseInt(text) || 0 })}
            />
          </View>
          <View>
            <Text className="text-sm text-gray-600 mb-1">Transport</Text>
            <TextInput
              className="border border-gray-300 rounded-lg px-4 py-3 text-base"
              placeholder="20000"
              keyboardType="numeric"
              value={wizardData.expenses.transport.toString()}
              onChangeText={(text) => updateWizardData('expenses', { transport: parseInt(text) || 0 })}
            />
          </View>
        </View>
      </View>
    </ScrollView>
  );

  const renderGoalsStep = () => (
    <ScrollView className="flex-1 px-6">
      <Text className="text-xl font-bold text-gray-900 mb-6">Objectifs Financiers</Text>
      
      <View className="mb-6">
        <Text className="text-base font-semibold text-gray-900 mb-3">
          Combien veux-tu épargner par mois ? (FCFA)
        </Text>
        <TextInput
          className="border border-gray-300 rounded-lg px-4 py-3 text-base"
          placeholder="50000"
          keyboardType="numeric"
          value={wizardData.goals.savings_target.toString()}
          onChangeText={(text) => updateWizardData('goals', { savings_target: parseInt(text) || 0 })}
        />
      </View>
    </ScrollView>
  );

  const renderHabitsStep = () => (
    <ScrollView className="flex-1 px-6">
      <Text className="text-xl font-bold text-gray-900 mb-6">Organisation & Habitudes</Text>
      
      <View className="mb-6">
        <Text className="text-base font-semibold text-gray-900 mb-3">
          Quand préfères-tu organiser ta journée ?
        </Text>
        <View className="space-y-2">
          {['Matin ☀️', 'Midi 🌤', 'Soir 🌙'].map((time) => (
            <TouchableOpacity
              key={time}
              className={`p-3 rounded-lg border ${
                wizardData.habits.planning_time === time.split(' ')[0]
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-300'
              }`}
              onPress={() => updateWizardData('habits', { planning_time: time.split(' ')[0] })}
            >
              <Text className={`${
                wizardData.habits.planning_time === time.split(' ')[0]
                  ? 'text-blue-600 font-medium'
                  : 'text-gray-600'
              }`}>
                {time}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>
    </ScrollView>
  );

  const renderSummary = () => (
    <ScrollView className="flex-1 px-6">
      <Text className="text-xl font-bold text-gray-900 mb-6">Résumé & Confirmation</Text>
      
      <View className="space-y-6">
        <View className="bg-blue-50 p-4 rounded-lg">
          <Text className="text-lg font-semibold text-blue-900 mb-3">💰 Revenus</Text>
          <Text className="text-sm text-blue-800">
            Total mensuel: {wizardData.income.monthly_total.toLocaleString()} FCFA
          </Text>
        </View>

        <View className="bg-green-50 p-4 rounded-lg">
          <Text className="text-lg font-semibold text-green-900 mb-3">💸 Dépenses</Text>
          <Text className="text-sm text-green-800">
            Budget total: {(wizardData.expenses.food + wizardData.expenses.transport).toLocaleString()} FCFA
          </Text>
        </View>

        <View className="bg-purple-50 p-4 rounded-lg">
          <Text className="text-lg font-semibold text-purple-900 mb-3">🎯 Objectifs</Text>
          <Text className="text-sm text-purple-800">
            Épargne: {wizardData.goals.savings_target.toLocaleString()} FCFA/mois
          </Text>
        </View>
      </View>

      <Text className="text-base text-gray-600 text-center mt-6 mb-6">
        Tout est prêt ! Nous allons personnaliser ton expérience.
      </Text>
    </ScrollView>
  );

  const renderStep = () => {
    switch (currentStep) {
      case 0:
        return renderWelcome();
      case 1:
        return renderIncomeStep();
      case 2:
        return renderExpensesStep();
      case 3:
        return renderGoalsStep();
      case 4:
        return renderHabitsStep();
      case 5:
        return renderSummary();
      default:
        return renderWelcome();
    }
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      presentationStyle="pageSheet"
      onRequestClose={onClose}
    >
      <KeyboardAvoidingView 
        style={{ flex: 1 }} 
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      >
        <View className="flex-1 bg-white">
          {currentStep > 0 && (
            <View className="flex-row items-center justify-between p-4 border-b border-gray-200">
              <TouchableOpacity onPress={prevStep}>
                <Icon name="arrow-back" size={24} color="#374151" />
              </TouchableOpacity>
              <View className="flex-1 mx-4">
                <View className="w-full bg-gray-200 rounded-full h-2">
                  <View 
                    className="bg-blue-500 h-2 rounded-full" 
                    style={{ width: `${((currentStep - 1) / (totalSteps - 2)) * 100}%` }} 
                  />
                </View>
                <Text className="text-xs text-gray-500 mt-1 text-center">
                  Étape {currentStep} sur {totalSteps - 1}
                </Text>
              </View>
              <TouchableOpacity onPress={onClose}>
                <Icon name="close" size={24} color="#374151" />
              </TouchableOpacity>
            </View>
          )}

          <View className="flex-1">
            {renderStep()}
          </View>

          {currentStep > 0 && currentStep < totalSteps - 1 && (
            <View className="p-4 border-t border-gray-200">
              <Button
                title="Suivant"
                onPress={nextStep}
                variant="primary"
                size="large"
                fullWidth
              />
            </View>
          )}

          {currentStep === totalSteps - 1 && (
            <View className="p-4 border-t border-gray-200">
              <Button
                title="Lancer mon tableau de bord 🚀"
                onPress={handleComplete}
                disabled={isLoading}
                loading={isLoading}
                variant="primary"
                size="large"
                fullWidth
              />
            </View>
          )}
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
}; 