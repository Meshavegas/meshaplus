-- Migration pour mettre à jour la table users avec la nouvelle structure
-- Exécutez ce script dans votre base de données PostgreSQL

-- Supprimer la table users existante si elle existe
DROP TABLE IF EXISTS users CASCADE;

-- Créer la nouvelle table users avec la structure mise à jour
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    avatar TEXT DEFAULT '',
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Créer un index sur l'email pour les recherches rapides
CREATE INDEX idx_users_email ON users(email);

-- Créer un index sur deleted_at pour le soft delete
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Créer un index sur created_at pour le tri
CREATE INDEX idx_users_created_at ON users(created_at);

-- Créer un trigger pour mettre à jour automatiquement updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insérer quelques utilisateurs de test (optionnel)
-- INSERT INTO users (name, email, password_hash) VALUES 
-- ('Admin User', 'admin@example.com', '$2a$10$example_hash_here'),
-- ('Test User', 'test@example.com', '$2a$10$example_hash_here'); 