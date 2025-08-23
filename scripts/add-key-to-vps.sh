#!/bin/bash

# Script pour ajouter automatiquement la clé publique au VPS
# Usage: ./add-key-to-vps.sh [host] [username] [key_path]

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
PUBLIC_KEY_PATH="${KEY_PATH}.pub"

log_info "🔑 Ajout de la clé publique au VPS"
log_info "Host: $HOST"
log_info "Username: $USERNAME"
log_info "Key: $KEY_PATH"
echo

# 1. Vérifier que la clé publique existe
if [ ! -f "$PUBLIC_KEY_PATH" ]; then
    log_error "Clé publique introuvable: $PUBLIC_KEY_PATH"
    exit 1
fi

log_success "Clé publique trouvée: $PUBLIC_KEY_PATH"
echo

# 2. Lire la clé publique
PUBLIC_KEY=$(cat "$PUBLIC_KEY_PATH")
log_info "Clé publique à ajouter:"
echo "$PUBLIC_KEY"
echo

# 3. Demander confirmation
log_warning "⚠️  Cette action va ajouter la clé publique au VPS"
read -p "Voulez-vous continuer? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_info "Action annulée"
    exit 0
fi

# 4. Ajouter la clé au VPS
log_info "1. Ajout de la clé publique au VPS..."

# Méthode 1: Utiliser ssh-copy-id si disponible
if command -v ssh-copy-id &> /dev/null; then
    log_info "Utilisation de ssh-copy-id..."
    ssh-copy-id -i "$PUBLIC_KEY_PATH" "$USERNAME@$HOST" || {
        log_warning "ssh-copy-id a échoué, tentative manuelle..."
        # Méthode 2: Ajout manuel
        ssh "$USERNAME@$HOST" "mkdir -p ~/.ssh && echo '$PUBLIC_KEY' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && chmod 700 ~/.ssh"
    }
else
    log_info "ssh-copy-id non disponible, ajout manuel..."
    ssh "$USERNAME@$HOST" "mkdir -p ~/.ssh && echo '$PUBLIC_KEY' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && chmod 700 ~/.ssh"
fi

log_success "Clé publique ajoutée au VPS"
echo

# 5. Tester la connexion avec la nouvelle clé
log_info "2. Test de connexion avec la nouvelle clé..."
if ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" "echo 'Connexion SSH réussie avec la nouvelle clé!'"; then
    log_success "✅ Connexion SSH réussie!"
else
    log_error "❌ Échec de la connexion SSH"
    log_info "Vérifiez que:"
    log_info "1. La clé privée correspond à la clé publique"
    log_info "2. L'utilisateur $USERNAME existe sur le VPS"
    log_info "3. SSH est configuré pour accepter les clés publiques"
    exit 1
fi
echo

# 6. Vérifier la configuration sur le VPS
log_info "3. Vérification de la configuration sur le VPS..."
ssh -o StrictHostKeyChecking=no -i "$KEY_PATH" "$USERNAME@$HOST" << 'EOF'
echo "📋 Contenu de ~/.ssh/authorized_keys:"
cat ~/.ssh/authorized_keys
echo
echo "📋 Permissions:"
ls -la ~/.ssh/
echo
echo "📋 Configuration SSH:"
sudo sshd -T | grep -E "(pubkey|authorized)" || echo "Impossible de vérifier la configuration SSH"
EOF

echo
log_success "🎉 Configuration SSH terminée avec succès!"
log_info "Le VPS est maintenant prêt pour le déploiement automatique" 