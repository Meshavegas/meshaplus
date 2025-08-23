#!/bin/bash

# 🔧 Script d'Installation des Services Locaux - MeshaPlus Backend
# Ce script installe et configure PostgreSQL et Redis sur le VPS

set -e

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

# Mise à jour du système
update_system() {
    log "🔄 Mise à jour du système..."
    
    sudo apt update
    sudo apt upgrade -y
    
    log_success "Système mis à jour"
}

# Installation de PostgreSQL
install_postgresql() {
    log "🐘 Installation de PostgreSQL..."
    
    # Installer PostgreSQL
    sudo apt install -y postgresql postgresql-contrib
    
    # Démarrer et activer le service
    sudo systemctl start postgresql
    sudo systemctl enable postgresql
    
    # Vérifier le statut
    if sudo systemctl is-active --quiet postgresql; then
        log_success "PostgreSQL installé et démarré"
    else
        log_error "Échec du démarrage de PostgreSQL"
        exit 1
    fi
}

# Configuration de PostgreSQL
configure_postgresql() {
    log "⚙️ Configuration de PostgreSQL..."
    
    # Créer la base de données et l'utilisateur
    sudo -u postgres psql << EOF
-- Créer la base de données
CREATE DATABASE meshaplus;

-- Créer l'utilisateur avec mot de passe
CREATE USER meshaplus WITH PASSWORD 'meshaplus_password';

-- Donner les privilèges
GRANT ALL PRIVILEGES ON DATABASE meshaplus TO meshaplus;

-- Se connecter à la base de données
\c meshaplus

-- Donner les privilèges sur le schéma public
GRANT ALL ON SCHEMA public TO meshaplus;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO meshaplus;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO meshaplus;

-- Configurer les privilèges par défaut
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO meshaplus;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO meshaplus;

\q
EOF

    log_success "PostgreSQL configuré"
}

# Installation de Redis
install_redis() {
    log "🔴 Installation de Redis..."
    
    # Installer Redis
    sudo apt install -y redis-server
    
    # Démarrer et activer le service
    sudo systemctl start redis-server
    sudo systemctl enable redis-server
    
    # Vérifier le statut
    if sudo systemctl is-active --quiet redis-server; then
        log_success "Redis installé et démarré"
    else
        log_error "Échec du démarrage de Redis"
        exit 1
    fi
}

# Configuration de Redis
configure_redis() {
    log "⚙️ Configuration de Redis..."
    
    # Sauvegarder la configuration originale
    sudo cp /etc/redis/redis.conf /etc/redis/redis.conf.backup
    
    # Modifier la configuration pour la production
    sudo sed -i 's/bind 127.0.0.1/bind 127.0.0.1/' /etc/redis/redis.conf
    sudo sed -i 's/# maxmemory <bytes>/maxmemory 256mb/' /etc/redis/redis.conf
    sudo sed -i 's/# maxmemory-policy noeviction/maxmemory-policy allkeys-lru/' /etc/redis/redis.conf
    
    # Redémarrer Redis
    sudo systemctl restart redis-server
    
    # Tester la connexion
    if redis-cli ping | grep -q "PONG"; then
        log_success "Redis configuré et fonctionnel"
    else
        log_error "Échec de la configuration de Redis"
        exit 1
    fi
}

# Installation de Nginx (optionnel)
install_nginx() {
    log "🌐 Installation de Nginx..."
    
    # Installer Nginx
    sudo apt install -y nginx
    
    # Démarrer et activer le service
    sudo systemctl start nginx
    sudo systemctl enable nginx
    
    # Configurer le firewall
    sudo ufw allow 'Nginx Full'
    
    # Vérifier le statut
    if sudo systemctl is-active --quiet nginx; then
        log_success "Nginx installé et démarré"
    else
        log_error "Échec du démarrage de Nginx"
        exit 1
    fi
}

# Configuration du firewall
configure_firewall() {
    log "🔥 Configuration du firewall..."
    
    # Installer UFW si pas déjà installé
    if ! command -v ufw &> /dev/null; then
        sudo apt install -y ufw
    fi
    
    # Activer UFW
    sudo ufw --force enable
    
    # Autoriser SSH
    sudo ufw allow ssh
    
    # Autoriser HTTP et HTTPS
    sudo ufw allow 80
    sudo ufw allow 443
    
    # Autoriser le port de l'API (optionnel, car derrière Nginx)
    # sudo ufw allow 8080
    
    log_success "Firewall configuré"
}

# Installation des outils de monitoring
install_monitoring_tools() {
    log "📊 Installation des outils de monitoring..."
    
    # Installer htop pour le monitoring système
    sudo apt install -y htop
    
    # Installer net-tools pour les commandes réseau
    sudo apt install -y net-tools
    
    # Installer curl pour les tests
    sudo apt install -y curl
    
    log_success "Outils de monitoring installés"
}

# Test des services
test_services() {
    log "🧪 Test des services..."
    
    # Test PostgreSQL
    if pg_isready -h localhost -p 5432 >/dev/null 2>&1; then
        log_success "PostgreSQL est accessible"
    else
        log_error "PostgreSQL n'est pas accessible"
        return 1
    fi
    
    # Test Redis
    if redis-cli ping | grep -q "PONG"; then
        log_success "Redis est accessible"
    else
        log_error "Redis n'est pas accessible"
        return 1
    fi
    
    # Test Nginx
    if curl -f http://localhost >/dev/null 2>&1; then
        log_success "Nginx est accessible"
    else
        log_warning "Nginx n'est pas accessible (normal si pas encore configuré)"
    fi
    
    log_success "Tous les services testés"
}

# Affichage des informations de configuration
show_configuration_info() {
    log "📋 Informations de configuration..."
    
    echo ""
    echo "🎉 Installation des services terminée !"
    echo ""
    echo "📋 Services installés :"
    echo "• PostgreSQL: localhost:5432"
    echo "  - Base de données: meshaplus"
    echo "  - Utilisateur: meshaplus"
    echo "  - Mot de passe: meshaplus_password"
    echo ""
    echo "• Redis: localhost:6379"
    echo "  - Pas de mot de passe"
    echo "  - Mémoire max: 256MB"
    echo ""
    echo "• Nginx: localhost:80"
    echo "  - Prêt pour la configuration de proxy"
    echo ""
    echo "🔍 Commandes utiles :"
    echo "• Statut PostgreSQL: sudo systemctl status postgresql"
    echo "• Statut Redis: sudo systemctl status redis-server"
    echo "• Statut Nginx: sudo systemctl status nginx"
    echo "• Logs PostgreSQL: sudo tail -f /var/log/postgresql/postgresql-*.log"
    echo "• Logs Redis: sudo tail -f /var/log/redis/redis-server.log"
    echo "• Logs Nginx: sudo tail -f /var/log/nginx/access.log"
    echo ""
    echo "⚠️  IMPORTANT :"
    echo "• Changez le mot de passe PostgreSQL en production"
    echo "• Configurez Redis avec un mot de passe en production"
    echo "• Mettez à jour la configuration de l'application"
    echo ""
}

# Fonction principale
main() {
    log "🔧 Début de l'installation des services locaux"
    echo "============================================"
    
    # Mettre à jour le système
    update_system
    
    # Installer PostgreSQL
    install_postgresql
    configure_postgresql
    
    # Installer Redis
    install_redis
    configure_redis
    
    # Installer Nginx
    install_nginx
    
    # Configurer le firewall
    configure_firewall
    
    # Installer les outils de monitoring
    install_monitoring_tools
    
    # Tester les services
    if test_services; then
        # Afficher les informations
        show_configuration_info
        
        log_success "Installation terminée avec succès !"
        exit 0
    else
        log_error "Échec des tests des services"
        exit 1
    fi
}

# Gestion des erreurs
trap 'log_error "Erreur survenue à la ligne $LINENO"; exit 1' ERR

# Exécution du script principal
main "$@" 