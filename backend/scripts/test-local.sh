#!/bin/bash

# 🧪 Script de Test Local - MeshaPlus Backend
# Ce script teste les commandes du workflow GitHub Actions localement

set -e

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Test Local - MeshaPlus Backend${NC}"
echo "====================================="
echo ""

# Fonction pour afficher les étapes
print_step() {
    echo -e "${YELLOW}📋 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Vérification des prérequis
check_prerequisites() {
    print_step "Vérification des prérequis..."
    
    # Vérifier que nous sommes dans le bon répertoire
    if [ ! -f "go.mod" ]; then
        print_error "go.mod non trouvé. Assurez-vous d'être dans le répertoire backend"
        exit 1
    fi
    
    # Vérifier Go
    if ! command -v go &> /dev/null; then
        print_error "Go n'est pas installé"
        exit 1
    fi
    
    print_info "Go version: $(go version)"
    print_success "Prérequis vérifiés"
}

# Test des dépendances
test_dependencies() {
    print_step "Test des dépendances..."
    
    # Télécharger les dépendances
    go mod download
    print_success "Dépendances téléchargées"
    
    # Nettoyer les dépendances
    go mod tidy
    print_success "Dépendances nettoyées"
    
    # Installer les outils
    go install github.com/swaggo/swag/cmd/swag@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    print_success "Outils installés"
}

# Test du linting
test_linting() {
    print_step "Test du linting..."
    
    if golangci-lint run ./...; then
        print_success "Linting réussi"
    else
        print_error "Linting échoué"
        exit 1
    fi
}

# Test des tests
test_tests() {
    print_step "Test des tests..."
    
    # Exécuter les tests
    if go test -v -race -coverprofile=coverage.out ./...; then
        print_success "Tests réussis"
    else
        print_error "Tests échoués"
        exit 1
    fi
    
    # Afficher la couverture
    go tool cover -func=coverage.out
}

# Test de la génération Swagger
test_swagger() {
    print_step "Test de la génération Swagger..."
    
    # Exporter le PATH
    export PATH=$PATH:$(go env GOPATH)/bin
    
    # Générer la documentation Swagger
    if swag init -g cmd/api/main.go -o ./docs; then
        print_success "Documentation Swagger générée"
    else
        print_error "Génération Swagger échouée"
        exit 1
    fi
}

# Test de la compilation
test_build() {
    print_step "Test de la compilation..."
    
    # Créer le répertoire bin s'il n'existe pas
    mkdir -p bin
    
    # Compiler l'application
    if CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags='-w -s -extldflags "-static"' \
        -a -installsuffix cgo \
        -o bin/api \
        ./cmd/api; then
        print_success "Compilation réussie"
        
        # Vérifier que le binaire existe
        if [ -f "bin/api" ]; then
            print_info "Binaire créé: bin/api"
            ls -la bin/api
        fi
    else
        print_error "Compilation échouée"
        exit 1
    fi
}

# Test de l'exécution
test_run() {
    print_step "Test de l'exécution..."
    
    # Vérifier que le binaire existe
    if [ ! -f "bin/api" ]; then
        print_error "Binaire non trouvé. Exécutez d'abord le test de compilation"
        exit 1
    fi
    
    print_info "Démarrage de l'application en arrière-plan..."
    
    # Démarrer l'application en arrière-plan
    ./bin/api &
    APP_PID=$!
    
    # Attendre que l'application démarre
    sleep 5
    
    # Tester l'endpoint de santé
    if curl -f http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        print_success "Application démarrée et accessible"
    else
        print_error "Application non accessible"
        kill $APP_PID 2>/dev/null || true
        exit 1
    fi
    
    # Tester Swagger
    if curl -f http://localhost:8080/swagger/index.html > /dev/null 2>&1; then
        print_success "Swagger accessible"
    else
        print_warning "Swagger non accessible"
    fi
    
    # Arrêter l'application
    kill $APP_PID 2>/dev/null || true
    print_info "Application arrêtée"
}

# Nettoyage
cleanup() {
    print_step "Nettoyage..."
    
    # Supprimer les fichiers temporaires
    rm -f coverage.out
    rm -f bin/api
    
    print_success "Nettoyage terminé"
}

# Affichage du résumé
show_summary() {
    print_step "Test local terminé !"
    echo ""
    echo -e "${GREEN}🎉 Tous les tests sont passés !${NC}"
    echo ""
    echo -e "${BLUE}📋 Résumé des tests :${NC}"
    echo "✅ Dépendances téléchargées et nettoyées"
    echo "✅ Linting réussi"
    echo "✅ Tests exécutés avec succès"
    echo "✅ Documentation Swagger générée"
    echo "✅ Application compilée"
    echo "✅ Application testée et fonctionnelle"
    echo ""
    echo -e "${BLUE}🚀 Votre application est prête pour le déploiement !${NC}"
}

# Fonction principale
main() {
    echo -e "${BLUE}🧪 Démarrage des tests locaux...${NC}"
    echo ""
    
    check_prerequisites
    echo ""
    
    test_dependencies
    echo ""
    
    test_linting
    echo ""
    
    test_tests
    echo ""
    
    test_swagger
    echo ""
    
    test_build
    echo ""
    
    test_run
    echo ""
    
    cleanup
    echo ""
    
    show_summary
}

# Gestion des erreurs
trap 'print_error "Erreur survenue. Arrêt du script."; exit 1' ERR

# Exécution
main "$@" 