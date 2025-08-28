# Fonctionnalités de la page de détail des comptes

## Vue d'ensemble

Une page de détail complète a été créée pour chaque compte, accessible en cliquant sur un compte dans la page des finances. Cette page offre une vue détaillée des transactions, des statistiques et des actions possibles sur le compte.

## Fonctionnalités implémentées

### 1. Navigation et interface
- **Navigation** : Bouton retour vers la page des finances
- **Header** : Affichage du nom du compte et du type
- **Actions rapides** : Boutons d'édition et de suppression du compte

### 2. Résumé du compte
- **Solde actuel** : Affichage du solde avec l'icône et la couleur du compte
- **Statistiques** : 
  - Total des revenus
  - Total des dépenses
  - Nombre de transactions

### 3. Filtres temporels
- **Périodes** : Semaine, Mois, Année
- **Navigation par onglets** : Transactions et Statistiques

### 4. Actions rapides
- **Nouvelle transaction** : Ajouter une transaction au compte
- **Transfert** : Effectuer un transfert entre comptes (à implémenter)
- **Export** : Exporter les données du compte (à implémenter)
- **Partage** : Partager les informations du compte (à implémenter)

### 5. Section Transactions
- **Liste des transactions** : Affichage chronologique avec icônes et couleurs
- **Détails** : Montant, catégorie, date pour chaque transaction
- **Navigation** : Clic sur une transaction pour voir ses détails
- **État vide** : Message et bouton d'action quand aucune transaction

### 6. Section Statistiques
- **Graphique d'évolution du solde** : Ligne temporelle du solde
- **Graphique revenus vs dépenses** : Comparaison en barres
- **Graphique par catégorie** : Répartition en camembert
- **États vides** : Messages appropriés quand pas de données

### 7. Détails du compte
- **Informations** : Type de compte, devise, date de création
- **Activité** : Dernière transaction effectuée

### 8. Page de détail des transactions
- **Vue complète** : Toutes les informations de la transaction
- **Actions** : Modifier, dupliquer, marquer comme récurrent, supprimer
- **Navigation** : Retour vers la page du compte

## Structure des fichiers

```
app/dashboard/
├── account/
│   └── [id].tsx              # Page de détail du compte
└── transaction/
    └── [id].tsx              # Page de détail de la transaction

src/components/
├── charts/
│   └── AccountChart.tsx      # Composant de graphiques
└── AccountActions.tsx        # Composant d'actions rapides

src/services/
├── accountService/
│   └── accountApi.ts         # API des comptes (ajout deleteAccount)
├── transactionService/
│   └── transactionApi.ts     # API des transactions (ajout getTransactionsByAccount)
└── categoryService/
    └── categoryApi.ts        # API des catégories (ajout getCategory)
```

## Fonctionnalités à implémenter

### Actions rapides
- [ ] Transfert entre comptes
- [ ] Export des données (PDF, CSV)
- [ ] Partage des informations

### Gestion des transactions
- [ ] Modification d'une transaction
- [ ] Duplication d'une transaction
- [ ] Gestion des transactions récurrentes
- [ ] Filtrage avancé des transactions

### Statistiques avancées
- [ ] Filtrage par période personnalisée
- [ ] Comparaison avec les périodes précédentes
- [ ] Tendances et prévisions

### Interface utilisateur
- [ ] Animations de transition
- [ ] Mode sombre
- [ ] Personnalisation des couleurs
- [ ] Notifications pour les actions importantes

## Utilisation

1. **Accéder à un compte** : Cliquer sur un compte dans la page des finances
2. **Voir les transactions** : Onglet "Transactions" par défaut
3. **Voir les statistiques** : Basculer vers l'onglet "Statistiques"
4. **Ajouter une transaction** : Utiliser le bouton "Nouvelle Transaction"
5. **Voir les détails d'une transaction** : Cliquer sur une transaction dans la liste
6. **Modifier le compte** : Utiliser le bouton d'édition dans le header
7. **Supprimer le compte** : Utiliser le bouton de suppression (avec confirmation)

## Technologies utilisées

- **React Native** : Interface utilisateur
- **Expo Router** : Navigation entre les pages
- **react-native-chart-kit** : Graphiques et visualisations
- **Ionicons** : Icônes
- **Tailwind CSS** : Styles et mise en page
- **TypeScript** : Typage statique

## API Endpoints utilisés

- `GET /accounts/{id}` : Récupérer les détails d'un compte
- `DELETE /accounts/{id}` : Supprimer un compte
- `GET /transactions?accountId={id}` : Récupérer les transactions d'un compte
- `GET /transactions/{id}` : Récupérer les détails d'une transaction
- `DELETE /transactions/{id}` : Supprimer une transaction
- `GET /categories/{id}` : Récupérer les détails d'une catégorie
- `POST /transactions` : Créer une nouvelle transaction 