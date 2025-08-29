# Fonctionnalité des budgets

## Vue d'ensemble

La fonctionnalité des budgets permet aux utilisateurs de créer, gérer et suivre leurs budgets par catégorie de dépenses. Elle offre une vue claire de l'utilisation des budgets avec des indicateurs visuels et des alertes.

## Fonctionnalités implémentées

### 1. **Création de budgets**
- **Formulaire complet** avec validation
- **Sélection de catégorie** parmi les catégories de dépenses
- **Montant prévu** avec validation
- **Période** (hebdomadaire, mensuel, annuel)
- **Description optionnelle**

### 2. **Affichage des budgets**
- **Cartes de budget** avec statuts visuels
- **Barre de progression** colorée selon l'utilisation
- **Indicateurs de statut** (Actif, Dépassé, Terminé)
- **Montants restants** et jours restants

### 3. **Gestion des budgets**
- **Modification** des budgets existants
- **Suppression** avec confirmation
- **Actualisation automatique** des données

### 4. **Statuts des budgets**
- **Actif** : Budget en cours d'utilisation
- **Dépassé** : Montant dépensé > montant prévu
- **Terminé** : Période écoulée

## Composants créés

### 1. **BudgetForm** (`src/components/forms/BudgetForm.tsx`)
Formulaire modal pour créer et modifier les budgets.

**Fonctionnalités :**
- Validation des champs obligatoires
- Sélection de catégorie avec icônes
- Choix de période (hebdomadaire/mensuel/annuel)
- Résumé en temps réel
- Gestion des erreurs

**Utilisation :**
```typescript
<BudgetForm
  visible={showBudgetForm}
  onClose={() => setShowBudgetForm(false)}
  onSubmit={handleCreateBudget}
  categories={categories}
  initialData={budgetToEdit} // Optionnel pour la modification
/>
```

### 2. **BudgetCard** (`src/components/BudgetCard.tsx`)
Carte d'affichage d'un budget avec toutes ses informations.

**Fonctionnalités :**
- Affichage du statut avec couleur
- Barre de progression dynamique
- Actions d'édition et suppression
- Alertes pour les budgets dépassés
- Formatage des montants

**Utilisation :**
```typescript
<BudgetCard
  budget={budget}
  onEdit={() => handleEditBudget(budget.id)}
  onDelete={() => handleDeleteBudget(budget.id)}
/>
```

### 3. **BudgetService** (`src/services/budgetService/budgetApi.ts`)
Service API pour gérer les budgets.

**Méthodes disponibles :**
- `getBudgets()` : Récupérer tous les budgets
- `getBudget(id)` : Récupérer un budget spécifique
- `createBudget(data)` : Créer un nouveau budget
- `updateBudget(id, data)` : Modifier un budget
- `deleteBudget(id)` : Supprimer un budget
- `getBudgetsWithStatus()` : Récupérer les budgets avec statuts
- `getBudgetsByCategory(categoryId)` : Récupérer les budgets par catégorie

## Interface utilisateur

### 1. **Section Budgets dans Finances**
- **Bouton "Ajouter"** pour créer un nouveau budget
- **État vide** avec invitation à créer le premier budget
- **Liste des budgets** avec cartes détaillées

### 2. **Formulaire de création**
- **Modal plein écran** avec navigation
- **Champs validés** avec messages d'erreur
- **Sélection visuelle** des catégories
- **Résumé en temps réel**

### 3. **Cartes de budget**
- **Header** avec nom et statut
- **Barre de progression** colorée
- **Détails** (montants, jours restants)
- **Actions** (édition, suppression)
- **Alertes** pour budgets dépassés

## Types de données

### Budget
```typescript
interface Budget {
  id: string;
  userId: string;
  categoryId: string;
  name: string;
  amountPlanned: number;
  amountSpent: number;
  period: 'weekly' | 'monthly' | 'yearly';
  description?: string;
  status: 'active' | 'completed' | 'overdue';
  percentageUsed: number;
  remainingAmount: number;
  daysRemaining: number;
  createdAt: Date;
  updatedAt: Date;
}
```

### CreateBudgetRequest
```typescript
interface CreateBudgetRequest {
  name: string;
  categoryId: string;
  amountPlanned: number;
  period: 'weekly' | 'monthly' | 'yearly';
  description?: string;
}
```

## Logique métier

### 1. **Calcul des statuts**
- **Actif** : Période en cours et montant dépensé < montant prévu
- **Dépassé** : Montant dépensé > montant prévu
- **Terminé** : Période écoulée

### 2. **Couleurs de progression**
- **Vert** : < 75% utilisé
- **Orange** : 75-90% utilisé
- **Rouge** : > 90% utilisé

### 3. **Validation**
- **Nom requis** : Non vide
- **Catégorie requise** : Sélection obligatoire
- **Montant valide** : > 0
- **Période valide** : weekly, monthly, yearly

## Intégration

### 1. **Page des finances**
- **Section budgets** intégrée
- **Gestion des états** (chargement, erreur, vide)
- **Actualisation automatique** après modifications

### 2. **API**
- **Endpoints REST** pour CRUD des budgets
- **Gestion des erreurs** avec messages utilisateur
- **Types TypeScript** pour la sécurité

### 3. **Navigation**
- **Modal pour création/modification**
- **Retour automatique** après actions
- **Confirmation** pour suppressions

## Utilisation

### 1. **Créer un budget**
1. Cliquer sur "Ajouter" dans la section Budgets
2. Remplir le formulaire (nom, catégorie, montant, période)
3. Ajouter une description optionnelle
4. Valider la création

### 2. **Modifier un budget**
1. Cliquer sur l'icône d'édition sur une carte de budget
2. Modifier les champs souhaités
3. Sauvegarder les modifications

### 3. **Supprimer un budget**
1. Cliquer sur l'icône de suppression
2. Confirmer la suppression
3. Le budget est supprimé définitivement

### 4. **Suivre l'utilisation**
- **Barre de progression** : Utilisation en temps réel
- **Statut** : Indicateur visuel de l'état
- **Alertes** : Notifications pour budgets dépassés
- **Montants** : Dépensé vs prévu vs restant

## Fonctionnalités à venir

### 1. **Notifications**
- **Alertes push** pour budgets dépassés
- **Rappels** avant la fin de période
- **Résumés** hebdomadaires/mensuels

### 2. **Analytics**
- **Graphiques** d'évolution des budgets
- **Comparaisons** entre périodes
- **Tendances** d'utilisation

### 3. **Automatisation**
- **Budgets récurrents** automatiques
- **Ajustements** basés sur l'historique
- **Suggestions** de montants

### 4. **Partage**
- **Budgets partagés** entre utilisateurs
- **Collaboration** sur les objectifs
- **Responsabilités** par catégorie

## Structure des fichiers

```
src/
├── components/
│   ├── forms/
│   │   └── BudgetForm.tsx          # Formulaire de création/modification
│   └── BudgetCard.tsx              # Carte d'affichage d'un budget
├── services/
│   └── budgetService/
│       └── budgetApi.ts            # Service API des budgets
└── types/
    └── budget.d.ts                 # Types TypeScript (à créer)

app/dashboard/tabs/
└── finances.tsx                    # Page principale avec section budgets
```

## Tests recommandés

### 1. **Tests unitaires**
- Validation du formulaire
- Calculs des statuts
- Formatage des données

### 2. **Tests d'intégration**
- Création/modification/suppression
- Gestion des erreurs API
- Actualisation des données

### 3. **Tests utilisateur**
- Parcours de création
- Gestion des cas d'erreur
- Performance avec beaucoup de budgets

## Notes techniques

1. **Performance** : Les budgets sont chargés avec les autres données financières
2. **Sécurité** : Validation côté client et serveur
3. **Accessibilité** : Support des lecteurs d'écran et navigation clavier
4. **Responsive** : Adaptation aux différentes tailles d'écran
5. **Offline** : Support du mode hors ligne (à implémenter) 