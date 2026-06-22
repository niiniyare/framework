-- Settings module: config definitions, tenant configs, templates, audit

-- System-wide config definitions (schema for all possible config keys)
CREATE TABLE IF NOT EXISTS config_definitions (
    module_name           VARCHAR(50)  NOT NULL,
    config_key            VARCHAR(100) NOT NULL,
    config_type           VARCHAR(20)  NOT NULL CHECK (config_type IN ('string','integer','boolean','decimal','json')),
    default_value         JSONB,
    validation_rules      JSONB,
    description           TEXT,
    required_feature_flag VARCHAR(100),
    created_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (module_name, config_key)
);

-- Seed common config definitions
INSERT INTO config_definitions (module_name, config_key, config_type, default_value, description) VALUES
    ('iam',       'session_ttl_hours',       'integer', '8',     'Session TTL in hours'),
    ('iam',       'max_concurrent_sessions', 'integer', '5',     'Max active sessions per user'),
    ('iam',       'password_min_length',     'integer', '8',     'Minimum password length'),
    ('finance',   'invoice_prefix',          'string',  '"INV-"','Invoice naming prefix'),
    ('finance',   'auto_approval_limit',     'decimal', '0',     'Auto-approve journal entries below this amount'),
    ('finance',   'fiscal_year_start_month', 'integer', '1',     'Month fiscal year starts (1=January)'),
    ('inventory', 'valuation_method',        'string',  '"FIFO"','Default stock valuation method'),
    ('inventory', 'reorder_notifications',   'boolean', 'true',  'Send reorder level notifications'),
    ('hr',        'overtime_threshold',      'decimal', '40',    'Weekly hours before overtime kicks in'),
    ('hr',        'overtime_multiplier',     'decimal', '1.5',   'Overtime pay multiplier')
ON CONFLICT (module_name, config_key) DO NOTHING;

-- Per-tenant configuration overrides (flat module.key JSONB map)
CREATE TABLE IF NOT EXISTS tenant_configurations_settings (
    tenant_id             UUID        PRIMARY KEY,
    settings              JSONB       NOT NULL DEFAULT '{}'::jsonb,
    settings_version      INTEGER     NOT NULL DEFAULT 0,
    last_template_applied UUID,
    template_applied_at   TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed empty configuration for this tenant
INSERT INTO tenant_configurations_settings (tenant_id)
SELECT gen_random_uuid()
WHERE NOT EXISTS (SELECT 1 FROM tenant_configurations_settings)
ON CONFLICT DO NOTHING;

-- Configuration templates
CREATE TABLE IF NOT EXISTS configuration_templates (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    scope                   VARCHAR(10) NOT NULL DEFAULT 'TENANT'
                                CHECK (scope IN ('SYSTEM', 'TENANT')),
    name                    VARCHAR(255) NOT NULL,
    category                VARCHAR(50),
    description             TEXT,
    version                 VARCHAR(20)  NOT NULL DEFAULT '1.0',
    configurations          JSONB        NOT NULL DEFAULT '{}'::jsonb,
    applicable_tenant_types TEXT[],
    required_feature_flags  TEXT[],
    conflict_resolution     VARCHAR(20)  NOT NULL DEFAULT 'merge'
                                CHECK (conflict_resolution IN ('merge', 'replace', 'preserve')),
    is_active               BOOLEAN      NOT NULL DEFAULT TRUE,
    created_by              UUID,
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Configuration change audit (append-only)
CREATE TABLE IF NOT EXISTS configuration_audit (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id       UUID,
    module_name     VARCHAR(50)  NOT NULL,
    config_key_name VARCHAR(100) NOT NULL,
    old_value       JSONB,
    new_value       JSONB        NOT NULL,
    source          VARCHAR(20),   -- 'tenant' | 'entity' | 'template'
    operation       VARCHAR(20),   -- 'SET' | 'DELETE' | 'TEMPLATE_APPLY'
    user_id         UUID         NOT NULL,
    session_id      UUID,
    correlation_id  UUID,
    applied_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_config_audit_module
    ON configuration_audit(module_name, config_key_name, applied_at DESC);
CREATE INDEX IF NOT EXISTS idx_config_audit_entity
    ON configuration_audit(entity_id, applied_at DESC) WHERE entity_id IS NOT NULL;
