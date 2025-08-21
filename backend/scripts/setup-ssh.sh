#!/bin/bash

# 🔑 Script de Configuration SSH pour VPS - MeshaPlus
# Ce script génère les clés SSH et guide l'utilisateur

set -e

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔑 Configuration SSH pour VPS MeshaPlus${NC}"
echo "============================================="
echo ""

# Variables
SSH_KEY_PATH="$HOME/.ssh/vps_key"
VPS_IP=""
VPS_USER="ubuntu"

# Fonction pour demander l'IP du VPS
get_vps_info() {
    echo -e "${YELLOW}📋 Informations VPS${NC}"
    read -p "Entrez l'adresse IP de votre VPS: " VPS_IP
    read -p "Nom d'utilisateur SSH [$VPS_USER]: " input_user
    VPS_USER=${input_user:-$VPS_USER}
    echo ""
}

# Fonction pour générer les clés SSH
generate_ssh_keys() {
    echo -e "${YELLOW}🔐 Génération des clés SSH${NC}"
    
    if [ -f "$SSH_KEY_PATH" ]; then
        echo -e "${YELLOW}⚠️  Une clé SSH existe déjà à $SSH_KEY_PATH${NC}"
        read -p "Voulez-vous la remplacer ? (y/N): " replace
        if [[ ! $replace =~ ^[Yy]$ ]]; then
            echo "Utilisation de la clé existante..."
            return
        fi
        rm -f "$SSH_KEY_PATH" "$SSH_KEY_PATH.pub"
    fi
    
    # Générer la clé SSH
    ssh-keygen -t ed25519 -C "meshaplus-vps@$(hostname)" -f "$SSH_KEY_PATH" -N ""
    
    # Sécuriser les permissions
    chmod 600 "$SSH_KEY_PATH"
    chmod 644 "$SSH_KEY_PATH.pub"
    
    echo -e "${GREEN}✅ Clés SSH générées avec succès${NC}"
    echo "  - Clé privée: $SSH_KEY_PATH"
    echo "  - Clé publique: $SSH_KEY_PATH.pub"
    echo ""
}

# Fonction pour afficher la clé publique
show_public_key() {
    echo -e "${YELLOW}📋 Clé publique à installer sur le VPS${NC}"
    echo "========================================"
    cat "$SSH_KEY_PATH.pub"
    echo "========================================"
    echo ""
}

# Fonction pour installer la clé sur le VPS
install_key_on_vps() {
    echo -e "${YELLOW}🚀 Installation de la clé sur le VPS${NC}"
    echo ""
    echo "Méthode 1: Automatique (si connexion par mot de passe possible)"
    read -p "Voulez-vous essayer l'installation automatique ? (y/N): " auto_install
    
    if [[ $auto_install =~ ^[Yy]$ ]]; then
        echo "Installation automatique de la clé..."
        if ssh-copy-id -i "$SSH_KEY_PATH.pub" "$VPS_USER@$VPS_IP"; then
            echo -e "${GREEN}✅ Clé installée automatiquement${NC}"
            return
        else
            echo -e "${RED}❌ Installation automatique échouée${NC}"
        fi
    fi
    
    echo ""
    echo -e "${YELLOW}Méthode 2: Installation manuelle${NC}"
    echo "1. Connectez-vous à votre VPS :"
    echo "   ssh $VPS_USER@$VPS_IP"
    echo ""
    echo "2. Exécutez ces commandes sur le VPS :"
    echo "   mkdir -p ~/.ssh"
    echo "   chmod 700 ~/.ssh"
    echo "   echo \"$(cat "$SSH_KEY_PATH.pub")\" >> ~/.ssh/authorized_keys"
    echo "   chmod 600 ~/.ssh/authorized_keys"
    echo ""
    read -p "Appuyez sur Entrée une fois l'installation manuelle terminée..."
}

# Fonction pour tester la connexion
test_connection() {
    echo -e "${YELLOW}🧪 Test de la connexion SSH${NC}"
    
    if ssh -i "$SSH_KEY_PATH" -o ConnectTimeout=10 -o StrictHostKeyChecking=no "$VPS_USER@$VPS_IP" "echo 'Connexion SSH réussie !'"; then
        echo -e "${GREEN}✅ Connexion SSH fonctionnelle${NC}"
        return 0
    else
        echo -e "${RED}❌ Connexion SSH échouée${NC}"
        return 1
    fi
}

# Fonction pour afficher la clé privée pour GitHub
show_private_key_for_github() {
    echo -e "${YELLOW}📋 Clé privée pour GitHub Secret (VPS_SSH_KEY)${NC}"
    echo "=================================================="
    echo -e "${BLUE}Copiez TOUT le contenu ci-dessous :${NC}"
    echo ""
    cat "$SSH_KEY_PATH"
    echo ""
    echo "=================================================="
    echo ""
    echo -e "${YELLOW}Instructions GitHub :${NC}"
    echo "1. Allez sur votre repository GitHub"
    echo "2. Settings → Secrets and variables → Actions"
    echo "3. New repository secret"
    echo "4. Name: VPS_SSH_KEY"
    echo "5. Value: Collez le contenu ci-dessus (avec BEGIN et END)"
    echo ""
}

# Fonction pour afficher le résumé
show_summary() {
    echo -e "${GREEN}🎉 Configuration SSH Terminée !${NC}"
    echo ""
    echo -e "${BLUE}📋 Récapitulatif des secrets GitHub :${NC}"
    echo "======================================"
    echo "VPS_HOST=$VPS_IP"
    echo "VPS_USERNAME=$VPS_USER"
    echo "VPS_SSH_KEY=[Clé privée affichée ci-dessus]"
    echo "VPS_DOCKER_COMPOSE_PATH=/opt/meshaplus"
    echo ""
    echo -e "${BLUE}📝 Prochaines étapes :${NC}"
    echo "1. ✅ Ajouter tous les secrets dans GitHub"
    echo "2. ✅ Exécuter le script setup-vps.sh sur votre VPS"
    echo "3. ✅ Configurer les autres secrets (base de données, Telegram)"
    echo "4. ✅ Tester le déploiement"
    echo ""
    echo -e "${YELLOW}💡 Commande pour exécuter setup-vps.sh :${NC}"
    echo "ssh -i ~/.ssh/vps_key $VPS_USER@$VPS_IP"
    echo "wget https://raw.githubusercontent.com/votre-username/meshaplus/main/backend/scripts/setup-vps.sh"
    echo "chmod +x setup-vps.sh"
    echo "sudo ./setup-vps.sh"
    echo ""
}

# Fonction principale
main() {
    get_vps_info
    generate_ssh_keys
    show_public_key
    install_key_on_vps
    
    if test_connection; then
        echo ""
        show_private_key_for_github
        show_summary
    else
        echo ""
        echo -e "${RED}❌ La connexion SSH ne fonctionne pas${NC}"
        echo "Vérifiez que la clé publique est bien installée sur le VPS"
        echo ""
        echo -e "${YELLOW}Pour débugger :${NC}"
        echo "ssh -i $SSH_KEY_PATH -v $VPS_USER@$VPS_IP"
    fi
}

# Exécution
main "$@"