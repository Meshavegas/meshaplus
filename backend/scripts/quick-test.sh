#!/bin/bash

# 🚀 Test Rapide - MeshaPlus Backend
# Script de test rapide pour vérifier les commandes de base

set -e

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}🚀 Test Rapide - MeshaPlus Backend${NC}"
echo "====================================="
echo ""

# Test 1: Vérifier le répertoire
echo -e "${YELLOW}📋 Test 1: Vérification du répertoire${NC}"
if [ -f "go.mod" ]; then
    echo -e "${GREEN}✅ go.mod trouvé${NC}"
else
    echo -e "${RED}❌ go.mod non trouvé${NC}"
    exit 1
fi

# Test 2: Télécharger les dépendances
echo -e "${YELLOW}📋 Test 2: Téléchargement des dépendances${NC}"
if go mod download; then
    echo -e "${GREEN}✅ Dépendances téléchargées${NC}"
else
    echo -e "${RED}❌ Erreur lors du téléchargement${NC}"
    exit 1
fi

# Test 3: Nettoyer les dépendances
echo -e "${YELLOW}📋 Test 3: Nettoyage des dépendances${NC}"
if go mod tidy; then
    echo -e "${GREEN}✅ Dépendances nettoyées${NC}"
else
    echo -e "${RED}❌ Erreur lors du nettoyage${NC}"
    exit 1
fi

# Test 4: Vérifier la compilation
echo -e "${YELLOW}📋 Test 4: Test de compilation${NC}"
if go build -o /tmp/test-api ./cmd/api; then
    echo -e "${GREEN}✅ Compilation réussie${NC}"
    rm -f /tmp/test-api
else
    echo -e "${RED}❌ Erreur de compilation${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 Tous les tests sont passés !${NC}"
echo -e "${GREEN}✅ Le workflow GitHub Actions devrait maintenant fonctionner${NC}" 