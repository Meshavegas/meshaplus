import { View, Text, ScrollView, TouchableOpacity, StatusBar } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { SyncStatus } from '@/src/components/ui/SyncStatus';
import CircularProgress from '@/src/components/ui/CircularProgress';
import { useAuthStore } from '@/src';

export default function Overview() {
  const { user } = useAuthStore()
  
  return (
    <ScrollView className="flex-1 bg-gray-50 mt-safe">
      <StatusBar barStyle="dark-content" backgroundColor="#ffffff" />
      <View className="p-4">
        <View className='flex-row items-center justify-between gap-2'>
        <View className='flex-row items-center gap-2'>
          <View className='h-12 w-12 bg-primary items-center justify-center rounded-full'> 
            <Text className="text-2xl font-bold text-white">MV</Text>
          </View>
          <View className='flex-col'>
          <Text className="text-md font-bold text-gray-900">{user?.name}</Text>
          <Text className="text-xs text-gray-500">Bienvenue sur votre tableau de bord</Text>
          </View>
        </View>
          <SyncStatus />
        </View>

        {/* Daily Budget Status */}
        <View className="bg-white  rounded-lg shadow-sm mb-6 mt-3  overflow-hidden">
         
          
            <View className='flex-row items-center bg-primary/30 '>
              <CircularProgress progressPercent={30} size={60} strokeWidth={10} text='3/5' textSize={20}/>
              <View className=' flex flex-col gap-2'>
              <Text className="text-md font-bold text-gray-900 ">Taches</Text>
              <Text>Ranger la chambre</Text>
              </View>
            </View>
          
          <View className="bg-teal/70 px-4 py-4">
            <View className="flex-row justify-between items-center">
              <Text className="text-sm text-white">Budget quotidien</Text>
              <Text className="text-white font-semibold">$67 / $100</Text>
            </View>
            <View className="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
              <View className="bg-teal-500 h-3 rounded-full" style={{ width: '69%' }} />
            </View>
            <Text className="text-xs text-white mt-1">$33 restant aujourd'hui</Text>
          </View>
        </View>

        
        <View>
          <Text className="text-md text-black">Mois en cours</Text>
          <View className="space-y-4 mb-2 mt-2 flex flex-row gap-2">
          
          <View className="bg-white px-4 py-2 flex flex-col justify-between rounded-lg shadow-sm flex-1">
            <Text className="text-[12px] text-gray-500">Solde Total</Text>
            <Text className="text-md font-bold text-green-600">$2,450.00</Text>
            <Text className="text-xs text-green-600">+$150 ce mois</Text>
          </View>
          
          <View className="bg-white px-4 py-2 flex flex-col justify-between rounded-lg shadow-sm flex-1">
            <Text className="text-[12px] text-gray-500">Dépenses</Text>
            <Text className="text-md font-bold text-red-600">$1,200.00</Text>
            <Text className="text-xs text-red-600">60% du budget</Text>
          </View>
          
          <View className="bg-white px-4 py-2 flex flex-col justify-between rounded-lg shadow-sm flex-1">
            <Text className="text-[12px] text-gray-500">Revenus</Text>
            <Text className="text-md font-bold text-blue-600">$3,650.00</Text>
            <Text className="text-xs text-blue-600">+12.5%</Text>
          </View>
        </View>
        </View>

        <View className='w-full h-[200px] bg-white rounded-lg shadow-sm mb-6 items-center justify-center'>
          <Text className="text-md text-black">Statistique</Text>
        </View>

        {/* Daily Tasks */}
        <View className="bg-white p-4 rounded-lg shadow-sm mb-6">
          <Text className="text-lg font-semibold text-gray-900 mb-4">Tâches du Jour</Text>
          <View className="space-y-3">
            <View className="flex-row items-center">
              <View className="w-5 h-5 bg-green-100 rounded-full items-center justify-center mr-3">
                <Ionicons name="checkmark" size={12} color="#10b981" />
              </View>
              <Text className="text-gray-700 flex-1">Saisir les dépenses du jour</Text>
            </View>
            <View className="flex-row items-center">
              <View className="w-5 h-5 bg-gray-200 rounded-full items-center justify-center mr-3">
                <Ionicons name="time" size={12} color="#9ca3af" />
              </View>
              <Text className="text-gray-700 flex-1">Vérifier le budget alimentation</Text>
            </View>
            <View className="flex-row items-center">
              <View className="w-5 h-5 bg-gray-200 rounded-full items-center justify-center mr-3">
                <Ionicons name="time" size={12} color="#9ca3af" />
              </View>
              <Text className="text-gray-700 flex-1">Mettre à jour les objectifs</Text>
            </View>
          </View>
        </View>

        {/* Quick Actions */}
        <View className="space-y-3 mb-6">
          <Text className="text-lg font-semibold text-gray-900">Actions Rapides</Text>
          <View className="flex-row gap-2">
            <TouchableOpacity className="flex-1 bg-[#dc2626] p-4 rounded-lg items-center">
              <Ionicons name="add" size={20} color="white" />
              <Text className="text-white font-medium mt-1">Dépense</Text>
            </TouchableOpacity>
            <TouchableOpacity className="flex-1 bg-green-600 p-4 rounded-lg items-center">
              <Ionicons name="add" size={20} color="white" />
              <Text className="text-white font-medium mt-1">Revenu</Text>
            </TouchableOpacity>
            <TouchableOpacity className="flex-1 bg-blue-600 p-4 rounded-lg items-center">
              <Ionicons name="flag" size={20} color="white" />
              <Text className="text-white font-medium mt-1">Taches</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Recent Activity */}
        <View className="bg-white p-4 rounded-lg shadow-sm">
          <Text className="text-lg font-semibold text-gray-900 mb-4">Activité Récente</Text>
          <View className="space-y-3">
            <View className="flex-row justify-between items-center">
              <View className="flex-row items-center">
                <View className="w-8 h-8 bg-red-100 rounded-full items-center justify-center mr-3">
                  <Ionicons name="restaurant" size={16} color="#dc2626" />
                </View>
                <View>
                  <Text className="text-gray-900 font-medium">Courses alimentaires</Text>
                  <Text className="text-xs text-gray-500">Il y a 2h</Text>
                </View>
              </View>
              <Text className="text-red-600 font-semibold">-$85.50</Text>
            </View>
            <View className="flex-row justify-between items-center">
              <View className="flex-row items-center">
                <View className="w-8 h-8 bg-green-100 rounded-full items-center justify-center mr-3">
                  <Ionicons name="briefcase" size={16} color="#16a34a" />
                </View>
                <View>
                  <Text className="text-gray-900 font-medium">Salaire</Text>
                  <Text className="text-xs text-gray-500">Hier</Text>
                </View>
              </View>
              <Text className="text-green-600 font-semibold">+$2,500.00</Text>
            </View>
            <View className="flex-row justify-between items-center">
              <View className="flex-row items-center">
                <View className="w-8 h-8 bg-blue-100 rounded-full items-center justify-center mr-3">
                  <Ionicons name="car" size={16} color="#2563eb" />
                </View>
                <View>
                  <Text className="text-gray-900 font-medium">Essence</Text>
                  <Text className="text-xs text-gray-500">Il y a 1 jour</Text>
                </View>
              </View>
              <Text className="text-red-600 font-semibold">-$45.00</Text>
            </View>
          </View>
        </View>
      </View>
    </ScrollView>
  );
} 