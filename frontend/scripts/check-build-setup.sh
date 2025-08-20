#!/bin/bash

echo "🔍 Vérification de la configuration de build Android..."

# Vérifier si EAS CLI est installé
if ! command -v eas &> /dev/null; then
    echo "❌ EAS CLI n'est pas installé"
    echo "Installez-le avec: npm install -g @expo/eas-cli"
    exit 1
else
    echo "✅ EAS CLI est installé"
fi

# Vérifier si l'utilisateur est connecté
if ! eas whoami &> /dev/null; then
    echo "❌ Vous n'êtes pas connecté à Expo"
    echo "Connectez-vous avec: eas login"
    exit 1
else
    echo "✅ Connecté à Expo"
fi

# Vérifier les fichiers de configuration
if [ ! -f "eas.json" ]; then
    echo "❌ eas.json manquant"
    exit 1
else
    echo "✅ eas.json trouvé"
fi

if [ ! -f "app.json" ]; then
    echo "❌ app.json manquant"
    exit 1
else
    echo "✅ app.json trouvé"
fi

# Vérifier les permissions Android
if grep -q "package.*com.meshaplus.app" app.json; then
    echo "✅ Package Android configuré"
else
    echo "❌ Package Android non configuré"
fi

# Vérifier les scripts npm
if grep -q "build:dev:android" package.json; then
    echo "✅ Scripts de build configurés"
else
    echo "❌ Scripts de build manquants"
fi

echo ""
echo "🎉 Configuration de build vérifiée !"
echo "Vous pouvez maintenant lancer: npm run build:dev:android" 