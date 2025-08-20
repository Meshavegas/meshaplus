#!/bin/bash

# Script de test pour les routes Expense
# Usage: ./scripts/test_expense.sh

set -e

echo "🧪 Test des routes Expense"
echo "=========================="

# Configuration
BASE_URL="http://localhost:8080"
EMAIL="expense-test@example.com"
PASSWORD="password123"

echo "📝 1. Création d'un utilisateur de test..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test User\",\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

if echo "$REGISTER_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Utilisateur créé avec succès"
else
    echo "❌ Erreur création utilisateur: $REGISTER_RESPONSE"
    exit 1
fi

# Extraction du token
TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
echo "🔑 Token extrait: ${TOKEN:0:20}..."

echo ""
echo "📝 2. Test création de dépenses..."

# Création de plusieurs dépenses
EXPENSE1_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/expenses" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount":50.00,"description":"Achat alimentaire","date":"2025-07-26T00:00:00Z"}')

if echo "$EXPENSE1_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépense 1 créée avec succès"
    EXPENSE1_ID=$(echo "$EXPENSE1_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
else
    echo "❌ Erreur création dépense 1: $EXPENSE1_RESPONSE"
    exit 1
fi

EXPENSE2_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/expenses" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount":25.50,"description":"Transport","date":"2025-07-26T00:00:00Z"}')

if echo "$EXPENSE2_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépense 2 créée avec succès"
    EXPENSE2_ID=$(echo "$EXPENSE2_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
else
    echo "❌ Erreur création dépense 2: $EXPENSE2_RESPONSE"
    exit 1
fi

EXPENSE3_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/expenses" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount":100.00,"description":"Loisirs","date":"2025-07-25T00:00:00Z"}')

if echo "$EXPENSE3_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépense 3 créée avec succès"
    EXPENSE3_ID=$(echo "$EXPENSE3_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
else
    echo "❌ Erreur création dépense 3: $EXPENSE3_RESPONSE"
    exit 1
fi

echo ""
echo "📝 3. Test récupération de toutes les dépenses..."
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/expenses" \
  -H "Authorization: Bearer $TOKEN")

if echo "$LIST_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Liste des dépenses récupérée avec succès"
    EXPENSE_COUNT=$(echo "$LIST_RESPONSE" | grep -o '"id":"[^"]*"' | wc -l)
    echo "📊 Nombre de dépenses: $EXPENSE_COUNT"
else
    echo "❌ Erreur récupération liste: $LIST_RESPONSE"
    exit 1
fi

echo ""
echo "📝 4. Test récupération d'une dépense par ID..."
GET_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/expenses/$EXPENSE1_ID" \
  -H "Authorization: Bearer $TOKEN")

if echo "$GET_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépense récupérée par ID avec succès"
    EXPENSE_AMOUNT=$(echo "$GET_RESPONSE" | grep -o '"amount":[0-9.]*' | cut -d':' -f2)
    echo "💰 Montant: $EXPENSE_AMOUNT"
else
    echo "❌ Erreur récupération par ID: $GET_RESPONSE"
    exit 1
fi

echo ""
echo "📝 5. Test calcul du total des dépenses..."
TOTAL_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/expenses/total" \
  -H "Authorization: Bearer $TOKEN")

if echo "$TOTAL_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Total calculé avec succès"
    TOTAL_AMOUNT=$(echo "$TOTAL_RESPONSE" | grep -o '"data":[0-9.]*' | cut -d':' -f2)
    echo "💰 Total des dépenses: $TOTAL_AMOUNT"
else
    echo "❌ Erreur calcul total: $TOTAL_RESPONSE"
    exit 1
fi

echo ""
echo "📝 6. Test récupération par période..."
DATE_RANGE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/expenses/date-range?start_date=2025-07-25&end_date=2025-07-26" \
  -H "Authorization: Bearer $TOKEN")

if echo "$DATE_RANGE_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépenses par période récupérées avec succès"
    RANGE_COUNT=$(echo "$DATE_RANGE_RESPONSE" | grep -o '"id":"[^"]*"' | wc -l)
    echo "📊 Nombre de dépenses sur la période: $RANGE_COUNT"
else
    echo "❌ Erreur récupération par période: $DATE_RANGE_RESPONSE"
    exit 1
fi

echo ""
echo "📝 7. Test mise à jour d'une dépense..."
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/v1/expenses/$EXPENSE1_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount":60.00,"description":"Achat alimentaire mis à jour"}')

if echo "$UPDATE_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépense mise à jour avec succès"
    NEW_AMOUNT=$(echo "$UPDATE_RESPONSE" | grep -o '"amount":[0-9.]*' | cut -d':' -f2)
    echo "💰 Nouveau montant: $NEW_AMOUNT"
else
    echo "❌ Erreur mise à jour: $UPDATE_RESPONSE"
    exit 1
fi

echo ""
echo "📝 8. Test suppression des dépenses..."

# Suppression de la première dépense
DELETE1_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/expenses/$EXPENSE1_ID" \
  -H "Authorization: Bearer $TOKEN")

if echo "$DELETE1_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépense 1 supprimée avec succès"
else
    echo "❌ Erreur suppression dépense 1: $DELETE1_RESPONSE"
    exit 1
fi

# Suppression de la deuxième dépense
DELETE2_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/expenses/$EXPENSE2_ID" \
  -H "Authorization: Bearer $TOKEN")

if echo "$DELETE2_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépense 2 supprimée avec succès"
else
    echo "❌ Erreur suppression dépense 2: $DELETE2_RESPONSE"
    exit 1
fi

# Suppression de la troisième dépense
DELETE3_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/v1/expenses/$EXPENSE3_ID" \
  -H "Authorization: Bearer $TOKEN")

if echo "$DELETE3_RESPONSE" | grep -q '"success":true'; then
    echo "✅ Dépense 3 supprimée avec succès"
else
    echo "❌ Erreur suppression dépense 3: $DELETE3_RESPONSE"
    exit 1
fi

echo ""
echo "📝 9. Vérification que toutes les dépenses sont supprimées..."
FINAL_LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/expenses" \
  -H "Authorization: Bearer $TOKEN")

if echo "$FINAL_LIST_RESPONSE" | grep -q '"data":null'; then
    echo "✅ Toutes les dépenses ont été supprimées"
else
    echo "❌ Erreur: il reste des dépenses après suppression"
    exit 1
fi

echo ""
echo "🎉 Tous les tests des routes Expense sont passés avec succès !"
echo "==============================================================" 