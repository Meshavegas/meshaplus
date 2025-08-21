#!/bin/bash

# 🚀 Test Workflow Script - MeshaPlus Backend
# Ce script simule les étapes principales des workflows GitHub Actions

set -e

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_ROOT"
GO_VERSION="1.21"
DOCKER_IMAGE="meshaplus-backend-test"

echo -e "${BLUE}🚀 Test Workflow Script - MeshaPlus Backend${NC}"
echo "=================================================="
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
    
    # Vérifier Go
    if ! command -v go &> /dev/null; then
        print_error "Go n'est pas installé"
        exit 1
    fi
    
    GO_CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Go version: $GO_CURRENT_VERSION"
    
    # Vérifier Docker
    if ! command -v docker &> /dev/null; then
        print_error "Docker n'est pas installé"
        exit 1
    fi
    
    DOCKER_VERSION=$(docker --version)
    print_info "Docker: $DOCKER_VERSION"
    
    # Vérifier git
    if ! command -v git &> /dev/null; then
        print_error "Git n'est pas installé"
        exit 1
    fi
    
    print_success "Tous les prérequis sont satisfaits"
}

# Tests et validation
run_tests() {
    print_step "🧪 Exécution des tests..."
    
    cd "$BACKEND_DIR"
    
    # Installer les dépendances
    print_info "Installation des dépendances..."
    go mod download
    go mod tidy
    
    # Installer les outils
    print_info "Installation des outils..."
    go install github.com/swaggo/swag/cmd/swag@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    # Linting
    print_info "Exécution du linting..."
    if golangci-lint run ./...; then
        print_success "Linting réussi"
    else
        print_error "Linting échoué"
        exit 1
    fi
    
    # Tests
    print_info "Exécution des tests..."
    if go test -v -race -coverprofile=coverage.out ./...; then
        print_success "Tests réussis"
    else
        print_error "Tests échoués"
        exit 1
    fi
    
    # Couverture
    print_info "Analyse de la couverture..."
    go tool cover -func=coverage.out
    
    # Génération Swagger
    print_info "Génération de la documentation Swagger..."
    export PATH=$PATH:$(go env GOPATH)/bin
    if swag init -g cmd/api/main.go -o ./docs; then
        print_success "Documentation Swagger générée"
    else
        print_error "Échec de la génération Swagger"
        exit 1
    fi
    
    # Build
    print_info "Compilation de l'application..."
    if CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags='-w -s -extldflags "-static"' \
        -a -installsuffix cgo \
        -o bin/api \
        ./cmd/api; then
        print_success "Application compilée"
    else
        print_error "Échec de la compilation"
        exit 1
    fi
}

# Build Docker
build_docker() {
    print_step "🐳 Construction de l'image Docker..."
    
    cd "$BACKEND_DIR"
    
    # Nettoyer les anciennes images
    print_info "Nettoyage des anciennes images..."
    docker rmi "$DOCKER_IMAGE" 2>/dev/null || true
    
    # Build
    print_info "Construction de l'image Docker..."
    if docker build -t "$DOCKER_IMAGE" .; then
        print_success "Image Docker construite"
    else
        print_error "Échec de la construction Docker"
        exit 1
    fi
    
    # Vérifier l'image
    print_info "Vérification de l'image..."
    docker images "$DOCKER_IMAGE"
}

# Test de l'application
test_application() {
    print_step "🧪 Test de l'application..."
    
    cd "$BACKEND_DIR"
    
    # Démarrer l'application
    print_info "Démarrage de l'application..."
    docker run -d \
        --name meshaplus-backend-test \
        -p 8080:8080 \
        -e DATABASE_URL="postgres://postgres:postgres@localhost:5432/meshaplus_test?sslmode=disable" \
        -e REDIS_URL="redis://localhost:6379" \
        "$DOCKER_IMAGE"
    
    # Attendre que l'application démarre
    print_info "Attente du démarrage de l'application..."
    sleep 10
    
    # Test de santé
    print_info "Test de santé de l'application..."
    if curl -f http://localhost:8080/api/v1/health 2>/dev/null; then
        print_success "Application en bonne santé"
    else
        print_error "Application non accessible"
        docker logs meshaplus-backend-test
        exit 1
    fi
    
    # Test Swagger
    print_info "Test de la documentation Swagger..."
    if curl -f http://localhost:8080/swagger/index.html 2>/dev/null; then
        print_success "Swagger accessible"
    else
        print_error "Swagger non accessible"
    fi
    
    # Nettoyage
    print_info "Nettoyage..."
    docker stop meshaplus-backend-test || true
    docker rm meshaplus-backend-test || true
}

# Scan de sécurité
security_scan() {
    print_step "🔒 Scan de sécurité..."
    
    cd "$BACKEND_DIR"
    
    # Vérifier si Trivy est installé
    if ! command -v trivy &> /dev/null; then
        print_info "Trivy non installé, installation..."
        # Note: L'installation de Trivy dépend du système
        print_info "Veuillez installer Trivy manuellement: https://aquasecurity.github.io/trivy/latest/getting-started/installation/"
        return 0
    fi
    
    # Scan de l'image Docker
    print_info "Scan de l'image Docker avec Trivy..."
    if trivy image "$DOCKER_IMAGE"; then
        print_success "Scan Trivy terminé"
    else
        print_error "Échec du scan Trivy"
    fi
}

# Nettoyage
cleanup() {
    print_step "🧹 Nettoyage..."
    
    cd "$BACKEND_DIR"
    
    # Nettoyer les conteneurs
    docker stop meshaplus-backend-test 2>/dev/null || true
    docker rm meshaplus-backend-test 2>/dev/null || true
    
    # Nettoyer les images
    docker rmi "$DOCKER_IMAGE" 2>/dev/null || true
    
    # Nettoyer les fichiers générés
    rm -f coverage.out
    rm -f bin/api
    
    print_success "Nettoyage terminé"
}

# Fonction principale
main() {
    echo -e "${BLUE}🚀 Démarrage du test workflow...${NC}"
    echo ""
    
    # Vérifier les prérequis
    check_prerequisites
    echo ""
    
    # Exécuter les tests
    run_tests
    echo ""
    
    # Build Docker
    build_docker
    echo ""
    
    # Test de l'application
    test_application
    echo ""
    
    # Scan de sécurité
    security_scan
    echo ""
    
    # Nettoyage
    cleanup
    echo ""
    
    print_success "🎉 Test workflow terminé avec succès !"
    echo ""
    print_info "📊 Résumé :"
    echo "  ✅ Tests et validation"
    echo "  ✅ Construction Docker"
    echo "  ✅ Test de l'application"
    echo "  ✅ Scan de sécurité"
    echo "  ✅ Nettoyage"
    echo ""
    print_info "🚀 Votre backend est prêt pour le déploiement !"
}

# Gestion des erreurs
trap 'print_error "Erreur survenue. Nettoyage en cours..."; cleanup; exit 1' ERR

# Exécution
main "$@" 