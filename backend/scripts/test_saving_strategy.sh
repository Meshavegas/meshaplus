#!/bin/bash

# Script de test pour les routes SavingStrategy
# Ce script teste toutes les fonctionnalités CRUD et spécialisées des stratégies d'épargne

set -e

# Configuration
BASE_URL="http://localhost:8080/api/v1"
EMAIL="test@example.com"
PASSWORD="password123"

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Tests des routes SavingStrategy${NC}"
echo "=================================="

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

# 1. Inscription d'un utilisateur
echo -e "\n${YELLOW}1. Inscription d'un utilisateur...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\",
        \"name\": \"Test User\"
    }")

echo "Réponse d'inscription: $REGISTER_RESPONSE"

# 2. Connexion de l'utilisateur
echo -e "\n${YELLOW}2. Connexion de l'utilisateur...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

echo "Réponse de connexion: $LOGIN_RESPONSE"

# Extraction du token
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Impossible d'extraire le token${NC}"
    exit 1
fi

echo -e "${GREEN}Token extrait: ${TOKEN:0:20}...${NC}"

# 3. Création d'une stratégie d'épargne par pourcentage
echo -e "\n${YELLOW}3. Création d'une stratégie d'épargne par pourcentage...${NC}"
PERCENTAGE_STRATEGY_RESPONSE=$(curl -s -X POST "$BASE_URL/saving-strategies" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"strategy_name\": \"Épargne 20% du salaire\",
        \"type\": \"percentage\",
        \"amount\": 20.0,
        \"frequency\": \"monthly\"
    }")

echo "Réponse création stratégie pourcentage: $PERCENTAGE_STRATEGY_RESPONSE"

PERCENTAGE_STRATEGY_ID=$(echo $PERCENTAGE_STRATEGY_RESPONSE | jq -r '.data.id')
if [ "$PERCENTAGE_STRATEGY_ID" = "null" ] || [ -z "$PERCENTAGE_STRATEGY_ID" ]; then
    echo -e "${RED}❌ Impossible d'extraire l'ID de la stratégie pourcentage${NC}"
    exit 1
fi

print_result $? "Stratégie d'épargne par pourcentage créée"

# 4. Création d'une stratégie d'épargne fixe
echo -e "\n${YELLOW}4. Création d'une stratégie d'épargne fixe...${NC}"
FIXED_STRATEGY_RESPONSE=$(curl -s -X POST "$BASE_URL/saving-strategies" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"strategy_name\": \"Épargne fixe 500€\",
        \"type\": \"fixed\",
        \"amount\": 500.0,
        \"frequency\": \"monthly\"
    }")

echo "Réponse création stratégie fixe: $FIXED_STRATEGY_RESPONSE"

FIXED_STRATEGY_ID=$(echo $FIXED_STRATEGY_RESPONSE | jq -r '.data.id')
if [ "$FIXED_STRATEGY_ID" = "null" ] || [ -z "$FIXED_STRATEGY_ID" ]; then
    echo -e "${RED}❌ Impossible d'extraire l'ID de la stratégie fixe${NC}"
    exit 1
fi

print_result $? "Stratégie d'épargne fixe créée"

# 5. Création d'une stratégie d'épargne basée sur un objectif
echo -e "\n${YELLOW}5. Création d'une stratégie d'épargne basée sur un objectif...${NC}"
GOAL_STRATEGY_RESPONSE=$(curl -s -X POST "$BASE_URL/saving-strategies" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"strategy_name\": \"Épargne pour vacances\",
        \"type\": \"goal_based\",
        \"amount\": 1000.0,
        \"frequency\": \"monthly\",
        \"target_goal_id\": \"123e4567-e89b-12d3-a456-426614174000\"
    }")

echo "Réponse création stratégie objectif: $GOAL_STRATEGY_RESPONSE"

GOAL_STRATEGY_ID=$(echo $GOAL_STRATEGY_RESPONSE | jq -r '.data.id')
if [ "$GOAL_STRATEGY_ID" = "null" ] || [ -z "$GOAL_STRATEGY_ID" ]; then
    echo -e "${RED}❌ Impossible d'extraire l'ID de la stratégie objectif${NC}"
    exit 1
fi

print_result $? "Stratégie d'épargne basée sur objectif créée"

# 6. Récupération de toutes les stratégies d'épargne
echo -e "\n${YELLOW}6. Récupération de toutes les stratégies d'épargne...${NC}"
GET_ALL_RESPONSE=$(curl -s -X GET "$BASE_URL/saving-strategies" \
    -H "Authorization: Bearer $TOKEN")

echo "Réponse récupération toutes les stratégies: $GET_ALL_RESPONSE"

STRATEGIES_COUNT=$(echo $GET_ALL_RESPONSE | jq '.data | length')
if [ "$STRATEGIES_COUNT" -ge 3 ]; then
    print_result 0 "Toutes les stratégies d'épargne récupérées ($STRATEGIES_COUNT stratégies)"
else
    print_result 1 "Nombre de stratégies incorrect: $STRATEGIES_COUNT"
fi

# 7. Récupération d'une stratégie d'épargne par ID
echo -e "\n${YELLOW}7. Récupération d'une stratégie d'épargne par ID...${NC}"
GET_ONE_RESPONSE=$(curl -s -X GET "$BASE_URL/saving-strategies/$PERCENTAGE_STRATEGY_ID" \
    -H "Authorization: Bearer $TOKEN")

echo "Réponse récupération stratégie par ID: $GET_ONE_RESPONSE"

RETRIEVED_ID=$(echo $GET_ONE_RESPONSE | jq -r '.data.id')
if [ "$RETRIEVED_ID" = "$PERCENTAGE_STRATEGY_ID" ]; then
    print_result 0 "Stratégie d'épargne récupérée par ID"
else
    print_result 1 "ID de stratégie incorrect: $RETRIEVED_ID"
fi

# 8. Récupération des stratégies d'épargne par type (pourcentage)
echo -e "\n${YELLOW}8. Récupération des stratégies d'épargne par type (pourcentage)...${NC}"
GET_BY_TYPE_RESPONSE=$(curl -s -X GET "$BASE_URL/saving-strategies/type?type=percentage" \
    -H "Authorization: Bearer $TOKEN")

echo "Réponse récupération par type: $GET_BY_TYPE_RESPONSE"

PERCENTAGE_COUNT=$(echo $GET_BY_TYPE_RESPONSE | jq '.data | length')
if [ "$PERCENTAGE_COUNT" -eq 1 ]; then
    print_result 0 "Stratégies d'épargne par type récupérées ($PERCENTAGE_COUNT stratégie pourcentage)"
else
    print_result 1 "Nombre de stratégies pourcentage incorrect: $PERCENTAGE_COUNT"
fi

# 9. Calcul du montant d'épargne pour une stratégie
echo -e "\n${YELLOW}9. Calcul du montant d'épargne pour une stratégie...${NC}"
CALCULATE_RESPONSE=$(curl -s -X GET "$BASE_URL/saving-strategies/$PERCENTAGE_STRATEGY_ID/calculate?base_amount=3000" \
    -H "Authorization: Bearer $TOKEN")

echo "Réponse calcul montant d'épargne: $CALCULATE_RESPONSE"

CALCULATED_AMOUNT=$(echo $CALCULATE_RESPONSE | jq -r '.data.saving_amount')
if [ "$CALCULATED_AMOUNT" = "600" ]; then
    print_result 0 "Montant d'épargne calculé correctement: $CALCULATED_AMOUNT€"
else
    print_result 1 "Montant d'épargne incorrect: $CALCULATED_AMOUNT€"
fi

# 10. Mise à jour d'une stratégie d'épargne
echo -e "\n${YELLOW}10. Mise à jour d'une stratégie d'épargne...${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/saving-strategies/$FIXED_STRATEGY_ID" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"strategy_name\": \"Épargne fixe 600€ (mise à jour)\",
        \"amount\": 600.0
    }")

echo "Réponse mise à jour: $UPDATE_RESPONSE"

UPDATED_AMOUNT=$(echo $UPDATE_RESPONSE | jq -r '.data.amount')
if [ "$UPDATED_AMOUNT" = "600" ]; then
    print_result 0 "Stratégie d'épargne mise à jour correctement"
else
    print_result 1 "Montant mis à jour incorrect: $UPDATED_AMOUNT"
fi

# 11. Suppression d'une stratégie d'épargne
echo -e "\n${YELLOW}11. Suppression d'une stratégie d'épargne...${NC}"
DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/saving-strategies/$GOAL_STRATEGY_ID" \
    -H "Authorization: Bearer $TOKEN")

echo "Réponse suppression: $DELETE_RESPONSE"

if echo $DELETE_RESPONSE | jq -e '.success' > /dev/null; then
    print_result 0 "Stratégie d'épargne supprimée avec succès"
else
    print_result 1 "Erreur lors de la suppression"
fi

# 12. Vérification que la stratégie a bien été supprimée
echo -e "\n${YELLOW}12. Vérification de la suppression...${NC}"
GET_DELETED_RESPONSE=$(curl -s -X GET "$BASE_URL/saving-strategies/$GOAL_STRATEGY_ID" \
    -H "Authorization: Bearer $TOKEN")

echo "Réponse récupération stratégie supprimée: $GET_DELETED_RESPONSE"

if echo $GET_DELETED_RESPONSE | jq -e '.error' > /dev/null; then
    print_result 0 "Stratégie d'épargne bien supprimée (404 retourné)"
else
    print_result 1 "Stratégie d'épargne toujours accessible après suppression"
fi

# 13. Vérification finale du nombre de stratégies
echo -e "\n${YELLOW}13. Vérification finale du nombre de stratégies...${NC}"
FINAL_COUNT_RESPONSE=$(curl -s -X GET "$BASE_URL/saving-strategies" \
    -H "Authorization: Bearer $TOKEN")

FINAL_COUNT=$(echo $FINAL_COUNT_RESPONSE | jq '.data | length')
if [ "$FINAL_COUNT" -eq 2 ]; then
    print_result 0 "Nombre final de stratégies correct: $FINAL_COUNT"
else
    print_result 1 "Nombre final de stratégies incorrect: $FINAL_COUNT"
fi

echo -e "\n${GREEN}🎉 Tous les tests des routes SavingStrategy sont passés avec succès !${NC}"
echo "==================================" 