#!/bin/bash

# Script pour nettoyer les données de test
# ATTENTION: Ce script supprime TOUS les utilisateurs de test

echo "🧹 Nettoyage des données de test"
echo "================================"

# Variables de configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="meshaplus"
DB_USER="postgres"
DB_PASSWORD="postgres"

# Couleurs pour l'affichage
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Demander confirmation
echo -e "${YELLOW}⚠️  ATTENTION: Ce script va supprimer tous les utilisateurs avec des emails @example.com${NC}"
read -p "Êtes-vous sûr de vouloir continuer ? (y/N): " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}❌ Nettoyage annulé${NC}"
    exit 1
fi

# Exécuter le script SQL
echo -e "${YELLOW}🗑️  Suppression des utilisateurs de test...${NC}"
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/clean_test_data.sql

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Nettoyage terminé avec succès!${NC}"
else
    echo -e "${RED}❌ Erreur lors du nettoyage${NC}"
    exit 1
fi

echo -e "\n${GREEN}🎉 Base de données nettoyée !${NC}"
echo -e "${GREEN}Vous pouvez maintenant relancer les tests avec: make test-auth${NC}" 