#!/bin/bash

# Script de diagnostic SSH pour MeshaPlus
# Usage: ./ssh-diagnostic.sh [host] [username] [ssh_key_path]

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
HOST=${1:-"your-vps-host"}
USERNAME=${2:-"your-username"}
SSH_KEY_PATH=${3:-"~/.ssh/vps_key"}

log_info "🔍 Diagnostic SSH pour MeshaPlus"
log_info "Host: $HOST"
log_info "Username: $USERNAME"
log_info "SSH Key: $SSH_KEY_PATH"
echo

# 1. Vérifier que la clé SSH existe
log_info "1. Vérification de la clé SSH..."
if [ -f "$SSH_KEY_PATH" ]; then
    log_success "Clé SSH trouvée: $SSH_KEY_PATH"
    log_info "Permissions: $(ls -la "$SSH_KEY_PATH")"
    log_info "Taille: $(wc -c < "$SSH_KEY_PATH") bytes"
    
    # Vérifier le format de la clé
    if grep -q "BEGIN.*PRIVATE KEY" "$SSH_KEY_PATH"; then
        log_success "Format de clé privée détecté"
    else
        log_warning "Format de clé non reconnu"
    fi
else
    log_error "Clé SSH introuvable: $SSH_KEY_PATH"
    exit 1
fi
echo

# 2. Vérifier les permissions de la clé
log_info "2. Vérification des permissions..."
PERMISSIONS=$(stat -c "%a" "$SSH_KEY_PATH")
if [ "$PERMISSIONS" = "600" ]; then
    log_success "Permissions correctes: $PERMISSIONS"
else
    log_warning "Permissions incorrectes: $PERMISSIONS (devrait être 600)"
    log_info "Correction des permissions..."
    chmod 600 "$SSH_KEY_PATH"
    log_success "Permissions corrigées"
fi
echo

# 3. Test de connectivité réseau
log_info "3. Test de connectivité réseau..."
if ping -c 1 "$HOST" > /dev/null 2>&1; then
    log_success "Host accessible via ping"
else
    log_warning "Host non accessible via ping (peut être normal si ICMP est bloqué)"
fi

# Test de connectivité sur le port SSH
if nc -z "$HOST" 22 2>/dev/null; then
    log_success "Port SSH (22) accessible"
else
    log_error "Port SSH (22) non accessible"
    exit 1
fi
echo

# 4. Test de connexion SSH basique
log_info "4. Test de connexion SSH..."
if ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=no -i "$SSH_KEY_PATH" "$USERNAME@$HOST" "echo 'Connexion SSH réussie'" 2>/dev/null; then
    log_success "Connexion SSH réussie!"
else
    log_error "Échec de la connexion SSH"
    echo
    log_info "🔍 Tentative de diagnostic détaillé..."
    
    # Test avec verbose
    log_info "Test avec mode verbose:"
    ssh -v -o ConnectTimeout=10 -o StrictHostKeyChecking=no -i "$SSH_KEY_PATH" "$USERNAME@$HOST" "echo 'test'" 2>&1 | head -20
    
    echo
    log_info "🔧 Solutions possibles:"
    log_info "1. Vérifiez que la clé publique est ajoutée à ~/.ssh/authorized_keys sur le VPS"
    log_info "2. Vérifiez que l'utilisateur $USERNAME existe sur le VPS"
    log_info "3. Vérifiez que SSH est configuré pour accepter les clés publiques"
    log_info "4. Vérifiez les logs SSH sur le VPS: sudo journalctl -u ssh"
    exit 1
fi
echo

# 5. Test de commandes sur le VPS
log_info "5. Test des commandes sur le VPS..."
log_info "Vérification de l'utilisateur:"
ssh -o StrictHostKeyChecking=no -i "$SSH_KEY_PATH" "$USERNAME@$HOST" "whoami && pwd && id"

log_info "Vérification de Docker:"
if ssh -o StrictHostKeyChecking=no -i "$SSH_KEY_PATH" "$USERNAME@$HOST" "docker --version" 2>/dev/null; then
    log_success "Docker installé"
else
    log_warning "Docker non installé ou non accessible"
fi

log_info "Vérification de l'espace disque:"
ssh -o StrictHostKeyChecking=no -i "$SSH_KEY_PATH" "$USERNAME@$HOST" "df -h /"

echo
log_success "🎉 Diagnostic SSH terminé avec succès!"
log_info "Le VPS est prêt pour le déploiement MeshaPlus" 