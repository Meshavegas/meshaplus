import React, { useState } from 'react'
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
import Icon from '../Icons'
import ColorsSelector from './ColorsSelector'
import IconsSelector from './IconsSelector'

interface IAccountRequest {
    id?: string
    Name: string
    Type: string
    Balance: number
    Currency: string
    icon?: string
    color?: string
}

interface AccountFormProps {
    visible: boolean
    onClose: () => void
    onSubmit: (account: IAccountRequest) => Promise<void>
    errors?: string
}

const ACCOUNT_TYPES = [
    { id: 'checking', name: 'Compte Courant', icon: 'card', color: '#2563eb' },
    { id: 'savings', name: 'Épargne', icon: 'wallet', color: '#16a34a' },
    { id: 'mobile_money', name: 'Mobile Money', icon: 'phone-portrait', color: '#f59e0b' },
    { id: 'credit_card', name: 'Carte de Crédit', icon: 'card-outline', color: '#dc2626' },
    { id: 'investment', name: 'Investissement', icon: 'trending-up', color: '#9333ea' },
]

const CURRENCIES = [
    { code: 'XAF', symbol: 'FCFA', name: 'Franc CFA' },
    { code: 'EUR', symbol: '€', name: 'Euro' },
    { code: 'USD', symbol: '$', name: 'Dollar US' },
    { code: 'GBP', symbol: '£', name: 'Livre Sterling' },
]

export const AccountForm: React.FC<AccountFormProps> = ({
    visible,
    onClose,
    onSubmit,
    errors
}) => {
    const [name, setName] = useState('')
    const [balance, setBalance] = useState('')
    const [selectedType, setSelectedType] = useState(ACCOUNT_TYPES[0])
    const [color, setColor] = useState<string>()
    const [selectedCurrency, setSelectedCurrency] = useState(CURRENCIES[0])
    const [isLoading, setIsLoading] = useState(false)
    const [formErrors, setFormErrors] = useState<Record<string, string>>({})
    const [selectedIcon, setSelectedIcon] = useState<string>()
    const [showIconsSelector, setShowIconsSelector] = useState<boolean>(false)

    const [showColorsSelector, SetShowColorsSelector] = useState<boolean>(false)

    const validateForm = () => {
        const newErrors: Record<string, string> = {}

        if (!name.trim()) {
            newErrors.name = 'Le nom du compte est requis'
        }

        if (!balance || parseFloat(balance) < 0) {
            newErrors.balance = 'Le solde doit être un nombre positif'
        }

        setFormErrors(newErrors)
        return Object.keys(newErrors).length === 0
    }

    const handleSubmit = async () => {
        if (!validateForm()) {
            return
        }

        setIsLoading(true)
        try {
            const account: IAccountRequest = {
                Name: name.trim(),
                Type: selectedType.id,
                Balance: parseFloat(balance) || 0,
                Currency: selectedCurrency.code,
                icon: selectedIcon ?? "fa:card",
                color: color ?? selectedType.color,
            }

            await onSubmit(account)

            // Reset form
            setName('')
            setBalance('')
            setSelectedType(ACCOUNT_TYPES[0])
            setSelectedCurrency(CURRENCIES[0])
            setFormErrors({})

            onClose()
        } catch (error) {
            Alert.alert(
                'Erreur',
                error instanceof Error ? error.message : 'Erreur lors de la création du compte'
            )
        } finally {
            setIsLoading(false)
        }
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
                            Nouveau Compte
                        </Text>
                        <View style={{ width: 24 }} />
                    </View>

                    <ScrollView className="flex-1 px-4 py-6">
                        {/* Nom du compte */}
                        <View className="mb-6">
                            <Text className="text-sm font-medium mb-2 text-gray-900">
                                Nom du compte *
                            </Text>
                            <View className={`border rounded-lg px-4 py-3 ${formErrors.name ? 'border-red-500' : 'border-gray-300'
                                } bg-gray-50`}>
                                <TextInput
                                    className="text-base text-gray-900"
                                    placeholder="Ex: Compte Courant Principal"
                                    placeholderTextColor={colors.text.tertiary}
                                    value={name}
                                    onChangeText={(text) => {
                                        setName(text)
                                        if (formErrors.name) {
                                            setFormErrors(prev => ({ ...prev, name: '' }))
                                        }
                                    }}
                                />
                            </View>
                            {formErrors.name && (
                                <Text className="text-sm mt-1 text-red-500">{formErrors.name}</Text>
                            )}
                        </View>

                        {/* Type de compte */}
                        <View className="mb-6">
                            <Text className="text-sm font-medium mb-3 text-gray-900">
                                Type de compte *
                            </Text>
                            <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                                <View className="flex-row gap-2">
                                    {ACCOUNT_TYPES.map((type) => (
                                        <Pressable
                                            key={type.id}
                                            className={`flex-row items-center gap-2 py-4 px-4 rounded-xl overflow-hidden`}
                                            style={{
                                                backgroundColor:
                                                    selectedType.id === type.id
                                                        ? `#${type.color.replace('#', '')}`
                                                        : `#${type.color.replace('#', '')}87`,
                                                borderRadius: 10,
                                                borderWidth: 1,
                                                borderColor:
                                                    selectedType.id === type.id
                                                        ? `#${type.color.replace('#', '')}`
                                                        : `#${type.color.replace('#', '')}87`,
                                            }}
                                            onPress={() => setSelectedType(type)}
                                        >
                                            <View className="flex-row items-center gap-2">
                                                <Icon
                                                    name={type.icon as any}
                                                    size={20}
                                                    color={selectedType.id === type.id ? '#fff' : type.color}
                                                />
                                                <Text className="text-white">{type.name}</Text>
                                            </View>
                                        </Pressable>
                                    ))}
                                </View>
                            </ScrollView>
                        </View>
                        <View className='flex flex-row w-full gap-4'>
                            <View className='flex-1'>
                                <Text className=' text-sm font-medium text-gray-900'>Colors</Text>
                                <Pressable onPress={() => SetShowColorsSelector(true)} style={{
                                    backgroundColor: color ?? "teal"
                                }} className='w-full h-10 rounded-md items-center justify-center rounded-lg' >
                                    <Text className=' text-white text-xl font-medium'>{color?.toUpperCase()}</Text>
                                </Pressable>
                            </View>
                            <View className='w-20 justify-center items-center' >
                                <Text>Icon</Text>
                                <Pressable onPress={() => setShowIconsSelector(true)}>
                                    <Icon name={selectedIcon ?? "fa:user"} color={color ?? "teal"} />
                                </Pressable>
                            </View>
                        </View>

                        {/* Solde initial */}
                        <View className="mb-6">
                            <Text className="text-sm font-medium mb-2 text-gray-900">
                                Solde initial *
                            </Text>
                            <View className={`flex-row items-center border rounded-lg px-4 py-3 ${formErrors.balance ? 'border-red-500' : 'border-gray-300'
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
                                    value={balance}
                                    onChangeText={(text) => {
                                        setBalance(text.replace(/[^0-9.]/g, ''))
                                        if (formErrors.balance) {
                                            setFormErrors(prev => ({ ...prev, balance: '' }))
                                        }
                                    }}
                                    keyboardType="numeric"
                                />
                                <Text className="text-lg font-semibold ml-2 text-gray-900">
                                    {selectedCurrency.symbol}
                                </Text>
                            </View>
                            {formErrors.balance && (
                                <Text className="text-sm mt-1 text-red-500">{formErrors.balance}</Text>
                            )}
                        </View>

                        {/* Devise */}
                        <View className="mb-6">
                            <Text className="text-sm font-medium mb-3 text-gray-900">
                                Devise *
                            </Text>
                            <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                                <View className="flex-row gap-2">
                                    {CURRENCIES.map((currency) => (
                                        <Pressable
                                            key={currency.code}
                                            className={`flex-row items-center gap-2 py-4 px-4 rounded-xl overflow-hidden`}
                                            style={{
                                                backgroundColor:
                                                    selectedCurrency.code === currency.code
                                                        ? '#2563eb'
                                                        : '#2563eb87',
                                                borderRadius: 10,
                                                borderWidth: 1,
                                                borderColor:
                                                    selectedCurrency.code === currency.code
                                                        ? '#2563eb'
                                                        : '#2563eb87',
                                            }}
                                            onPress={() => setSelectedCurrency(currency)}
                                        >
                                            <View className="flex-row items-center gap-2">
                                                <Text className="text-white font-medium">
                                                    {currency.symbol}
                                                </Text>
                                                <Text className="text-white">{currency.name}</Text>
                                            </View>
                                        </Pressable>
                                    ))}
                                </View>
                            </ScrollView>
                        </View>

                        {/* Aperçu du compte */}
                        <View className="mb-6">
                            <Text className="text-sm font-medium mb-3 text-gray-900">
                                Aperçu
                            </Text>
                            <View className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                                <View className="flex-row items-center mb-3">
                                    <View
                                        className="w-10 h-10 rounded-full items-center justify-center mr-3"
                                        style={{ backgroundColor: `${color}30` }}
                                    >
                                        <Icon name={selectedType.icon as any} size={20} color={color} />
                                    </View>
                                    <View className="flex-1">
                                        <Text className="text-gray-900 font-medium">{name || 'Nom du compte'}</Text>
                                        <Text className="text-sm text-gray-500">{selectedType.name}</Text>
                                    </View>
                                </View>
                                <View className="flex-row justify-between items-center">
                                    <Text className="text-lg font-bold text-gray-900">
                                        {balance || '0.00'} {selectedCurrency.symbol}
                                    </Text>
                                    <Text className="text-sm text-gray-500">{selectedCurrency.code}</Text>
                                </View>
                            </View>
                        </View>
                    </ScrollView>

                    {/* Footer */}
                    <View className="p-4 border-t border-gray-200">
                        <Button
                            title="Créer le compte"
                            onPress={handleSubmit}
                            disabled={isLoading}
                            loading={isLoading}
                            fullWidth
                            size="large"
                            variant="primary"
                        />
                    </View>
                </View>
                <ColorsSelector
                    onClose={() => { SetShowColorsSelector(false) }}
                    onSelect={(val) => {
                        setColor(val)
                        SetShowColorsSelector(false)
                    }}
                    visible={showColorsSelector}
                />
                <IconsSelector
                    onClose={() => { setShowIconsSelector(false) }}
                    onSelect={(val) => {
                        setSelectedIcon(val)
                        setShowIconsSelector(false)
                    }}
                    visible={showIconsSelector}
                />
            </KeyboardAvoidingView>
        </Modal>
    )
}