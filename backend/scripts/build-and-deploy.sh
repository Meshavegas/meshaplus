#!/bin/bash

# 🚀 Script de Build et Déploiement Binaire - MeshaPlus Backend
# Ce script compile le binaire et le déploie sur le VPS

set -e

# Configuration
APP_NAME="meshaplus-backend"
VPS_HOST=""
VPS_USER=""
VPS_PORT="22"
DEPLOY_SCRIPT="deploy-binary-vps.sh"
LOCAL_BUILD_DIR="bin"
REMOTE_DEPLOY_PATH="/opt/meshaplus"

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
    
    # Vérifier Go
    if ! command -v go &> /dev/null; then
        log_error "Go n'est pas installé"
        exit 1
    fi
    
    # Vérifier SSH
    if ! command -v ssh &> /dev/null; then
        log_error "SSH n'est pas installé"
        exit 1
    fi
    
    # Vérifier SCP
    if ! command -v scp &> /dev/null; then
        log_error "SCP n'est pas installé"
        exit 1
    fi
    
    # Vérifier les variables d'environnement
    if [ -z "$VPS_HOST" ] || [ -z "$VPS_USER" ]; then
        log_error "VPS_HOST et VPS_USER doivent être définis"
        log "Exemple: VPS_HOST=your-vps.com VPS_USER=root ./build-and-deploy.sh"
        exit 1
    fi
    
    log_success "Prérequis vérifiés"
}

# Build du binaire
build_binary() {
    log "🏗️ Compilation du binaire..."
    
    # Nettoyer le répertoire de build
    rm -rf "$LOCAL_BUILD_DIR"
    mkdir -p "$LOCAL_BUILD_DIR"
    
    # Compiler pour Linux AMD64
    GOOS=linux GOARCH=amd64 go build -o "$LOCAL_BUILD_DIR/api" ./cmd/api
    
    if [ $? -eq 0 ]; then
        log_success "Binaire compilé: $LOCAL_BUILD_DIR/api"
    else
        log_error "Échec de la compilation"
        exit 1
    fi
}

# Test de connexion SSH
test_ssh_connection() {
    log "🔌 Test de connexion SSH..."
    
    if ssh -p "$VPS_PORT" -o ConnectTimeout=10 -o BatchMode=yes "$VPS_USER@$VPS_HOST" "echo 'SSH connection successful'" 2>/dev/null; then
        log_success "Connexion SSH établie"
    else
        log_error "Impossible de se connecter au VPS"
        log "Vérifiez:"
        log "• L'adresse IP/domaine du VPS"
        log "• Les identifiants SSH"
        log "• La clé SSH (si utilisée)"
        exit 1
    fi
}

# Transfert des fichiers
transfer_files() {
    log "📤 Transfert des fichiers..."
    
    # Créer un répertoire temporaire sur le VPS
    ssh -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" "mkdir -p /tmp/meshaplus-deploy"
    
    # Transférer le binaire
    scp -P "$VPS_PORT" "$LOCAL_BUILD_DIR/api" "$VPS_USER@$VPS_HOST:/tmp/meshaplus-deploy/"
    
    # Transférer la configuration
    scp -P "$VPS_PORT" "configs/config.prod.yaml" "$VPS_USER@$VPS_HOST:/tmp/meshaplus-deploy/"
    
    # Transférer le script de déploiement
    scp -P "$VPS_PORT" "scripts/$DEPLOY_SCRIPT" "$VPS_USER@$VPS_HOST:/tmp/meshaplus-deploy/"
    
    log_success "Fichiers transférés"
}

# Exécution du déploiement
execute_deployment() {
    log "🚀 Exécution du déploiement..."
    
    # Rendre le script exécutable et l'exécuter
    ssh -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" << EOF
        cd /tmp/meshaplus-deploy
        chmod +x $DEPLOY_SCRIPT
        ./$DEPLOY_SCRIPT
EOF

    if [ $? -eq 0 ]; then
        log_success "Déploiement terminé"
    else
        log_error "Échec du déploiement"
        exit 1
    fi
}

# Nettoyage des fichiers temporaires
cleanup_temp_files() {
    log "🧹 Nettoyage des fichiers temporaires..."
    
    # Nettoyer sur le VPS
    ssh -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" "rm -rf /tmp/meshaplus-deploy"
    
    # Nettoyer localement
    rm -rf "$LOCAL_BUILD_DIR"
    
    log_success "Nettoyage terminé"
}

# Vérification post-déploiement
verify_deployment() {
    log "🔍 Vérification post-déploiement..."
    
    # Attendre que l'application soit prête
    sleep 10
    
    # Test de santé
    for i in {1..5}; do
        if ssh -p "$VPS_PORT" "$VPS_USER@$VPS_HOST" "curl -f http://localhost:8080/api/v1/health" >/dev/null 2>&1; then
            log_success "Application déployée et fonctionnelle"
            return 0
        fi
        
        log_warning "Tentative $i/5 - Application pas encore prête"
        sleep 10
    done
    
    log_error "L'application n'est pas accessible après déploiement"
    return 1
}

# Affichage des informations de déploiement
show_deployment_info() {
    log "📊 Informations de déploiement..."
    
    echo ""
    echo "🎉 Déploiement terminé avec succès !"
    echo ""
    echo "📋 Informations :"
    echo "• VPS: $VPS_HOST"
    echo "• Utilisateur: $VPS_USER"
    echo "• Application: $APP_NAME"
    echo "• Version: $(date +%Y-%m-%d\ %H:%M:%S)"
    echo ""
    echo "🌐 Accès :"
    echo "• API: http://$VPS_HOST:8080"
    echo "• Health: http://$VPS_HOST:8080/api/v1/health"
    echo "• Swagger: http://$VPS_HOST:8080/swagger/index.html"
    echo ""
    echo "🔍 Commandes utiles :"
    echo "• Connexion SSH: ssh -p $VPS_PORT $VPS_USER@$VPS_HOST"
    echo "• Statut service: ssh $VPS_USER@$VPS_HOST 'sudo systemctl status $APP_NAME'"
    echo "• Logs service: ssh $VPS_USER@$VPS_HOST 'sudo journalctl -u $APP_NAME -f'"
    echo ""
}

# Fonction d'aide
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --host HOST     Adresse IP ou domaine du VPS"
    echo "  -u, --user USER     Utilisateur SSH du VPS"
    echo "  -p, --port PORT     Port SSH (défaut: 22)"
    echo "  --help              Affiche cette aide"
    echo ""
    echo "Exemples:"
    echo "  VPS_HOST=192.168.1.100 VPS_USER=root $0"
    echo "  $0 -h your-vps.com -u root -p 2222"
    echo ""
}

# Parsing des arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--host)
                VPS_HOST="$2"
                shift 2
                ;;
            -u|--user)
                VPS_USER="$2"
                shift 2
                ;;
            -p|--port)
                VPS_PORT="$2"
                shift 2
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                log_error "Option inconnue: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Fonction principale
main() {
    log "🚀 Début du build et déploiement MeshaPlus Backend"
    echo "================================================"
    
    # Parser les arguments
    parse_arguments "$@"
    
    # Vérifier les prérequis
    check_prerequisites
    
    # Tester la connexion SSH
    test_ssh_connection
    
    # Build du binaire
    build_binary
    
    # Transférer les fichiers
    transfer_files
    
    # Exécuter le déploiement
    execute_deployment
    
    # Vérifier le déploiement
    if verify_deployment; then
        # Nettoyer les fichiers temporaires
        cleanup_temp_files
        
        # Afficher les informations
        show_deployment_info
        
        log_success "Build et déploiement terminés avec succès !"
        exit 0
    else
        log_error "Échec de la vérification post-déploiement"
        exit 1
    fi
}

# Gestion des erreurs
trap 'log_error "Erreur survenue à la ligne $LINENO"; exit 1' ERR

# Exécution du script principal
main "$@" 