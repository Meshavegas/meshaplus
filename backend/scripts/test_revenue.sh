#!/bin/bash

# Script de test pour les routes Revenue
# Ce script teste toutes les fonctionnalités CRUD et avancées des revenus

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

echo -e "${BLUE}=== Tests des routes Revenue ===${NC}"

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
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
        \"first_name\": \"Test\",
        \"last_name\": \"User\"
    }")

echo "$REGISTER_RESPONSE" | jq .
print_result $? "Inscription utilisateur"

# 2. Connexion et récupération du token
echo -e "\n${YELLOW}2. Connexion et récupération du token...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')
echo "$LOGIN_RESPONSE" | jq .
print_result $? "Connexion utilisateur"

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}Erreur: Token non récupéré${NC}"
    exit 1
fi

echo -e "${GREEN}Token récupéré: ${TOKEN:0:20}...${NC}"

# 3. Création d'un revenu ponctuel
echo -e "\n${YELLOW}3. Création d'un revenu ponctuel...${NC}"
REVENUE1_RESPONSE=$(curl -s -X POST "$BASE_URL/revenues" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"amount\": 2500.00,
        \"source\": \"Salaire principal\",
        \"date\": \"2025-01-15T00:00:00Z\",
        \"notes\": \"Salaire de janvier\",
        \"type\": \"one_time\"
    }")

echo "$REVENUE1_RESPONSE" | jq .
print_result $? "Création revenu ponctuel"

REVENUE1_ID=$(echo "$REVENUE1_RESPONSE" | jq -r '.id')

# 4. Création d'un revenu récurrent
echo -e "\n${YELLOW}4. Création d'un revenu récurrent...${NC}"
REVENUE2_RESPONSE=$(curl -s -X POST "$BASE_URL/revenues" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"amount\": 3000.00,
        \"source\": \"Salaire mensuel\",
        \"date\": \"2025-01-25T00:00:00Z\",
        \"notes\": \"Salaire mensuel récurrent\",
        \"type\": \"recurring\",
        \"frequency\": \"monthly\",
        \"payment_day\": 25,
        \"start_date\": \"2025-01-01T00:00:00Z\",
        \"is_active\": true
    }")

echo "$REVENUE2_RESPONSE" | jq .
print_result $? "Création revenu récurrent"

REVENUE2_ID=$(echo "$REVENUE2_RESPONSE" | jq -r '.id')

# 5. Création d'un bonus
echo -e "\n${YELLOW}5. Création d'un bonus...${NC}"
REVENUE3_RESPONSE=$(curl -s -X POST "$BASE_URL/revenues" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"amount\": 500.00,
        \"source\": \"Bonus de fin d'année\",
        \"date\": \"2025-12-15T00:00:00Z\",
        \"notes\": \"Bonus exceptionnel\",
        \"type\": \"bonus\"
    }")

echo "$REVENUE3_RESPONSE" | jq .
print_result $? "Création bonus"

REVENUE3_ID=$(echo "$REVENUE3_RESPONSE" | jq -r '.id')

# 6. Récupération de tous les revenus
echo -e "\n${YELLOW}6. Récupération de tous les revenus...${NC}"
GET_ALL_RESPONSE=$(curl -s -X GET "$BASE_URL/revenues" \
    -H "Authorization: Bearer $TOKEN")

echo "$GET_ALL_RESPONSE" | jq .
print_result $? "Récupération tous les revenus"

# 7. Récupération d'un revenu par ID
echo -e "\n${YELLOW}7. Récupération d'un revenu par ID...${NC}"
GET_ONE_RESPONSE=$(curl -s -X GET "$BASE_URL/revenues/$REVENUE1_ID" \
    -H "Authorization: Bearer $TOKEN")

echo "$GET_ONE_RESPONSE" | jq .
print_result $? "Récupération revenu par ID"

# 8. Calcul du total des revenus
echo -e "\n${YELLOW}8. Calcul du total des revenus...${NC}"
TOTAL_RESPONSE=$(curl -s -X GET "$BASE_URL/revenues/total" \
    -H "Authorization: Bearer $TOKEN")

echo "$TOTAL_RESPONSE" | jq .
print_result $? "Calcul total revenus"

# 9. Récupération des revenus par type
echo -e "\n${YELLOW}9. Récupération des revenus par type (recurring)...${NC}"
TYPE_RESPONSE=$(curl -s -X GET "$BASE_URL/revenues/type?type=recurring" \
    -H "Authorization: Bearer $TOKEN")

echo "$TYPE_RESPONSE" | jq .
print_result $? "Récupération revenus par type"

# 10. Récupération des revenus actifs
echo -e "\n${YELLOW}10. Récupération des revenus actifs...${NC}"
ACTIVE_RESPONSE=$(curl -s -X GET "$BASE_URL/revenues/active" \
    -H "Authorization: Bearer $TOKEN")

echo "$ACTIVE_RESPONSE" | jq .
print_result $? "Récupération revenus actifs"

# 11. Récupération des revenus par période
echo -e "\n${YELLOW}11. Récupération des revenus par période...${NC}"
DATE_RANGE_RESPONSE=$(curl -s -X GET "$BASE_URL/revenues/date-range?start_date=2025-01-01T00:00:00Z&end_date=2025-12-31T23:59:59Z" \
    -H "Authorization: Bearer $TOKEN")

echo "$DATE_RANGE_RESPONSE" | jq .
print_result $? "Récupération revenus par période"

# 12. Mise à jour d'un revenu
echo -e "\n${YELLOW}12. Mise à jour d'un revenu...${NC}"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/revenues/$REVENUE1_ID" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
        \"amount\": 2600.00,
        \"notes\": \"Salaire de janvier mis à jour\"
    }")

echo "$UPDATE_RESPONSE" | jq .
print_result $? "Mise à jour revenu"

# 13. Vérification de la mise à jour
echo -e "\n${YELLOW}13. Vérification de la mise à jour...${NC}"
VERIFY_UPDATE_RESPONSE=$(curl -s -X GET "$BASE_URL/revenues/$REVENUE1_ID" \
    -H "Authorization: Bearer $TOKEN")

echo "$VERIFY_UPDATE_RESPONSE" | jq .
print_result $? "Vérification mise à jour"

# 14. Suppression des revenus
echo -e "\n${YELLOW}14. Suppression des revenus...${NC}"

# Suppression du premier revenu
curl -s -X DELETE "$BASE_URL/revenues/$REVENUE1_ID" \
    -H "Authorization: Bearer $TOKEN"
print_result $? "Suppression revenu 1"

# Suppression du deuxième revenu
curl -s -X DELETE "$BASE_URL/revenues/$REVENUE2_ID" \
    -H "Authorization: Bearer $TOKEN"
print_result $? "Suppression revenu 2"

# Suppression du troisième revenu
curl -s -X DELETE "$BASE_URL/revenues/$REVENUE3_ID" \
    -H "Authorization: Bearer $TOKEN"
print_result $? "Suppression revenu 3"

# 15. Vérification que la liste est vide
echo -e "\n${YELLOW}15. Vérification que la liste est vide...${NC}"
FINAL_CHECK_RESPONSE=$(curl -s -X GET "$BASE_URL/revenues" \
    -H "Authorization: Bearer $TOKEN")

echo "$FINAL_CHECK_RESPONSE" | jq .
REVENUE_COUNT=$(echo "$FINAL_CHECK_RESPONSE" | jq 'length')

if [ "$REVENUE_COUNT" -eq 0 ]; then
    print_result 0 "Vérification liste vide"
else
    print_result 1 "Vérification liste vide (attendu: 0, obtenu: $REVENUE_COUNT)"
fi

echo -e "\n${GREEN}=== Tous les tests Revenue sont passés avec succès ! ===${NC}" 