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
import { Button } from '@/src/components/ui/Button'
import { colors, spacing } from '@/src/theme'
import { CreateTransactionRequest, Category, Account } from '@/src/types'
import Select from '../Select'
import Icon from '../Icons'
import { isEmpty } from '@/src/utils/stringUtils'

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
              <Icon name="close" size={24} color={colors.text.primary} />
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
                    <Icon 
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
                    <Icon 
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
                <Icon 
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
            <ScrollView horizontal showsHorizontalScrollIndicator={false}>
              <View className="flex-row gap-2">

              {(accounts??[]).map((account) => (
                <Pressable key={account.id} className={`flex-row items-center gap-2 py-4 px-4 rounded-xl overflow-hidden`} style={{ backgroundColor: selectedAccount?.id === account.id ? `#${(isEmpty(account.color)?"#000":account.color).replace('#', '')}` : `#${account.color.replace('#', '')}87`, borderRadius: 10, borderWidth: 1, borderColor: selectedAccount?.id === account.id ? `#${account.color.replace('#', '')}` : `#${account.color.replace('#', '')}87` }} onPress={() => setSelectedAccount(account)}>
                 <View className="flex-row items-center gap-2">

                  <Text>
                    <Icon name={account.icon as any} size={20} color={`${ selectedAccount?.id !== account.id ? account.color : '#fff' }`} />
                    </Text>
                  <Text className="text-white">{account.name}</Text> 
                 </View>
                </Pressable>
              ))}
              </View>
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
             <ScrollView horizontal showsHorizontalScrollIndicator={false}>
              <View className="flex-row gap-2">
              {categories?.filter(category => category.type === (type==="income"?"revenue":type))?.map(category => (
                <Pressable
                  key={category.id}
                  className={`flex-row items-center gap-2 py-4 px-4 rounded-xl overflow-hidden`}
                  style={{
                    backgroundColor:
                      selectedCategory?.id === category.id
                        ? `#${category.color.replace('#', '')}`
                        : `#${category.color.replace('#', '')}87`,
                    borderRadius: 10,
                    borderWidth: 1,
                    borderColor:
                      selectedCategory?.id === category.id
                        ? `#${category.color.replace('#', '')}`
                        : `#${category.color.replace('#', '')}87`,
                  }}
                  onPress={() => {
                    if (selectedCategory?.id === category.id) {
                      setSelectedCategory(null);
                    } else {
                      setSelectedCategory(category);
                    }
                  }}
                >
                  <View className="flex-row items-center gap-2">
                    <Text>
                      <Icon name={category.icon as any} size={20} color={`${ selectedCategory?.id !== category.id ? category.color : '#fff' }`} />
                    </Text>
                    <Text className="text-white">{category.name}</Text>
                  </View>
                </Pressable>
              ))}
              </View>
             </ScrollView>
            </View>

            {/* Date */}
            <View className="mb-6">
              <Text className="text-sm font-medium mb-2 text-gray-900">
                Date
              </Text>
              <View className="flex-row items-center justify-between border border-gray-300 rounded-lg px-4 py-3 bg-gray-50">
                <View className="flex-row items-center">
                  <Icon 
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