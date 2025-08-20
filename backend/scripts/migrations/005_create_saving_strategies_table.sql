-- Migration: Création de la table saving_strategies
-- Description: Table pour stocker les stratégies d'épargne des utilisateurs
-- Date: 2025-01-27

-- Création de la table saving_strategies
CREATE TABLE IF NOT EXISTS saving_strategies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    strategy_name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('percentage', 'fixed', 'goal_based')),
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    frequency VARCHAR(20) NOT NULL CHECK (frequency IN ('weekly', 'monthly', 'yearly')),
    target_goal_id UUID REFERENCES goals(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Contraintes supplémentaires
ALTER TABLE saving_strategies 
ADD CONSTRAINT saving_strategies_percentage_check 
CHECK (
    (type = 'percentage' AND amount <= 100) OR 
    (type IN ('fixed', 'goal_based'))
);

ALTER TABLE saving_strategies 
ADD CONSTRAINT saving_strategies_goal_based_check 
CHECK (
    (type = 'goal_based' AND target_goal_id IS NOT NULL) OR 
    (type IN ('percentage', 'fixed'))
);

-- Index pour les performances
CREATE INDEX IF NOT EXISTS idx_saving_strategies_user_id ON saving_strategies(user_id);
CREATE INDEX IF NOT EXISTS idx_saving_strategies_type ON saving_strategies(type);
CREATE INDEX IF NOT EXISTS idx_saving_strategies_frequency ON saving_strategies(frequency);
CREATE INDEX IF NOT EXISTS idx_saving_strategies_target_goal_id ON saving_strategies(target_goal_id);
CREATE INDEX IF NOT EXISTS idx_saving_strategies_created_at ON saving_strategies(created_at);

-- Index composite pour les requêtes fréquentes
CREATE INDEX IF NOT EXISTS idx_saving_strategies_user_type ON saving_strategies(user_id, type);
CREATE INDEX IF NOT EXISTS idx_saving_strategies_user_frequency ON saving_strategies(user_id, frequency);

-- Commentaires pour la documentation
COMMENT ON TABLE saving_strategies IS 'Stratégies d''épargne définies par les utilisateurs';
COMMENT ON COLUMN saving_strategies.id IS 'Identifiant unique de la stratégie d''épargne';
COMMENT ON COLUMN saving_strategies.user_id IS 'Identifiant de l''utilisateur propriétaire';
COMMENT ON COLUMN saving_strategies.strategy_name IS 'Nom de la stratégie d''épargne';
COMMENT ON COLUMN saving_strategies.type IS 'Type de stratégie: percentage (%), fixed (montant fixe), goal_based (basé sur objectif)';
COMMENT ON COLUMN saving_strategies.amount IS 'Montant ou pourcentage de la stratégie';
COMMENT ON COLUMN saving_strategies.frequency IS 'Fréquence d''application: weekly, monthly, yearly';
COMMENT ON COLUMN saving_strategies.target_goal_id IS 'Objectif cible pour les stratégies goal_based';
COMMENT ON COLUMN saving_strategies.created_at IS 'Date de création de la stratégie';

-- Trigger pour mettre à jour updated_at (si nécessaire à l'avenir)
-- CREATE OR REPLACE FUNCTION update_saving_strategies_updated_at()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = NOW();
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

-- CREATE TRIGGER trigger_update_saving_strategies_updated_at
--     BEFORE UPDATE ON saving_strategies
--     FOR EACH ROW
--     EXECUTE FUNCTION update_saving_strategies_updated_at(); 