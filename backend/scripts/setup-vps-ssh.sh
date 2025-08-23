#!/bin/bash

# 🔧 Script de Configuration SSH VPS - MeshaPlus Backend
# Ce script configure SSH sur le VPS pour le déploiement

set -e

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔧 Configuration SSH VPS - MeshaPlus Backend${NC}"
echo "============================================="
echo ""

# Demander les informations de connexion
echo -e "${YELLOW}📋 Entrez les informations de connexion VPS :${NC}"
echo ""

read -p "Adresse IP du VPS: " VPS_HOST
read -p "Nom d'utilisateur SSH: " VPS_USERNAME
read -p "Mot de passe SSH (si nécessaire): " -s VPS_PASSWORD
echo ""

# Vérifier que les informations sont fournies
if [ -z "$VPS_HOST" ] || [ -z "$VPS_USERNAME" ]; then
    echo -e "${RED}❌ VPS_HOST et VPS_USERNAME sont requis${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}📊 Informations de connexion :${NC}"
echo "• VPS_HOST: $VPS_HOST"
echo "• VPS_USERNAME: $VPS_USERNAME"
echo ""

# Étape 1: Vérifier la clé SSH locale
echo -e "${YELLOW}📋 Étape 1: Vérification de la clé SSH locale${NC}"

if [ ! -f ~/.ssh/id_rsa ]; then
    echo -e "${RED}❌ Clé privée non trouvée${NC}"
    echo -e "${BLUE}💡 Génération d'une nouvelle clé SSH...${NC}"
    
    ssh-keygen -t rsa -b 4096 -C "meshaplus-deployment@$(hostname)" -f ~/.ssh/id_rsa -N ""
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Nouvelle clé SSH générée${NC}"
    else
        echo -e "${RED}❌ Erreur lors de la génération de la clé${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✅ Clé privée trouvée: ~/.ssh/id_rsa${NC}"
fi

# Étape 2: Test de connexion initiale
echo ""
echo -e "${YELLOW}📋 Étape 2: Test de connexion initiale${NC}"

if [ -n "$VPS_PASSWORD" ]; then
    echo -e "${BLUE}🔗 Test de connexion avec mot de passe...${NC}"
    
    # Test avec sshpass si disponible
    if command -v sshpass &> /dev/null; then
        if sshpass -p "$VPS_PASSWORD" ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 "$VPS_USERNAME@$VPS_HOST" "echo 'Connexion avec mot de passe réussie'" 2>/dev/null; then
            echo -e "${GREEN}✅ Connexion avec mot de passe réussie${NC}"
            CONNECTION_METHOD="password"
        else
            echo -e "${RED}❌ Échec de la connexion avec mot de passe${NC}"
            CONNECTION_METHOD="none"
        fi
    else
        echo -e "${YELLOW}⚠️  sshpass non disponible, tentative de connexion manuelle${NC}"
        CONNECTION_METHOD="manual"
    fi
else
    echo -e "${BLUE}🔗 Test de connexion sans mot de passe...${NC}"
    
    if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "echo 'Connexion SSH réussie'" 2>/dev/null; then
        echo -e "${GREEN}✅ Connexion SSH réussie (clé déjà configurée)${NC}"
        CONNECTION_METHOD="key"
    else
        echo -e "${YELLOW}⚠️  Connexion SSH échouée, configuration nécessaire${NC}"
        CONNECTION_METHOD="none"
    fi
fi

# Étape 3: Configuration SSH sur le VPS
if [ "$CONNECTION_METHOD" != "key" ]; then
    echo ""
    echo -e "${YELLOW}📋 Étape 3: Configuration SSH sur le VPS${NC}"
    
    if [ "$CONNECTION_METHOD" = "password" ] && command -v sshpass &> /dev/null; then
        echo -e "${BLUE}🔧 Configuration SSH avec mot de passe...${NC}"
        
        # Copier la clé publique avec sshpass
        if sshpass -p "$VPS_PASSWORD" ssh-copy-id -i ~/.ssh/id_rsa.pub -o StrictHostKeyChecking=no "$VPS_USERNAME@$VPS_HOST"; then
            echo -e "${GREEN}✅ Clé publique copiée avec succès${NC}"
        else
            echo -e "${RED}❌ Échec de la copie de la clé publique${NC}"
            echo -e "${YELLOW}💡 Tentative de copie manuelle...${NC}"
            
            # Copie manuelle
            PUBKEY=$(cat ~/.ssh/id_rsa.pub)
            sshpass -p "$VPS_PASSWORD" ssh -o StrictHostKeyChecking=no "$VPS_USERNAME@$VPS_HOST" "mkdir -p ~/.ssh && echo '$PUBKEY' >> ~/.ssh/authorized_keys && chmod 700 ~/.ssh && chmod 600 ~/.ssh/authorized_keys"
            
            if [ $? -eq 0 ]; then
                echo -e "${GREEN}✅ Clé publique copiée manuellement${NC}"
            else
                echo -e "${RED}❌ Échec de la copie manuelle${NC}"
                exit 1
            fi
        fi
    else
        echo -e "${YELLOW}⚠️  Configuration manuelle requise${NC}"
        echo ""
        echo -e "${BLUE}📋 Instructions de configuration manuelle :${NC}"
        echo ""
        echo "1. Connectez-vous manuellement au VPS :"
        echo "   ssh $VPS_USERNAME@$VPS_HOST"
        echo ""
        echo "2. Créez le répertoire SSH :"
        echo "   mkdir -p ~/.ssh"
        echo ""
        echo "3. Ajoutez votre clé publique :"
        echo "   echo '$(cat ~/.ssh/id_rsa.pub)' >> ~/.ssh/authorized_keys"
        echo ""
        echo "4. Configurez les permissions :"
        echo "   chmod 700 ~/.ssh"
        echo "   chmod 600 ~/.ssh/authorized_keys"
        echo ""
        echo "5. Testez la connexion :"
        echo "   ssh $VPS_USERNAME@$VPS_HOST"
        echo ""
        
        read -p "Appuyez sur Enter une fois la configuration manuelle terminée..."
    fi
fi

# Étape 4: Test de connexion finale
echo ""
echo -e "${YELLOW}📋 Étape 4: Test de connexion finale${NC}"

if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "echo 'Connexion SSH finale réussie'" 2>/dev/null; then
    echo -e "${GREEN}✅ Connexion SSH finale réussie${NC}"
else
    echo -e "${RED}❌ Échec de la connexion SSH finale${NC}"
    echo ""
    echo -e "${YELLOW}🔍 Debug de la connexion SSH :${NC}"
    ssh -v -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "echo 'test'" 2>&1 | tail -10
    exit 1
fi

# Étape 5: Configuration des prérequis sur le VPS
echo ""
echo -e "${YELLOW}📋 Étape 5: Configuration des prérequis sur le VPS${NC}"

echo -e "${BLUE}🔍 Vérification de Docker...${NC}"
if ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "command -v docker >/dev/null 2>&1 && echo 'Docker installé' || echo 'Docker non installé'" 2>/dev/null; then
    echo -e "${GREEN}✅ Docker disponible${NC}"
else
    echo -e "${YELLOW}⚠️  Docker non disponible${NC}"
    echo -e "${BLUE}💡 Installation de Docker...${NC}"
    
    ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" << 'EOF'
        # Installation de Docker
        sudo apt-get update
        sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
        echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
        sudo apt-get update
        sudo apt-get install -y docker-ce docker-ce-cli containerd.io
        sudo usermod -aG docker $USER
        echo "✅ Docker installé"
EOF
fi

echo -e "${BLUE}🔍 Vérification de Docker Compose...${NC}"
if ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" "command -v docker-compose >/dev/null 2>&1 && echo 'Docker Compose installé' || echo 'Docker Compose non installé'" 2>/dev/null; then
    echo -e "${GREEN}✅ Docker Compose disponible${NC}"
else
    echo -e "${YELLOW}⚠️  Docker Compose non disponible${NC}"
    echo -e "${BLUE}💡 Installation de Docker Compose...${NC}"
    
    ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" << 'EOF'
        # Installation de Docker Compose
        sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        sudo chmod +x /usr/local/bin/docker-compose
        echo "✅ Docker Compose installé"
EOF
fi

echo -e "${BLUE}🔍 Création du répertoire de déploiement...${NC}"
ssh -o StrictHostKeyChecking=no -i ~/.ssh/id_rsa "$VPS_USERNAME@$VPS_HOST" << 'EOF'
    # Créer le répertoire de déploiement
    sudo mkdir -p /opt/meshaplus
    sudo chown $USER:$USER /opt/meshaplus
    echo "✅ Répertoire de déploiement créé: /opt/meshaplus"
EOF

# Résumé final
echo ""
echo -e "${GREEN}🎉 Configuration SSH VPS terminée avec succès !${NC}"
echo ""
echo -e "${BLUE}📋 Résumé de la configuration :${NC}"
echo "✅ Clé SSH locale vérifiée"
echo "✅ Connexion SSH configurée"
echo "✅ Docker installé/configuré"
echo "✅ Docker Compose installé/configuré"
echo "✅ Répertoire de déploiement créé"
echo ""
echo -e "${BLUE}📋 Secrets GitHub à configurer :${NC}"
echo "• VPS_HOST: $VPS_HOST"
echo "• VPS_USERNAME: $VPS_USERNAME"
echo "• VPS_SSH_KEY: [contenu de ~/.ssh/id_rsa]"
echo "• VPS_DEPLOY_PATH: /opt/meshaplus (optionnel)"
echo ""
echo -e "${BLUE}💡 Pour récupérer la clé privée :${NC}"
echo "cat ~/.ssh/id_rsa"
echo ""
echo -e "${BLUE}🚀 Prêt pour le déploiement !${NC}" 