#!/bin/bash

# 🔍 Script de Debug SSH - MeshaPlus Backend
# Ce script aide à diagnostiquer les problèmes de connexion SSH

set -e

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 Debug SSH - MeshaPlus Backend${NC}"
echo "================================"
echo ""

# Test 1: Vérifier les variables d'environnement
echo -e "${YELLOW}📋 Test 1: Vérification des variables d'environnement${NC}"

if [ -z "$VPS_HOST" ]; then
    echo -e "${RED}❌ VPS_HOST non défini${NC}"
else
    echo -e "${GREEN}✅ VPS_HOST: $VPS_HOST${NC}"
fi

if [ -z "$VPS_USERNAME" ]; then
    echo -e "${RED}❌ VPS_USERNAME non défini${NC}"
else
    echo -e "${GREEN}✅ VPS_USERNAME: $VPS_USERNAME${NC}"
fi

if [ -z "$VPS_SSH_KEY" ]; then
    echo -e "${RED}❌ VPS_SSH_KEY non défini${NC}"
else
    echo -e "${GREEN}✅ VPS_SSH_KEY défini (longueur: ${#VPS_SSH_KEY} caractères)${NC}"
fi

# Test 2: Vérifier la clé SSH
echo ""
echo -e "${YELLOW}📋 Test 2: Vérification de la clé SSH${NC}"

if [ -n "$VPS_SSH_KEY" ]; then
    # Créer le répertoire SSH
    mkdir -p ~/.ssh
    
    # Écrire la clé
    echo "$VPS_SSH_KEY" > ~/.ssh/vps_key
    chmod 600 ~/.ssh/vps_key
    
    echo -e "${GREEN}✅ Clé SSH écrite dans ~/.ssh/vps_key${NC}"
    
    # Vérifier le format de la clé
    if grep -q "BEGIN.*PRIVATE KEY" ~/.ssh/vps_key; then
        echo -e "${GREEN}✅ Format de clé privée détecté${NC}"
    else
        echo -e "${RED}❌ Format de clé privée non reconnu${NC}"
    fi
    
    # Vérifier les permissions
    if [ -r ~/.ssh/vps_key ]; then
        echo -e "${GREEN}✅ Clé lisible${NC}"
    else
        echo -e "${RED}❌ Clé non lisible${NC}"
    fi
else
    echo -e "${RED}❌ Aucune clé SSH fournie${NC}"
fi

# Test 3: Test de connexion SSH
echo ""
echo -e "${YELLOW}📋 Test 3: Test de connexion SSH${NC}"

if [ -n "$VPS_HOST" ] && [ -n "$VPS_USERNAME" ] && [ -f ~/.ssh/vps_key ]; then
    echo -e "${BLUE}🔗 Test de connexion à $VPS_USERNAME@$VPS_HOST...${NC}"
    
    # Test avec verbose pour voir les détails
    if ssh -v -o StrictHostKeyChecking=no -i ~/.ssh/vps_key $VPS_USERNAME@$VPS_HOST "echo 'Connexion SSH réussie'" 2>&1; then
        echo -e "${GREEN}✅ Connexion SSH réussie${NC}"
    else
        echo -e "${RED}❌ Échec de la connexion SSH${NC}"
        echo ""
        echo -e "${YELLOW}🔍 Suggestions de résolution :${NC}"
        echo "1. Vérifiez que la clé publique est installée sur le VPS"
        echo "2. Vérifiez que l'utilisateur a les bonnes permissions"
        echo "3. Vérifiez que SSH est activé sur le VPS"
        echo "4. Vérifiez les logs SSH sur le VPS: sudo tail -f /var/log/auth.log"
    fi
else
    echo -e "${RED}❌ Impossible de tester la connexion (variables manquantes)${NC}"
fi

# Test 4: Vérifier les secrets GitHub
echo ""
echo -e "${YELLOW}📋 Test 4: Vérification des secrets GitHub${NC}"

echo -e "${BLUE}📋 Secrets requis pour le déploiement :${NC}"
echo "• VPS_HOST: Adresse IP ou domaine du VPS"
echo "• VPS_USERNAME: Nom d'utilisateur SSH"
echo "• VPS_SSH_KEY: Clé SSH privée (format OpenSSH)"
echo "• VPS_DEPLOY_PATH: Chemin de déploiement sur le VPS"
echo "• PROD_DATABASE_URL: URL de la base de données production"
echo "• PROD_JWT_SECRET: Clé secrète JWT"
echo "• TELEGRAM_BOT_TOKEN: Token du bot Telegram"
echo "• TELEGRAM_CHAT_ID: ID du chat Telegram"
echo ""

# Test 5: Instructions de configuration
echo -e "${YELLOW}📋 Test 5: Instructions de configuration${NC}"

echo -e "${BLUE}🔧 Pour configurer SSH sur le VPS :${NC}"
echo ""
echo "1. Générer une paire de clés SSH :"
echo "   ssh-keygen -t rsa -b 4096 -C 'your-email@example.com'"
echo ""
echo "2. Copier la clé publique sur le VPS :"
echo "   ssh-copy-id $VPS_USERNAME@$VPS_HOST"
echo ""
echo "3. Ou manuellement :"
echo "   cat ~/.ssh/id_rsa.pub | ssh $VPS_USERNAME@$VPS_HOST 'mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys'"
echo ""
echo "4. Configurer les permissions :"
echo "   ssh $VPS_USERNAME@$VPS_HOST 'chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys'"
echo ""

# Nettoyage
echo -e "${YELLOW}🧹 Nettoyage...${NC}"
rm -f ~/.ssh/vps_key
echo -e "${GREEN}✅ Nettoyage terminé${NC}"

echo ""
echo -e "${BLUE}📋 Résumé du debug :${NC}"
echo "✅ Variables d'environnement vérifiées"
echo "✅ Clé SSH analysée"
echo "✅ Test de connexion effectué"
echo "✅ Instructions de configuration fournies"
echo ""
echo -e "${BLUE}🚀 Prochaines étapes :${NC}"
echo "1. Configurer SSH sur le VPS"
echo "2. Ajouter les secrets GitHub"
echo "3. Tester la connexion manuellement"
echo "4. Relancer le déploiement" 