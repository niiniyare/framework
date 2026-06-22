-- Tenant schema template: applied to every new tenant's schema.
-- This file is run against "tenant_{slug}" schema.

-- Custom entity records (JSONB-backed, tenant-defined entities)
CREATE TABLE IF NOT EXISTS custom_records (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_name TEXT NOT NULL,
    tenant_id   TEXT NOT NULL,
    data        JSONB NOT NULL DEFAULT '{}',
    doc_status  TEXT NOT NULL DEFAULT 'Draft' CHECK (doc_status IN ('Draft', 'Submitted', 'Cancelled')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- GIN index for JSONB field queries
CREATE INDEX IF NOT EXISTS custom_records_data_gin ON custom_records USING GIN (data);
-- Composite index for tenant-scoped entity queries
CREATE INDEX IF NOT EXISTS custom_records_entity_tenant ON custom_records (entity_name, tenant_id, created_at DESC);

-- Custom field definitions (metadata for tenant-defined fields on any entity)
CREATE TABLE IF NOT EXISTS custom_field_defs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   TEXT NOT NULL,
    entity_name TEXT NOT NULL,
    field_key   TEXT NOT NULL,  -- the JSON key used in data column
    field_label TEXT NOT NULL,
    field_type  TEXT NOT NULL,  -- subset of framework field types
    required    BOOLEAN NOT NULL DEFAULT false,
    options     JSONB NOT NULL DEFAULT '{}', -- choices, max_len, regex, etc.
    sort_order  INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, entity_name, field_key)
);

-- Sessions stored per tenant (cross-tenant sessions live in system Redis)
-- This table is a fallback when Redis is unavailable.
CREATE TABLE IF NOT EXISTS sessions (
    id          TEXT PRIMARY KEY,
    tenant_id   TEXT NOT NULL,
    user_id     TEXT NOT NULL,
    roles       TEXT[] NOT NULL DEFAULT '{}',
    is_super    BOOLEAN NOT NULL DEFAULT false,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS sessions_tenant_user ON sessions (tenant_id, user_id);
CREATE INDEX IF NOT EXISTS sessions_expires ON sessions (expires_at);
