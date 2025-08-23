#!/bin/bash

# Script pour configurer le service systemd sur le VPS
# Usage: ./scripts/setup-vps-service.sh

set -e

echo "🚀 Configuration du service systemd pour MeshaPlus API..."

# Variables
SERVICE_NAME="meshaplus-api"
SERVICE_FILE="/etc/systemd/system/$SERVICE_NAME.service"
DEPLOY_DIR="/opt/$SERVICE_NAME"

# Créer les répertoires nécessaires
echo "📁 Création des répertoires..."
sudo mkdir -p "$DEPLOY_DIR/bin"
sudo mkdir -p "$DEPLOY_DIR/backups"
sudo mkdir -p "$DEPLOY_DIR/logs"

# Copier le fichier de service
echo "📋 Installation du fichier de service..."
sudo cp scripts/meshaplus-api.service "$SERVICE_FILE"

# Recharger systemd
echo "🔄 Rechargement de systemd..."
sudo systemctl daemon-reload

# Activer le service
echo "✅ Activation du service..."
sudo systemctl enable "$SERVICE_NAME"

# Vérifier la configuration
echo "🔍 Vérification de la configuration..."
sudo systemctl status "$SERVICE_NAME" --no-pager || true

echo "🎉 Configuration terminée!"
echo ""
echo "📋 Commandes utiles:"
echo "  sudo systemctl start $SERVICE_NAME    # Démarrer le service"
echo "  sudo systemctl stop $SERVICE_NAME     # Arrêter le service"
echo "  sudo systemctl restart $SERVICE_NAME  # Redémarrer le service"
echo "  sudo systemctl status $SERVICE_NAME   # Vérifier le statut"
echo "  sudo journalctl -u $SERVICE_NAME -f   # Voir les logs en temps réel" 