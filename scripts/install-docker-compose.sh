#!/bin/bash

# Script pour installer Docker Compose sur le VPS
# Usage: ./install-docker-compose.sh [host] [username] [key_path]

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
KEY_PATH=${3:-"~/.ssh/meshaplus_deploy_key"}

# Vérification des paramètres
if [ -z "$HOST" ] || [ -z "$USERNAME" ]; then
    log_error "Usage: $0 <host> <username> [key_path]"
    log_error "Example: $0 192.168.1.100 ubuntu ~/.ssh/meshaplus_deploy_key"
    exit 1
fi

KEY_PATH=$(eval echo "$KEY_PATH")

log_info "🐳 Installation de Docker Compose sur le VPS"
log_info "Host: $HOST"
log_info "Username: $USERNAME"
log_info "SSH Key: $KEY_PATH"
echo

# 1. Vérifier que la clé SSH existe
if [ ! -f "$KEY_PATH" ]; then
    log_error "Clé SSH introuvable: $KEY_PATH"
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
    exit 1
fi
echo

# 3. Vérifier l'état actuel de Docker
log_info "2. Vérification de l'état actuel de Docker..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🐳 Vérification de Docker:"
    docker --version || echo "❌ Docker non installé"
    
    echo "🐳 Vérification de Docker Compose:"
    docker-compose --version || echo "❌ Docker Compose non installé"
    
    echo "🐳 Vérification de Docker Buildx:"
    docker buildx version || echo "❌ Docker Buildx non installé"
    
    echo "📊 Espace disque disponible:"
    df -h /
    
    echo "📊 Mémoire disponible:"
    free -h
EOF
echo

# 4. Installer Docker Compose
log_info "3. Installation de Docker Compose..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🔧 Mise à jour des paquets..."
    sudo apt update
    
    echo "🔧 Installation de Docker Compose..."
    # Supprimer l'ancienne version si elle existe
    sudo rm -f /usr/local/bin/docker-compose
    
    # Télécharger la dernière version de Docker Compose
    DOCKER_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
    echo "📦 Version Docker Compose: $DOCKER_COMPOSE_VERSION"
    
    sudo curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    
    # Rendre exécutable
    sudo chmod +x /usr/local/bin/docker-compose
    
    # Créer un lien symbolique
    sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
    
    echo "✅ Docker Compose installé"
EOF
echo

# 5. Installer Docker Buildx
log_info "4. Installation de Docker Buildx..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🔧 Installation de Docker Buildx..."
    
    # Vérifier si buildx est déjà installé
    if docker buildx version > /dev/null 2>&1; then
        echo "✅ Docker Buildx déjà installé"
    else
        # Installer buildx
        mkdir -p ~/.docker/cli-plugins
        curl -SL https://github.com/docker/buildx/releases/download/v0.12.1/buildx-v0.12.1.linux-amd64 -o ~/.docker/cli-plugins/docker-buildx
        chmod +x ~/.docker/cli-plugins/docker-buildx
        
        # Créer une nouvelle instance buildx
        docker buildx create --name meshaplus-builder --use
        docker buildx inspect --bootstrap
        
        echo "✅ Docker Buildx installé et configuré"
    fi
EOF
echo

# 6. Configurer Docker pour éviter les rate limits
log_info "5. Configuration de Docker pour éviter les rate limits..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🔧 Configuration de Docker..."
    
    # Créer le répertoire de configuration Docker
    sudo mkdir -p /etc/docker
    
    # Créer ou modifier le fichier daemon.json
    sudo tee /etc/docker/daemon.json > /dev/null << 'DOCKER_CONFIG'
{
  "registry-mirrors": [
    "https://mirror.gcr.io",
    "https://registry-1.docker.io"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
DOCKER_CONFIG
    
    # Redémarrer Docker
    sudo systemctl restart docker
    
    echo "✅ Docker configuré et redémarré"
EOF
echo

# 7. Vérification finale
log_info "6. Vérification finale..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
    echo "🐳 Vérification de Docker:"
    docker --version
    
    echo "🐳 Vérification de Docker Compose:"
    docker-compose --version
    
    echo "🐳 Vérification de Docker Buildx:"
    docker buildx version
    
    echo "🐳 Test de construction d'image simple:"
    docker buildx build --platform linux/amd64 -t test-image:latest . << 'DOCKERFILE'
FROM alpine:latest
RUN echo "Test image built successfully"
DOCKERFILE
    
    echo "✅ Tous les composants Docker sont fonctionnels"
EOF
echo

log_success "🎉 Installation de Docker Compose terminée!"
log_info "Le VPS est maintenant prêt pour le déploiement MeshaPlus" 