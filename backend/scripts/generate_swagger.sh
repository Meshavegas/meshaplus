#!/bin/bash

# Script pour générer la documentation Swagger
echo "🔄 Génération de la documentation Swagger..."

# Vérifier si swag est installé
if ! command -v ~/go/bin/swag &> /dev/null; then
    echo "📦 Installation de swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Nettoyer les anciens fichiers de documentation
echo "🧹 Nettoyage des anciens fichiers..."
rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml

# Générer la nouvelle documentation
echo "📝 Génération de la documentation..."
~/go/bin/swag init -g cmd/api/main.go -o docs

# Vérifier si la génération a réussi
if [ $? -eq 0 ]; then
    echo "✅ Documentation Swagger générée avec succès !"
    echo "📁 Fichiers créés :"
    echo "   - docs/docs.go"
    echo "   - docs/swagger.json"
    echo "   - docs/swagger.yaml"
    echo ""
    echo "🌐 Accédez à la documentation : http://localhost:8080/swagger/"
    echo "📊 Interface Swagger UI disponible après démarrage du serveur"
else
    echo "❌ Erreur lors de la génération de la documentation"
    exit 1
fi 