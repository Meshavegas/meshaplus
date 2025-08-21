#!/bin/bash

# 🔧 Script de Configuration SSH Automatique - MeshaPlus Backend
# Ce script automatise la configuration SSH pour le déploiement

set -e

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔧 Configuration SSH Automatique - MeshaPlus Backend${NC}"
echo "=================================================="
echo ""

# Fonction d'aide
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Afficher cette aide"
    echo "  -g, --generate      Générer une nouvelle paire de clés SSH"
    echo "  -c, --copy HOST     Copier la clé publique sur le VPS"
    echo "  -t, --test HOST     Tester la connexion SSH"
    echo "  -s, --setup HOST    Configuration complète (générer + copier + tester)"
    echo ""
    echo "Exemples:"
    echo "  $0 --setup ubuntu@192.168.1.100"
    echo "  $0 --generate"
    echo "  $0 --copy ubuntu@192.168.1.100"
    echo "  $0 --test ubuntu@192.168.1.100"
}

# Fonction de génération de clés
generate_keys() {
    echo -e "${YELLOW}🔑 Génération d'une nouvelle paire de clés SSH...${NC}"
    
    # Vérifier si les clés existent déjà
    if [ -f ~/.ssh/id_rsa ] && [ -f ~/.ssh/id_rsa.pub ]; then
        echo -e "${YELLOW}⚠️  Des clés SSH existent déjà.${NC}"
        read -p "Voulez-vous les remplacer ? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}⏭️  Génération annulée${NC}"
            return 0
        fi
        
        # Sauvegarder les anciennes clés
        mv ~/.ssh/id_rsa ~/.ssh/id_rsa.backup.$(date +%Y%m%d-%H%M%S)
        mv ~/.ssh/id_rsa.pub ~/.ssh/id_rsa.pub.backup.$(date +%Y%m%d-%H%M%S)
        echo -e "${BLUE}💾 Anciennes clés sauvegardées${NC}"
    fi
    
    # Générer les nouvelles clés
    ssh-keygen -t rsa -b 4096 -C "meshaplus-deployment@$(hostname)" -f ~/.ssh/id_rsa -N ""
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Nouvelles clés SSH générées avec succès${NC}"
        echo -e "${BLUE}📁 Emplacement: ~/.ssh/id_rsa (privée) et ~/.ssh/id_rsa.pub (publique)${NC}"
    else
        echo -e "${RED}❌ Erreur lors de la génération des clés${NC}"
        return 1
    fi
}

# Fonction de copie de clé
copy_key() {
    local host=$1
    
    if [ -z "$host" ]; then
        echo -e "${RED}❌ Hôte non spécifié${NC}"
        echo "Usage: $0 --copy username@host"
        return 1
    fi
    
    echo -e "${YELLOW}📋 Copie de la clé publique sur $host...${NC}"
    
    # Vérifier que la clé publique existe
    if [ ! -f ~/.ssh/id_rsa.pub ]; then
        echo -e "${RED}❌ Clé publique non trouvée. Générez d'abord les clés.${NC}"
        return 1
    fi
    
    # Copier la clé
    if ssh-copy-id -i ~/.ssh/id_rsa.pub "$host"; then
        echo -e "${GREEN}✅ Clé publique copiée avec succès sur $host${NC}"
    else
        echo -e "${RED}❌ Erreur lors de la copie de la clé${NC}"
        echo -e "${YELLOW}💡 Essayez de copier manuellement :${NC}"
        echo "cat ~/.ssh/id_rsa.pub | ssh $host 'mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys'"
        return 1
    fi
}

# Fonction de test de connexion
test_connection() {
    local host=$1
    
    if [ -z "$host" ]; then
        echo -e "${RED}❌ Hôte non spécifié${NC}"
        echo "Usage: $0 --test username@host"
        return 1
    fi
    
    echo -e "${YELLOW}🧪 Test de connexion SSH vers $host...${NC}"
    
    # Test de connexion
    if ssh -o ConnectTimeout=10 -o BatchMode=yes "$host" "echo 'Connexion SSH réussie'" 2>/dev/null; then
        echo -e "${GREEN}✅ Connexion SSH réussie${NC}"
        return 0
    else
        echo -e "${RED}❌ Échec de la connexion SSH${NC}"
        echo ""
        echo -e "${YELLOW}🔍 Suggestions de dépannage :${NC}"
        echo "1. Vérifiez que la clé publique est installée sur le VPS"
        echo "2. Vérifiez les permissions : ssh $host 'chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys'"
        echo "3. Vérifiez les logs SSH : ssh $host 'sudo tail -f /var/log/auth.log'"
        echo "4. Testez avec verbose : ssh -v $host"
        return 1
    fi
}

# Fonction de configuration complète
setup_complete() {
    local host=$1
    
    if [ -z "$host" ]; then
        echo -e "${RED}❌ Hôte non spécifié${NC}"
        echo "Usage: $0 --setup username@host"
        return 1
    fi
    
    echo -e "${BLUE}🚀 Configuration SSH complète pour $host${NC}"
    echo "=========================================="
    
    # Étape 1: Générer les clés
    if ! generate_keys; then
        return 1
    fi
    
    echo ""
    
    # Étape 2: Copier la clé
    if ! copy_key "$host"; then
        return 1
    fi
    
    echo ""
    
    # Étape 3: Tester la connexion
    if ! test_connection "$host"; then
        return 1
    fi
    
    echo ""
    echo -e "${GREEN}🎉 Configuration SSH terminée avec succès !${NC}"
    echo ""
    echo -e "${BLUE}📋 Prochaines étapes :${NC}"
    echo "1. Ajoutez la clé privée dans GitHub Secrets :"
    echo "   cat ~/.ssh/id_rsa"
    echo ""
    echo "2. Ajoutez les autres secrets GitHub :"
    echo "   - VPS_HOST: $(echo $host | cut -d@ -f2)"
    echo "   - VPS_USERNAME: $(echo $host | cut -d@ -f1)"
    echo "   - VPS_SSH_KEY: [contenu de ~/.ssh/id_rsa]"
    echo ""
    echo "3. Testez le déploiement en poussant vers GitHub"
}

# Fonction d'affichage de la clé privée
show_private_key() {
    echo -e "${YELLOW}🔑 Affichage de la clé privée pour GitHub Secrets :${NC}"
    echo ""
    echo -e "${BLUE}📋 Copiez le contenu suivant dans le secret VPS_SSH_KEY :${NC}"
    echo "=========================================="
    cat ~/.ssh/id_rsa
    echo "=========================================="
    echo ""
    echo -e "${YELLOW}⚠️  Attention : Cette clé est sensible, ne la partagez pas !${NC}"
}

# Traitement des arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    -g|--generate)
        generate_keys
        show_private_key
        ;;
    -c|--copy)
        copy_key "$2"
        ;;
    -t|--test)
        test_connection "$2"
        ;;
    -s|--setup)
        setup_complete "$2"
        ;;
    *)
        echo -e "${RED}❌ Option non reconnue${NC}"
        echo ""
        show_help
        exit 1
        ;;
esac 