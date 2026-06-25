-- Casbin policy storage
CREATE TABLE IF NOT EXISTS casbin_rule (
    id     BIGSERIAL PRIMARY KEY,
    ptype  TEXT NOT NULL DEFAULT '',
    v0     TEXT NOT NULL DEFAULT '',
    v1     TEXT NOT NULL DEFAULT '',
    v2     TEXT NOT NULL DEFAULT '',
    v3     TEXT NOT NULL DEFAULT '',
    v4     TEXT NOT NULL DEFAULT '',
    v5     TEXT NOT NULL DEFAULT '',
    CONSTRAINT casbin_rule_unique UNIQUE (ptype, v0, v1, v2, v3, v4, v5)
);

CREATE INDEX IF NOT EXISTS idx_casbin_ptype ON casbin_rule (ptype);
CREATE INDEX IF NOT EXISTS idx_casbin_v0    ON casbin_rule (v0);
