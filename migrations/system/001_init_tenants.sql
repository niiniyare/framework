-- System schema: core tables shared across all tenants.
-- Applied once to the system (public) schema.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ── Tenants ──────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS tenants (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    slug             TEXT        NOT NULL UNIQUE,
    name             TEXT        NOT NULL,
    email            TEXT        NOT NULL,
    subdomain        TEXT        UNIQUE,
    status           TEXT        NOT NULL DEFAULT 'PENDING'
                         CHECK (status IN ('PENDING','ACTIVE','SUSPENDED','ARCHIVED')),
    plan             TEXT        NOT NULL DEFAULT 'Basic'
                         CHECK (plan IN ('Basic','Professional','Enterprise')),
    industry         TEXT,
    company_size     TEXT,
    currency_code    CHAR(3)     NOT NULL DEFAULT 'KES',
    timezone         TEXT        NOT NULL DEFAULT 'Africa/Nairobi',
    settings         JSONB       NOT NULL DEFAULT '{}',
    metadata         JSONB       NOT NULL DEFAULT '{}',
    -- compat: generated schema name for schema-per-tenant routing
    schema_name      TEXT        NOT NULL GENERATED ALWAYS AS ('tenant_' || slug) STORED,
    deleted_at       TIMESTAMPTZ,
    last_activity_at TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tenants_slug      ON tenants(slug)      WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tenants_subdomain ON tenants(subdomain) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tenants_status    ON tenants(status);
CREATE INDEX IF NOT EXISTS idx_tenants_email     ON tenants(email);

-- Resource limits per tenant (created automatically on INSERT via trigger or by lifecycle)
CREATE TABLE IF NOT EXISTS tenant_configurations (
    id                     UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              UUID    NOT NULL UNIQUE REFERENCES tenants(id) ON DELETE CASCADE,
    max_users              INTEGER NOT NULL DEFAULT 10,
    max_entities           INTEGER NOT NULL DEFAULT 5,
    max_transactions_month INTEGER NOT NULL DEFAULT 1000,
    storage_quota_mb       INTEGER NOT NULL DEFAULT 1024,
    allowed_modules        JSONB   NOT NULL DEFAULT '["finance"]'::jsonb,
    accounting_method      TEXT    NOT NULL DEFAULT 'FIFO',
    fiscal_year_start_month INTEGER NOT NULL DEFAULT 1,
    password_policy        JSONB   NOT NULL DEFAULT '{}'::jsonb,
    api_rate_limits        JSONB   NOT NULL DEFAULT '{}'::jsonb
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tenant_configurations_tenant
    ON tenant_configurations(tenant_id);

-- ── Users ─────────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS users (
    id                    UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email                 TEXT        NOT NULL,
    full_name             TEXT        NOT NULL DEFAULT '',
    password_hash         TEXT        NOT NULL DEFAULT '',
    user_type             TEXT        NOT NULL DEFAULT 'INTERNAL'
                              CHECK (user_type IN ('INTERNAL','SYSADMIN','CUSTOMER','PORTAL','API')),
    mfa_enabled           BOOLEAN     NOT NULL DEFAULT FALSE,
    mfa_secret            TEXT,                             -- AES-256-GCM encrypted TOTP secret
    is_active             BOOLEAN     NOT NULL DEFAULT TRUE,
    is_suspended          BOOLEAN     NOT NULL DEFAULT FALSE,
    is_super              BOOLEAN     NOT NULL DEFAULT FALSE,
    failed_login_attempts INTEGER     NOT NULL DEFAULT 0,
    locked_until          TIMESTAMPTZ,
    last_login_at         TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, email)
);

CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_email  ON users(email);

-- ── Roles (legacy simple model — Casbin is the authoritative RBAC store) ─────

CREATE TABLE IF NOT EXISTS roles (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    label      TEXT        NOT NULL DEFAULT '',
    is_system  BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, name)
);

CREATE INDEX IF NOT EXISTS idx_roles_tenant ON roles(tenant_id);

-- ── Feature flags (system-level definitions) ─────────────────────────────────

CREATE TABLE IF NOT EXISTS feature_flags (
    id            UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_key      TEXT    NOT NULL UNIQUE,
    label         TEXT    NOT NULL,
    description   TEXT,
    default_value BOOLEAN NOT NULL DEFAULT FALSE,
    is_system     BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_feature_flags_key ON feature_flags(flag_key);

-- Per-tenant overrides of system feature flags
CREATE TABLE IF NOT EXISTS tenant_flag_overrides (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    flag_id    UUID        NOT NULL REFERENCES feature_flags(id) ON DELETE CASCADE,
    flag_key   TEXT        NOT NULL,
    enabled    BOOLEAN     NOT NULL,
    set_by     UUID        REFERENCES users(id),
    set_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, flag_id)
);

CREATE INDEX IF NOT EXISTS idx_tenant_flag_overrides_tenant ON tenant_flag_overrides(tenant_id);

-- Seed built-in module-level feature flags
INSERT INTO feature_flags (flag_key, label, default_value) VALUES
    ('finance',   'Finance Module',   TRUE),
    ('inventory', 'Inventory Module', FALSE),
    ('hr',        'HR Module',        FALSE),
    ('crm',       'CRM Module',       FALSE),
    ('forecourt', 'Forecourt Module', FALSE)
ON CONFLICT (flag_key) DO NOTHING;

-- ── Audit log (system-level) ─────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS audit_log (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID        REFERENCES tenants(id),
    user_id    UUID        REFERENCES users(id),
    action     TEXT        NOT NULL,
    entity     TEXT,
    record_id  UUID,
    diff       JSONB,
    ip_address INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_tenant ON audit_log(tenant_id, created_at DESC)
    WHERE tenant_id IS NOT NULL;
