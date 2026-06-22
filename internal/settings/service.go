// Package settings implements hierarchical configuration resolution.
// Precedence (highest wins): entity (2) > tenant (1) > system default (0).
// Deletions at a level fall back to the level below.
package settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"awo.so/framework/internal/redis"
)

const (
	tenantCachePrefix = "config:"
	tenantCacheTTL    = 10 * time.Minute
)

// ConfigValue is the resolved value for a single config key.
type ConfigValue struct {
	ModuleName string
	ConfigKey  string
	Value      any
	Source     string // "system" | "tenant" | "entity"
	ConfigType string // "string" | "integer" | "boolean" | "decimal" | "json"
}

// UpdateParams sets a single config key at tenant or entity level.
type UpdateParams struct {
	TenantID  uuid.UUID
	EntityID  *uuid.UUID // nil = tenant-level
	Module    string
	Key       string
	Value     any
	UserID    uuid.UUID
	SessionID *uuid.UUID
}

// Service provides config resolution and writes.
type Service struct {
	db    *pgxpool.Pool
	redis *redis.Client
	log   *slog.Logger
}

// New creates a settings Service.
func New(db *pgxpool.Pool, r *redis.Client, log *slog.Logger) *Service {
	return &Service{db: db, redis: r, log: log}
}

// ResolveForTenant returns all effective settings for a tenant as a flat
// "module.key" → string map for embedding into sessions.
func (s *Service) ResolveForTenant(ctx context.Context, tenantID uuid.UUID) (map[string]string, error) {
	cacheKey := tenantCachePrefix + tenantID.String() + ":tenant"

	if data, err := s.redis.Get(ctx, cacheKey); err == nil {
		var m map[string]string
		if json.Unmarshal([]byte(data), &m) == nil {
			return m, nil
		}
	}

	// Load all system defaults
	sysRows, err := s.db.Query(ctx, `
		SELECT module_name, config_key, default_value, config_type
		FROM config_definitions
		WHERE default_value IS NOT NULL`)
	if err != nil {
		return nil, fmt.Errorf("settings: load defaults: %w", err)
	}
	defer sysRows.Close()

	out := make(map[string]string)
	for sysRows.Next() {
		var mod, key, typ string
		var val []byte
		if err := sysRows.Scan(&mod, &key, &val, &typ); err != nil {
			return nil, err
		}
		out[mod+"."+key] = jsonToString(val, typ)
	}
	if err := sysRows.Err(); err != nil {
		return nil, err
	}

	// Load tenant overrides (JSONB map) and apply
	row := s.db.QueryRow(ctx, `
		SELECT settings FROM tenant_configurations_settings WHERE tenant_id = $1`, tenantID)
	var settingsJSON []byte
	if err := row.Scan(&settingsJSON); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("settings: load tenant: %w", err)
	}
	if settingsJSON != nil {
		var overrides map[string]any
		if err := json.Unmarshal(settingsJSON, &overrides); err == nil {
			for k, v := range overrides {
				out[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	// Cache
	if data, err := json.Marshal(out); err == nil {
		_ = s.redis.SetEX(ctx, cacheKey, string(data), tenantCacheTTL)
	}
	return out, nil
}

// Get resolves a single config key respecting entity > tenant > system precedence.
func (s *Service) Get(ctx context.Context, tenantID uuid.UUID, entityID *uuid.UUID, module, key string) (*ConfigValue, error) {
	// Try entity level
	if entityID != nil {
		if v := s.getEntityValue(ctx, *entityID, module, key); v != nil {
			return v, nil
		}
	}
	// Try tenant level
	if v := s.getTenantValue(ctx, tenantID, module, key); v != nil {
		return v, nil
	}
	// System default
	return s.getSystemDefault(ctx, module, key)
}

// Set writes a config key at tenant or entity level and records audit.
func (s *Service) Set(ctx context.Context, p UpdateParams) error {
	fullKey := p.Module + "." + p.Key
	valJSON, err := json.Marshal(p.Value)
	if err != nil {
		return fmt.Errorf("settings set: marshal value: %w", err)
	}

	// Read old value for audit
	old, _ := s.Get(ctx, p.TenantID, p.EntityID, p.Module, p.Key)

	var source string
	if p.EntityID != nil {
		// Entity-level: jsonb_set on entities.settings
		_, err = s.db.Exec(ctx, `
			UPDATE entities
			SET settings = jsonb_set(
			    COALESCE(settings, '{}'::jsonb), $1, $2, true),
			    updated_at = NOW()
			WHERE id = $3`,
			"{"+fullKey+"}", valJSON, *p.EntityID)
		source = "entity"
	} else {
		// Tenant-level: jsonb_set on tenant_configurations_settings.settings
		_, err = s.db.Exec(ctx, `
			INSERT INTO tenant_configurations_settings (tenant_id, settings, settings_version)
			VALUES ($1, jsonb_build_object($2, $3::jsonb), 1)
			ON CONFLICT (tenant_id) DO UPDATE
			SET settings         = jsonb_set(
			        COALESCE(tenant_configurations_settings.settings, '{}'::jsonb),
			        $4, $3::jsonb, true),
			    settings_version = tenant_configurations_settings.settings_version + 1,
			    updated_at       = NOW()`,
			p.TenantID, fullKey, string(valJSON), "{"+fullKey+"}")
		source = "tenant"
	}
	if err != nil {
		return fmt.Errorf("settings set: %w", err)
	}

	// Invalidate cache
	_ = s.redis.Del(ctx, tenantCachePrefix+p.TenantID.String()+":tenant")

	// Audit
	s.writeAudit(ctx, p, old, string(valJSON), source, "SET")
	return nil
}

// Delete removes a config key at tenant or entity level (falls back to lower level).
func (s *Service) Delete(ctx context.Context, p UpdateParams) error {
	fullKey := p.Module + "." + p.Key
	old, _ := s.Get(ctx, p.TenantID, p.EntityID, p.Module, p.Key)

	var err error
	var source string
	if p.EntityID != nil {
		_, err = s.db.Exec(ctx, `
			UPDATE entities SET settings = settings - $1, updated_at = NOW()
			WHERE id = $2`, fullKey, *p.EntityID)
		source = "entity"
	} else {
		_, err = s.db.Exec(ctx, `
			UPDATE tenant_configurations_settings
			SET settings         = settings - $1,
			    settings_version = settings_version + 1,
			    updated_at       = NOW()
			WHERE tenant_id = $2`, fullKey, p.TenantID)
		source = "tenant"
	}
	if err != nil {
		return fmt.Errorf("settings delete: %w", err)
	}
	_ = s.redis.Del(ctx, tenantCachePrefix+p.TenantID.String()+":tenant")
	s.writeAudit(ctx, p, old, "null", source, "DELETE")
	return nil
}

// --- internal helpers ---

func (s *Service) getEntityValue(ctx context.Context, entityID uuid.UUID, module, key string) *ConfigValue {
	fullKey := module + "." + key
	row := s.db.QueryRow(ctx, `
		SELECT settings->$1 FROM entities WHERE id = $2 AND settings ? $1`,
		fullKey, entityID)
	var raw []byte
	if err := row.Scan(&raw); err != nil || raw == nil {
		return nil
	}
	typ := s.getConfigType(ctx, module, key)
	return &ConfigValue{
		ModuleName: module, ConfigKey: key,
		Value:      jsonToAny(raw, typ),
		Source:     "entity", ConfigType: typ,
	}
}

func (s *Service) getTenantValue(ctx context.Context, tenantID uuid.UUID, module, key string) *ConfigValue {
	fullKey := module + "." + key
	row := s.db.QueryRow(ctx, `
		SELECT settings->$1 FROM tenant_configurations_settings
		WHERE tenant_id = $2 AND settings ? $1`, fullKey, tenantID)
	var raw []byte
	if err := row.Scan(&raw); err != nil || raw == nil {
		return nil
	}
	typ := s.getConfigType(ctx, module, key)
	return &ConfigValue{
		ModuleName: module, ConfigKey: key,
		Value:      jsonToAny(raw, typ),
		Source:     "tenant", ConfigType: typ,
	}
}

func (s *Service) getSystemDefault(ctx context.Context, module, key string) (*ConfigValue, error) {
	row := s.db.QueryRow(ctx, `
		SELECT default_value, config_type FROM config_definitions
		WHERE module_name = $1 AND config_key = $2`, module, key)
	var raw []byte
	var typ string
	if err := row.Scan(&raw, &typ); err != nil {
		return nil, nil // unknown key
	}
	return &ConfigValue{
		ModuleName: module, ConfigKey: key,
		Value:      jsonToAny(raw, typ),
		Source:     "system", ConfigType: typ,
	}, nil
}

func (s *Service) getConfigType(ctx context.Context, module, key string) string {
	row := s.db.QueryRow(ctx, `
		SELECT config_type FROM config_definitions
		WHERE module_name = $1 AND config_key = $2`, module, key)
	var typ string
	_ = row.Scan(&typ)
	if typ == "" {
		return "string"
	}
	return typ
}

func (s *Service) writeAudit(ctx context.Context, p UpdateParams, old *ConfigValue, newVal, source, op string) {
	var oldVal string
	if old != nil {
		b, _ := json.Marshal(old.Value)
		oldVal = string(b)
	}
	_, err := s.db.Exec(ctx, `
		INSERT INTO configuration_audit
		    (entity_id, module_name, config_key_name, old_value, new_value,
		     source, operation, user_id, session_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		p.EntityID, p.Module, p.Key,
		nullJSON(oldVal), []byte(newVal),
		source, op, p.UserID, p.SessionID)
	if err != nil {
		s.log.Warn("settings: audit write failed", slog.Any("err", err))
	}
}

func jsonToString(raw []byte, typ string) string {
	if raw == nil {
		return ""
	}
	// JSONB values from PG come with quotes for strings; strip them
	s := string(raw)
	if typ == "string" {
		v, err := strconv.Unquote(s)
		if err == nil {
			return v
		}
	}
	return s
}

func jsonToAny(raw []byte, typ string) any {
	if raw == nil {
		return nil
	}
	var v any
	_ = json.Unmarshal(raw, &v)
	return v
}

func nullJSON(s string) any {
	if s == "" {
		return nil
	}
	return []byte(s)
}
