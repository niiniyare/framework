-- System IAM: Casbin rule store, role assignments, user sessions, API keys

-- Casbin rule store (ptype=p for policies, ptype=g for role assignments)
CREATE TABLE IF NOT EXISTS casbin_rule (
    id    UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    ptype VARCHAR(10)  NOT NULL,
    v0    VARCHAR(256) NOT NULL DEFAULT '',
    v1    VARCHAR(256) NOT NULL DEFAULT '',
    v2    VARCHAR(256) NOT NULL DEFAULT '',
    v3    VARCHAR(256) NOT NULL DEFAULT '',
    v4    VARCHAR(256) NOT NULL DEFAULT '',
    v5    VARCHAR(256) NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_casbin_rule_unique
    ON casbin_rule(ptype, v0, v1, v2, v3, v4, v5);
CREATE INDEX IF NOT EXISTS idx_casbin_rule_ptype ON casbin_rule(ptype);

-- Role assignment audit table (Casbin g-rules are the source of truth;
-- this table records who assigned what and when for compliance.)
CREATE TABLE IF NOT EXISTS role_assignments (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID         REFERENCES tenants(id) ON DELETE CASCADE,  -- NULL for platform roles
    subject      VARCHAR(256) NOT NULL,
    role_name    VARCHAR(100) NOT NULL,
    domain       VARCHAR(256) NOT NULL,
    assigned_by  VARCHAR(256),
    delegated_by VARCHAR(256),
    expires_at   TIMESTAMPTZ,
    is_active    BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT role_assignments_unique UNIQUE (subject, role_name, domain)
);

CREATE INDEX IF NOT EXISTS idx_role_assignments_subject ON role_assignments(subject, domain);
CREATE INDEX IF NOT EXISTS idx_role_assignments_tenant  ON role_assignments(tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_role_assignments_expires ON role_assignments(expires_at) WHERE expires_at IS NOT NULL;

-- User sessions with pre-computed entity scope and configuration
CREATE TABLE IF NOT EXISTS user_sessions (
    id             UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID         NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id        UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_type      VARCHAR(20)  NOT NULL DEFAULT 'INTERNAL',  -- INTERNAL|SYSADMIN|CUSTOMER|PORTAL|API
    session_token  TEXT         NOT NULL UNIQUE,              -- SHA-256 hex of raw token; raw never stored
    principal_id   UUID,                                      -- non-null for portal users only
    entity_scope   JSONB        NOT NULL DEFAULT '{"type":"all"}'::jsonb,
    configuration  JSONB        NOT NULL DEFAULT '{}'::jsonb, -- {flags:{},settings:{},prefs:{}}
    ip_address     INET,
    user_agent     TEXT,
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    expires_at     TIMESTAMPTZ  NOT NULL,
    last_seen_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user    ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_tenant  ON user_sessions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token   ON user_sessions(session_token) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires ON user_sessions(expires_at) WHERE is_active = TRUE;

-- API keys — SHA-256 hash stored; raw key returned once on creation
CREATE TABLE IF NOT EXISTS api_keys (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    key_hash    TEXT        NOT NULL UNIQUE,  -- SHA-256 hex of raw key
    scopes      TEXT[]      NOT NULL DEFAULT '{}',
    created_by  UUID        REFERENCES users(id),
    expires_at  TIMESTAMPTZ,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_keys_tenant   ON api_keys(tenant_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);
