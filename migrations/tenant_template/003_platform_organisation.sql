-- Platform module: Organisation and Department tables (per-tenant schema)

CREATE TABLE IF NOT EXISTS departments (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                TEXT NOT NULL,
    code                TEXT UNIQUE,
    parent_department_id UUID REFERENCES departments(id),
    head_of_department_id UUID REFERENCES users(id),
    cost_centre_code    TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add FK from users.department_id now that departments table exists
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS department_id UUID REFERENCES departments(id);

CREATE TABLE IF NOT EXISTS organisations (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    legal_name          TEXT NOT NULL,
    trading_name        TEXT,
    registration_number TEXT,
    tax_pin             TEXT,          -- KRA PIN, encrypted at app layer
    vat_number          TEXT,
    address_line1       TEXT,
    address_line2       TEXT,
    city                TEXT,
    country             TEXT NOT NULL DEFAULT 'KE',
    phone               TEXT,
    email               TEXT,
    website             TEXT,
    logo                TEXT,          -- object storage key
    currency            TEXT NOT NULL DEFAULT 'KES',
    fiscal_year_start   TEXT NOT NULL DEFAULT 'January',
    timezone            TEXT NOT NULL DEFAULT 'Africa/Nairobi',
    date_format         TEXT NOT NULL DEFAULT 'DD/MM/YYYY',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed a placeholder organisation row (singleton — populated via onboarding wizard)
INSERT INTO organisations (legal_name, currency) VALUES ('My Organisation', 'KES')
ON CONFLICT DO NOTHING;
