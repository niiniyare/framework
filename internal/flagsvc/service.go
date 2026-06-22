// Package flagsvc implements the feature flag service.
// Flags are resolved as: tenant_override COALESCE default_value.
// Tenant sessions cache the full resolved flag map for their TTL.
package flagsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"awo.so/framework/internal/redis"
)

const (
	flagCachePrefix = "flags:"
	flagCacheTTL    = 5 * time.Minute
)

// FlagDefinition is a system-wide feature flag specification.
type FlagDefinition struct {
	ID           uuid.UUID
	FlagKey      string
	Label        string
	Description  string
	DefaultValue bool
	IsSystem     bool // only platform operators can toggle
}

// TenantFlag is a per-tenant override of a flag definition.
type TenantFlag struct {
	FlagID  uuid.UUID
	FlagKey string
	Enabled bool
	SetBy   *uuid.UUID
	SetAt   time.Time
}

// SetFlagParams sets or clears a tenant flag override.
type SetFlagParams struct {
	TenantID uuid.UUID
	FlagKey  string
	Enabled  bool
	SetBy    uuid.UUID
}

// SessionInvalidator evicts sessions when flags change.
type SessionInvalidator interface {
	InvalidateByTenant(ctx context.Context, tenantID uuid.UUID) error
}

// Service provides flag resolution and management.
type Service struct {
	db          *pgxpool.Pool
	redis       *redis.Client
	invalidator SessionInvalidator
	log         *slog.Logger
}

// New creates a flag Service.
func New(db *pgxpool.Pool, r *redis.Client, inv SessionInvalidator, log *slog.Logger) *Service {
	return &Service{db: db, redis: r, invalidator: inv, log: log}
}

// ResolveForTenant returns the effective flag map for a tenant.
// Redis cache: "flags:{tenantID}" → JSON map[string]bool, TTL 5 min.
func (s *Service) ResolveForTenant(ctx context.Context, tenantID uuid.UUID) (map[string]bool, error) {
	cacheKey := flagCachePrefix + tenantID.String()

	// Cache-aside
	if data, err := s.redis.Get(ctx, cacheKey); err == nil {
		var flags map[string]bool
		if json.Unmarshal([]byte(data), &flags) == nil {
			return flags, nil
		}
	}

	// DB: COALESCE(tenant_override, system_default)
	rows, err := s.db.Query(ctx, `
		SELECT ffd.flag_key,
		       COALESCE(tff.enabled, ffd.default_value) AS effective_value
		FROM feature_flag_definitions ffd
		LEFT JOIN tenant_feature_flags tff ON tff.flag_id = ffd.id
		ORDER BY ffd.flag_key`)
	if err != nil {
		return nil, fmt.Errorf("flagsvc: resolve for tenant %s: %w", tenantID, err)
	}
	defer rows.Close()

	flags := make(map[string]bool)
	for rows.Next() {
		var key string
		var val bool
		if err := rows.Scan(&key, &val); err != nil {
			return nil, fmt.Errorf("flagsvc: scan: %w", err)
		}
		flags[key] = val
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache
	if data, err := json.Marshal(flags); err == nil {
		_ = s.redis.SetEX(ctx, cacheKey, string(data), flagCacheTTL)
	}
	return flags, nil
}

// Set creates or updates a tenant flag override, then invalidates sessions.
func (s *Service) Set(ctx context.Context, p SetFlagParams) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO tenant_feature_flags (flag_id, flag_key, enabled, set_by, set_at)
		SELECT id, flag_key, $2, $3, NOW()
		FROM feature_flag_definitions
		WHERE flag_key = $1
		ON CONFLICT (flag_id) DO UPDATE
		    SET enabled = EXCLUDED.enabled,
		        set_by  = EXCLUDED.set_by,
		        set_at  = EXCLUDED.set_at`,
		p.FlagKey, p.Enabled, p.SetBy)
	if err != nil {
		return fmt.Errorf("flagsvc: set %s: %w", p.FlagKey, err)
	}

	// Evict flag cache
	_ = s.redis.Del(ctx, flagCachePrefix+p.TenantID.String())

	// Evict all tenant sessions (force re-login to pick up new flags)
	if s.invalidator != nil {
		if err := s.invalidator.InvalidateByTenant(ctx, p.TenantID); err != nil {
			s.log.Warn("flagsvc: session invalidation failed after flag change",
				slog.String("tenant_id", p.TenantID.String()),
				slog.String("flag_key", p.FlagKey),
				slog.Any("err", err))
		}
	}
	return nil
}

// Reset removes a tenant override, restoring the system default.
func (s *Service) Reset(ctx context.Context, tenantID uuid.UUID, flagKey string) error {
	_, err := s.db.Exec(ctx, `
		DELETE FROM tenant_feature_flags
		WHERE flag_key = $1`, flagKey)
	if err != nil {
		return fmt.Errorf("flagsvc: reset %s: %w", flagKey, err)
	}
	_ = s.redis.Del(ctx, flagCachePrefix+tenantID.String())
	return nil
}

// ListDefinitions returns all defined flags (system catalogue).
func (s *Service) ListDefinitions(ctx context.Context) ([]*FlagDefinition, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, flag_key, label, COALESCE(description,''), default_value, is_system
		FROM feature_flag_definitions
		ORDER BY flag_key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var defs []*FlagDefinition
	for rows.Next() {
		var d FlagDefinition
		if err := rows.Scan(&d.ID, &d.FlagKey, &d.Label, &d.Description, &d.DefaultValue, &d.IsSystem); err != nil {
			return nil, err
		}
		defs = append(defs, &d)
	}
	return defs, rows.Err()
}
