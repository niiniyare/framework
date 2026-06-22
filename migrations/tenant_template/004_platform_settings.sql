-- Platform module: Settings singleton and Feature Flag overrides (per-tenant schema)

CREATE TABLE IF NOT EXISTS settings (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Locale
    timezone                        TEXT NOT NULL DEFAULT 'Africa/Nairobi',
    date_format                     TEXT NOT NULL DEFAULT 'DD/MM/YYYY',
    time_format                     TEXT NOT NULL DEFAULT 'HH:mm',
    language                        TEXT NOT NULL DEFAULT 'en',
    currency                        TEXT NOT NULL DEFAULT 'KES',
    number_format                   TEXT NOT NULL DEFAULT '1,234.56',

    -- Module enablement
    enable_finance                  BOOLEAN NOT NULL DEFAULT true,
    enable_inventory                BOOLEAN NOT NULL DEFAULT false,
    enable_hr                       BOOLEAN NOT NULL DEFAULT false,
    enable_crm                      BOOLEAN NOT NULL DEFAULT false,
    enable_forecourt                BOOLEAN NOT NULL DEFAULT false,

    -- Security
    session_timeout_hours           INTEGER NOT NULL DEFAULT 24,
    max_concurrent_sessions         INTEGER NOT NULL DEFAULT 5,
    require_mfa                     BOOLEAN NOT NULL DEFAULT false,
    password_min_length             INTEGER NOT NULL DEFAULT 8,
    password_requires_uppercase     BOOLEAN NOT NULL DEFAULT false,
    password_requires_number        BOOLEAN NOT NULL DEFAULT true,

    -- KRA eTIMS
    etims_enabled                   BOOLEAN NOT NULL DEFAULT false,
    etims_taxpayer_pin              TEXT,
    etims_device_serial             TEXT,
    etims_environment               TEXT NOT NULL DEFAULT 'sandbox',

    -- NEMA
    nema_tank_variance_threshold_pct NUMERIC(5,2) NOT NULL DEFAULT 0.50,

    -- Rate limits
    api_rate_limit_per_minute       INTEGER NOT NULL DEFAULT 300,

    -- Branding
    primary_color                   TEXT NOT NULL DEFAULT '#1890ff',
    logo_url                        TEXT,

    created_at                      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed the singleton settings row
INSERT INTO settings DEFAULT VALUES
ON CONFLICT DO NOTHING;

-- Feature flag per-tenant overrides
CREATE TABLE IF NOT EXISTS tenant_flag_overrides (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_name        TEXT NOT NULL UNIQUE,  -- references system feature_flags.name
    boolean_value    BOOLEAN NOT NULL DEFAULT false,
    string_value     TEXT,
    percentage_value INTEGER NOT NULL DEFAULT 0 CHECK (percentage_value BETWEEN 0 AND 100),
    reason           TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
