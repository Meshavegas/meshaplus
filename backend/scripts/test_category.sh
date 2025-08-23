#!/bin/bash

# Script de test pour la route de catégorisation
# Assurez-vous que le serveur est démarré sur localhost:8080

BASE_URL="http://localhost:8080/api/v1"
TOKEN="" # À remplir avec un token valide

echo "🧪 Test de la route de catégorisation"
echo "======================================"

# Test 1: Catégoriser un item de type expense
echo ""
echo "📝 Test 1: Catégoriser 'Pizza Margherita' comme expense"
curl -X POST "$BASE_URL/categories/categorize" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "item": "Pizza Margherita",
    "categoryType": "expense"
  }' | jq '.'

# Test 2: Catégoriser un item de type revenue
echo ""
echo "📝 Test 2: Catégoriser 'Salaire mensuel' comme revenue"
curl -X POST "$BASE_URL/categories/categorize" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "item": "Salaire mensuel",
    "categoryType": "revenue"
  }' | jq '.'

# Test 3: Catégoriser un item de type task
echo ""
echo "📝 Test 3: Catégoriser 'Réunion équipe' comme task"
curl -X POST "$BASE_URL/categories/categorize" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "item": "Réunion équipe",
    "categoryType": "task"
  }' | jq '.'

echo ""
echo "✅ Tests terminés" 