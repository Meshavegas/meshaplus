#!/bin/bash

# 🔍 Test de Linting - MeshaPlus Backend
# Script de test spécifique pour golangci-lint

set -e

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 Test de Linting - MeshaPlus Backend${NC}"
echo "====================================="
echo ""

# Vérifier que golangci-lint est installé
check_golangci_lint() {
    echo -e "${YELLOW}📋 Vérification de golangci-lint${NC}"
    
    if ! command -v golangci-lint &> /dev/null; then
        echo -e "${YELLOW}⚠️  golangci-lint non trouvé, installation...${NC}"
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    fi
    
    echo -e "${GREEN}✅ golangci-lint version: $(golangci-lint version)${NC}"
}

# Test de linting avec exclusion du répertoire docs
test_linting() {
    echo -e "${YELLOW}📋 Test de linting (excluant docs/)${NC}"
    
    if golangci-lint run --exclude-dirs docs ./...; then
        echo -e "${GREEN}✅ Linting réussi${NC}"
    else
        echo -e "${RED}❌ Linting échoué${NC}"
        exit 1
    fi
}

# Test de linting complet (pour debug)
test_linting_debug() {
    echo -e "${YELLOW}📋 Test de linting avec debug${NC}"
    
    if golangci-lint run --exclude-dirs docs --verbose ./...; then
        echo -e "${GREEN}✅ Linting debug réussi${NC}"
    else
        echo -e "${RED}❌ Linting debug échoué${NC}"
        exit 1
    fi
}

# Vérifier la configuration
check_config() {
    echo -e "${YELLOW}📋 Vérification de la configuration${NC}"
    
    if golangci-lint config; then
        echo -e "${GREEN}✅ Configuration valide${NC}"
    else
        echo -e "${RED}❌ Configuration invalide${NC}"
        exit 1
    fi
}

# Nettoyer les fichiers temporaires
cleanup() {
    echo -e "${YELLOW}📋 Nettoyage${NC}"
    
    # Supprimer les fichiers de cache si nécessaire
    rm -rf .golangci-cache 2>/dev/null || true
    
    echo -e "${GREEN}✅ Nettoyage terminé${NC}"
}

# Affichage du résumé
show_summary() {
    echo ""
    echo -e "${GREEN}🎉 Test de linting terminé !${NC}"
    echo ""
    echo -e "${BLUE}📋 Résumé :${NC}"
    echo "✅ golangci-lint installé et configuré"
    echo "✅ Configuration valide"
    echo "✅ Linting réussi (excluant docs/)"
    echo ""
    echo -e "${BLUE}🚀 Le workflow GitHub Actions devrait maintenant passer !${NC}"
}

# Fonction principale
main() {
    check_golangci_lint
    echo ""
    
    check_config
    echo ""
    
    test_linting
    echo ""
    
    cleanup
    echo ""
    
    show_summary
}

# Exécution
main "$@" 