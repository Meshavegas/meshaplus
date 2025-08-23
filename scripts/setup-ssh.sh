#!/bin/bash

# Script de configuration SSH pour MeshaPlus
# Usage: ./setup-ssh.sh [host] [username] [key_name]

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
KEY_NAME=${3:-"meshaplus_key"}

# Vérification des paramètres
if [ -z "$HOST" ] || [ -z "$USERNAME" ]; then
    log_error "Usage: $0 <host> <username> [key_name]"
    log_error "Example: $0 192.168.1.100 ubuntu meshaplus_key"
    exit 1
fi

log_info "🔧 Configuration SSH pour MeshaPlus"
log_info "Host: $HOST"
log_info "Username: $USERNAME"
log_info "Key name: $KEY_NAME"
echo

# 1. Générer une nouvelle paire de clés SSH
log_info "1. Génération d'une nouvelle paire de clés SSH..."
if [ -f "$HOME/.ssh/$KEY_NAME" ]; then
    log_warning "Clé SSH existante détectée: $HOME/.ssh/$KEY_NAME"
    read -p "Voulez-vous la remplacer? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Configuration annulée"
        exit 0
    fi
    rm -f "$HOME/.ssh/$KEY_NAME" "$HOME/.ssh/$KEY_NAME.pub"
fi

ssh-keygen -t ed25519 -f "$HOME/.ssh/$KEY_NAME" -N "" -C "meshaplus-deployment-$HOST"
log_success "Paire de clés générée: $HOME/.ssh/$KEY_NAME"
echo

# 2. Configurer les permissions
log_info "2. Configuration des permissions..."
chmod 600 "$HOME/.ssh/$KEY_NAME"
chmod 644 "$HOME/.ssh/$KEY_NAME.pub"
log_success "Permissions configurées"
echo

# 3. Tester la connexion SSH (si possible)
log_info "3. Test de connexion SSH..."
if ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no -i "$HOME/.ssh/$KEY_NAME" "$USERNAME@$HOST" "echo 'Test de connexion'" 2>/dev/null; then
    log_success "Connexion SSH réussie avec la nouvelle clé!"
else
    log_warning "Connexion SSH échouée - la clé publique doit être ajoutée au VPS"
    echo
    log_info "🔧 Pour ajouter la clé publique au VPS:"
    log_info "1. Copiez la clé publique ci-dessous:"
    echo
    echo "=== CLÉ PUBLIQUE À COPIER ==="
    cat "$HOME/.ssh/$KEY_NAME.pub"
    echo "=============================="
    echo
    log_info "2. Connectez-vous au VPS et ajoutez-la à ~/.ssh/authorized_keys:"
    log_info "   ssh $USERNAME@$HOST"
    log_info "   mkdir -p ~/.ssh"
    log_info "   echo 'CLÉ_PUBLIQUE_CI_DESSUS' >> ~/.ssh/authorized_keys"
    log_info "   chmod 600 ~/.ssh/authorized_keys"
    echo
fi
echo

# 4. Afficher les informations pour GitHub Secrets
log_info "4. Configuration GitHub Secrets"
echo
log_info "📋 Ajoutez ces secrets dans GitHub (Settings > Secrets and variables > Actions):"
echo
log_info "VPS_HOST:"
echo "$HOST"
echo
log_info "VPS_USERNAME:"
echo "$USERNAME"
echo
log_info "VPS_SSH_KEY:"
echo "=== CLÉ PRIVÉE À COPIER ==="
cat "$HOME/.ssh/$KEY_NAME"
echo "============================="
echo

# 5. Créer un fichier de configuration SSH
log_info "5. Création du fichier de configuration SSH..."
SSH_CONFIG="$HOME/.ssh/config"
mkdir -p "$(dirname "$SSH_CONFIG")"

# Ajouter la configuration si elle n'existe pas déjà
if ! grep -q "Host meshaplus-$HOST" "$SSH_CONFIG" 2>/dev/null; then
    cat >> "$SSH_CONFIG" << EOF

# Configuration MeshaPlus pour $HOST
Host meshaplus-$HOST
    HostName $HOST
    User $USERNAME
    IdentityFile ~/.ssh/$KEY_NAME
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
EOF
    log_success "Configuration SSH ajoutée à $SSH_CONFIG"
    log_info "Vous pouvez maintenant vous connecter avec: ssh meshaplus-$HOST"
else
    log_warning "Configuration SSH déjà présente pour $HOST"
fi
echo

# 6. Test final
log_info "6. Test final de configuration..."
if [ -f "$SSH_CONFIG" ] && grep -q "Host meshaplus-$HOST" "$SSH_CONFIG"; then
    log_success "Configuration SSH complète!"
    log_info "Test de connexion avec la configuration:"
    if ssh -o ConnectTimeout=10 "meshaplus-$HOST" "echo 'Configuration SSH réussie!'" 2>/dev/null; then
        log_success "🎉 Configuration SSH terminée avec succès!"
    else
        log_warning "⚠️  Connexion échouée - vérifiez que la clé publique est ajoutée au VPS"
    fi
else
    log_error "❌ Configuration SSH incomplète"
fi
echo

# 7. Instructions finales
log_info "📝 Instructions finales:"
log_info "1. Copiez la clé privée ci-dessus dans GitHub Secrets (VPS_SSH_KEY)"
log_info "2. Ajoutez la clé publique au VPS si ce n'est pas déjà fait"
log_info "3. Testez la connexion: ssh meshaplus-$HOST"
log_info "4. Relancez le workflow GitHub Actions"
echo

log_success "🔧 Configuration SSH terminée!" 