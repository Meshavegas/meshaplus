#!/bin/bash

# Script pour configurer les secrets GitHub pour le déploiement VPS
# Usage: ./scripts/setup-github-secrets.sh

echo "🔐 Configuration des secrets GitHub pour le déploiement VPS"
echo ""
echo "Ce script vous aide à configurer les secrets nécessaires dans GitHub."
echo "Vous devez avoir les droits d'administration sur le repository."
echo ""

# Vérifier si gh CLI est installé
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) n'est pas installé."
    echo "Installez-le depuis: https://cli.github.com/"
    exit 1
fi

# Vérifier si l'utilisateur est connecté
if ! gh auth status &> /dev/null; then
    echo "❌ Vous n'êtes pas connecté à GitHub CLI."
    echo "Connectez-vous avec: gh auth login"
    exit 1
fi

echo "📋 Secrets nécessaires:"
echo "  - VPS_HOST: L'adresse IP ou le nom de domaine de votre VPS"
echo "  - VPS_USER: Le nom d'utilisateur pour SSH (généralement 'root')"
echo "  - VPS_SSH_KEY: La clé SSH privée pour se connecter au VPS"
echo "  - VPS_PORT: Le port SSH (optionnel, défaut: 22)"
echo ""

read -p "Voulez-vous configurer ces secrets maintenant? (y/N): " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    
    # VPS_HOST
    read -p "VPS_HOST (adresse IP/domaine): " vps_host
    if [ ! -z "$vps_host" ]; then
        echo "🔧 Configuration de VPS_HOST..."
        gh secret set VPS_HOST --body "$vps_host"
        echo "✅ VPS_HOST configuré"
    fi
    
    # VPS_USER
    read -p "VPS_USER (utilisateur SSH): " vps_user
    if [ ! -z "$vps_user" ]; then
        echo "🔧 Configuration de VPS_USER..."
        gh secret set VPS_USER --body "$vps_user"
        echo "✅ VPS_USER configuré"
    fi
    
    # VPS_PORT
    read -p "VPS_PORT (port SSH, défaut 22): " vps_port
    if [ ! -z "$vps_port" ]; then
        echo "🔧 Configuration de VPS_PORT..."
        gh secret set VPS_PORT --body "$vps_port"
        echo "✅ VPS_PORT configuré"
    fi
    
    # VPS_SSH_KEY
    echo ""
    echo "🔑 Pour VPS_SSH_KEY, vous devez fournir le contenu de votre clé SSH privée."
    echo "Généralement située dans ~/.ssh/id_rsa ou ~/.ssh/id_ed25519"
    read -p "Chemin vers votre clé SSH privée: " ssh_key_path
    
    if [ -f "$ssh_key_path" ]; then
        echo "🔧 Configuration de VPS_SSH_KEY..."
        gh secret set VPS_SSH_KEY < "$ssh_key_path"
        echo "✅ VPS_SSH_KEY configuré"
    else
        echo "❌ Fichier de clé SSH non trouvé: $ssh_key_path"
    fi
    
    echo ""
    echo "🎉 Configuration terminée!"
    echo ""
    echo "📋 Vérification des secrets configurés:"
    gh secret list | grep VPS_ || echo "Aucun secret VPS trouvé"
    
else
    echo "Configuration annulée."
fi

echo ""
echo "📚 Documentation:"
echo "  - Workflow: .github/workflows/vps-deploy.yml"
echo "  - Service systemd: scripts/meshaplus-api.service"
echo "  - Script de configuration: scripts/setup-vps-service.sh" 