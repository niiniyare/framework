-- Platform module: IAM tables (per-tenant schema)

-- Users within this tenant
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL,
    full_name       TEXT NOT NULL,
    password_hash   TEXT NOT NULL DEFAULT '',
    is_active       BOOLEAN NOT NULL DEFAULT true,
    is_super        BOOLEAN NOT NULL DEFAULT false,
    phone           TEXT,
    avatar_url      TEXT,
    last_login_at   TIMESTAMPTZ,
    department_id   UUID, -- FK to departments added after that table exists
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (email)
);

-- Roles
CREATE TABLE IF NOT EXISTS roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    label       TEXT NOT NULL,
    description TEXT,
    is_system   BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed system roles
INSERT INTO roles (name, label, is_system) VALUES
    ('admin',           'Administrator',     true),
    ('viewer',          'Viewer',            true),
    ('finance_manager', 'Finance Manager',   true),
    ('accountant',      'Accountant',        true),
    ('hr_manager',      'HR Manager',        true),
    ('inventory_clerk', 'Inventory Clerk',   true),
    ('cashier',         'Cashier',           true)
ON CONFLICT (name) DO NOTHING;

-- Role assignments (user ↔ role)
CREATE TABLE IF NOT EXISTS role_assignments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id      UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    granted_by   UUID REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, role_id)
);

-- Tenant-level permission overrides (allow/deny per role+entity+action)
CREATE TABLE IF NOT EXISTS permission_rules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id     UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    entity_name TEXT NOT NULL,
    action      TEXT NOT NULL CHECK (action IN ('read','create','update','delete','submit','cancel','amend')),
    effect      TEXT NOT NULL DEFAULT 'allow' CHECK (effect IN ('allow','deny')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (role_id, entity_name, action)
);

CREATE INDEX IF NOT EXISTS permission_rules_role_entity ON permission_rules (role_id, entity_name);
