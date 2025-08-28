-- Script pour réinitialiser la base de données
-- À exécuter en cas de problème avec les migrations

-- Supprimer toutes les tables pour éviter les conflits de structure
DROP TABLE IF EXISTS budgets CASCADE;
DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS saving_goals CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS moods CASCADE;
DROP TABLE IF EXISTS reminders CASCADE;
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS motivations CASCADE;
DROP TABLE IF EXISTS life_notes CASCADE;
DROP TABLE IF EXISTS weekly_summaries CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS preferences CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Les tables seront recréées automatiquement par les migrations au prochain démarrage