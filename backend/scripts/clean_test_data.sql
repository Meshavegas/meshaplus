-- Script pour nettoyer les données de test
-- ATTENTION: Ce script supprime TOUS les utilisateurs de test

-- Supprimer les utilisateurs de test
DELETE FROM users WHERE email LIKE '%@example.com';

-- Réinitialiser la séquence si nécessaire
-- (Pas nécessaire avec UUID, mais gardé pour référence)

-- Afficher le nombre d'utilisateurs restants
SELECT COUNT(*) as total_users FROM users;

-- Afficher les utilisateurs restants
SELECT id, name, email, created_at FROM users ORDER BY created_at DESC; 