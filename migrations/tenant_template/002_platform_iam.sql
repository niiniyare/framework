-- Platform module: IAM tables (per-tenant schema)

-- Users within this tenant
CREATE TABLE IF NOT EXISTS users (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID,       -- NULL for system-schema users; set for tenant-schema copies
    email                 TEXT        NOT NULL UNIQUE,
    full_name             TEXT        NOT NULL DEFAULT '',
    password_hash         TEXT        NOT NULL DEFAULT '',
    user_type             TEXT        NOT NULL DEFAULT 'INTERNAL'
                              CHECK (user_type IN ('INTERNAL','SYSADMIN','CUSTOMER','PORTAL','API')),
    mfa_enabled           BOOLEAN     NOT NULL DEFAULT FALSE,
    mfa_secret            TEXT,
    is_active             BOOLEAN     NOT NULL DEFAULT TRUE,
    is_suspended          BOOLEAN     NOT NULL DEFAULT FALSE,
    is_super              BOOLEAN     NOT NULL DEFAULT FALSE,
    failed_login_attempts INTEGER     NOT NULL DEFAULT 0,
    locked_until          TIMESTAMPTZ,
    last_login_at         TIMESTAMPTZ,
    phone                 TEXT,
    avatar_url            TEXT,
    department_id         UUID,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Roles
CREATE TABLE IF NOT EXISTS roles (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL UNIQUE,
    label       TEXT        NOT NULL,
    description TEXT,
    is_system   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed system roles
INSERT INTO roles (name, label, is_system) VALUES
    ('admin',           'Administrator',   TRUE),
    ('viewer',          'Viewer',          TRUE),
    ('finance_manager', 'Finance Manager', TRUE),
    ('accountant',      'Accountant',      TRUE),
    ('hr_manager',      'HR Manager',      TRUE),
    ('inventory_clerk', 'Inventory Clerk', TRUE),
    ('cashier',         'Cashier',         TRUE)
ON CONFLICT (name) DO NOTHING;

-- Role assignments (user ↔ role)
CREATE TABLE IF NOT EXISTS role_assignments (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id    UUID        NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    granted_by UUID        REFERENCES users(id),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_role_assignments_user ON role_assignments(user_id);

-- Tenant-level permission overrides (allow/deny per role+entity+action)
CREATE TABLE IF NOT EXISTS permission_rules (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id     UUID        NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    entity_name TEXT        NOT NULL,
    action      TEXT        NOT NULL
                    CHECK (action IN ('read','create','update','delete','submit','cancel','*')),
    effect      TEXT        NOT NULL DEFAULT 'allow' CHECK (effect IN ('allow','deny')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (role_id, entity_name, action)
);

CREATE INDEX IF NOT EXISTS idx_permission_rules_role_entity
    ON permission_rules(role_id, entity_name);

-- User sessions (per-tenant copy — system schema is authoritative)
CREATE TABLE IF NOT EXISTS user_sessions (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_type     TEXT        NOT NULL DEFAULT 'INTERNAL',
    session_token TEXT        NOT NULL UNIQUE,
    entity_scope  JSONB       NOT NULL DEFAULT '{"type":"all"}'::jsonb,
    configuration JSONB       NOT NULL DEFAULT '{}'::jsonb,
    ip_address    INET,
    user_agent    TEXT,
    is_active     BOOLEAN     NOT NULL DEFAULT TRUE,
    expires_at    TIMESTAMPTZ NOT NULL,
    last_seen_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_token
    ON user_sessions(session_token) WHERE is_active = TRUE;
