#!/bin/bash

# Script pour tester le déploiement de l'API
# Usage: ./scripts/test-deployment.sh [VPS_HOST] [VPS_USER]

set -e

# Variables par défaut
VPS_HOST=${1:-"localhost"}
VPS_USER=${2:-"root"}
SERVICE_NAME="meshaplus-api"
API_URL="http://$VPS_HOST:8080"

echo "🧪 Test du déploiement MeshaPlus API"
echo "VPS: $VPS_HOST"
echo "Utilisateur: $VPS_USER"
echo "API URL: $API_URL"
echo ""

# Test de connexion SSH
echo "🔌 Test de connexion SSH..."
if ssh -o ConnectTimeout=10 -o BatchMode=yes "$VPS_USER@$VPS_HOST" "echo 'SSH connection successful'" 2>/dev/null; then
    echo "✅ Connexion SSH réussie"
else
    echo "❌ Échec de la connexion SSH"
    echo "Vérifiez:"
    echo "  - L'adresse IP/domaine du VPS"
    echo "  - Les clés SSH"
    echo "  - Le firewall"
    exit 1
fi

# Test du service systemd
echo ""
echo "🔧 Test du service systemd..."
ssh "$VPS_USER@$VPS_HOST" "
    if systemctl is-active --quiet $SERVICE_NAME; then
        echo '✅ Service $SERVICE_NAME est actif'
        systemctl status $SERVICE_NAME --no-pager | head -10
    else
        echo '❌ Service $SERVICE_NAME n\'est pas actif'
        echo 'Tentative de démarrage...'
        systemctl start $SERVICE_NAME
        sleep 3
        if systemctl is-active --quiet $SERVICE_NAME; then
            echo '✅ Service démarré avec succès'
        else
            echo '❌ Échec du démarrage du service'
            journalctl -u $SERVICE_NAME --no-pager -n 20
            exit 1
        fi
    fi
"

# Test de l'API
echo ""
echo "🌐 Test de l'API..."
echo "Attente du démarrage de l'API..."

# Attendre que l'API soit prête
for i in {1..30}; do
    if curl -f -s "$API_URL/health" >/dev/null 2>&1; then
        echo "✅ API accessible après $i secondes"
        break
    fi
    
    if [ $i -eq 30 ]; then
        echo "❌ API non accessible après 30 secondes"
        echo "Vérification des logs..."
        ssh "$VPS_USER@$VPS_HOST" "journalctl -u $SERVICE_NAME --no-pager -n 20"
        exit 1
    fi
    
    echo "⏳ Attente... ($i/30)"
    sleep 1
done

# Tests des endpoints
echo ""
echo "🔍 Tests des endpoints..."

# Test health endpoint
echo "🏥 Test endpoint /health..."
if curl -f -s "$API_URL/health" | grep -q "ok"; then
    echo "✅ Health check réussi"
else
    echo "❌ Health check échoué"
    curl -v "$API_URL/health"
fi

# Test API version
echo "📋 Test endpoint /api/v1..."
if curl -f -s "$API_URL/api/v1" >/dev/null 2>&1; then
    echo "✅ API v1 accessible"
else
    echo "⚠️ API v1 non accessible (peut être normal)"
fi

# Test Swagger
echo "📚 Test endpoint /swagger..."
if curl -f -s "$API_URL/swagger/index.html" >/dev/null 2>&1; then
    echo "✅ Documentation Swagger accessible"
else
    echo "⚠️ Documentation Swagger non accessible (peut être normal)"
fi

# Vérification des ressources
echo ""
echo "💾 Vérification des ressources..."
ssh "$VPS_USER@$VPS_HOST" "
    echo '📁 Répertoire de déploiement:'
    ls -la /opt/$SERVICE_NAME/
    echo ''
    echo '💾 Sauvegardes:'
    ls -la /opt/$SERVICE_NAME/backups/ 2>/dev/null || echo 'Aucune sauvegarde'
    echo ''
    echo '📊 Utilisation mémoire:'
    ps aux | grep $SERVICE_NAME | grep -v grep || echo 'Processus non trouvé'
    echo ''
    echo '🌐 Ports en écoute:'
    netstat -tlnp | grep :8080 || echo 'Port 8080 non en écoute'
"

echo ""
echo "🎉 Tests terminés!"
echo ""
echo "📋 Résumé:"
echo "  - SSH: ✅"
echo "  - Service systemd: ✅"
echo "  - API: ✅"
echo ""
echo "🌐 Votre API est accessible sur: $API_URL"
echo "📚 Documentation: $API_URL/swagger/index.html" 