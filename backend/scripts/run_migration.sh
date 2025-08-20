#!/bin/bash

# Script pour exécuter la migration de la base de données
# Assurez-vous que PostgreSQL est en cours d'exécution

echo "🚀 Exécution de la migration de la base de données..."

# Variables de configuration (ajustez selon votre configuration)
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="meshaplus"
DB_USER="postgres"
DB_PASSWORD="postgres"

# Exécuter la migration
echo "📝 Application de la migration users..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/migrate_users.sql

if [ $? -eq 0 ]; then
    echo "✅ Migration terminée avec succès!"
    echo "🎉 La table users a été mise à jour avec la nouvelle structure."
else
    echo "❌ Erreur lors de la migration."
    echo "💡 Vérifiez que PostgreSQL est en cours d'exécution et que les paramètres de connexion sont corrects."
    exit 1
fi 