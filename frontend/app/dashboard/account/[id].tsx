import React, { useEffect, useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Alert, RefreshControl } from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import Icon from '@/src/components/Icons';
import { TransactionForm } from '@/src/components/forms/TransactionForm';
import AccountChart from '@/src/components/charts/AccountChart';
import AccountActions from '@/src/components/AccountActions';
import accountApi from '@/src/services/accountService/accountApi';
import transactionApi from '@/src/services/transactionService/transactionApi';
import categoryApi from '@/src/services/categoryService/categoryApi';

interface Transaction {
    id: string;
    amount: number;
    description: string;
    category: string;
    date: string;
    type: 'income' | 'expense';
    accountId: string;
}

interface Account {
    id: string;
    userId: string;
    name: string;
    type: 'checking' | 'savings' | 'mobile_money';
    balance: number;
    icon: string;
    color: string;
    currency: string;
    createdAt: Date;
    updatedAt: Date;
}

export default function AccountDetail() {
    const { id } = useLocalSearchParams();
    const router = useRouter();
    const [account, setAccount] = useState<AccountDetails | null>(null);
    const [transactions, setTransactions] = useState<Transaction[]>([]);
    const [categories, setCategories] = useState<Category[]>([]);
    const [showTransactionForm, setShowTransactionForm] = useState(false);
    const [refreshing, setRefreshing] = useState(false);
    const [selectedPeriod, setSelectedPeriod] = useState('month');
    const [activeTab, setActiveTab] = useState<'transactions' | 'stats'>('transactions');

    useEffect(() => {
        loadAccountData();
    }, [id]);

    const loadAccountData = async () => {
        try {
            const [accountData, categoriesData] = await Promise.all([
                accountApi.getAccount(id as string),
                categoryApi.getCategories()
            ]);

            setAccount(accountData);
            setCategories(categoriesData);
        } catch (error) {
            console.error('Error loading account data:', error);
            Alert.alert('Erreur', 'Impossible de charger les données du compte');
        }
    };

    const onRefresh = async () => {
        setRefreshing(true);
        await loadAccountData();
        setRefreshing(false);
    };

    const handleCreateTransaction = async (transactionData: any) => {
        try {
            await transactionApi.createTransaction({
                ...transactionData,
                accountId: id
            });

            // Refresh data
            await loadAccountData();
            setShowTransactionForm(false);
            Alert.alert('Succès', 'Transaction créée avec succès');
        } catch (error) {
            Alert.alert('Erreur', 'Erreur lors de la création de la transaction');
        }
    };

    const handleEditAccount = () => {
        Alert.alert('Modifier le compte', 'Fonctionnalité à implémenter');
    };

    const handleDeleteAccount = () => {
        Alert.alert(
            'Supprimer le compte',
            'Êtes-vous sûr de vouloir supprimer ce compte ? Cette action est irréversible.',
            [
                { text: 'Annuler', style: 'cancel' },
                {
                    text: 'Supprimer',
                    style: 'destructive',
                    onPress: async () => {
                        try {
                            await accountApi.deleteAccount(id as string);
                            router.back();
                        } catch (error) {
                            Alert.alert('Erreur', 'Impossible de supprimer le compte');
                        }
                    }
                }
            ]
        );
    };

    const handleTransfer = () => {
        Alert.alert('Transfert', 'Fonctionnalité de transfert à implémenter');
    };

    const handleExport = () => {
        Alert.alert('Export', 'Fonctionnalité d\'export à implémenter');
    };

    const handleShare = () => {
        Alert.alert('Partager', 'Fonctionnalité de partage à implémenter');
    };

    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('fr-FR', {
            style: 'currency',
            currency: account?.account.currency || 'EUR'
        }).format(amount);
    };

         const formatDate = (date: Date | string | undefined) => {
         if (!date) return '';
         
         const dateObj = typeof date === 'string' ? new Date(date) : date;
         
         if (isNaN(dateObj.getTime())) {
             return 'Date invalide';
         }
         
         return dateObj.toLocaleDateString('fr-FR', {
             day: 'numeric',
             month: 'short',
             year: 'numeric'
         });
     };

    const getTransactionIcon = (category: string) => {
        const categoryMap: { [key: string]: string } = {
            'alimentation': 'restaurant',
            'transport': 'car',
            'logement': 'home',
            'sante': 'medical',
            'loisirs': 'game-controller',
            'shopping': 'bag',
            'default': 'card'
        };
        return categoryMap[category.toLowerCase()] || categoryMap.default;
    };

    const getTransactionColor = (type: string) => {
        return type === 'income' ? '#16a34a' : '#dc2626';
    };

    const getTransactionBgColor = (type: string) => {
        return type === 'income' ? '#dcfce7' : '#fef2f2';
    };

    if (!account) {
        return (
            <View className="flex-1 bg-gray-50 justify-center items-center">
                <Text className="text-gray-500">Chargement...</Text>
            </View>
        );
    }

    return (
        <View className="flex-1 bg-gray-50">
            {/* Header */}
            <View className="bg-white shadow-sm">
                <View className="flex-row items-center p-4 pt-12">
                    <TouchableOpacity onPress={() => router.back()} className="mr-4">
                        <Ionicons name="arrow-back" size={24} color="#374151" />
                    </TouchableOpacity>
                    <View className="flex-1">
                        <Text className="text-xl font-bold text-gray-900">{account.account.name}</Text>
                        <Text className="text-sm text-gray-500">Compte {account.account.type}</Text>
                    </View>
                    <TouchableOpacity onPress={handleEditAccount} className="mr-2">
                        <Ionicons name="create-outline" size={24} color="#374151" />
                    </TouchableOpacity>
                    <TouchableOpacity onPress={handleDeleteAccount}>
                        <Ionicons name="trash-outline" size={24} color="#dc2626" />
                    </TouchableOpacity>
                </View>
            </View>

            <ScrollView
                className="p-4"
                refreshControl={
                    <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
                }
            >
                {/* Account Summary Card */}
                <View className="bg-white p-6 rounded-lg shadow-sm mb-6">
                    <View className="flex-row items-center mb-4">
                        <View
                            className="w-12 h-12 rounded-full items-center justify-center mr-4"
                            style={{ backgroundColor: `${account.account.color}20` }}
                        >
                            <Icon name={account.account.icon} size={24} color={account.account.color} />
                        </View>
                        <View className="flex-1">
                            <Text className="text-2xl font-bold text-gray-900">
                                {formatCurrency(account.account.balance)}
                            </Text>
                            <Text className="text-sm text-gray-500">Solde actuel</Text>
                        </View>
                    </View>

                    <View className="flex-row justify-between py-3 border-t border-gray-100">
                                                 <View className="items-center flex-1">
                             <Text className="text-lg font-semibold text-green-600">
                                 {formatCurrency(account.transactions?.filter(t => t.type === 'income').reduce((sum, t) => sum + t.amount, 0) || 0)}
                             </Text>
                             <Text className="text-xs text-gray-500">Revenus</Text>
                         </View>
                         <View className="items-center flex-1">
                             <Text className="text-lg font-semibold text-red-600">
                                 {formatCurrency(account.transactions?.filter(t => t.type === 'expense').reduce((sum, t) => sum + t.amount, 0) || 0)}
                             </Text>
                             <Text className="text-xs text-gray-500">Dépenses</Text>
                         </View>
                                                 <View className="items-center flex-1">
                             <Text className="text-lg font-semibold text-gray-900">
                                 {account.transactions?.length || 0}
                             </Text>
                             <Text className="text-xs text-gray-500">Transactions</Text>
                         </View>
                    </View>
                </View>

                {/* Period Filter */}
                <View className="flex-row bg-white rounded-lg shadow-sm mb-6 p-1">
                    {['week', 'month', 'year'].map((period) => (
                        <TouchableOpacity
                            key={period}
                            className={`flex-1 py-2 px-4 rounded-md ${selectedPeriod === period ? 'bg-blue-100' : ''
                                }`}
                            onPress={() => setSelectedPeriod(period)}
                        >
                            <Text className={`text-center font-medium ${selectedPeriod === period ? 'text-blue-600' : 'text-gray-600'
                                }`}>
                                {period === 'week' ? 'Semaine' : period === 'month' ? 'Mois' : 'Année'}
                            </Text>
                        </TouchableOpacity>
                    ))}
                </View>

                {/* Tab Navigation */}
                <View className="flex-row bg-white rounded-lg shadow-sm mb-6 p-1">
                    <TouchableOpacity
                        className={`flex-1 py-3 px-4 rounded-md ${activeTab === 'transactions' ? 'bg-blue-100' : ''
                            }`}
                        onPress={() => setActiveTab('transactions')}
                    >
                        <Text className={`text-center font-medium ${activeTab === 'transactions' ? 'text-blue-600' : 'text-gray-600'
                            }`}>
                            Transactions
                        </Text>
                    </TouchableOpacity>
                    <TouchableOpacity
                        className={`flex-1 py-3 px-4 rounded-md ${activeTab === 'stats' ? 'bg-blue-100' : ''
                            }`}
                        onPress={() => setActiveTab('stats')}
                    >
                        <Text className={`text-center font-medium ${activeTab === 'stats' ? 'text-blue-600' : 'text-gray-600'
                            }`}>
                            Statistiques
                        </Text>
                    </TouchableOpacity>
                </View>

                {/* Quick Actions */}
                <AccountActions
                    onAddTransaction={() => setShowTransactionForm(true)}
                    onTransfer={handleTransfer}
                    onExport={handleExport}
                    onShare={handleShare}
                />

                {/* Content based on active tab */}
                {activeTab === 'transactions' ? (
                    /* Transactions Section */
                    <View className="mb-6">
                        <View className="flex-row justify-between items-center mb-4">
                            <Text className="text-xl font-bold text-gray-900">Transactions</Text>
                            <TouchableOpacity onPress={() => Alert.alert('Filtrer', 'Fonctionnalité à implémenter')}>
                                <Ionicons name="filter" size={20} color="#374151" />
                            </TouchableOpacity>
                        </View>

                                                 {!account.transactions || account.transactions.length === 0 ? (
                            <View className="bg-white p-8 rounded-lg shadow-sm items-center">
                                <Ionicons name="receipt-outline" size={48} color="#9ca3af" />
                                <Text className="text-gray-500 mt-2 text-center">Aucune transaction pour le moment</Text>
                                <TouchableOpacity
                                    className="mt-4 bg-blue-600 px-4 py-2 rounded-lg"
                                    onPress={() => setShowTransactionForm(true)}
                                >
                                    <Text className="text-white font-medium">Ajouter une transaction</Text>
                                </TouchableOpacity>
                            </View>
                        ) : (
                                                         <View className="gap-2">
                                 {account.transactions?.map((transaction) => (
                                    <TouchableOpacity
                                        key={transaction.id}
                                        className="bg-white p-4 rounded-lg shadow-sm"
                                        onPress={() => router.push(`/dashboard/transaction/${transaction.id}`)}
                                    >
                                        <View className="flex-row justify-between items-center">
                                            <View className="flex-row items-center flex-1">
                                                <View
                                                    className="w-10 h-10 rounded-full items-center justify-center mr-3"
                                                    style={{ backgroundColor: transaction.category.color ?? '#000' }}
                                                >
                                                    <Icon
                                                        name={transaction.category.icon ?? 'io:card'}
                                                        size={18}
                                                        color={transaction.category.color ?? '#000'}
                                                    />
                                                </View>
                                                <View className="flex-1">
                                                    <Text className="text-gray-900 font-medium">{transaction.description}</Text>
                                                    <Text className="text-sm text-gray-500">
                                                        {transaction.category.name} • {formatDate(transaction.created_at)}
                                                    </Text>
                                                </View>
                                            </View>
                                            <View className="items-end">
                                                <Text
                                                    className={`font-semibold ${transaction.type === 'income' ? 'text-green-600' : 'text-red-600'
                                                        }`}
                                                >
                                                    {transaction.type === 'income' ? '+' : '-'}{formatCurrency(transaction.amount)}
                                                </Text>
                                            </View>
                                        </View>
                                    </TouchableOpacity>
                                ))}
                            </View>
                        )}
                    </View>
                ) : (
                                         /* Statistics Section */
                     <View className="mb-6">
                         <AccountChart
                             transactions={account.transactions || []}
                             period={selectedPeriod as 'week' | 'month' | 'year'}
                         />
                     </View>
                )}

                {/* Account Details */}
                <View className="bg-white p-6 rounded-lg shadow-sm mb-6">
                    <Text className="text-lg font-bold text-gray-900 mb-4">Détails du compte</Text>
                    <View className="space-y-3">
                        <View className="flex-row justify-between">
                            <Text className="text-gray-600">Type de compte</Text>
                            <Text className="text-gray-900 font-medium">{account.account.type}</Text>
                        </View>
                        <View className="flex-row justify-between">
                            <Text className="text-gray-600">Devise</Text>
                            <Text className="text-gray-900 font-medium">{account.account.currency}</Text>
                        </View>
                        <View className="flex-row justify-between">
                            <Text className="text-gray-600">Créé le</Text>
                            <Text className="text-gray-900 font-medium">{formatDate(account.account.created_at)}</Text>
                        </View>
                        <View className="flex-row justify-between">
                            <Text className="text-gray-600">Dernière activité</Text>
                                                         <Text className="text-gray-900 font-medium">
                                 {account.transactions && account.transactions.length > 0 ? formatDate(account.transactions[0].created_at) : 'Aucune'}
                             </Text>
                        </View>
                    </View>
                </View>
            </ScrollView>

            {/* Transaction Form Modal */}
            <TransactionForm
                visible={showTransactionForm}
                onClose={() => setShowTransactionForm(false)}
                onSubmit={handleCreateTransaction}
                categories={categories}
                accounts={[account.account]}
            />
        </View>
    );
} 