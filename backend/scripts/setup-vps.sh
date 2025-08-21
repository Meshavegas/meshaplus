#!/bin/bash

# 🚀 Script de Configuration VPS Ubuntu - MeshaPlus Backend
# Ce script configure un VPS Ubuntu pour le déploiement automatique

set -e

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables de configuration
DOCKER_COMPOSE_PATH="/opt/meshaplus"
NGINX_CONF_PATH="/etc/nginx/sites-available"
NGINX_ENABLED_PATH="/etc/nginx/sites-enabled"
SSL_PATH="/etc/nginx/ssl"

echo -e "${BLUE}🚀 Configuration VPS Ubuntu - MeshaPlus Backend${NC}"
echo "======================================================"
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
    
    # Vérifier que nous sommes sur Ubuntu
    if ! grep -q "Ubuntu" /etc/os-release; then
        print_error "Ce script est conçu pour Ubuntu"
        exit 1
    fi
    
    print_info "Système: Ubuntu $(lsb_release -rs)"
    
    # Vérifier les permissions root
    if [ "$EUID" -ne 0 ]; then
        print_error "Ce script doit être exécuté en tant que root (sudo)"
        exit 1
    fi
    
    print_success "Prérequis vérifiés"
}

# Mise à jour du système
update_system() {
    print_step "Mise à jour du système..."
    
    apt update
    apt upgrade -y
    
    print_success "Système mis à jour"
}

# Installation de Docker
install_docker() {
    print_step "Installation de Docker..."
    
    # Supprimer les anciennes versions
    apt remove -y docker docker-engine docker.io containerd runc 2>/dev/null || true
    
    # Installer les dépendances
    apt install -y apt-transport-https ca-certificates curl gnupg lsb-release
    
    # Ajouter la clé GPG officielle de Docker
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    
    # Ajouter le repository Docker
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Installer Docker
    apt update
    apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
    
    # Démarrer et activer Docker
    systemctl start docker
    systemctl enable docker
    
    # Ajouter l'utilisateur au groupe docker
    usermod -aG docker $SUDO_USER
    
    print_success "Docker installé"
}

# Installation de Nginx
install_nginx() {
    print_step "Installation de Nginx..."
    
    apt install -y nginx
    
    # Démarrer et activer Nginx
    systemctl start nginx
    systemctl enable nginx
    
    # Configurer le firewall
    ufw allow 'Nginx Full'
    ufw allow OpenSSH
    
    print_success "Nginx installé"
}

# Installation de Certbot (SSL)
install_certbot() {
    print_step "Installation de Certbot pour SSL..."
    
    apt install -y certbot python3-certbot-nginx
    
    print_success "Certbot installé"
}

# Configuration des répertoires
setup_directories() {
    print_step "Configuration des répertoires..."
    
    # Créer le répertoire principal
    mkdir -p $DOCKER_COMPOSE_PATH
    mkdir -p $SSL_PATH
    
    # Créer les répertoires pour les données
    mkdir -p /var/lib/meshaplus/postgres
    mkdir -p /var/lib/meshaplus/redis
    
    # Définir les permissions
    chown -R $SUDO_USER:$SUDO_USER $DOCKER_COMPOSE_PATH
    chown -R $SUDO_USER:$SUDO_USER /var/lib/meshaplus
    
    print_success "Répertoires configurés"
}

# Configuration de Nginx
setup_nginx() {
    print_step "Configuration de Nginx..."
    
    # Créer la configuration par défaut
    cat > /etc/nginx/sites-available/meshaplus << 'EOF'
server {
    listen 80;
    server_name _;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location /swagger/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF
    
    # Activer le site
    ln -sf /etc/nginx/sites-available/meshaplus /etc/nginx/sites-enabled/
    
    # Supprimer la configuration par défaut
    rm -f /etc/nginx/sites-enabled/default
    
    # Tester la configuration
    nginx -t
    
    # Redémarrer Nginx
    systemctl reload nginx
    
    print_success "Nginx configuré"
}

# Configuration du firewall
setup_firewall() {
    print_step "Configuration du firewall..."
    
    # Activer UFW
    ufw --force enable
    
    # Règles de base
    ufw default deny incoming
    ufw default allow outgoing
    ufw allow ssh
    ufw allow 'Nginx Full'
    
    print_success "Firewall configuré"
}

# Configuration des logs
setup_logging() {
    print_step "Configuration des logs..."
    
    # Créer le répertoire de logs
    mkdir -p /var/log/meshaplus
    
    # Configuration logrotate
    cat > /etc/logrotate.d/meshaplus << 'EOF'
/var/log/meshaplus/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 root root
}
EOF
    
    print_success "Logs configurés"
}

# Configuration des sauvegardes
setup_backups() {
    print_step "Configuration des sauvegardes..."
    
    # Créer le répertoire de sauvegardes
    mkdir -p /var/backups/meshaplus
    
    # Script de sauvegarde
    cat > /usr/local/bin/backup-meshaplus.sh << 'EOF'
#!/bin/bash

BACKUP_DIR="/var/backups/meshaplus"
DATE=$(date +%Y%m%d_%H%M%S)

# Sauvegarder la base de données
docker exec meshaplus-postgres pg_dump -U postgres meshaplus > $BACKUP_DIR/db_backup_$DATE.sql

# Sauvegarder les fichiers de configuration
tar -czf $BACKUP_DIR/config_backup_$DATE.tar.gz -C /opt/meshaplus .

# Supprimer les sauvegardes de plus de 30 jours
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete

echo "Sauvegarde terminée: $DATE"
EOF
    
    chmod +x /usr/local/bin/backup-meshaplus.sh
    
    # Ajouter au crontab (sauvegarde quotidienne à 2h du matin)
    (crontab -l 2>/dev/null; echo "0 2 * * * /usr/local/bin/backup-meshaplus.sh") | crontab -
    
    print_success "Sauvegardes configurées"
}

# Configuration du monitoring
setup_monitoring() {
    print_step "Configuration du monitoring..."
    
    # Installer htop pour le monitoring
    apt install -y htop
    
    # Script de monitoring
    cat > /usr/local/bin/monitor-meshaplus.sh << 'EOF'
#!/bin/bash

echo "=== MeshaPlus Monitoring ==="
echo "Date: $(date)"
echo ""

echo "=== Docker Containers ==="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""

echo "=== System Resources ==="
echo "CPU: $(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1)%"
echo "Memory: $(free -m | awk 'NR==2{printf "%.1f%%", $3*100/$2}')"
echo "Disk: $(df -h / | awk 'NR==2{print $5}')"
echo ""

echo "=== Application Health ==="
curl -s http://localhost:8080/api/v1/health || echo "Application not responding"
echo ""
EOF
    
    chmod +x /usr/local/bin/monitor-meshaplus.sh
    
    print_success "Monitoring configuré"
}

# Configuration des variables d'environnement
setup_environment() {
    print_step "Configuration des variables d'environnement..."
    
    # Créer le fichier .env
    cat > $DOCKER_COMPOSE_PATH/.env << 'EOF'
# Configuration de l'environnement
NODE_ENV=production

# Base de données
POSTGRES_DB=meshaplus
POSTGRES_USER=meshaplus_user
POSTGRES_PASSWORD=change_this_password

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=change_this_jwt_secret

# Application
APP_PORT=8080
APP_HOST=0.0.0.0
EOF
    
    print_info "Fichier .env créé à $DOCKER_COMPOSE_PATH/.env"
    print_info "⚠️  N'oubliez pas de modifier les mots de passe !"
    
    print_success "Variables d'environnement configurées"
}

# Configuration des scripts utilitaires
setup_utilities() {
    print_step "Configuration des scripts utilitaires..."
    
    # Script de redémarrage
    cat > /usr/local/bin/restart-meshaplus.sh << 'EOF'
#!/bin/bash

cd /opt/meshaplus
docker-compose down
docker-compose up -d
echo "MeshaPlus redémarré"
EOF
    
    chmod +x /usr/local/bin/restart-meshaplus.sh
    
    # Script de logs
    cat > /usr/local/bin/logs-meshaplus.sh << 'EOF'
#!/bin/bash

cd /opt/meshaplus
docker-compose logs -f
EOF
    
    chmod +x /usr/local/bin/logs-meshaplus.sh
    
    # Script de mise à jour
    cat > /usr/local/bin/update-meshaplus.sh << 'EOF'
#!/bin/bash

cd /opt/meshaplus
docker-compose pull
docker-compose up -d
echo "MeshaPlus mis à jour"
EOF
    
    chmod +x /usr/local/bin/update-meshaplus.sh
    
    print_success "Scripts utilitaires configurés"
}

# Affichage des informations finales
show_final_info() {
    print_step "Configuration terminée !"
    echo ""
    echo -e "${GREEN}🎉 VPS Ubuntu configuré avec succès !${NC}"
    echo ""
    echo -e "${BLUE}📋 Informations importantes :${NC}"
    echo "• Répertoire principal: $DOCKER_COMPOSE_PATH"
    echo "• Configuration Nginx: /etc/nginx/sites-available/meshaplus"
    echo "• Logs: /var/log/meshaplus"
    echo "• Sauvegardes: /var/backups/meshaplus"
    echo ""
    echo -e "${BLUE}🔧 Scripts disponibles :${NC}"
    echo "• /usr/local/bin/restart-meshaplus.sh - Redémarrer l'application"
    echo "• /usr/local/bin/logs-meshaplus.sh - Voir les logs"
    echo "• /usr/local/bin/update-meshaplus.sh - Mettre à jour l'application"
    echo "• /usr/local/bin/monitor-meshaplus.sh - Monitoring"
    echo "• /usr/local/bin/backup-meshaplus.sh - Sauvegarde"
    echo ""
    echo -e "${BLUE}⚠️  Actions requises :${NC}"
    echo "1. Modifier le fichier .env avec vos vraies valeurs"
    echo "2. Configurer votre domaine dans Nginx"
    echo "3. Obtenir un certificat SSL avec Certbot"
    echo "4. Configurer les secrets GitHub Actions"
    echo ""
    echo -e "${BLUE}🔐 Sécurité :${NC}"
    echo "• Firewall UFW activé"
    echo "• Seuls les ports SSH et HTTP/HTTPS ouverts"
    echo "• Sauvegardes automatiques configurées"
    echo ""
    echo -e "${GREEN}🚀 Votre VPS est prêt pour le déploiement !${NC}"
}

# Fonction principale
main() {
    echo -e "${BLUE}🚀 Démarrage de la configuration VPS...${NC}"
    echo ""
    
    check_prerequisites
    echo ""
    
    update_system
    echo ""
    
    install_docker
    echo ""
    
    install_nginx
    echo ""
    
    install_certbot
    echo ""
    
    setup_directories
    echo ""
    
    setup_nginx
    echo ""
    
    setup_firewall
    echo ""
    
    setup_logging
    echo ""
    
    setup_backups
    echo ""
    
    setup_monitoring
    echo ""
    
    setup_environment
    echo ""
    
    setup_utilities
    echo ""
    
    show_final_info
}

# Exécution
main "$@" 