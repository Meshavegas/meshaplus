-- Migration pour créer la table revenues
-- Date: 2025-01-26
-- Description: Table pour gérer les revenus des utilisateurs

CREATE TABLE IF NOT EXISTS revenues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    source VARCHAR(255) NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    notes TEXT,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('one_time', 'recurring', 'bonus', 'commission')),
    recurrence_id UUID REFERENCES revenue_recurrences(id) ON DELETE SET NULL,
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    frequency VARCHAR(20) CHECK (frequency IN ('weekly', 'monthly', 'yearly')),
    payment_day INTEGER CHECK (payment_day >= 1 AND payment_day <= 31),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index pour améliorer les performances
CREATE INDEX IF NOT EXISTS idx_revenues_user_id ON revenues(user_id);
CREATE INDEX IF NOT EXISTS idx_revenues_date ON revenues(date);
CREATE INDEX IF NOT EXISTS idx_revenues_type ON revenues(type);
CREATE INDEX IF NOT EXISTS idx_revenues_category_id ON revenues(category_id);
CREATE INDEX IF NOT EXISTS idx_revenues_is_active ON revenues(is_active);
CREATE INDEX IF NOT EXISTS idx_revenues_user_date ON revenues(user_id, date);

-- Contrainte pour vérifier que end_date > start_date si les deux sont définis
ALTER TABLE revenues ADD CONSTRAINT check_date_range 
    CHECK (end_date IS NULL OR start_date IS NULL OR end_date > start_date);

-- Contrainte pour vérifier que payment_day est défini si frequency est 'monthly'
ALTER TABLE revenues ADD CONSTRAINT check_payment_day_monthly
    CHECK (frequency != 'monthly' OR payment_day IS NOT NULL);

-- Contrainte pour vérifier que start_date est défini si type est 'recurring'
ALTER TABLE revenues ADD CONSTRAINT check_recurring_start_date
    CHECK (type != 'recurring' OR start_date IS NOT NULL);

-- Contrainte pour vérifier que frequency est défini si type est 'recurring'
ALTER TABLE revenues ADD CONSTRAINT check_recurring_frequency
    CHECK (type != 'recurring' OR frequency IS NOT NULL);

-- Trigger pour mettre à jour updated_at automatiquement
CREATE OR REPLACE FUNCTION update_revenues_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_revenues_updated_at
    BEFORE UPDATE ON revenues
    FOR EACH ROW
    EXECUTE FUNCTION update_revenues_updated_at();

-- Commentaires sur la table et les colonnes
COMMENT ON TABLE revenues IS 'Table pour stocker les revenus des utilisateurs';
COMMENT ON COLUMN revenues.id IS 'Identifiant unique du revenu';
COMMENT ON COLUMN revenues.user_id IS 'Identifiant de l''utilisateur propriétaire';
COMMENT ON COLUMN revenues.amount IS 'Montant du revenu (doit être positif)';
COMMENT ON COLUMN revenues.source IS 'Source du revenu (ex: Salaire, Freelance, etc.)';
COMMENT ON COLUMN revenues.date IS 'Date du revenu';
COMMENT ON COLUMN revenues.notes IS 'Notes optionnelles sur le revenu';
COMMENT ON COLUMN revenues.category_id IS 'Catégorie associée au revenu';
COMMENT ON COLUMN revenues.type IS 'Type de revenu: one_time, recurring, bonus, commission';
COMMENT ON COLUMN revenues.recurrence_id IS 'Référence vers la récurrence si applicable';
COMMENT ON COLUMN revenues.start_date IS 'Date de début pour les revenus récurrents';
COMMENT ON COLUMN revenues.end_date IS 'Date de fin pour les revenus récurrents';
COMMENT ON COLUMN revenues.frequency IS 'Fréquence pour les revenus récurrents: weekly, monthly, yearly';
COMMENT ON COLUMN revenues.payment_day IS 'Jour de paiement pour les revenus mensuels (1-31)';
COMMENT ON COLUMN revenues.is_active IS 'Indique si le revenu est actif';
COMMENT ON COLUMN revenues.created_at IS 'Date de création du revenu';
COMMENT ON COLUMN revenues.updated_at IS 'Date de dernière modification du revenu'; 