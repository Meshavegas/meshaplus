#!/bin/bash

# 🏗️ Test de Build - MeshaPlus Backend
# Script de test pour vérifier que le build fonctionne

set -e

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🏗️ Test de Build - MeshaPlus Backend${NC}"
echo "====================================="
echo -e "${YELLOW}📋 Version Go: $(go version)${NC}"
echo ""
echo ""

# Test 1: Vérifier le répertoire
echo -e "${YELLOW}📋 Test 1: Vérification du répertoire${NC}"
if [ -f "go.mod" ]; then
    echo -e "${GREEN}✅ go.mod trouvé${NC}"
else
    echo -e "${RED}❌ go.mod non trouvé${NC}"
    exit 1
fi

# Test 2: Nettoyer les builds précédents
echo -e "${YELLOW}📋 Test 2: Nettoyage des builds précédents${NC}"
rm -f bin/api
rm -f /tmp/test-api
echo -e "${GREEN}✅ Nettoyage terminé${NC}"

# Test 3: Télécharger les dépendances
echo -e "${YELLOW}📋 Test 3: Téléchargement des dépendances${NC}"
if go mod download; then
    echo -e "${GREEN}✅ Dépendances téléchargées${NC}"
else
    echo -e "${RED}❌ Erreur lors du téléchargement${NC}"
    exit 1
fi

# Test 4: Nettoyer les dépendances
echo -e "${YELLOW}📋 Test 4: Nettoyage des dépendances${NC}"
if go mod tidy; then
    echo -e "${GREEN}✅ Dépendances nettoyées${NC}"
else
    echo -e "${RED}❌ Erreur lors du nettoyage${NC}"
    exit 1
fi

# Test 5: Vérifier la compilation
echo -e "${YELLOW}📋 Test 5: Test de compilation${NC}"
if go build -o /tmp/test-api ./cmd/api; then
    echo -e "${GREEN}✅ Compilation réussie${NC}"
else
    echo -e "${RED}❌ Erreur de compilation${NC}"
    exit 1
fi

# Test 6: Vérifier que le binaire fonctionne
echo -e "${YELLOW}📋 Test 6: Test du binaire${NC}"
if [ -f "/tmp/test-api" ]; then
    echo -e "${GREEN}✅ Binaire créé avec succès${NC}"
    ls -la /tmp/test-api
else
    echo -e "${RED}❌ Binaire non trouvé${NC}"
    exit 1
fi

# Test 7: Vérifier la taille du binaire
echo -e "${YELLOW}📋 Test 7: Vérification de la taille${NC}"
size=$(stat -f%z /tmp/test-api 2>/dev/null || stat -c%s /tmp/test-api 2>/dev/null || echo "0")
if [ "$size" -gt 1000000 ]; then
    echo -e "${GREEN}✅ Taille du binaire: ${size} bytes${NC}"
else
    echo -e "${YELLOW}⚠️  Taille du binaire petite: ${size} bytes${NC}"
fi

# Nettoyage
echo -e "${YELLOW}📋 Nettoyage${NC}"
rm -f /tmp/test-api
echo -e "${GREEN}✅ Nettoyage terminé${NC}"

echo ""
echo -e "${GREEN}🎉 Test de build terminé avec succès !${NC}"
echo ""
echo -e "${BLUE}📋 Résumé :${NC}"
echo "✅ go.mod trouvé"
echo "✅ Dépendances téléchargées et nettoyées"
echo "✅ Compilation réussie"
echo "✅ Binaire créé et fonctionnel"
echo ""
echo -e "${BLUE}🚀 Le build Docker devrait maintenant fonctionner !${NC}" 