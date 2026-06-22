-- Feature flag definitions and per-tenant overrides

CREATE TABLE IF NOT EXISTS feature_flag_definitions (
    id            UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_key      TEXT    UNIQUE NOT NULL,   -- 'finance' | 'finance.transactions'
    label         TEXT    NOT NULL,
    description   TEXT,
    default_value BOOLEAN NOT NULL DEFAULT FALSE,
    is_system     BOOLEAN NOT NULL DEFAULT FALSE  -- only platform operators can toggle
);

-- Seed built-in feature flags (module-level)
INSERT INTO feature_flag_definitions (flag_key, label, default_value) VALUES
    ('finance',    'Finance Module',    TRUE),
    ('inventory',  'Inventory Module',  FALSE),
    ('hr',         'HR Module',         FALSE),
    ('crm',        'CRM Module',        FALSE),
    ('forecourt',  'Forecourt Module',  FALSE)
ON CONFLICT (flag_key) DO NOTHING;

-- Per-tenant overrides
CREATE TABLE IF NOT EXISTS tenant_feature_flags (
    id        UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_id   UUID        NOT NULL REFERENCES feature_flag_definitions(id) ON DELETE CASCADE,
    flag_key  TEXT        NOT NULL,   -- denormalised for fast lookup
    enabled   BOOLEAN     NOT NULL,
    set_by    UUID,
    set_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (flag_id)
);

CREATE INDEX IF NOT EXISTS idx_tenant_flags_flag ON tenant_feature_flags(flag_id);
