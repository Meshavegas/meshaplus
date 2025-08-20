import React, { useState } from 'react';
import { View, Text, ScrollView, TouchableOpacity, Switch, Alert } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { SyncStatus } from '@/src/components/ui/SyncStatus';

export default function Profile() {
  const [notificationsEnabled, setNotificationsEnabled] = useState(true);
  const [biometricEnabled, setBiometricEnabled] = useState(false);
  const [darkModeEnabled, setDarkModeEnabled] = useState(false);

  const handleExport = () => {
    Alert.alert(
      'Exporter les Données',
      'Voulez-vous exporter toutes vos données financières ?',
      [
        { text: 'Annuler', style: 'cancel' },
        { text: 'Exporter', onPress: () => console.log('Export started') }
      ]
    );
  };

  const handleBackup = () => {
    Alert.alert(
      'Sauvegarde',
      'Voulez-vous créer une sauvegarde de vos données ?',
      [
        { text: 'Annuler', style: 'cancel' },
        { text: 'Sauvegarder', onPress: () => console.log('Backup started') }
      ]
    );
  };

  const renderSettingsSection = () => (
    <View className="space-y-4">
      <Text className="text-lg font-semibold text-gray-900">Paramètres</Text>
      
      {/* Notifications */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="notifications" size={20} color="#3b82f6" />
            <Text className="text-gray-900 font-medium ml-3">Notifications</Text>
          </View>
          <Switch
            value={notificationsEnabled}
            onValueChange={setNotificationsEnabled}
            trackColor={{ false: '#d1d5db', true: '#93c5fd' }}
            thumbColor={notificationsEnabled ? '#3b82f6' : '#f3f4f6'}
          />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Alertes budget et rappels de factures
        </Text>
      </View>

      {/* Sécurité */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="finger-print" size={20} color="#10b981" />
            <Text className="text-gray-900 font-medium ml-3">Déverrouillage biométrique</Text>
          </View>
          <Switch
            value={biometricEnabled}
            onValueChange={setBiometricEnabled}
            trackColor={{ false: '#d1d5db', true: '#86efac' }}
            thumbColor={biometricEnabled ? '#10b981' : '#f3f4f6'}
          />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Utiliser l'empreinte digitale ou Face ID
        </Text>
      </View>

      {/* Thème */}
      <View className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="moon" size={20} color="#8b5cf6" />
            <Text className="text-gray-900 font-medium ml-3">Mode sombre</Text>
          </View>
          <Switch
            value={darkModeEnabled}
            onValueChange={setDarkModeEnabled}
            trackColor={{ false: '#d1d5db', true: '#c4b5fd' }}
            thumbColor={darkModeEnabled ? '#8b5cf6' : '#f3f4f6'}
          />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Interface en mode sombre
        </Text>
      </View>
    </View>
  );

  const renderPersonalizationSection = () => (
    <View className="space-y-4">
      <Text className="text-lg font-semibold text-gray-900">Personnalisation</Text>
      
      {/* Catégories personnalisées */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="color-palette" size={20} color="#f59e0b" />
            <Text className="text-gray-900 font-medium ml-3">Catégories personnalisées</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Créer et modifier vos catégories
        </Text>
      </TouchableOpacity>

      {/* Couleurs par catégorie */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="brush" size={20} color="#ec4899" />
            <Text className="text-gray-900 font-medium ml-3">Couleurs par catégorie</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Personnaliser les couleurs de vos catégories
        </Text>
      </TouchableOpacity>

      {/* Widgets */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="grid" size={20} color="#06b6d4" />
            <Text className="text-gray-900 font-medium ml-3">Widgets personnalisables</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Organiser votre dashboard
        </Text>
      </TouchableOpacity>
    </View>
  );

  const renderDataSection = () => (
    <View className="space-y-4">
      <Text className="text-lg font-semibold text-gray-900">Données</Text>
      
      {/* Export */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm" onPress={handleExport}>
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="download" size={20} color="#059669" />
            <Text className="text-gray-900 font-medium ml-3">Exporter les données</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Exporter en PDF ou Excel
        </Text>
      </TouchableOpacity>

      {/* Sauvegarde */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm" onPress={handleBackup}>
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="cloud-upload" size={20} color="#3b82f6" />
            <Text className="text-gray-900 font-medium ml-3">Sauvegarde cloud</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Sauvegarder vos données en ligne
        </Text>
      </TouchableOpacity>

      {/* Nettoyage cache */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="trash" size={20} color="#dc2626" />
            <Text className="text-gray-900 font-medium ml-3">Nettoyer le cache</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Libérer de l'espace de stockage
        </Text>
      </TouchableOpacity>
    </View>
  );

  const renderSupportSection = () => (
    <View className="space-y-4">
      <Text className="text-lg font-semibold text-gray-900">Support</Text>
      
      {/* Aide */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="help-circle" size={20} color="#8b5cf6" />
            <Text className="text-gray-900 font-medium ml-3">Centre d'aide</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Tutoriels et FAQ
        </Text>
      </TouchableOpacity>

      {/* Contact */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="mail" size={20} color="#f59e0b" />
            <Text className="text-gray-900 font-medium ml-3">Nous contacter</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Support technique et suggestions
        </Text>
      </TouchableOpacity>

      {/* À propos */}
      <TouchableOpacity className="bg-white p-4 rounded-lg shadow-sm">
        <View className="flex-row justify-between items-center">
          <View className="flex-row items-center">
            <Ionicons name="information-circle" size={20} color="#06b6d4" />
            <Text className="text-gray-900 font-medium ml-3">À propos</Text>
          </View>
          <Ionicons name="chevron-forward" size={20} color="#9ca3af" />
        </View>
        <Text className="text-sm text-gray-500 mt-1">
          Version 1.0.0 - Mother AI
        </Text>
      </TouchableOpacity>
    </View>
  );

  return (
    <ScrollView className="flex-1 bg-gray-50">
      <View className="p-4">
        <Text className="text-2xl font-bold text-gray-900 mb-6">Profil</Text>
        
        {/* Sync Status */}
        <View className="mb-6">
          <SyncStatus showDetails={true} />
        </View>

        {/* User Info */}
        <View className="bg-white p-4 rounded-lg shadow-sm mb-6">
          <View className="flex-row items-center">
            <View className="w-16 h-16 bg-blue-100 rounded-full items-center justify-center">
              <Ionicons name="person" size={32} color="#3b82f6" />
            </View>
            <View className="ml-4 flex-1">
              <Text className="text-lg font-semibold text-gray-900">John Doe</Text>
              <Text className="text-sm text-gray-500">john.doe@example.com</Text>
            </View>
            <TouchableOpacity>
              <Ionicons name="pencil" size={20} color="#6b7280" />
            </TouchableOpacity>
          </View>
        </View>

        {/* Sections */}
        {renderSettingsSection()}
        {renderPersonalizationSection()}
        {renderDataSection()}
        {renderSupportSection()}

        {/* Logout */}
        <TouchableOpacity className="bg-red-50 border border-red-200 rounded-lg p-4 mt-6">
          <View className="flex-row items-center justify-center">
            <Ionicons name="log-out" size={20} color="#dc2626" />
            <Text className="text-red-700 font-semibold ml-2">Se déconnecter</Text>
          </View>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
} 