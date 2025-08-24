# Composant Select

Un composant Select réutilisable et personnalisable pour React Native, basé sur le système d'icônes unifié.

## Installation

Le composant Select utilise le système d'icônes unifié. Assurez-vous que les dépendances sont installées :

```bash
npm install @expo/vector-icons
```

## Utilisation de base

```tsx
import React, { useState } from 'react';
import { View } from 'react-native';
import Select, { SelectOption } from '@/src/components/Select';

export default function MyComponent() {
  const [selectedValue, setSelectedValue] = useState<string>();

  const options: SelectOption[] = [
    { label: 'Option 1', value: 'option1' },
    { label: 'Option 2', value: 'option2' },
    { label: 'Option 3', value: 'option3' },
  ];

  return (
    <View>
      <Select
        options={options}
        value={selectedValue}
        onValueChange={setSelectedValue}
        placeholder="Sélectionner une option"
        label="Mon Select"
      />
    </View>
  );
}
```

## Options avec icônes et descriptions

```tsx
const options: SelectOption[] = [
  { 
    label: 'Alimentation', 
    value: 'food', 
    icon: 'fa6:utensils',
    description: 'Dépenses alimentaires'
  },
  { 
    label: 'Transport', 
    value: 'transport', 
    icon: 'fa6:car',
    description: 'Frais de transport'
  },
  { 
    label: 'Logement', 
    value: 'housing', 
    icon: 'fa6:house',
    description: 'Loyer et charges'
  },
];
```

## Fonctionnalités

### 🔍 Recherche
Activez la recherche pour filtrer les options :

```tsx
<Select
  options={options}
  value={selectedValue}
  onValueChange={setSelectedValue}
  searchable={true}
  placeholder="Rechercher et sélectionner..."
/>
```

### 📏 Tailles
Trois tailles disponibles : `sm`, `md`, `lg`

```tsx
<Select size="sm" />  // 32px de hauteur
<Select size="md" />  // 40px de hauteur (défaut)
<Select size="lg" />  // 48px de hauteur
```

### 🎨 Variants
Trois styles différents : `default`, `outline`, `filled`

```tsx
<Select variant="default" />  // Bordure grise, fond blanc
<Select variant="outline" />  // Bordure bleue, fond transparent
<Select variant="filled" />   // Fond gris, pas de bordure
```

### 🚫 États désactivés et erreurs

```tsx
// Désactivé
<Select
  options={options}
  disabled={true}
  error="Ce champ est temporairement indisponible"
/>

// Avec erreur
<Select
  options={options}
  error="Veuillez sélectionner une option valide"
/>
```

## API Reference

### SelectProps

| Prop | Type | Défaut | Description |
|------|------|--------|-------------|
| `options` | `SelectOption[]` | - | **Requis** - Liste des options disponibles |
| `value` | `T` | - | Valeur actuellement sélectionnée |
| `onValueChange` | `(value: T) => void` | - | Callback appelé lors de la sélection |
| `placeholder` | `string` | "Sélectionner une option" | Texte affiché quand aucune option n'est sélectionnée |
| `disabled` | `boolean` | `false` | Désactive le composant |
| `searchable` | `boolean` | `false` | Active la recherche dans les options |
| `multiple` | `boolean` | `false` | Permet la sélection multiple (non implémenté) |
| `size` | `'sm' \| 'md' \| 'lg'` | `'md'` | Taille du composant |
| `variant` | `'default' \| 'outline' \| 'filled'` | `'default'` | Style du composant |
| `width` | `number \| string` | `'100%'` | Largeur du composant |
| `label` | `string` | - | Label affiché au-dessus du select |
| `error` | `string` | - | Message d'erreur |
| `helperText` | `string` | - | Texte d'aide |
| `onOpen` | `() => void` | - | Callback appelé à l'ouverture |
| `onClose` | `() => void` | - | Callback appelé à la fermeture |

### SelectOption

| Prop | Type | Description |
|------|------|-------------|
| `label` | `string` | **Requis** - Texte affiché pour l'option |
| `value` | `T` | **Requis** - Valeur de l'option |
| `disabled` | `boolean` | Désactive cette option |
| `icon` | `string` | Nom de l'icône (ex: 'fa6:house') |
| `description` | `string` | Description optionnelle |

## Exemples d'utilisation

### Select pour catégories de dépenses

```tsx
const categoryOptions: SelectOption[] = [
  { label: 'Alimentation', value: 'food', icon: 'fa6:utensils' },
  { label: 'Transport', value: 'transport', icon: 'fa6:car' },
  { label: 'Logement', value: 'housing', icon: 'fa6:house' },
  { label: 'Loisirs', value: 'entertainment', icon: 'fa6:gamepad' },
  { label: 'Santé', value: 'health', icon: 'fa6:heart-pulse' },
  { label: 'Shopping', value: 'shopping', icon: 'fa6:shopping-bag' },
];

<Select
  options={categoryOptions}
  value={selectedCategory}
  onValueChange={setSelectedCategory}
  placeholder="Choisir une catégorie"
  label="Catégorie de dépense"
  searchable={true}
  size="md"
  variant="outline"
/>
```

### Select pour devises

```tsx
const currencyOptions: SelectOption[] = [
  { label: 'Euro (€)', value: 'EUR', icon: 'fa6:euro-sign' },
  { label: 'Dollar US ($)', value: 'USD', icon: 'fa6:dollar-sign' },
  { label: 'Livre Sterling (£)', value: 'GBP', icon: 'fa6:sterling-sign' },
  { label: 'Franc Suisse (CHF)', value: 'CHF', icon: 'fa6:franc-sign' },
];

<Select
  options={currencyOptions}
  value={selectedCurrency}
  onValueChange={setSelectedCurrency}
  placeholder="Sélectionner une devise"
  label="Devise"
  helperText="Choisissez la devise pour vos transactions"
/>
```

### Select pour comptes bancaires

```tsx
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
    icon: 'fa6:piggy-bank',
    description: 'Compte d\'épargne avec intérêts'
  },
];

<Select
  options={accountOptions}
  value={selectedAccount}
  onValueChange={setSelectedAccount}
  placeholder="Choisir un compte"
  label="Compte bancaire"
  searchable={true}
  size="lg"
  variant="filled"
/>
```

## Avantages

- ✅ **TypeScript** - Typage complet et autocomplétion
- ✅ **Icônes** - Support des icônes FontAwesome 6
- ✅ **Recherche** - Filtrage en temps réel des options
- ✅ **Accessibilité** - Support des états désactivés et erreurs
- ✅ **Personnalisable** - Tailles, variants et styles flexibles
- ✅ **Réutilisable** - API simple et cohérente
- ✅ **Performance** - Rendu optimisé avec React Native

## Migration depuis d'autres composants

Le composant Select est conçu pour remplacer les Picker natifs et autres composants de sélection avec une API plus moderne et flexible. 