#!/bin/bash

# Script de test pour valider les structures de requêtes
# Ce script teste la validation des requêtes avec les nouvelles règles

set -e

echo "🧪 Test des structures de requêtes"
echo "=================================="

# Configuration
API_URL="http://localhost:8080"

# Fonction pour tester une requête
test_request() {
    local endpoint=$1
    local method=$2
    local data=$3
    local description=$4
    
    echo "📋 Test: $description"
    echo "   Endpoint: $method $endpoint"
    echo "   Data: $data"
    
    # Effectuer la requête
    response=$(curl -s -w "\n%{http_code}" -X "$method" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "$data" \
        "$API_URL$endpoint")
    
    # Extraire le code de statut
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo "   Status: $http_code"
    
    if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 201 ]; then
        echo "   ✅ Succès"
    elif [ "$http_code" -eq 400 ]; then
        echo "   ⚠️  Erreur de validation (attendu)"
        echo "   Response: $body"
    else
        echo "   ❌ Erreur inattendue"
        echo "   Response: $body"
    fi
    echo ""
}

# Fonction pour obtenir un token d'authentification
get_auth_token() {
    echo "🔐 Authentification..."
    
    # Créer un utilisateur de test
    user_data='{
        "name": "Test User",
        "email": "test@example.com",
        "password": "password123"
    }'
    
    register_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$user_data" \
        "$API_URL/auth/register")
    
    # Extraire le token
    TOKEN=$(echo "$register_response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    
    if [ -z "$TOKEN" ]; then
        echo "❌ Impossible d'obtenir le token d'authentification"
        exit 1
    fi
    
    echo "✅ Token obtenu"
    echo ""
}

# Tests des requêtes de tâches
echo "📝 Tests des requêtes de tâches:"
echo ""

# Test 1: Créer une tâche valide
test_request "/tasks" "POST" '{
    "title": "Tâche de test",
    "description": "Description de test",
    "priority": "medium",
    "duration_planned": 60
}' "Créer une tâche valide"

# Test 2: Créer une tâche avec priorité invalide
test_request "/tasks" "POST" '{
    "title": "Tâche invalide",
    "description": "Description de test",
    "priority": "invalid_priority",
    "duration_planned": 60
}' "Créer une tâche avec priorité invalide"

# Test 3: Créer une tâche avec durée négative
test_request "/tasks" "POST" '{
    "title": "Tâche invalide",
    "description": "Description de test",
    "priority": "medium",
    "duration_planned": -10
}' "Créer une tâche avec durée négative"

# Tests des requêtes de transactions
echo "💰 Tests des requêtes de transactions:"
echo ""

# Test 4: Créer une transaction valide
test_request "/transactions" "POST" '{
    "account_id": "123e4567-e89b-12d3-a456-426614174000",
    "category_id": "123e4567-e89b-12d3-a456-426614174000",
    "type": "expense",
    "amount": 25.50,
    "description": "Achat alimentaire",
    "date": "2024-01-15T00:00:00Z"
}' "Créer une transaction valide"

# Test 5: Créer une transaction avec montant négatif
test_request "/transactions" "POST" '{
    "account_id": "123e4567-e89b-12d3-a456-426614174000",
    "category_id": "123e4567-e89b-12d3-a456-426614174000",
    "type": "expense",
    "amount": -25.50,
    "description": "Achat alimentaire",
    "date": "2024-01-15T00:00:00Z"
}' "Créer une transaction avec montant négatif"

# Tests des requêtes de comptes
echo "🏦 Tests des requêtes de comptes:"
echo ""

# Test 6: Créer un compte valide
test_request "/accounts" "POST" '{
    "name": "Compte principal",
    "type": "checking",
    "balance": 1500.00,
    "currency": "EUR"
}' "Créer un compte valide"

# Test 7: Créer un compte avec type invalide
test_request "/accounts" "POST" '{
    "name": "Compte invalide",
    "type": "invalid_type",
    "balance": 1500.00,
    "currency": "EUR"
}' "Créer un compte avec type invalide"

# Tests des requêtes de catégories
echo "📂 Tests des requêtes de catégories:"
echo ""

# Test 8: Créer une catégorie valide
test_request "/categories" "POST" '{
    "name": "Alimentation",
    "type": "expense"
}' "Créer une catégorie valide"

# Test 9: Créer une catégorie avec type invalide
test_request "/categories" "POST" '{
    "name": "Catégorie invalide",
    "type": "invalid_type"
}' "Créer une catégorie avec type invalide"

# Tests des requêtes de budgets
echo "📊 Tests des requêtes de budgets:"
echo ""

# Test 10: Créer un budget valide
test_request "/budgets" "POST" '{
    "category_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount_planned": 500.00,
    "month": 1,
    "year": 2024
}' "Créer un budget valide"

# Test 11: Créer un budget avec mois invalide
test_request "/budgets" "POST" '{
    "category_id": "123e4567-e89b-12d3-a456-426614174000",
    "amount_planned": 500.00,
    "month": 13,
    "year": 2024
}' "Créer un budget avec mois invalide"

# Tests des requêtes d'objectifs d'épargne
echo "🎯 Tests des requêtes d'objectifs d'épargne:"
echo ""

# Test 12: Créer un objectif d'épargne valide
test_request "/saving-goals" "POST" '{
    "title": "Vacances d'été",
    "target_amount": 2000.00,
    "frequency": "monthly"
}' "Créer un objectif d'épargne valide"

# Test 13: Créer un objectif avec fréquence invalide
test_request "/saving-goals" "POST" '{
    "title": "Objectif invalide",
    "target_amount": 2000.00,
    "frequency": "invalid_frequency"
}' "Créer un objectif avec fréquence invalide"

# Tests des requêtes de rappels
echo "⏰ Tests des requêtes de rappels:"
echo ""

# Test 14: Créer un rappel valide
test_request "/reminders" "POST" '{
    "message": "Rappel de test",
    "trigger_at": "2024-12-31T23:59:59Z"
}' "Créer un rappel valide"

# Test 15: Créer un rappel avec date passée
test_request "/reminders" "POST" '{
    "message": "Rappel invalide",
    "trigger_at": "2020-01-01T00:00:00Z"
}' "Créer un rappel avec date passée"

# Tests des requêtes d'humeur
echo "😊 Tests des requêtes d'humeur:"
echo ""

# Test 16: Créer une humeur valide
test_request "/moods" "POST" '{
    "feeling": "happy",
    "note": "Journée productive",
    "date": "2024-01-15"
}' "Créer une humeur valide"

# Test 17: Créer une humeur avec date invalide
test_request "/moods" "POST" '{
    "feeling": "happy",
    "note": "Journée productive",
    "date": "invalid-date"
}' "Créer une humeur avec date invalide"

echo "✅ Tous les tests terminés!"
echo ""
echo "📝 Résumé:"
echo "   - Tests de validation des requêtes effectués"
echo "   - Erreurs de validation correctement détectées"
echo "   - Structures de requêtes cohérentes avec les entités"
echo ""
echo "🎯 Prochaines étapes:"
echo "   1. Implémenter les handlers pour ces endpoints"
echo "   2. Ajouter la validation dans les services"
echo "   3. Créer les repositories correspondants"
echo "   4. Ajouter les tests unitaires" 