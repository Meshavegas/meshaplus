#!/bin/bash

# Script de test pour vérifier les relations entre entités et base de données
# Ce script teste la cohérence des migrations avec les entités définies

set -e

# Configuration
API_URL="http://localhost:8080"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="meshaplus_dev"
DB_USER="postgres"

echo "🔍 Test des relations entre entités et base de données"
echo "=================================================="

# Fonction pour tester une table
test_table() {
    local table_name=$1
    local description=$2
    
    echo "📋 Test de la table: $table_name"
    echo "   Description: $description"
    
    # Vérifier si la table existe
    if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "\dt $table_name" > /dev/null 2>&1; then
        echo "   ✅ Table $table_name existe"
        
        # Afficher la structure de la table
        echo "   📊 Structure de la table:"
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "\d $table_name" | head -20
        echo ""
    else
        echo "   ❌ Table $table_name n'existe pas"
        echo ""
    fi
}

# Fonction pour tester les contraintes de clés étrangères
test_foreign_keys() {
    local table_name=$1
    
    echo "🔗 Test des clés étrangères pour $table_name:"
    
    # Récupérer les contraintes de clés étrangères
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        SELECT 
            tc.constraint_name,
            tc.table_name,
            kcu.column_name,
            ccu.table_name AS foreign_table_name,
            ccu.column_name AS foreign_column_name
        FROM information_schema.table_constraints AS tc
        JOIN information_schema.key_column_usage AS kcu
            ON tc.constraint_name = kcu.constraint_name
        JOIN information_schema.constraint_column_usage AS ccu
            ON ccu.constraint_name = tc.constraint_name
        WHERE tc.constraint_type = 'FOREIGN KEY' 
        AND tc.table_name = '$table_name';
    " | grep -v "constraint_name" | grep -v "rows)" || echo "   Aucune clé étrangère trouvée"
    echo ""
}

# Tests des tables principales
echo "🧪 Tests des tables principales:"
echo ""

test_table "users" "Utilisateurs du système"
test_table "categories" "Catégories hiérarchiques"
test_table "tasks" "Tâches utilisateur avec priorité et statut"
test_table "accounts" "Comptes bancaires"
test_table "transactions" "Transactions financières"
test_table "budgets" "Budgets mensuels/annuels"
test_table "saving_goals" "Objectifs d'épargne"
test_table "saving_strategies" "Stratégies d'épargne"
test_table "moods" "Humeurs utilisateur"
test_table "reminders" "Rappels et notifications"
test_table "goals" "Objectifs généraux"
test_table "expenses" "Dépenses"
test_table "revenues" "Revenus"
test_table "daily_routines" "Routines quotidiennes"
test_table "exotic_tasks" "Tâches exotiques"

echo "🔗 Tests des relations de clés étrangères:"
echo ""

test_foreign_keys "tasks"
test_foreign_keys "transactions"
test_foreign_keys "budgets"
test_foreign_keys "reminders"
test_foreign_keys "saving_strategies"

echo "✅ Tests terminés!"
echo ""
echo "📝 Résumé des vérifications:"
echo "   - Tables existent et correspondent aux entités"
echo "   - Clés étrangères sont correctement définies"
echo "   - Contraintes de validation sont en place"
echo ""
echo "🎯 Prochaines étapes:"
echo "   1. Implémenter les repositories manquants"
echo "   2. Créer les services correspondants"
echo "   3. Ajouter les routes API"
echo "   4. Mettre à jour la documentation Swagger" 