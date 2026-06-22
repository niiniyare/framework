-- Module-Resource-Action registry (system-wide, not per-tenant)
-- Provides the permission key taxonomy: {module}.{resource}.{action}

CREATE TABLE IF NOT EXISTS modules (
    id        UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    slug      TEXT    UNIQUE NOT NULL,
    label     TEXT    NOT NULL,
    icon      TEXT,
    nav_order INT     NOT NULL DEFAULT 999,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS resources (
    id        UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID    NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    slug      TEXT    NOT NULL,
    label     TEXT    NOT NULL,
    nav_url   TEXT,
    nav_order INT     NOT NULL DEFAULT 999,
    UNIQUE (module_id, slug)
);

CREATE TABLE IF NOT EXISTS actions (
    id          UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    resource_id UUID    NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    slug        TEXT    NOT NULL,
    label       TEXT    NOT NULL,
    http_method TEXT,
    UNIQUE (resource_id, slug)
);

-- Seed built-in modules
INSERT INTO modules (slug, label, nav_order) VALUES
    ('platform',   'Platform',   0),
    ('finance',    'Finance',    10),
    ('inventory',  'Inventory',  20),
    ('hr',         'HR',         30),
    ('crm',        'CRM',        40),
    ('forecourt',  'Forecourt',  50)
ON CONFLICT (slug) DO NOTHING;

-- Domain events outbox (written in same tx as operation; worker delivers async)
CREATE TABLE IF NOT EXISTS domain_events (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID        REFERENCES tenants(id),  -- NULL for platform events
    event_type   TEXT        NOT NULL,
    payload      JSONB       NOT NULL DEFAULT '{}',
    occurred_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_domain_events_unprocessed
    ON domain_events(occurred_at) WHERE processed_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_domain_events_tenant
    ON domain_events(tenant_id) WHERE tenant_id IS NOT NULL;

-- System-level audit events (append-only, no UPDATE/DELETE)
CREATE TABLE IF NOT EXISTS audit_events (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID        REFERENCES tenants(id),
    event_type     TEXT        NOT NULL,
    actor_subject  TEXT        NOT NULL,
    actor_domain   TEXT        NOT NULL,
    target_subject TEXT,
    resource       TEXT,
    action         TEXT,
    outcome        TEXT,        -- 'allowed' | 'denied'
    metadata       JSONB       NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_events_tenant ON audit_events(tenant_id, created_at DESC)
    WHERE tenant_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_audit_events_type   ON audit_events(event_type, created_at DESC);
