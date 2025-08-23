#!/bin/bash

# 🚀 Script de Déploiement Binaire VPS - MeshaPlus Backend
# Ce script déploie directement le binaire sur le VPS sans Docker

set -e

# Configuration
APP_NAME="meshaplus-backend"
DEPLOY_PATH="/opt/meshaplus"
BINARY_PATH="/opt/meshaplus/bin"
CONFIG_PATH="/opt/meshaplus/configs"
LOG_PATH="/var/log/meshaplus"
SERVICE_NAME="meshaplus-backend"
USER_NAME="meshaplus"

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
    
    # Vérifier que le binaire existe
    if [ ! -f "bin/api" ]; then
        log_error "Le binaire bin/api n'existe pas"
        exit 1
    fi
    
    # Vérifier que les configs existent
    if [ ! -f "configs/config.prod.yaml" ]; then
        log_error "Le fichier de configuration configs/config.prod.yaml n'existe pas"
        exit 1
    fi
    
    log_success "Prérequis vérifiés"
}

# Création des répertoires et utilisateur
setup_environment() {
    log "📁 Configuration de l'environnement..."
    
    # Créer l'utilisateur si il n'existe pas
    if ! id "$USER_NAME" &>/dev/null; then
        sudo useradd -r -s /bin/false -d "$DEPLOY_PATH" "$USER_NAME"
        log_success "Utilisateur $USER_NAME créé"
    fi
    
    # Créer les répertoires
    sudo mkdir -p "$DEPLOY_PATH"
    sudo mkdir -p "$BINARY_PATH"
    sudo mkdir -p "$CONFIG_PATH"
    sudo mkdir -p "$LOG_PATH"
    sudo mkdir -p "$DEPLOY_PATH/uploads"
    
    # Changer les permissions
    sudo chown -R "$USER_NAME:$USER_NAME" "$DEPLOY_PATH"
    sudo chown -R "$USER_NAME:$USER_NAME" "$LOG_PATH"
    
    log_success "Environnement configuré"
}

# Sauvegarde de l'ancienne version
backup_previous_version() {
    log "💾 Sauvegarde de l'ancienne version..."
    
    if [ -f "$BINARY_PATH/api" ]; then
        # Arrêter le service s'il tourne
        if sudo systemctl is-active --quiet "$SERVICE_NAME"; then
            sudo systemctl stop "$SERVICE_NAME"
            log_success "Ancien service arrêté"
        fi
        
        # Créer une sauvegarde
        sudo cp "$BINARY_PATH/api" "$BINARY_PATH/api.backup.$(date +%Y%m%d-%H%M%S)"
        log_success "Ancienne version sauvegardée"
    else
        log_warning "Aucune version précédente trouvée"
    fi
}

# Déploiement du binaire
deploy_binary() {
    log "📦 Déploiement du binaire..."
    
    # Copier le binaire
    sudo cp bin/api "$BINARY_PATH/"
    sudo chown "$USER_NAME:$USER_NAME" "$BINARY_PATH/api"
    sudo chmod +x "$BINARY_PATH/api"
    
    # Copier les configurations
    sudo cp configs/config.prod.yaml "$CONFIG_PATH/"
    sudo chown "$USER_NAME:$USER_NAME" "$CONFIG_PATH/config.prod.yaml"
    
    log_success "Binaire déployé"
}

# Création du service systemd
create_systemd_service() {
    log "🔧 Création du service systemd..."
    
    cat << EOF | sudo tee /etc/systemd/system/$SERVICE_NAME.service
[Unit]
Description=MeshaPlus Backend API
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
User=$USER_NAME
Group=$USER_NAME
WorkingDirectory=$DEPLOY_PATH
ExecStart=$BINARY_PATH/api
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=$SERVICE_NAME

# Variables d'environnement
Environment=CONFIG_FILE=$CONFIG_PATH/config.prod.yaml
Environment=PORT=8080

# Sécurité
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$DEPLOY_PATH $LOG_PATH

[Install]
WantedBy=multi-user.target
EOF

    # Recharger systemd
    sudo systemctl daemon-reload
    sudo systemctl enable "$SERVICE_NAME"
    
    log_success "Service systemd créé"
}

# Vérification des services locaux
check_local_services() {
    log "🔍 Vérification des services locaux..."
    
    # Vérifier PostgreSQL
    if ! sudo systemctl is-active --quiet postgresql; then
        log_warning "PostgreSQL n'est pas en cours d'exécution"
        log "Pour installer PostgreSQL: sudo apt update && sudo apt install postgresql postgresql-contrib"
        log "Pour démarrer PostgreSQL: sudo systemctl start postgresql && sudo systemctl enable postgresql"
    else
        log_success "PostgreSQL est en cours d'exécution"
    fi
    
    # Vérifier Redis (optionnel)
    if command -v redis-server &> /dev/null; then
        if ! sudo systemctl is-active --quiet redis-server; then
            log_warning "Redis n'est pas en cours d'exécution"
            log "Pour démarrer Redis: sudo systemctl start redis-server && sudo systemctl enable redis-server"
        else
            log_success "Redis est en cours d'exécution"
        fi
    else
        log_warning "Redis n'est pas installé"
        log "Pour installer Redis: sudo apt update && sudo apt install redis-server"
    fi
}

# Démarrage de l'application
start_application() {
    log "🚀 Démarrage de l'application..."
    
    # Démarrer le service
    sudo systemctl start "$SERVICE_NAME"
    
    if [ $? -eq 0 ]; then
        log_success "Application démarrée"
    else
        log_error "Échec du démarrage de l'application"
        sudo systemctl status "$SERVICE_NAME"
        exit 1
    fi
}

# Vérification de la santé de l'application
health_check() {
    log "🏥 Vérification de la santé de l'application..."
    
    # Attendre que l'application soit prête
    sleep 10
    
    # Test de santé
    for i in {1..10}; do
        if curl -f http://localhost:8080/api/v1/health >/dev/null 2>&1; then
            log_success "Application en bonne santé"
            return 0
        fi
        
        log_warning "Tentative $i/10 - Application pas encore prête"
        sleep 5
    done
    
    log_error "L'application n'est pas en bonne santé après 10 tentatives"
    return 1
}

# Configuration de Nginx (optionnel)
setup_nginx() {
    log "🌐 Configuration de Nginx..."
    
    if command -v nginx &> /dev/null; then
        cat << EOF | sudo tee /etc/nginx/sites-available/$SERVICE_NAME
server {
    listen 80;
    server_name _;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

        # Activer le site
        sudo ln -sf /etc/nginx/sites-available/$SERVICE_NAME /etc/nginx/sites-enabled/
        sudo nginx -t && sudo systemctl reload nginx
        
        log_success "Nginx configuré"
    else
        log_warning "Nginx n'est pas installé"
        log "Pour installer Nginx: sudo apt update && sudo apt install nginx"
    fi
}

# Nettoyage des anciennes sauvegardes
cleanup_old_backups() {
    log "🧹 Nettoyage des anciennes sauvegardes..."
    
    # Garder seulement les 5 dernières sauvegardes
    find "$BINARY_PATH" -name "api.backup.*" -type f -printf '%T@ %p\n' | \
    sort -nr | tail -n +6 | cut -d' ' -f2- | xargs -r sudo rm -f
    
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
    echo "• Binaire: $BINARY_PATH/api"
    echo "• Config: $CONFIG_PATH/config.prod.yaml"
    echo "• Logs: $LOG_PATH"
    echo "• Service: $SERVICE_NAME"
    echo ""
    echo "🔍 Commandes utiles :"
    echo "• Statut: sudo systemctl status $SERVICE_NAME"
    echo "• Logs: sudo journalctl -u $SERVICE_NAME -f"
    echo "• Redémarrer: sudo systemctl restart $SERVICE_NAME"
    echo "• Arrêter: sudo systemctl stop $SERVICE_NAME"
    echo "• Démarrer: sudo systemctl start $SERVICE_NAME"
    echo ""
    echo "🌐 Accès :"
    echo "• API: http://localhost:8080"
    echo "• Health: http://localhost:8080/api/v1/health"
    echo "• Swagger: http://localhost:8080/swagger/index.html"
    echo ""
}

# Fonction principale
main() {
    log "🚀 Début du déploiement binaire MeshaPlus Backend"
    echo "=============================================="
    
    # Vérifier les prérequis
    check_prerequisites
    
    # Configurer l'environnement
    setup_environment
    
    # Sauvegarder l'ancienne version
    backup_previous_version
    
    # Déployer le binaire
    deploy_binary
    
    # Créer le service systemd
    create_systemd_service
    
    # Vérifier les services locaux
    check_local_services
    
    # Configurer Nginx
    setup_nginx
    
    # Démarrer l'application
    start_application
    
    # Vérifier la santé
    if health_check; then
        # Nettoyer les anciennes sauvegardes
        cleanup_old_backups
        
        # Afficher les informations
        show_deployment_info
        
        log_success "Déploiement terminé avec succès !"
        exit 0
    else
        log_error "Échec du déploiement - L'application n'est pas en bonne santé"
        
        # Rollback automatique
        log "🔄 Tentative de rollback..."
        sudo systemctl stop "$SERVICE_NAME"
        if [ -f "$BINARY_PATH/api.backup.$(date +%Y%m%d-%H%M%S)" ]; then
            sudo cp "$BINARY_PATH/api.backup.$(date +%Y%m%d-%H%M%S)" "$BINARY_PATH/api"
            sudo systemctl start "$SERVICE_NAME"
        fi
        
        exit 1
    fi
}

# Gestion des erreurs
trap 'log_error "Erreur survenue à la ligne $LINENO"; exit 1' ERR

# Exécution du script principal
main "$@" 