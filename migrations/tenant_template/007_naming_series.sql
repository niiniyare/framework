-- Naming series counters: one row per unique series key (e.g. "SE-2025-01-{SEQ}")
-- Incremented atomically via INSERT ... ON CONFLICT DO UPDATE.

CREATE TABLE IF NOT EXISTS naming_series_counters (
    series_key    TEXT    PRIMARY KEY,
    current_value BIGINT  NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
