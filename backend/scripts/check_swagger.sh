#!/bin/bash

# Script pour vérifier que la documentation Swagger fonctionne
# Assurez-vous que le serveur est en cours d'exécution sur localhost:8080

echo "🔍 Vérification de la documentation Swagger"
echo "=========================================="

BASE_URL="http://localhost:8080"

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

# Test 1: Vérifier que le serveur répond
echo -e "\n${YELLOW}1. Test de connectivité du serveur${NC}"
if curl -s -o /dev/null -w "%{http_code}" "$BASE_URL" | grep -q "200\|404"; then
    print_result 0 "Serveur accessible"
else
    print_result 1 "Serveur inaccessible"
    echo "   Assurez-vous que le serveur est démarré avec: make run"
    exit 1
fi

# Test 2: Vérifier que Swagger UI est accessible
echo -e "\n${YELLOW}2. Test d'accès à Swagger UI${NC}"
SWAGGER_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/swagger/index.html")
if [ "$SWAGGER_RESPONSE" = "200" ]; then
    print_result 0 "Swagger UI accessible"
else
    print_result 1 "Swagger UI inaccessible (HTTP $SWAGGER_RESPONSE)"
fi

# Test 3: Vérifier que le fichier swagger.json est accessible
echo -e "\n${YELLOW}3. Test d'accès au fichier swagger.json${NC}"
JSON_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/swagger/doc.json")
if [ "$JSON_RESPONSE" = "200" ]; then
    print_result 0 "swagger.json accessible"
else
    print_result 1 "swagger.json inaccessible (HTTP $JSON_RESPONSE)"
fi

# Test 4: Vérifier que les routes d'authentification sont présentes
echo -e "\n${YELLOW}4. Vérification des routes d'authentification${NC}"
SWAGGER_CONTENT=$(curl -s "$BASE_URL/swagger/doc.json")

# Vérifier chaque route d'authentification
ROUTES=(
    "/api/v1/auth/register"
    "/api/v1/auth/login"
    "/api/v1/auth/refresh"
    "/api/v1/auth/me"
)

ALL_ROUTES_PRESENT=true
for route in "${ROUTES[@]}"; do
    if echo "$SWAGGER_CONTENT" | grep -q "$route"; then
        echo -e "   ${GREEN}✅ $route${NC}"
    else
        echo -e "   ${RED}❌ $route manquante${NC}"
        ALL_ROUTES_PRESENT=false
    fi
done

if [ "$ALL_ROUTES_PRESENT" = true ]; then
    print_result 0 "Toutes les routes d'authentification présentes"
else
    print_result 1 "Routes d'authentification manquantes"
fi

# Test 5: Vérifier que les tags sont présents
echo -e "\n${YELLOW}5. Vérification des tags${NC}"
if echo "$SWAGGER_CONTENT" | grep -q '"auth"'; then
    print_result 0 "Tag 'auth' présent"
else
    print_result 1 "Tag 'auth' manquant"
fi

# Test 6: Vérifier que les modèles sont présents
echo -e "\n${YELLOW}6. Vérification des modèles${NC}"
MODELS=(
    "service.RegisterRequest"
    "service.LoginRequest"
    "service.RefreshTokenRequest"
    "response.Response"
    "response.ErrorResponse"
)

ALL_MODELS_PRESENT=true
for model in "${MODELS[@]}"; do
    if echo "$SWAGGER_CONTENT" | grep -q "$model"; then
        echo -e "   ${GREEN}✅ $model${NC}"
    else
        echo -e "   ${RED}❌ $model manquant${NC}"
        ALL_MODELS_PRESENT=false
    fi
done

if [ "$ALL_MODELS_PRESENT" = true ]; then
    print_result 0 "Tous les modèles présents"
else
    print_result 1 "Modèles manquants"
fi

# Test 7: Vérifier que la sécurité BearerAuth est définie
echo -e "\n${YELLOW}7. Vérification de la sécurité${NC}"
if echo "$SWAGGER_CONTENT" | grep -q "BearerAuth"; then
    print_result 0 "Sécurité BearerAuth définie"
else
    print_result 1 "Sécurité BearerAuth manquante"
fi

echo -e "\n${YELLOW}🎉 Vérification terminée !${NC}"
echo -e "\n${GREEN}Documentation Swagger disponible sur: $BASE_URL/swagger/index.html${NC}"
echo -e "${GREEN}Serveur: $BASE_URL${NC}"

# Afficher un résumé des routes disponibles
echo -e "\n${YELLOW}📋 Routes d'authentification disponibles:${NC}"
echo -e "${GREEN}  POST /api/v1/auth/register${NC} - Inscription"
echo -e "${GREEN}  POST /api/v1/auth/login${NC} - Connexion"
echo -e "${GREEN}  POST /api/v1/auth/refresh${NC} - Rafraîchissement de token"
echo -e "${GREEN}  GET  /api/v1/auth/me${NC} - Utilisateur actuel (protégé)" 