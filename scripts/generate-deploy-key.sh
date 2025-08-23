#!/bin/bash

# Script pour générer une clé SSH de déploiement sans phrase de passe
# Usage: ./generate-deploy-key.sh [key_name]

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
KEY_NAME=${1:-"meshaplus_deploy_key"}

log_info "🔑 Génération d'une clé SSH de déploiement"
log_info "Nom de la clé: $KEY_NAME"
echo

# 1. Vérifier si la clé existe déjà
if [ -f "$HOME/.ssh/$KEY_NAME" ]; then
    log_warning "Clé SSH existante détectée: $HOME/.ssh/$KEY_NAME"
    read -p "Voulez-vous la remplacer? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Génération annulée"
        exit 0
    fi
    rm -f "$HOME/.ssh/$KEY_NAME" "$HOME/.ssh/$KEY_NAME.pub"
fi

# 2. Générer la nouvelle clé SSH sans phrase de passe
log_info "1. Génération de la clé SSH..."
ssh-keygen -t ed25519 -f "$HOME/.ssh/$KEY_NAME" -N "" -C "meshaplus-deployment-$(date +%Y%m%d)"
log_success "Clé SSH générée: $HOME/.ssh/$KEY_NAME"
echo

# 3. Configurer les permissions
log_info "2. Configuration des permissions..."
chmod 600 "$HOME/.ssh/$KEY_NAME"
chmod 644 "$HOME/.ssh/$KEY_NAME.pub"
log_success "Permissions configurées"
echo

# 4. Afficher les informations
log_info "3. Informations de la clé générée"
echo
log_info "📋 Clé publique (à ajouter au VPS):"
echo "=== CLÉ PUBLIQUE ==="
cat "$HOME/.ssh/$KEY_NAME.pub"
echo "==================="
echo

log_info "🔐 Clé privée (pour GitHub Secrets VPS_SSH_KEY):"
echo "=== CLÉ PRIVÉE ==="
cat "$HOME/.ssh/$KEY_NAME"
echo "=================="
echo

# 5. Instructions
log_info "📝 Instructions de configuration:"
echo
log_info "1. 🔑 Ajoutez la clé publique au VPS:"
log_info "   - Connectez-vous à votre VPS"
log_info "   - Ajoutez la clé publique ci-dessus à ~/.ssh/authorized_keys"
log_info "   - Ou utilisez: ssh-copy-id -i $HOME/.ssh/$KEY_NAME.pub user@vps-host"
echo
log_info "2. 🔐 Configurez GitHub Secrets:"
log_info "   - Allez dans votre repository GitHub"
log_info "   - Settings > Secrets and variables > Actions"
log_info "   - Ajoutez le secret VPS_SSH_KEY avec la clé privée ci-dessus"
echo
log_info "3. 🧪 Testez la connexion:"
log_info "   - Utilisez: ssh -i $HOME/.ssh/$KEY_NAME user@vps-host"
echo

log_success "🎉 Clé SSH de déploiement générée avec succès!"
log_info "La clé est prête pour le déploiement automatique via GitHub Actions" 