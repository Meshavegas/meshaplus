#!/bin/bash

# Script pour créer le répertoire de déploiement sur le VPS
# Usage: ./setup-deploy-directory.sh [host] [username] [deploy_path]

set -e

# Couleurs pour l'affichage
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Fonction pour afficher les messages
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Paramètres
HOST=${1:-""}
USERNAME=${2:-""}
DEPLOY_PATH=${3:-""}
KEY_PATH=${4:-"~/.ssh/meshaplus_deploy_key"}

# Vérification des paramètres
if [ -z "$HOST" ] || [ -z "$USERNAME" ]; then
    log_error "Usage: $0 <host> <username> [deploy_path] [key_path]"
    log_error "Example: $0 192.168.1.100 ubuntu /home/ubuntu/meshaplus"
    exit 1
fi

# Chemin de déploiement par défaut
if [ -z "$DEPLOY_PATH" ]; then
    DEPLOY_PATH="/home/$USERNAME/meshaplus"
fi

KEY_PATH=$(eval echo "$KEY_PATH")

log_info "📁 Configuration du répertoire de déploiement"
log_info "Host: $HOST"
log_info "Username: $USERNAME"
log_info "Deploy Path: $DEPLOY_PATH"
log_info "SSH Key: $KEY_PATH"
echo

# 1. Vérifier que la clé SSH existe
if [ ! -f "$KEY_PATH" ]; then
    log_error "Clé SSH introuvable: $KEY_PATH"
    log_info "Générez d'abord une clé avec: ./scripts/generate-deploy-key.sh"
    exit 1
fi

log_success "Clé SSH trouvée: $KEY_PATH"
echo

# 2. Tester la connexion SSH
log_info "1. Test de connexion SSH..."
if ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" "echo 'Connexion SSH réussie'"; then
    log_success "Connexion SSH réussie"
else
    log_error "Échec de la connexion SSH"
    log_info "Vérifiez que la clé publique est ajoutée au VPS"
    exit 1
fi
echo

# 3. Créer le répertoire de déploiement
log_info "2. Création du répertoire de déploiement..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << EOF
    echo "📁 Création du répertoire: $DEPLOY_PATH"
    mkdir -p "$DEPLOY_PATH"
    
    echo "📋 Vérification du répertoire:"
    ls -la "$(dirname "$DEPLOY_PATH")"
    
    echo "📋 Contenu du répertoire de déploiement:"
    ls -la "$DEPLOY_PATH"
    
    echo "🔧 Configuration des permissions:"
    chmod 755 "$DEPLOY_PATH"
    
    echo "📊 Espace disque disponible:"
    df -h "$(dirname "$DEPLOY_PATH")"
EOF

log_success "Répertoire de déploiement créé: $DEPLOY_PATH"
echo

# 4. Vérifier l'environnement Docker
log_info "3. Vérification de l'environnement Docker..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🐳 Vérification de Docker:"
    docker --version || echo "❌ Docker non installé"
    
    echo "🐳 Vérification de Docker Compose:"
    docker-compose --version || echo "❌ Docker Compose non installé"
    
    echo "🐳 Vérification des conteneurs en cours:"
    docker ps -a || echo "❌ Impossible de lister les conteneurs"
    
    echo "📊 Utilisation des ressources:"
    free -h
    df -h /
EOF

echo

# 5. Test de copie de fichier
log_info "4. Test de copie de fichier..."
echo "Test file content" > /tmp/test_deploy.txt
if scp -o StrictHostKeyChecking=no -i "$KEY_PATH" /tmp/test_deploy.txt "$USERNAME@$HOST:$DEPLOY_PATH/"; then
    log_success "Test de copie réussi"
    ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" "rm -f $DEPLOY_PATH/test_deploy.txt"
else
    log_error "Test de copie échoué"
    exit 1
fi
rm -f /tmp/test_deploy.txt
echo

# 6. Instructions finales
log_info "📝 Instructions finales:"
echo
log_info "✅ Le répertoire de déploiement est prêt: $DEPLOY_PATH"
log_info "🔧 Pour configurer GitHub Secrets, utilisez:"
echo
log_info "VPS_DEPLOY_PATH: $DEPLOY_PATH"
echo
log_info "🚀 Vous pouvez maintenant relancer le workflow GitHub Actions"
echo

log_success "🎉 Configuration du répertoire de déploiement terminée!"
log_info "Le VPS est prêt pour recevoir les fichiers de déploiement" 