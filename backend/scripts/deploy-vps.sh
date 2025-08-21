#!/bin/bash

# 🚀 Script de Déploiement VPS - MeshaPlus Backend
# Ce script est exécuté sur le VPS pour déployer l'application

set -e

# Configuration
APP_NAME="meshaplus-backend"
DEPLOY_PATH="/opt/meshaplus"
BACKUP_PATH="/opt/meshaplus/backups"
LOG_PATH="/var/log/meshaplus"

# Couleurs pour les logs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Fonction de logging
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Vérification des prérequis
check_prerequisites() {
    log "🔍 Vérification des prérequis..."
    
    # Vérifier Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker n'est pas installé"
        exit 1
    fi
    
    # Vérifier Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose n'est pas installé"
        exit 1
    fi
    
    # Vérifier que le binaire existe
    if [ ! -f "bin/api" ]; then
        log_error "Le binaire bin/api n'existe pas"
        exit 1
    fi
    
    log_success "Prérequis vérifiés"
}

# Création des répertoires nécessaires
setup_directories() {
    log "📁 Création des répertoires..."
    
    mkdir -p "$DEPLOY_PATH"
    mkdir -p "$BACKUP_PATH"
    mkdir -p "$LOG_PATH"
    mkdir -p "$DEPLOY_PATH/configs"
    mkdir -p "$DEPLOY_PATH/scripts"
    
    log_success "Répertoires créés"
}

# Sauvegarde de l'ancienne version
backup_previous_version() {
    log "💾 Sauvegarde de l'ancienne version..."
    
    if docker ps -q -f name="$APP_NAME" | grep -q .; then
        # Arrêter l'ancienne version
        docker-compose down || true
        
        # Créer une sauvegarde
        if [ -f "docker-compose.yml" ]; then
            cp docker-compose.yml "$BACKUP_PATH/docker-compose-$(date +%Y%m%d-%H%M%S).yml"
        fi
        
        log_success "Ancienne version sauvegardée"
    else
        log_warning "Aucune version précédente en cours d'exécution"
    fi
}

# Construction de l'image Docker
build_docker_image() {
    log "🏗️ Construction de l'image Docker..."
    
    # Supprimer l'ancienne image
    docker rmi "$APP_NAME:latest" 2>/dev/null || true
    
    # Construire la nouvelle image
    docker build -t "$APP_NAME:latest" . 2>&1 | tee "$LOG_PATH/build-$(date +%Y%m%d-%H%M%S).log"
    
    if [ $? -eq 0 ]; then
        log_success "Image Docker construite avec succès"
    else
        log_error "Échec de la construction de l'image Docker"
        exit 1
    fi
}

# Démarrage de l'application
start_application() {
    log "🚀 Démarrage de l'application..."
    
    # Démarrer les conteneurs
    docker-compose up -d 2>&1 | tee "$LOG_PATH/start-$(date +%Y%m%d-%H%M%S).log"
    
    if [ $? -eq 0 ]; then
        log_success "Application démarrée"
    else
        log_error "Échec du démarrage de l'application"
        exit 1
    fi
}

# Vérification de la santé de l'application
health_check() {
    log "🏥 Vérification de la santé de l'application..."
    
    # Attendre que l'application soit prête
    sleep 30
    
    # Test de santé
    for i in {1..10}; do
        if curl -f http://localhost:8080/api/v1/health >/dev/null 2>&1; then
            log_success "Application en bonne santé"
            return 0
        fi
        
        log_warning "Tentative $i/10 - Application pas encore prête"
        sleep 10
    done
    
    log_error "L'application n'est pas en bonne santé après 10 tentatives"
    return 1
}

# Nettoyage des anciennes images
cleanup_old_images() {
    log "🧹 Nettoyage des anciennes images..."
    
    # Supprimer les images non utilisées
    docker image prune -f
    
    # Supprimer les conteneurs arrêtés
    docker container prune -f
    
    log_success "Nettoyage terminé"
}

# Affichage des informations de déploiement
show_deployment_info() {
    log "📊 Informations de déploiement..."
    
    echo ""
    echo "🎉 Déploiement terminé avec succès !"
    echo ""
    echo "📋 Informations :"
    echo "• Application: $APP_NAME"
    echo "• Version: $(date +%Y-%m-%d\ %H:%M:%S)"
    echo "• Port: 8080"
    echo "• Logs: $LOG_PATH"
    echo "• Backups: $BACKUP_PATH"
    echo ""
    echo "🔍 Commandes utiles :"
    echo "• Voir les logs: docker-compose logs -f"
    echo "• Statut: docker-compose ps"
    echo "• Arrêter: docker-compose down"
    echo "• Redémarrer: docker-compose restart"
    echo ""
}

# Fonction principale
main() {
    log "🚀 Début du déploiement MeshaPlus Backend"
    echo "====================================="
    
    # Vérifier les prérequis
    check_prerequisites
    
    # Créer les répertoires
    setup_directories
    
    # Sauvegarder l'ancienne version
    backup_previous_version
    
    # Construire l'image Docker
    build_docker_image
    
    # Démarrer l'application
    start_application
    
    # Vérifier la santé
    if health_check; then
        # Nettoyer les anciennes images
        cleanup_old_images
        
        # Afficher les informations
        show_deployment_info
        
        log_success "Déploiement terminé avec succès !"
        exit 0
    else
        log_error "Échec du déploiement - L'application n'est pas en bonne santé"
        
        # Rollback automatique
        log "🔄 Tentative de rollback..."
        docker-compose down
        docker rmi "$APP_NAME:latest" 2>/dev/null || true
        
        exit 1
    fi
}

# Gestion des erreurs
trap 'log_error "Erreur survenue à la ligne $LINENO"; exit 1' ERR

# Exécution du script principal
main "$@" 