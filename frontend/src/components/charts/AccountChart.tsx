import React from 'react';
import { View, Text, Dimensions } from 'react-native';
import { LineChart, BarChart, PieChart } from 'react-native-chart-kit';

 
interface AccountChartProps {
  transactions: AccountTransaction[];
  period: 'week' | 'month' | 'year';
}

const screenWidth = Dimensions.get('window').width;

export default function AccountChart({ transactions, period }: AccountChartProps) {
  // Calculer les données pour le graphique linéaire (évolution du solde)
  const getBalanceData = () => {
    const sortedTransactions = [...transactions].sort((a, b) => {
      const dateA = typeof a.date === 'string' ? new Date(a.date) : a.date;
      const dateB = typeof b.date === 'string' ? new Date(b.date) : b.date;
      return dateA.getTime() - dateB.getTime();
    });

    let balance = 0;
    const data = sortedTransactions.map(transaction => {
      if (transaction.type === 'income') {
        balance += transaction.amount;
      } else {
        balance -= transaction.amount;
      }
      return balance;
    });

    return {
      labels: sortedTransactions.map(t => {
        const date = typeof t.date === 'string' ? new Date(t.date) : t.date;
        return isNaN(date.getTime()) ? 'Date invalide' : date.toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit' });
      }),
      datasets: [{ data }]
    };
  };

  console.log(JSON.stringify(transactions, null, 2), "transactions")
  // Calculer les données pour le graphique en barres (revenus vs dépenses)
  const getIncomeExpenseData = () => {
    const income = transactions
      .filter(t => t.type === 'income')
      .reduce((sum, t) => sum + t.amount, 0);
    
    const expense = transactions
      .filter(t => t.type === 'expense')
      .reduce((sum, t) => sum + t.amount, 0);

    return {
      labels: ['Revenus', 'Dépenses'],
      datasets: [{ data: [income, expense] }]
    };
  };

  // Calculer les données pour le graphique circulaire (par catégorie)
  const getCategoryData = () => {
    const categoryMap = new Map<{label: string,color: string}, number>();
    
    transactions.forEach(transaction => {
      const categoryName = typeof transaction.category === 'string' 
        ? transaction.category 
        : transaction.category.name;
      const current = categoryMap.get({label: categoryName, color: transaction.category.color}) || 0;
      categoryMap.set({label: categoryName, color: transaction.category.color}, current + transaction.amount);
    });

    const colors = ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF', '#FF9F40'];
    
    return Array.from(categoryMap.entries()).map(([category, amount], index) => ({
      name: category.label,
      amount,
      color: category.color,
      legendFontColor: '#7F7F7F',
      legendFontSize: 12
    }));
  };

  const chartConfig = {
    backgroundColor: '#ffffff',
    backgroundGradientFrom: '#ffffff',
    backgroundGradientTo: '#ffffff',
    decimalPlaces: 2,
    color: (opacity = 1) => `rgba(59, 130, 246, ${opacity})`,
    labelColor: (opacity = 1) => `rgba(0, 0, 0, ${opacity})`,
    style: {
      borderRadius: 16
    },
    propsForDots: {
      r: '6',
      strokeWidth: '2',
      stroke: '#3b82f6'
    }
  };

  const barChartConfig = {
    ...chartConfig,
    color: (opacity = 1) => `rgba(34, 197, 94, ${opacity})`,
    fillShadowGradient: '#22c55e',
    fillShadowGradientOpacity: 0.8
  };

  return (
    <View className="space-y-6">
      {/* Évolution du solde */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <Text className="text-lg font-bold text-gray-900 mb-4">Évolution du solde</Text>
        {transactions.length > 0 ? (
          <LineChart
            data={getBalanceData()}
            width={screenWidth - 40}
            height={220}
            chartConfig={chartConfig}
            bezier
            style={{
              marginVertical: 8,
              borderRadius: 16
            }}
          />
        ) : (
          <View className="h-32 justify-center items-center">
            <Text className="text-gray-500">Aucune donnée disponible</Text>
          </View>
        )}
      </View>

      {/* Revenus vs Dépenses */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <Text className="text-lg font-bold text-gray-900 mb-4">Revenus vs Dépenses</Text>
        {transactions.length > 0 ? (
          <BarChart
            data={getIncomeExpenseData()}
            width={screenWidth - 40}
            height={220}
            chartConfig={barChartConfig}
            style={{
              marginVertical: 8,
              borderRadius: 16
            }}
            fromZero
            yAxisLabel=""
            yAxisSuffix=""
          />
        ) : (
          <View className="h-32 justify-center items-center">
            <Text className="text-gray-500">Aucune donnée disponible</Text>
          </View>
        )}
      </View>

      {/* Répartition par catégorie */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <Text className="text-lg font-bold text-gray-900 mb-4">Répartition par catégorie</Text>
        {transactions.length > 0 ? (
          <PieChart
            data={getCategoryData()}
            width={screenWidth - 40}
            height={220}
            chartConfig={chartConfig}
            accessor="amount"
            backgroundColor="transparent"
            paddingLeft="15"
            absolute
          />
        ) : (
          <View className="h-32 justify-center items-center">
            <Text className="text-gray-500">Aucune donnée disponible</Text>
          </View>
        )}
      </View>
    </View>
  );
} 