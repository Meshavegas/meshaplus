#!/bin/bash

# Script de test pour les routes de catégories
# Assurez-vous que le serveur est démarré sur localhost:8080

BASE_URL="http://localhost:8080/api/v1"
TOKEN="" # À remplir avec un token valide

echo "🧪 Test des routes de catégories"
echo "================================="

# Test 1: Récupérer les catégories de type expense
echo ""
echo "📝 Test 1: Récupérer les catégories de type 'expense'"
curl -X GET "$BASE_URL/categories?categoryType=expense" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# Test 2: Récupérer les catégories de type revenue
echo ""
echo "📝 Test 2: Récupérer les catégories de type 'revenue'"
curl -X GET "$BASE_URL/categories?categoryType=revenue" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# Test 3: Récupérer les catégories de type task
echo ""
echo "📝 Test 3: Récupérer les catégories de type 'task'"
curl -X GET "$BASE_URL/categories?categoryType=task" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

# Test 4: Catégoriser un item
echo ""
echo "📝 Test 4: Catégoriser 'Pizza Margherita' comme expense"
curl -X POST "$BASE_URL/categories/categorize" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "item": "Pizza Margherita",
    "categoryType": "expense"
  }' | jq '.'

echo ""
echo "✅ Tests terminés" 