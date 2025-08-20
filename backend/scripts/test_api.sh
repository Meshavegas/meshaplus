#!/bin/bash

# Script de test pour vérifier l'accessibilité des endpoints API
# Usage: ./scripts/test_api.sh

set -e

echo "🧪 Test de l'API MeshaPlus"
echo "=========================="

# Configuration
BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v1"

# Couleurs pour l'affichage
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour tester un endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local description=$3
    local expected_status=${4:-200}
    
    echo -n "Testing $method $endpoint... "
    
    # Effectuer la requête
    if [ "$method" = "GET" ]; then
        response=$(curl -s -o /dev/null -w "%{http_code}" "$endpoint" 2>/dev/null || echo "000")
    else
        response=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$endpoint" 2>/dev/null || echo "000")
    fi
    
    # Vérifier le statut
    if [ "$response" = "$expected_status" ]; then
        echo -e "${GREEN}✅ OK ($response)${NC}"
    else
        echo -e "${RED}❌ FAILED ($response, expected $expected_status)${NC}"
    fi
}

# Fonction pour tester un endpoint avec authentification (401 attendu)
test_auth_endpoint() {
    local method=$1
    local endpoint=$2
    local description=$3
    
    echo -n "Testing $method $endpoint (auth required)... "
    
    # Effectuer la requête
    if [ "$method" = "GET" ]; then
        response=$(curl -s -o /dev/null -w "%{http_code}" "$endpoint" 2>/dev/null || echo "000")
    else
        response=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$endpoint" 2>/dev/null || echo "000")
    fi
    
    # Vérifier le statut (401 = non autorisé, ce qui est normal sans token)
    if [ "$response" = "401" ]; then
        echo -e "${GREEN}✅ OK (401 - Auth required)${NC}"
    elif [ "$response" = "405" ]; then
        echo -e "${YELLOW}⚠️  Method Not Allowed (405)${NC}"
    else
        echo -e "${RED}❌ UNEXPECTED ($response)${NC}"
    fi
}

echo ""
echo "📋 Test des endpoints publics"
echo "----------------------------"

# Test Swagger
test_endpoint "GET" "$BASE_URL/swagger/index.html" "Swagger UI"
test_endpoint "GET" "$BASE_URL/swagger/doc.json" "Swagger JSON"

echo ""
echo "🔐 Test des endpoints d'authentification"
echo "---------------------------------------"

# Test des endpoints d'authentification (devraient retourner 400 ou 422 pour des données invalides)
test_endpoint "POST" "$API_BASE/auth/register" "Register (invalid data)" "400"
test_endpoint "POST" "$API_BASE/auth/login" "Login (invalid data)" "400"

echo ""
echo "🔒 Test des endpoints protégés (sans authentification)"
echo "-----------------------------------------------------"

# Test des endpoints protégés (devraient retourner 401)
test_auth_endpoint "GET" "$API_BASE/tasks" "Get tasks"
test_auth_endpoint "POST" "$API_BASE/tasks" "Create task"
test_auth_endpoint "GET" "$API_BASE/transactions" "Get transactions"
test_auth_endpoint "POST" "$API_BASE/transactions" "Create transaction"
test_auth_endpoint "GET" "$API_BASE/accounts" "Get accounts"
test_auth_endpoint "POST" "$API_BASE/accounts" "Create account"
test_auth_endpoint "GET" "$API_BASE/budgets" "Get budgets"
test_auth_endpoint "POST" "$API_BASE/budgets" "Create budget"
test_auth_endpoint "GET" "$API_BASE/saving-goals" "Get saving goals"
test_auth_endpoint "POST" "$API_BASE/saving-goals" "Create saving goal"

echo ""
echo "📊 Test des endpoints de statistiques"
echo "------------------------------------"

test_auth_endpoint "GET" "$API_BASE/tasks/stats" "Task stats"
test_auth_endpoint "GET" "$API_BASE/transactions/stats" "Transaction stats"
test_auth_endpoint "GET" "$API_BASE/budgets/stats" "Budget stats"

echo ""
echo "🎯 Test des endpoints avec paramètres"
echo "------------------------------------"

# Test des endpoints avec ID (devraient retourner 401)
test_auth_endpoint "GET" "$API_BASE/tasks/123e4567-e89b-12d3-a456-426614174000" "Get task by ID"
test_auth_endpoint "PUT" "$API_BASE/tasks/123e4567-e89b-12d3-a456-426614174000" "Update task"
test_auth_endpoint "DELETE" "$API_BASE/tasks/123e4567-e89b-12d3-a456-426614174000" "Delete task"

echo ""
echo "✅ Tests terminés !"
echo ""
echo "📖 Documentation Swagger disponible sur :"
echo "   $BASE_URL/swagger/index.html"
echo ""
echo "🚀 Pour tester avec authentification, utilisez :"
echo "   1. Créez un compte : POST $API_BASE/auth/register"
echo "   2. Connectez-vous : POST $API_BASE/auth/login"
echo "   3. Utilisez le token JWT dans l'en-tête Authorization" 