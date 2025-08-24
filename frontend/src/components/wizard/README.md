# Setup Wizard - Flow UX Étape par Étape

Un wizard complet pour configurer le profil utilisateur avec un flow UX guidé et une sortie JSON prête pour le backend.

## 🎯 Fonctionnalités

- **Flow guidé** - 5 étapes progressives avec navigation intuitive
- **Progression visuelle** - Barre de progression et indicateurs d'étape
- **Validation** - Vérification des données en temps réel
- **Sortie JSON** - Données structurées prêtes pour l'API
- **Interface native** - Modal avec animations fluides
- **Responsive** - Adaptation automatique au clavier

## 🚀 Utilisation

### Import et utilisation de base

```tsx
import React, { useState } from 'react';
import { SetupWizard, WizardData } from '@/src/components/wizard/SetupWizard';

export default function MyComponent() {
  const [showWizard, setShowWizard] = useState(false);

  const handleWizardComplete = async (data: WizardData) => {
    console.log('Configuration terminée:', data);
    
    // Envoyer les données au backend
    await apiService.saveUserSetup(data);
    
    setShowWizard(false);
  };

  return (
    <View>
      <TouchableOpacity onPress={() => setShowWizard(true)}>
        <Text>Configurer mon profil</Text>
      </TouchableOpacity>

      <SetupWizard
        visible={showWizard}
        onClose={() => setShowWizard(false)}
        onComplete={handleWizardComplete}
      />
    </View>
  );
}
```

## 📋 Étapes du Wizard

### 🟢 Écran 1 — Bienvenue
- Message d'accueil avec icône
- Bouton "Commencer"

### 🔹 Étape 1 — Revenus & Comptes
- **Montant total mensuel** (input numérique)
- **Dettes/crédits** (Oui/Non)

### 🔹 Étape 2 — Dépenses & Budget
- **Nourriture** (input numérique)
- **Transport** (input numérique)

### 🔹 Étape 3 — Objectifs Financiers
- **Épargne mensuelle** (input numérique)

### 🔹 Étape 4 — Organisation & Habitudes
- **Temps de planification** (Matin/Midi/Soir)

### 🟢 Écran Final — Résumé & Confirmation
- Affichage des données collectées
- Bouton "Lancer mon tableau de bord 🚀"

## 📊 Structure des Données

### WizardData Interface

```typescript
interface WizardData {
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
```

### Exemple de Sortie JSON

```json
{
  "income": {
    "sources": [],
    "monthly_total": 250000,
    "accounts": [],
    "has_debt": true,
    "debt_amount": 50000
  },
  "expenses": {
    "top_categories": [],
    "food": 60000,
    "transport": 20000,
    "housing": 0,
    "subscriptions": 0,
    "alerts_enabled": true,
    "auto_budget": true
  },
  "goals": {
    "main_goal": "",
    "savings_target": 50000,
    "deadline": "3 mois",
    "advice_enabled": true
  },
  "habits": {
    "planning_time": "Matin",
    "daily_focus_time": "30min",
    "summary_type": "Hebdomadaire"
  }
}
```

## 🎨 API Reference

### SetupWizardProps

| Prop | Type | Description |
|------|------|-------------|
| `visible` | `boolean` | Contrôle la visibilité du modal |
| `onClose` | `() => void` | Callback appelé à la fermeture |
| `onComplete` | `(data: WizardData) => Promise<void>` | Callback appelé à la fin du wizard |

### Méthodes

| Méthode | Description |
|---------|-------------|
| `nextStep()` | Passe à l'étape suivante |
| `prevStep()` | Retourne à l'étape précédente |
| `updateWizardData()` | Met à jour les données du wizard |

## 🎯 Intégration Backend

### Exemple d'envoi des données

```typescript
const handleWizardComplete = async (data: WizardData) => {
  try {
    // Envoyer au backend
    const response = await fetch('/api/user/setup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      // Succès
      Alert.alert('Succès', 'Configuration sauvegardée !');
    } else {
      throw new Error('Erreur serveur');
    }
  } catch (error) {
    Alert.alert('Erreur', 'Impossible de sauvegarder la configuration');
  }
};
```

## 🎨 Personnalisation

### Styles personnalisés

Le wizard utilise Tailwind CSS pour le styling. Vous pouvez personnaliser :

- **Couleurs** - Modifier les classes de couleur
- **Espacement** - Ajuster les marges et paddings
- **Typographie** - Changer les tailles de police
- **Animations** - Modifier les transitions

### Ajout d'étapes

Pour ajouter de nouvelles étapes :

1. Incrémenter `totalSteps`
2. Ajouter une nouvelle fonction `renderNewStep()`
3. Mettre à jour `renderStep()` avec le nouveau cas
4. Étendre l'interface `WizardData` si nécessaire

## 🚀 Avantages

- ✅ **UX guidée** - Flow étape par étape intuitif
- ✅ **Validation** - Vérification automatique des données
- ✅ **Progression** - Indicateurs visuels de progression
- ✅ **JSON prêt** - Sortie structurée pour l'API
- ✅ **Responsive** - Adaptation au clavier mobile
- ✅ **TypeScript** - Typage complet et autocomplétion
- ✅ **Réutilisable** - API simple et flexible

## 🔧 Développement

### Structure des fichiers

```
src/components/wizard/
├── SetupWizard.tsx    # Composant principal
├── example.tsx        # Exemple d'utilisation
└── README.md          # Documentation
```

### Tests

Pour tester le wizard :

```bash
# Lancer l'exemple
npm run start
# Ouvrir l'exemple dans l'app
```

Le wizard est maintenant prêt à être intégré dans votre application ! 🎉 