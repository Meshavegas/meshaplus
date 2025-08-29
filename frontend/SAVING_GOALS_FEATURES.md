# Fonctionnalité des Objectifs d'Épargne

## Vue d'ensemble

La fonctionnalité des objectifs d'épargne permet aux utilisateurs de créer, gérer et suivre leurs objectifs financiers avec des échéances spécifiques. Elle offre un système de suivi des progrès avec des contributions et des indicateurs visuels.

## Fonctionnalités implémentées

### 1. **Création d'objectifs d'épargne**
- **Formulaire complet** avec validation
- **Titre personnalisé** pour l'objectif
- **Montant cible** avec validation
- **Date limite** avec sélecteur de date
- **Fréquence de suivi** (hebdomadaire, mensuel, annuel)

### 2. **Affichage des objectifs**
- **Cartes d'objectifs** avec statuts visuels
- **Barre de progression** colorée selon l'avancement
- **Indicateurs de statut** (En cours, En retard, Atteint)
- **Montants restants** et jours restants

### 3. **Gestion des objectifs**
- **Modification** des objectifs existants
- **Suppression** avec confirmation
- **Actualisation automatique** des données

### 4. **Contributions**
- **Ajout de contributions** avec montants personnalisés
- **Suggestions rapides** de montants
- **Prévisualisation de l'impact** avant contribution
- **Validation** des montants

### 5. **Statuts des objectifs**
- **En cours** : Objectif actif avec échéance future
- **En retard** : Date limite dépassée
- **Atteint** : Montant cible atteint

## Structure des données

### Exemple de données
```json
{
  "deadline": "2024-06-30T00:00:00Z",
  "frequency": "monthly",
  "target_amount": 2000,
  "title": "Vacances d'été"
}
```

### Interface SavingGoal
```typescript
interface SavingGoal {
  id: string;
  userId: string;
  title: string;
  target_amount: number;
  current_amount: number;
  frequency: 'weekly' | 'monthly' | 'yearly';
  deadline: string;
  status: 'active' | 'completed' | 'overdue';
  percentage_achieved: number;
  remaining_amount: number;
  days_remaining: number;
  created_at: string;
  updated_at: string;
}
```

## Composants créés

### 1. **SavingGoalForm** (`src/components/forms/SavingGoalForm.tsx`)
Formulaire modal pour créer et modifier les objectifs d'épargne.

**Fonctionnalités :**
- Validation des champs obligatoires
- Sélecteur de date avec inputs séparés (jour/mois/année)
- Choix de fréquence de suivi
- Calcul automatique de la contribution mensuelle suggérée
- Conseils et motivation

**Utilisation :**
```typescript
<SavingGoalForm
  visible={showSavingGoalForm}
  onClose={() => setShowSavingGoalForm(false)}
  onSubmit={handleCreateSavingGoal}
  initialData={goalToEdit} // Optionnel pour la modification
/>
```

### 2. **SavingGoalCard** (`src/components/SavingGoalCard.tsx`)
Carte d'affichage d'un objectif d'épargne avec toutes ses informations.

**Fonctionnalités :**
- Affichage du statut avec couleur
- Barre de progression dynamique
- Actions d'édition, suppression et contribution
- Alertes pour les objectifs en retard
- Célébration pour les objectifs atteints
- Formatage des montants et dates

**Utilisation :**
```typescript
<SavingGoalCard
  goal={goal}
  onEdit={() => handleEditGoal(goal.id)}
  onDelete={() => handleDeleteGoal(goal.id)}
  onContribute={() => handleContribute(goal)}
/>
```

### 3. **ContributionForm** (`src/components/forms/ContributionForm.tsx`)
Formulaire pour ajouter des contributions aux objectifs d'épargne.

**Fonctionnalités :**
- Saisie du montant de contribution
- Suggestions rapides de montants (10%, 25%, 50%, 100%)
- Prévisualisation de l'impact de la contribution
- Validation des montants
- Messages de motivation

**Utilisation :**
```typescript
<ContributionForm
  visible={showContributionForm}
  onClose={() => setShowContributionForm(false)}
  onSubmit={handleAddContribution}
  goal={selectedGoal}
/>
```

### 4. **SavingGoalService** (`src/services/savingGoalService/savingGoalApi.ts`)
Service API pour gérer les objectifs d'épargne.

**Méthodes disponibles :**
- `getSavingGoals()` : Récupérer tous les objectifs
- `getSavingGoal(id)` : Récupérer un objectif spécifique
- `createSavingGoal(data)` : Créer un nouvel objectif
- `updateSavingGoal(id, data)` : Modifier un objectif
- `deleteSavingGoal(id)` : Supprimer un objectif
- `addContribution(id, amount)` : Ajouter une contribution
- `getSavingGoalsWithStatus()` : Récupérer les objectifs avec statuts

## Interface utilisateur

### 1. **Section Objectifs d'Épargne dans Finances**
- **Bouton "Ajouter"** pour créer un nouvel objectif
- **État vide** avec invitation à créer le premier objectif
- **Liste des objectifs** avec cartes détaillées

### 2. **Formulaire de création**
- **Modal plein écran** avec navigation
- **Champs validés** avec messages d'erreur
- **Sélecteur de date** avec inputs séparés
- **Résumé en temps réel** avec contribution suggérée

### 3. **Cartes d'objectifs**
- **Header** avec titre et statut
- **Barre de progression** colorée
- **Détails** (montants, jours restants, échéance)
- **Actions** (édition, suppression, contribution)
- **Alertes** pour objectifs en retard
- **Célébration** pour objectifs atteints

### 4. **Formulaire de contribution**
- **Informations de l'objectif** en cours
- **Saisie du montant** avec validation
- **Suggestions rapides** de montants
- **Prévisualisation de l'impact**
- **Messages de motivation**

## Logique métier

### 1. **Calcul des statuts**
- **En cours** : Échéance future et montant non atteint
- **En retard** : Date limite dépassée
- **Atteint** : Montant cible atteint ou dépassé

### 2. **Couleurs de progression**
- **Rouge** : < 50% atteint
- **Orange** : 50-75% atteint
- **Vert** : 75-99% atteint
- **Bleu** : 100% atteint

### 3. **Calcul de la contribution suggérée**
```typescript
const calculateMonthlyContribution = () => {
  const now = new Date();
  const monthsDiff = (deadline.getFullYear() - now.getFullYear()) * 12 + 
                    (deadline.getMonth() - now.getMonth());
  return monthsDiff > 0 ? target_amount / monthsDiff : target_amount;
};
```

### 4. **Validation**
- **Titre requis** : Non vide
- **Montant valide** : > 0
- **Date limite valide** : Dans le futur
- **Contribution valide** : ≤ montant restant

## Intégration

### 1. **Page des finances**
- **Section objectifs d'épargne** intégrée
- **Gestion des états** (chargement, erreur, vide)
- **Actualisation automatique** après modifications

### 2. **API**
- **Endpoints REST** pour CRUD des objectifs
- **Gestion des erreurs** avec messages utilisateur
- **Types TypeScript** pour la sécurité

### 3. **Navigation**
- **Modal pour création/modification**
- **Modal pour contributions**
- **Retour automatique** après actions
- **Confirmation** pour suppressions

## Utilisation

### 1. **Créer un objectif**
1. Cliquer sur "Ajouter" dans la section Objectifs d'Épargne
2. Remplir le formulaire (titre, montant, fréquence, date limite)
3. Valider la création

### 2. **Modifier un objectif**
1. Cliquer sur l'icône d'édition sur une carte d'objectif
2. Modifier les champs souhaités
3. Sauvegarder les modifications

### 3. **Supprimer un objectif**
1. Cliquer sur l'icône de suppression
2. Confirmer la suppression
3. L'objectif est supprimé définitivement

### 4. **Ajouter une contribution**
1. Cliquer sur "Contribuer" sur une carte d'objectif
2. Saisir le montant ou utiliser les suggestions
3. Voir l'impact de la contribution
4. Valider l'ajout

### 5. **Suivre les progrès**
- **Barre de progression** : Avancement en temps réel
- **Statut** : Indicateur visuel de l'état
- **Alertes** : Notifications pour objectifs en retard
- **Célébration** : Félicitations pour objectifs atteints
- **Montants** : Actuel vs cible vs restant

## Fonctionnalités à venir

### 1. **Notifications**
- **Alertes push** pour échéances approchantes
- **Rappels** de contributions
- **Célébrations** pour objectifs atteints

### 2. **Analytics**
- **Graphiques** d'évolution des objectifs
- **Comparaisons** entre objectifs
- **Tendances** de contributions

### 3. **Automatisation**
- **Contributions automatiques** récurrentes
- **Ajustements** basés sur l'historique
- **Suggestions** de nouveaux objectifs

### 4. **Partage**
- **Objectifs partagés** entre utilisateurs
- **Collaboration** sur les objectifs
- **Responsabilités** par objectif

## Structure des fichiers

```
src/
├── components/
│   ├── forms/
│   │   ├── SavingGoalForm.tsx          # Formulaire de création/modification
│   │   └── ContributionForm.tsx        # Formulaire de contribution
│   └── SavingGoalCard.tsx              # Carte d'affichage d'un objectif
├── services/
│   └── savingGoalService/
│       └── savingGoalApi.ts            # Service API des objectifs d'épargne
└── types/
    └── savingGoal.d.ts                 # Types TypeScript (à créer)

app/dashboard/tabs/
└── finances.tsx                        # Page principale avec section objectifs
```

## Tests recommandés

### 1. **Tests unitaires**
- Validation du formulaire
- Calculs des statuts
- Formatage des données

### 2. **Tests d'intégration**
- Création/modification/suppression
- Ajout de contributions
- Gestion des erreurs API
- Actualisation des données

### 3. **Tests utilisateur**
- Parcours de création d'objectif
- Ajout de contributions
- Gestion des cas d'erreur
- Performance avec beaucoup d'objectifs

## Notes techniques

1. **Performance** : Les objectifs sont chargés avec les autres données financières
2. **Sécurité** : Validation côté client et serveur
3. **Accessibilité** : Support des lecteurs d'écran et navigation clavier
4. **Responsive** : Adaptation aux différentes tailles d'écran
5. **Offline** : Support du mode hors ligne (à implémenter)
6. **Date Picker** : Implémentation avec inputs séparés pour éviter les problèmes de compatibilité

## Exemple d'utilisation complète

```typescript
// Créer un objectif
const newGoal = {
  title: "Vacances d'été",
  target_amount: 2000,
  frequency: "monthly",
  deadline: "2024-06-30T00:00:00Z"
};

// Ajouter une contribution
const contribution = 500;

// Résultat attendu
const updatedGoal = {
  ...newGoal,
  current_amount: 500,
  percentage_achieved: 25,
  remaining_amount: 1500,
  status: "active"
};
``` 