import React, { useState, useEffect } from 'react'
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ScrollView,
  Alert,
  Modal,
  KeyboardAvoidingView,
  Platform,
  Pressable,
} from 'react-native'
import { Ionicons } from '@expo/vector-icons'
import { Button } from '@/src/components/ui/Button'
import { colors, spacing } from '@/src/theme'
import { CreateTransactionRequest, Category, Account } from '@/src/types'

interface TransactionFormProps {
  visible: boolean
  onClose: () => void
  onSubmit: (transaction: CreateTransactionRequest) => Promise<void>
  categories: Category[]
  accounts: Account[]
  defaultType?: 'income' | 'expense'
}

export const TransactionForm: React.FC<TransactionFormProps> = ({
  visible,
  onClose,
  onSubmit,
  categories,
  accounts,
  defaultType = 'expense',
}) => {
  const [type, setType] = useState<'income' | 'expense'>(defaultType)
  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null)
  const [selectedAccount, setSelectedAccount] = useState<Account | null>(null)
  const [date, setDate] = useState(new Date())
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const filteredCategories = categories.filter(cat => cat.type === type)

  useEffect(() => {
    setSelectedCategory(null)
    setErrors({})
  }, [type])

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!amount || parseFloat(amount) <= 0) {
      newErrors.amount = 'Le montant doit être supérieur à 0'
    }

    if (!selectedAccount) {
      newErrors.account = 'Veuillez sélectionner un compte'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async () => {
    if (!validateForm()) {
      return
    }

    setIsLoading(true)
    try {
      const transaction: CreateTransactionRequest = {
        account_id: selectedAccount!.id,
        categoryId: selectedCategory?.id,
        type,
        amount: parseFloat(amount),
        description: description.trim() || undefined,
        date,
        recurring: false,
      }

      await onSubmit(transaction)
      
      setAmount('')
      setDescription('')
      setSelectedCategory(null)
      setDate(new Date())
      setErrors({})
      
      onClose()
    } catch (error) {
      Alert.alert(
        'Erreur',
        error instanceof Error ? error.message : 'Erreur lors de la création'
      )
    } finally {
      setIsLoading(false)
    }
  }

  const formatDate = (date: Date) => {
    return date.toLocaleDateString('fr-FR')
  }

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
          {/* Header */}
          <View className="flex-row items-center justify-between p-4 border-b border-gray-200">
            <TouchableOpacity onPress={onClose}>
              <Ionicons name="close" size={24} color={colors.text.primary} />
            </TouchableOpacity>
            <Text className="text-lg font-semibold text-gray-900">
              Nouvelle Transaction
            </Text>
            <View style={{ width: 24 }} />
          </View>

          <ScrollView className="flex-1 px-4 py-6">
            {/* Type de transaction */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-3 text-gray-900">
                Type de transaction
              </Text>
              <View className="flex-row space-x-3">
                <TouchableOpacity
                  className={`flex-1 p-3 rounded-lg border-2 ${
                    type === 'expense' ? 'border-red-500 bg-red-50' : 'border-gray-300'
                  }`}
                  onPress={() => setType('expense')}
                >
                  <View className="items-center">
                    <Ionicons 
                      name="arrow-down-circle" 
                      size={24} 
                      color={type === 'expense' ? colors.error[500] : colors.text.tertiary} 
                    />
                    <Text className={`text-sm font-medium mt-1 ${
                      type === 'expense' ? 'text-red-600' : 'text-gray-600'
                    }`}>
                      Dépense
                    </Text>
                  </View>
                </TouchableOpacity>

                <TouchableOpacity
                  className={`flex-1 p-3 rounded-lg border-2 ${
                    type === 'income' ? 'border-green-500 bg-green-50' : 'border-gray-300'
                  }`}
                  onPress={() => setType('income')}
                >
                  <View className="items-center">
                    <Ionicons 
                      name="arrow-up-circle" 
                      size={24} 
                      color={type === 'income' ? colors.success[500] : colors.text.tertiary} 
                    />
                    <Text className={`text-sm font-medium mt-1 ${
                      type === 'income' ? 'text-green-600' : 'text-gray-600'
                    }`}>
                      Revenu
                    </Text>
                  </View>
                </TouchableOpacity>
              </View>
            </View>

            {/* Montant */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Montant *
              </Text>
              <View className={`flex-row items-center border rounded-lg px-4 py-3 ${
                errors.amount ? 'border-red-500' : 'border-gray-300'
              } bg-gray-50`}>
                <Ionicons 
                  name="cash-outline" 
                  size={20} 
                  color={colors.text.tertiary} 
                  style={{ marginRight: spacing.sm }}
                />
                <TextInput
                  className="flex-1 text-lg font-semibold text-gray-900"
                  placeholder="0.00"
                  placeholderTextColor={colors.text.tertiary}
                  value={amount}
                  onChangeText={(text) => {
                    setAmount(text.replace(/[^0-9.]/g, ''))
                    if (errors.amount) {
                      setErrors(prev => ({ ...prev, amount: '' }))
                    }
                  }}
                  keyboardType="numeric"
                />
                <Text className="text-lg font-semibold ml-2 text-gray-900">€</Text>
              </View>
              {errors.amount && (
                <Text className="text-sm mt-1 text-red-500">{errors.amount}</Text>
              )}
            </View>

            {/* Compte */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Compte *
              </Text>
            <ScrollView horizontal>
              {(accounts??[]).map((account) => (
                <Pressable key={account.id} className="bg-primary p-2 rounded-lg px-4 py-4" onPress={() => setSelectedAccount(account)}>
                  <Text className="text-white">{account.name}</Text>
                </Pressable>
              ))}
            </ScrollView>
              {errors.account && (
                <Text className="text-sm mt-1 text-red-500">{errors.account}</Text>
              )}
            </View>

            {/* Catégorie (optionnelle) */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Catégorie (optionnel)
              </Text>
              <TouchableOpacity
                className="flex-row items-center justify-between border border-gray-300 rounded-lg px-4 py-3 bg-gray-50"
                onPress={() => {
                  if (filteredCategories.length === 0) {
                    Alert.alert('Aucune catégorie', 'Aucune catégorie disponible')
                    return
                  }
                }}
              >
                <View className="flex-row items-center">
                  <Ionicons 
                    name="pricetag-outline" 
                    size={20} 
                    color={colors.text.tertiary} 
                    style={{ marginRight: spacing.sm }}
                  />
                  <Text className={`text-base ${
                    selectedCategory ? 'text-gray-900' : 'text-gray-500'
                  }`}>
                    {selectedCategory ? selectedCategory.name : 'Sélectionner une catégorie'}
                  </Text>
                </View>
                <Ionicons name="chevron-down" size={20} color={colors.text.tertiary} />
              </TouchableOpacity>
            </View>

            {/* Date */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Date
              </Text>
              <View className="flex-row items-center justify-between border border-gray-300 rounded-lg px-4 py-3 bg-gray-50">
                <View className="flex-row items-center">
                  <Ionicons 
                    name="calendar-outline" 
                    size={20} 
                    color={colors.text.tertiary} 
                    style={{ marginRight: spacing.sm }}
                  />
                  <Text className="text-base text-gray-900">{formatDate(date)}</Text>
                </View>
              </View>
            </View>

            {/* Description (optionnelle) */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Description (optionnel)
              </Text>
              <View className="border border-gray-300 rounded-lg px-4 py-3 bg-gray-50">
                <TextInput
                  className="text-base text-gray-900"
                  placeholder="Ajouter une description..."
                  placeholderTextColor={colors.text.tertiary}
                  value={description}
                  onChangeText={setDescription}
                  multiline
                  numberOfLines={3}
                  textAlignVertical="top"
                />
              </View>
            </View>
          </ScrollView>

          {/* Footer */}
          <View className="p-4 border-t border-gray-200">
            <Button
              title="Créer la transaction"
              onPress={handleSubmit}
              disabled={isLoading}
              loading={isLoading}
              fullWidth
              size="large"
              variant={type === 'expense' ? 'danger' : 'primary'}
            />
          </View>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  )
} 