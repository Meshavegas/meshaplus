#!/bin/bash

# Script pour supprimer la phrase de passe d'une clé SSH existante
# Usage: ./remove-ssh-passphrase.sh [key_path]

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
KEY_PATH=${1:-"~/.ssh/vps_key"}
KEY_PATH=$(eval echo "$KEY_PATH")

log_info "🔓 Suppression de la phrase de passe SSH"
log_info "Clé: $KEY_PATH"
echo

# 1. Vérifier que la clé existe
if [ ! -f "$KEY_PATH" ]; then
    log_error "Clé SSH introuvable: $KEY_PATH"
    log_info "Usage: $0 [chemin_vers_cle_ssh]"
    exit 1
fi

log_success "Clé SSH trouvée: $KEY_PATH"
echo

# 2. Vérifier si la clé a une phrase de passe
log_info "1. Vérification de la phrase de passe..."
if ssh-keygen -y -f "$KEY_PATH" > /dev/null 2>&1; then
    log_success "La clé n'a pas de phrase de passe"
    log_info "Aucune action nécessaire"
    exit 0
else
    log_warning "La clé a une phrase de passe"
fi
echo

# 3. Créer une sauvegarde
log_info "2. Création d'une sauvegarde..."
BACKUP_PATH="${KEY_PATH}.backup.$(date +%Y%m%d_%H%M%S)"
cp "$KEY_PATH" "$BACKUP_PATH"
log_success "Sauvegarde créée: $BACKUP_PATH"
echo

# 4. Supprimer la phrase de passe
log_info "3. Suppression de la phrase de passe..."
log_warning "⚠️  Vous devrez saisir la phrase de passe actuelle"
log_info "Entrez la phrase de passe actuelle de la clé:"

# Créer une nouvelle clé sans phrase de passe
ssh-keygen -p -f "$KEY_PATH" -N ""
log_success "Phrase de passe supprimée"
echo

# 5. Vérifier le résultat
log_info "4. Vérification du résultat..."
if ssh-keygen -y -f "$KEY_PATH" > /dev/null 2>&1; then
    log_success "✅ La clé n'a plus de phrase de passe"
else
    log_error "❌ Erreur lors de la suppression de la phrase de passe"
    log_info "Restauration de la sauvegarde..."
    cp "$BACKUP_PATH" "$KEY_PATH"
    exit 1
fi
echo

# 6. Afficher les informations
log_info "5. Informations de la clé mise à jour"
echo
log_info "📋 Clé publique:"
echo "=== CLÉ PUBLIQUE ==="
ssh-keygen -y -f "$KEY_PATH"
echo "==================="
echo

log_info "🔐 Clé privée (pour GitHub Secrets VPS_SSH_KEY):"
echo "=== CLÉ PRIVÉE ==="
cat "$KEY_PATH"
echo "=================="
echo

# 7. Instructions
log_info "📝 Instructions:"
echo
log_info "1. 🔑 Vérifiez que la clé publique est toujours sur le VPS"
log_info "2. 🔐 Mettez à jour le secret GitHub VPS_SSH_KEY avec la nouvelle clé privée"
log_info "3. 🧪 Testez la connexion: ssh -i $KEY_PATH user@vps-host"
log_info "4. 🗑️  Supprimez la sauvegarde si tout fonctionne: rm $BACKUP_PATH"
echo

log_success "🎉 Phrase de passe supprimée avec succès!"
log_info "La clé est maintenant prête pour le déploiement automatique" 