#!/bin/bash

# 🧪 Script de Test de Connexion VPS - MeshaPlus Backend
# Ce script teste la connexion VPS avec des valeurs de test

set -e

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🧪 Test de Connexion VPS - MeshaPlus Backend${NC}"
echo "============================================="
echo ""

# Demander les informations de connexion
echo -e "${YELLOW}📋 Entrez les informations de connexion VPS :${NC}"
echo ""

read -p "Adresse IP du VPS: " VPS_HOST
read -p "Nom d'utilisateur SSH: " VPS_USERNAME
read -p "Chemin de déploiement (ex: /opt/meshaplus): " VPS_DEPLOY_PATH

# Vérifier que les informations sont fournies
if [ -z "$VPS_HOST" ] || [ -z "$VPS_USERNAME" ] || [ -z "$VPS_DEPLOY_PATH" ]; then
    echo -e "${RED}❌ Toutes les informations sont requises${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}📊 Informations de connexion :${NC}"
echo "• VPS_HOST: $VPS_HOST"
echo "• VPS_USERNAME: $VPS_USERNAME"
echo "• VPS_DEPLOY_PATH: $VPS_DEPLOY_PATH"
echo ""

# Test 1: Vérifier la clé SSH
echo -e "${YELLOW}📋 Test 1: Vérification de la clé SSH${NC}"

if [ -f ~/.ssh/id_rsa ]; then
    echo -e "${GREEN}✅ Clé privée trouvée: ~/.ssh/id_rsa${NC}"
    
    # Vérifier les permissions
    if [ -r ~/.ssh/id_rsa ]; then
        echo -e "${GREEN}✅ Clé privée lisible${NC}"
    else
        echo -e "${RED}❌ Clé privée non lisible${NC}"
        chmod 600 ~/.ssh/id_rsa
        echo -e "${YELLOW}⚠️  Permissions corrigées${NC}"
    fi
else
    echo -e "${RED}❌ Clé privée non trouvée${NC}"
    echo -e "${YELLOW}💡 Générer une nouvelle clé :${NC}"
    echo "ssh-keygen -t rsa -b 4096 -C 'meshaplus-deployment@example.com'"
    exit 1
fi

# Test 2: Test de connectivité réseau
echo ""
echo -e "${YELLOW}📋 Test 2: Test de connectivité réseau${NC}"

if ping -c 1 "$VPS_HOST" >/dev/null 2>&1; then
    echo -e "${GREEN}✅ VPS accessible via ping${NC}"
else
    echo -e "${YELLOW}⚠️  VPS non accessible via ping (peut être normal)${NC}"
fi

# Test 3: Test de connexion SSH
echo ""
echo -e "${YELLOW}📋 Test 3: Test de connexion SSH${NC}"

echo -e "${BLUE}🔗 Tentative de connexion SSH...${NC}"

if ssh -o ConnectTimeout=10 -o BatchMode=yes -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "echo 'Connexion SSH réussie'" 2>/dev/null; then
    echo -e "${GREEN}✅ Connexion SSH réussie${NC}"
    SSH_WORKING=true
else
    echo -e "${RED}❌ Échec de la connexion SSH${NC}"
    SSH_WORKING=false
    
    echo ""
    echo -e "${YELLOW}🔍 Debug de la connexion SSH :${NC}"
    echo "ssh -v -i ~/.ssh/id_rsa $VPS_USERNAME@$VPS_HOST"
    echo ""
    
    # Test avec verbose
    echo -e "${BLUE}🔍 Test avec verbose (dernières lignes) :${NC}"
    ssh -v -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "echo 'test'" 2>&1 | tail -10
fi

# Test 4: Test de copie de fichiers (si SSH fonctionne)
if [ "$SSH_WORKING" = true ]; then
    echo ""
    echo -e "${YELLOW}📋 Test 4: Test de copie de fichiers${NC}"
    
    # Créer un fichier de test
    echo "Test file created at $(date)" > test-file.txt
    
    echo -e "${BLUE}📁 Test de copie de fichier...${NC}"
    
    if scp -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa test-file.txt "$VPS_USERNAME@$VPS_HOST:$VPS_DEPLOY_PATH/" 2>/dev/null; then
        echo -e "${GREEN}✅ Copie de fichier réussie${NC}"
        
        # Vérifier que le fichier existe sur le VPS
        if ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "test -f $VPS_DEPLOY_PATH/test-file.txt && echo 'Fichier trouvé sur le VPS'" 2>/dev/null; then
            echo -e "${GREEN}✅ Fichier vérifié sur le VPS${NC}"
        else
            echo -e "${YELLOW}⚠️  Fichier non trouvé sur le VPS${NC}"
        fi
        
        # Nettoyer le fichier de test
        ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "rm -f $VPS_DEPLOY_PATH/test-file.txt" 2>/dev/null
    else
        echo -e "${RED}❌ Échec de la copie de fichier${NC}"
    fi
    
    # Nettoyer le fichier local
    rm -f test-file.txt
fi

# Test 5: Vérification des prérequis sur le VPS
if [ "$SSH_WORKING" = true ]; then
    echo ""
    echo -e "${YELLOW}📋 Test 5: Vérification des prérequis sur le VPS${NC}"
    
    echo -e "${BLUE}🔍 Vérification de Docker...${NC}"
    if ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "command -v docker >/dev/null 2>&1 && echo 'Docker installé' || echo 'Docker non installé'" 2>/dev/null; then
        echo -e "${GREEN}✅ Docker disponible${NC}"
    else
        echo -e "${RED}❌ Docker non disponible${NC}"
    fi
    
    echo -e "${BLUE}🔍 Vérification de Docker Compose...${NC}"
    if ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "command -v docker-compose >/dev/null 2>&1 && echo 'Docker Compose installé' || echo 'Docker Compose non installé'" 2>/dev/null; then
        echo -e "${GREEN}✅ Docker Compose disponible${NC}"
    else
        echo -e "${RED}❌ Docker Compose non disponible${NC}"
    fi
    
    echo -e "${BLUE}🔍 Vérification du répertoire de déploiement...${NC}"
    if ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "test -d $VPS_DEPLOY_PATH && echo 'Répertoire existe' || echo 'Répertoire inexistant'" 2>/dev/null; then
        echo -e "${GREEN}✅ Répertoire de déploiement existe${NC}"
    else
        echo -e "${YELLOW}⚠️  Répertoire de déploiement inexistant${NC}"
        echo -e "${BLUE}💡 Création du répertoire...${NC}"
        ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "sudo mkdir -p $VPS_DEPLOY_PATH && sudo chown $VPS_USERNAME:$VPS_USERNAME $VPS_DEPLOY_PATH" 2>/dev/null
        echo -e "${GREEN}✅ Répertoire créé${NC}"
    fi
fi

# Résumé
echo ""
echo -e "${BLUE}📋 Résumé du test :${NC}"

if [ "$SSH_WORKING" = true ]; then
    echo -e "${GREEN}✅ Connexion SSH fonctionnelle${NC}"
    echo -e "${GREEN}✅ Prêt pour le déploiement${NC}"
    echo ""
    echo -e "${BLUE}📋 Secrets GitHub à configurer :${NC}"
    echo "• VPS_HOST: $VPS_HOST"
    echo "• VPS_USERNAME: $VPS_USERNAME"
    echo "• VPS_DEPLOY_PATH: $VPS_DEPLOY_PATH"
    echo "• VPS_SSH_KEY: [contenu de ~/.ssh/id_rsa]"
    echo ""
    echo -e "${BLUE}💡 Pour récupérer la clé privée :${NC}"
    echo "cat ~/.ssh/id_rsa"
else
    echo -e "${RED}❌ Connexion SSH échouée${NC}"
    echo ""
    echo -e "${YELLOW}🔍 Suggestions de résolution :${NC}"
    echo "1. Vérifiez que la clé publique est installée sur le VPS"
    echo "2. Vérifiez les permissions SSH sur le VPS"
    echo "3. Vérifiez que SSH est activé sur le VPS"
    echo "4. Vérifiez les logs SSH : ssh $VPS_USERNAME@$VPS_HOST 'sudo tail -f /var/log/auth.log'"
fi 