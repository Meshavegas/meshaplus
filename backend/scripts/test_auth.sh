#!/bin/bash

# Script de test pour l'authentification JWT
# Assurez-vous que le serveur est en cours d'exécution sur localhost:8080

echo "🧪 Tests d'authentification JWT"
echo "================================"

BASE_URL="http://localhost:8080/api/v1/auth"

# Couleurs pour l'affichage
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Test 1: Inscription d'un utilisateur
echo -e "\n${YELLOW}1. Test d'inscription${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Test User",
        "email": "test@example.com",
        "password": "password123"
    }')

if echo "$REGISTER_RESPONSE" | grep -q '"success":true'; then
    print_result 0 "Inscription réussie"
    # Extraire le token d'accès
    ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)
    echo "   Token d'accès extrait"
else
    print_result 1 "Échec de l'inscription"
    echo "   Réponse: $REGISTER_RESPONSE"
fi

# Test 2: Connexion
echo -e "\n${YELLOW}2. Test de connexion${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "password": "password123"
    }')

if echo "$LOGIN_RESPONSE" | grep -q '"success":true'; then
    print_result 0 "Connexion réussie"
    # Extraire le nouveau token d'accès
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
    REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)
    echo "   Token d'accès extrait"
else
    print_result 1 "Échec de la connexion"
    echo "   Réponse: $LOGIN_RESPONSE"
fi

# Test 3: Accès à une ressource protégée
echo -e "\n${YELLOW}3. Test d'accès protégé${NC}"
if [ ! -z "$ACCESS_TOKEN" ]; then
    ME_RESPONSE=$(curl -s -X GET "$BASE_URL/me" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if echo "$ME_RESPONSE" | grep -q '"success":true'; then
        print_result 0 "Accès protégé réussi"
    else
        print_result 1 "Échec de l'accès protégé"
        echo "   Réponse: $ME_RESPONSE"
    fi
else
    print_result 1 "Impossible de tester l'accès protégé (pas de token)"
fi

# Test 4: Accès sans token (doit échouer)
echo -e "\n${YELLOW}4. Test d'accès sans token${NC}"
NO_TOKEN_RESPONSE=$(curl -s -X GET "$BASE_URL/me")

if echo "$NO_TOKEN_RESPONSE" | grep -q "Authorization header required"; then
    print_result 0 "Protection correcte (accès refusé sans token)"
else
    print_result 1 "Problème de protection (accès autorisé sans token)"
    echo "   Réponse: $NO_TOKEN_RESPONSE"
fi

# Test 5: Rafraîchissement de token
echo -e "\n${YELLOW}5. Test de rafraîchissement de token${NC}"
if [ ! -z "$REFRESH_TOKEN" ]; then
    REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/refresh" \
        -H "Content-Type: application/json" \
        -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}")
    
    if echo "$REFRESH_RESPONSE" | grep -q '"success":true'; then
        print_result 0 "Rafraîchissement de token réussi"
        # Extraire le nouveau token d'accès
        NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        echo "   Nouveau token d'accès généré"
    else
        print_result 1 "Échec du rafraîchissement de token"
        echo "   Réponse: $REFRESH_RESPONSE"
    fi
else
    print_result 1 "Impossible de tester le rafraîchissement (pas de refresh token)"
fi

# Test 6: Connexion avec mauvais mot de passe
echo -e "\n${YELLOW}6. Test de connexion avec mauvais mot de passe${NC}"
WRONG_PASSWORD_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "password": "wrongpassword"
    }')

if echo "$WRONG_PASSWORD_RESPONSE" | grep -q '"success":false'; then
    print_result 0 "Protection correcte (mauvais mot de passe rejeté)"
else
    print_result 1 "Problème de sécurité (mauvais mot de passe accepté)"
    echo "   Réponse: $WRONG_PASSWORD_RESPONSE"
fi

# Test 7: Inscription avec email existant
echo -e "\n${YELLOW}7. Test d'inscription avec email existant${NC}"
DUPLICATE_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Another User",
        "email": "test@example.com",
        "password": "password123"
    }')

if echo "$DUPLICATE_RESPONSE" | grep -q '"success":false'; then
    print_result 0 "Protection correcte (email dupliqué rejeté)"
else
    print_result 1 "Problème de contrainte (email dupliqué accepté)"
    echo "   Réponse: $DUPLICATE_RESPONSE"
fi

echo -e "\n${YELLOW}🎉 Tests terminés !${NC}"
echo -e "\n${GREEN}Documentation Swagger disponible sur: http://localhost:8080/swagger/index.html${NC}"
echo -e "${GREEN}Base de données: meshaplus${NC}"
echo -e "${GREEN}Serveur: http://localhost:8080${NC}" 