import React, { useState } from 'react';
import { View, Text, ScrollView } from 'react-native';
import Select, { SelectOption } from '../Select';

// Exemple d'utilisation du composant Select
export default function SelectExample() {
  const [selectedCategory, setSelectedCategory] = useState<string>();
  const [selectedCurrency, setSelectedCurrency] = useState<string>();
  const [selectedAccount, setSelectedAccount] = useState<string>();

  // Options pour les catégories
  const categoryOptions: SelectOption[] = [
    { label: 'Alimentation', value: 'food', icon: 'fa5:utensils' },
    { label: 'Transport', value: 'transport', icon: 'fa5:car' },
    { label: 'Logement', value: 'housing', icon: 'fa6:house' },
    { label: 'Loisirs', value: 'entertainment', icon: 'fa5:gamepad' },
    { label: 'Santé', value: 'health', icon: 'fa6:heart-pulse' },
    { label: 'Shopping', value: 'shopping', icon: 'fa5:shopping-bag' },
    { label: 'Éducation', value: 'education', icon: 'fa5:graduation-cap' },
    { label: 'Autres', value: 'other', icon: 'fa6:ellipsis' },
  ];

  // Options pour les devises
  const currencyOptions: SelectOption[] = [
    { label: 'Euro (€)', value: 'EUR', icon: 'fa5:euro-sign' },
    { label: 'Dollar US ($)', value: 'USD', icon: 'fa5:dollar-sign' },
    { label: 'Livre Sterling (£)', value: 'GBP', icon: 'fa5:sterling-sign' },
    { label: 'Franc Suisse (CHF)', value: 'CHF', icon: 'fa5:franc-sign' },
    { label: 'Yen Japonais (¥)', value: 'JPY', icon: 'fa5:yen-sign' },
  ];

  // Options pour les comptes bancaires
  const accountOptions: SelectOption[] = [
    { 
      label: 'Compte Courant Principal', 
      value: 'main-checking',
      icon: 'fa6:building-columns',
      description: 'Compte principal pour les dépenses quotidiennes'
    },
    { 
      label: 'Compte Épargne', 
      value: 'savings',
      icon: 'fa5:piggy-bank',
      description: 'Compte d\'épargne avec intérêts'
    },
    { 
      label: 'Compte Investissement', 
      value: 'investment',
      icon: 'fa5:chart-line',
      description: 'Compte pour les placements'
    },
    { 
      label: 'Carte de Crédit', 
      value: 'credit-card',
      icon: 'fa5:credit-card',
      description: 'Carte de crédit personnelle'
    },
  ];

  return (
    <ScrollView style={{ flex: 1, padding: 20, backgroundColor: '#f5f5f5' }}>
      <Text style={{ fontSize: 24, fontWeight: 'bold', marginBottom: 20 }}>
        Exemples de Select
      </Text>

      {/* Select simple */}
      <View style={{ marginBottom: 24 }}>
        <Text style={{ fontSize: 16, fontWeight: '600', marginBottom: 8 }}>
          Catégorie de dépense
        </Text>
        <Select
          options={categoryOptions}
          value={selectedCategory}
          onValueChange={setSelectedCategory}
          placeholder="Choisir une catégorie"
          label="Catégorie"
          helperText="Sélectionnez la catégorie de votre transaction"
        />
      </View>

      {/* Select avec recherche */}
      <View style={{ marginBottom: 24 }}>
        <Text style={{ fontSize: 16, fontWeight: '600', marginBottom: 8 }}>
          Devise
        </Text>
        <Select
          options={currencyOptions}
          value={selectedCurrency}
          onValueChange={setSelectedCurrency}
          placeholder="Sélectionner une devise"
          searchable={true}
          label="Devise"
          size="md"
          variant="outline"
        />
      </View>

      {/* Select avec descriptions */}
      <View style={{ marginBottom: 24 }}>
        <Text style={{ fontSize: 16, fontWeight: '600', marginBottom: 8 }}>
          Compte bancaire
        </Text>
        <Select
          options={accountOptions}
          value={selectedAccount}
          onValueChange={setSelectedAccount}
          placeholder="Choisir un compte"
          searchable={true}
          label="Compte"
          size="lg"
          variant="filled"
          helperText="Sélectionnez le compte à utiliser pour cette transaction"
        />
      </View>

      {/* Select désactivé */}
      <View style={{ marginBottom: 24 }}>
        <Text style={{ fontSize: 16, fontWeight: '600', marginBottom: 8 }}>
          Select désactivé
        </Text>
        <Select
          options={categoryOptions}
          value="food"
          onValueChange={() => {}}
          placeholder="Ce select est désactivé"
          disabled={true}
          label="Désactivé"
          error="Ce champ est temporairement indisponible"
        />
      </View>

      {/* Select avec erreur */}
      <View style={{ marginBottom: 24 }}>
        <Text style={{ fontSize: 16, fontWeight: '600', marginBottom: 8 }}>
          Select avec erreur
        </Text>
        <Select
          options={currencyOptions}
          value={selectedCurrency}
          onValueChange={setSelectedCurrency}
          placeholder="Sélectionner une devise"
          label="Devise"
          error="Veuillez sélectionner une devise valide"
        />
      </View>

      {/* Affichage des valeurs sélectionnées */}
      <View style={{ 
        backgroundColor: 'white', 
        padding: 16, 
        borderRadius: 8, 
        marginTop: 20 
      }}>
        <Text style={{ fontSize: 16, fontWeight: '600', marginBottom: 12 }}>
          Valeurs sélectionnées :
        </Text>
        <Text style={{ marginBottom: 4 }}>
          Catégorie : {selectedCategory || 'Aucune'}
        </Text>
        <Text style={{ marginBottom: 4 }}>
          Devise : {selectedCurrency || 'Aucune'}
        </Text>
        <Text style={{ marginBottom: 4 }}>
          Compte : {selectedAccount || 'Aucun'}
        </Text>
      </View>
    </ScrollView>
  );
} 