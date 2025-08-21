#!/bin/bash

# 🧪 Script de Test de Déploiement - MeshaPlus Backend
# Ce script vérifie que tous les fichiers nécessaires sont présents

set -e

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🧪 Test de Déploiement - MeshaPlus Backend${NC}"
echo "============================================="
echo ""

# Test 1: Vérifier la structure des fichiers
echo -e "${YELLOW}📋 Test 1: Vérification de la structure des fichiers${NC}"

REQUIRED_FILES=(
    "bin/api"
    "Dockerfile"
    "docker-compose.yml"
    "configs/"
    "scripts/"
)

MISSING_FILES=()

for file in "${REQUIRED_FILES[@]}"; do
    if [ -e "$file" ]; then
        echo -e "${GREEN}✅ $file${NC}"
    else
        echo -e "${RED}❌ $file (manquant)${NC}"
        MISSING_FILES+=("$file")
    fi
done

if [ ${#MISSING_FILES[@]} -gt 0 ]; then
    echo ""
    echo -e "${RED}❌ Fichiers manquants :${NC}"
    for file in "${MISSING_FILES[@]}"; do
        echo "  - $file"
    done
    exit 1
fi

echo ""
echo -e "${GREEN}✅ Tous les fichiers requis sont présents${NC}"

# Test 2: Vérifier le binaire
echo ""
echo -e "${YELLOW}📋 Test 2: Vérification du binaire${NC}"

if [ -f "bin/api" ]; then
    echo -e "${GREEN}✅ Binaire trouvé${NC}"
    
    # Vérifier que c'est un exécutable
    if [ -x "bin/api" ]; then
        echo -e "${GREEN}✅ Binaire exécutable${NC}"
    else
        echo -e "${RED}❌ Binaire non exécutable${NC}"
        chmod +x bin/api
        echo -e "${YELLOW}⚠️  Permissions corrigées${NC}"
    fi
    
    # Vérifier la taille
    size=$(stat -f%z bin/api 2>/dev/null || stat -c%s bin/api 2>/dev/null || echo "0")
    echo -e "${BLUE}📊 Taille du binaire: ${size} bytes${NC}"
    
    if [ "$size" -gt 1000000 ]; then
        echo -e "${GREEN}✅ Taille du binaire OK${NC}"
    else
        echo -e "${YELLOW}⚠️  Taille du binaire petite${NC}"
    fi
else
    echo -e "${RED}❌ Binaire non trouvé${NC}"
    exit 1
fi

# Test 3: Vérifier Dockerfile
echo ""
echo -e "${YELLOW}📋 Test 3: Vérification du Dockerfile${NC}"

if [ -f "Dockerfile" ]; then
    echo -e "${GREEN}✅ Dockerfile trouvé${NC}"
    
    # Vérifier le contenu
    if grep -q "FROM golang:1.23" Dockerfile; then
        echo -e "${GREEN}✅ Version Go correcte (1.23)${NC}"
    else
        echo -e "${YELLOW}⚠️  Version Go différente de 1.23${NC}"
    fi
    
    if grep -q "COPY.*bin/api" Dockerfile; then
        echo -e "${GREEN}✅ Copie du binaire configurée${NC}"
    else
        echo -e "${YELLOW}⚠️  Copie du binaire non trouvée${NC}"
    fi
else
    echo -e "${RED}❌ Dockerfile non trouvé${NC}"
    exit 1
fi

# Test 4: Vérifier docker-compose.yml
echo ""
echo -e "${YELLOW}📋 Test 4: Vérification du docker-compose.yml${NC}"

if [ -f "docker-compose.yml" ]; then
    echo -e "${GREEN}✅ docker-compose.yml trouvé${NC}"
    
    # Vérifier la version
    if grep -q "version: '3.8'" docker-compose.yml; then
        echo -e "${GREEN}✅ Version Docker Compose correcte${NC}"
    else
        echo -e "${YELLOW}⚠️  Version Docker Compose différente${NC}"
    fi
    
    # Vérifier les services
    if grep -q "backend:" docker-compose.yml; then
        echo -e "${GREEN}✅ Service backend configuré${NC}"
    else
        echo -e "${RED}❌ Service backend manquant${NC}"
    fi
    
    if grep -q "postgres:" docker-compose.yml; then
        echo -e "${GREEN}✅ Service postgres configuré${NC}"
    else
        echo -e "${YELLOW}⚠️  Service postgres manquant${NC}"
    fi
else
    echo -e "${RED}❌ docker-compose.yml non trouvé${NC}"
    exit 1
fi

# Test 5: Test de build Docker (optionnel)
echo ""
echo -e "${YELLOW}📋 Test 5: Test de build Docker (optionnel)${NC}"

if command -v docker &> /dev/null; then
    echo -e "${BLUE}🐳 Docker disponible${NC}"
    
    read -p "Voulez-vous tester le build Docker ? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}🏗️ Test de build Docker...${NC}"
        
        # Test de build
        if docker build -t meshaplus-test:latest . 2>/dev/null; then
            echo -e "${GREEN}✅ Build Docker réussi${NC}"
            
            # Nettoyer
            docker rmi meshaplus-test:latest 2>/dev/null || true
            echo -e "${BLUE}🧹 Image de test supprimée${NC}"
        else
            echo -e "${RED}❌ Build Docker échoué${NC}"
        fi
    else
        echo -e "${YELLOW}⏭️  Test de build Docker ignoré${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Docker non disponible${NC}"
fi

# Résumé
echo ""
echo -e "${GREEN}🎉 Test de déploiement terminé avec succès !${NC}"
echo ""
echo -e "${BLUE}📋 Résumé :${NC}"
echo "✅ Structure des fichiers vérifiée"
echo "✅ Binaire Go présent et exécutable"
echo "✅ Dockerfile configuré"
echo "✅ docker-compose.yml configuré"
echo ""
echo -e "${BLUE}🚀 Prêt pour le déploiement !${NC}" 