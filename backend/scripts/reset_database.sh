#!/bin/bash

# Script pour réinitialiser la base de données PostgreSQL
# Usage: ./scripts/reset_database.sh

echo "🔄 Réinitialisation de la base de données..."

# Configuration par défaut (peut être surchargée par les variables d'environnement)
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-meshaplus}
DB_USER=${DB_USER:-postgres}

echo "📊 Suppression des tables existantes..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/reset_database.sql

if [ $? -eq 0 ]; then
    echo "✅ Base de données réinitialisée avec succès"
    echo "🚀 Vous pouvez maintenant redémarrer l'application"
    echo "   Les tables seront recréées automatiquement par les migrations"
else
    echo "❌ Erreur lors de la réinitialisation de la base de données"
    exit 1
fi