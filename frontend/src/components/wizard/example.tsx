import React, { useState } from 'react';
import { View, Text, TouchableOpacity, ScrollView } from 'react-native';
import { SetupWizard, WizardData } from './SetupWizard';

// Exemple d'utilisation du SetupWizard
export default function WizardExample() {
  const [showWizard, setShowWizard] = useState(false);
  const [userData, setUserData] = useState<WizardData | null>(null);

  return (
    <ScrollView style={{ flex: 1, padding: 20, backgroundColor: '#f5f5f5' }}>
      <Text style={{ fontSize: 24, fontWeight: 'bold', marginBottom: 20 }}>
        Setup Wizard Example
      </Text>

      {/* Bouton pour ouvrir le wizard */}
      <TouchableOpacity
        style={{
          backgroundColor: '#3b82f6',
          padding: 16,
          borderRadius: 12,
          alignItems: 'center',
          marginBottom: 20,
        }}
        onPress={() => setShowWizard(true)}
      >
        <Text style={{ color: 'white', fontSize: 16, fontWeight: '600' }}>
          🚀 Lancer le Setup Wizard
        </Text>
      </TouchableOpacity>

      {/* Affichage des données collectées */}
      {userData && (
        <View style={{ marginBottom: 20 }}>
          <Text style={{ fontSize: 18, fontWeight: '600', marginBottom: 12 }}>
            Données collectées :
          </Text>
          
          <View style={{ backgroundColor: 'white', padding: 16, borderRadius: 12 }}>
            <Text style={{ fontSize: 14, fontWeight: '600', marginBottom: 8 }}>
              💰 Revenus
            </Text>
            <Text style={{ fontSize: 12, color: '#666', marginBottom: 4 }}>
              Total mensuel: {userData.income.monthly_total.toLocaleString()} FCFA
            </Text>
            <Text style={{ fontSize: 12, color: '#666', marginBottom: 4 }}>
              Dettes: {userData.income.has_debt ? 'Oui' : 'Non'}
            </Text>

            <Text style={{ fontSize: 14, fontWeight: '600', marginTop: 12, marginBottom: 8 }}>
              💸 Dépenses
            </Text>
            <Text style={{ fontSize: 12, color: '#666', marginBottom: 4 }}>
              Nourriture: {userData.expenses.food.toLocaleString()} FCFA
            </Text>
            <Text style={{ fontSize: 12, color: '#666', marginBottom: 4 }}>
              Transport: {userData.expenses.transport.toLocaleString()} FCFA
            </Text>

            <Text style={{ fontSize: 14, fontWeight: '600', marginTop: 12, marginBottom: 8 }}>
              🎯 Objectifs
            </Text>
            <Text style={{ fontSize: 12, color: '#666', marginBottom: 4 }}>
              Épargne mensuelle: {userData.goals.savings_target.toLocaleString()} FCFA
            </Text>

            <Text style={{ fontSize: 14, fontWeight: '600', marginTop: 12, marginBottom: 8 }}>
              📅 Habitudes
            </Text>
            <Text style={{ fontSize: 12, color: '#666', marginBottom: 4 }}>
              Planification: {userData.habits.planning_time}
            </Text>
          </View>
        </View>
      )}

      {/* JSON brut pour debug */}
      {userData && (
        <View style={{ backgroundColor: 'white', padding: 16, borderRadius: 12 }}>
          <Text style={{ fontSize: 14, fontWeight: '600', marginBottom: 8 }}>
            JSON envoyé au backend :
          </Text>
          <Text style={{ fontSize: 10, color: '#666', fontFamily: 'monospace' }}>
            {JSON.stringify(userData, null, 2)}
          </Text>
        </View>
      )}

      {/* Wizard */}
      
    </ScrollView>
  );
} 