# Composant TaskForm

Un composant de formulaire complet pour la création et la gestion de tâches, adapté pour les applications de gestion de projets et de productivité.

## Installation

Le composant TaskForm utilise les mêmes dépendances que le TransactionForm :

```bash
npm install @expo/vector-icons
```

## Utilisation de base

```tsx
import React, { useState } from 'react';
import { View } from 'react-native';
import { TaskForm, CreateTaskRequest } from '@/src/components/forms/TaskForm';

export default function MyComponent() {
  const [showTaskForm, setShowTaskForm] = useState(false);

  const handleCreateTask = async (task: CreateTaskRequest) => {
    console.log('Nouvelle tâche:', task);
    // Appel API pour sauvegarder la tâche
    // await taskService.createTask(task);
  };

  return (
    <View>
      <TaskForm
        visible={showTaskForm}
        onClose={() => setShowTaskForm(false)}
        onSubmit={handleCreateTask}
        defaultStatus="todo"
      />
    </View>
  );
}
```

## Fonctionnalités

### 📝 **Champs de base**
- **Titre** - Obligatoire
- **Description** - Optionnel, zone de texte multiligne
- **Statut** - À faire, En cours, Terminé
- **Priorité** - Basse, Moyenne, Haute

### 🏷️ **Organisation**
- **Catégories** - Sélection avec icônes et couleurs
- **Tags** - Système de tags personnalisables
- **Date d'échéance** - Optionnel avec gestion de date

### ⏱️ **Gestion du temps**
- **Temps estimé** - En minutes
- **Validation** - Vérification des données saisies

### 🎨 **Interface utilisateur**
- **Modal natif** - Interface utilisateur native
- **Validation en temps réel** - Feedback immédiat
- **États de chargement** - Indicateurs visuels
- **Gestion d'erreurs** - Messages d'erreur clairs

## API Reference

### TaskFormProps

| Prop | Type | Défaut | Description |
|------|------|--------|-------------|
| `visible` | `boolean` | - | **Requis** - Contrôle la visibilité du modal |
| `onClose` | `() => void` | - | **Requis** - Callback appelé à la fermeture |
| `onSubmit` | `(task: CreateTaskRequest) => Promise<void>` | - | **Requis** - Callback appelé lors de la soumission |
| `categories` | `Array<Category>` | `[]` | Liste des catégories disponibles |
| `defaultStatus` | `'todo' \| 'in_progress' \| 'completed'` | `'todo'` | Statut par défaut de la tâche |

### CreateTaskRequest

| Propriété | Type | Description |
|-----------|------|-------------|
| `title` | `string` | **Requis** - Titre de la tâche |
| `description` | `string \| undefined` | Description détaillée |
| `priority` | `'low' \| 'medium' \| 'high'` | Niveau de priorité |
| `status` | `'todo' \| 'in_progress' \| 'completed'` | Statut actuel |
| `dueDate` | `Date \| undefined` | Date d'échéance |
| `category` | `string \| undefined` | ID de la catégorie |
| `tags` | `string[] \| undefined` | Liste des tags |
| `estimatedTime` | `number \| undefined` | Temps estimé en minutes |

### Category

| Propriété | Type | Description |
|-----------|------|-------------|
| `id` | `string` | Identifiant unique |
| `name` | `string` | Nom de la catégorie |
| `color` | `string` | Couleur hexadécimale |
| `icon` | `string` | Nom de l'icône |

## Exemples d'utilisation

### Formulaire simple

```tsx
<TaskForm
  visible={showForm}
  onClose={() => setShowForm(false)}
  onSubmit={handleCreateTask}
/>
```

### Avec catégories

```tsx
const categories = [
  { id: 'work', name: 'Travail', color: '#3b82f6', icon: 'fa6:briefcase' },
  { id: 'personal', name: 'Personnel', color: '#10b981', icon: 'fa6:user' },
  { id: 'health', name: 'Santé', color: '#ef4444', icon: 'fa6:heart-pulse' },
];

<TaskForm
  visible={showForm}
  onClose={() => setShowForm(false)}
  onSubmit={handleCreateTask}
  categories={categories}
  defaultStatus="in_progress"
/>
```

### Gestion complète des tâches

```tsx
import React, { useState } from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import { TaskForm, CreateTaskRequest } from '@/src/components/forms/TaskForm';

export default function TaskManager() {
  const [showForm, setShowForm] = useState(false);
  const [tasks, setTasks] = useState<CreateTaskRequest[]>([]);

  const handleCreateTask = async (task: CreateTaskRequest) => {
    try {
      // Appel API
      const newTask = await taskService.createTask(task);
      
      // Mise à jour locale
      setTasks(prev => [...prev, newTask]);
      
      // Fermeture du formulaire
      setShowForm(false);
    } catch (error) {
      console.error('Erreur lors de la création:', error);
    }
  };

  return (
    <View>
      <TouchableOpacity onPress={() => setShowForm(true)}>
        <Text>+ Nouvelle Tâche</Text>
      </TouchableOpacity>

      <TaskForm
        visible={showForm}
        onClose={() => setShowForm(false)}
        onSubmit={handleCreateTask}
        categories={taskCategories}
      />
    </View>
  );
}
```

## Validation

Le formulaire inclut une validation automatique :

- **Titre** - Obligatoire, non vide
- **Temps estimé** - Doit être positif si renseigné
- **Date d'échéance** - Validation de format
- **Tags** - Pas de doublons

## États visuels

### Priorités
- **Basse** - Vert (`#10b981`)
- **Moyenne** - Orange (`#f59e0b`)
- **Haute** - Rouge (`#ef4444`)

### Statuts
- **À faire** - Gris (`#6b7280`)
- **En cours** - Bleu (`#3b82f6`)
- **Terminé** - Vert (`#10b981`)

## Avantages

- ✅ **TypeScript** - Typage complet et autocomplétion
- ✅ **Validation** - Vérification automatique des données
- ✅ **Accessibilité** - Support des états d'erreur et de chargement
- ✅ **Personnalisable** - Catégories et tags flexibles
- ✅ **Réutilisable** - API simple et cohérente
- ✅ **Performance** - Rendu optimisé avec React Native
- ✅ **UX moderne** - Interface utilisateur intuitive

## Migration depuis d'autres formulaires

Le TaskForm est conçu pour remplacer les formulaires de tâches basiques avec une API plus moderne et des fonctionnalités avancées comme les tags et la gestion du temps. 