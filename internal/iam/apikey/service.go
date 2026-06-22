// Package apikey manages API key creation, validation and revocation.
package apikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"awo.so/framework/internal/iam/domain"
	"awo.so/framework/internal/redis"
)

const (
	apiKeyPrefix    = "apikey:"
	apiKeyCacheTTL  = 5 * time.Minute
	rawKeyPrefix    = "eak_"
)

// FlagResolver resolves feature flags for a tenant.
type FlagResolver interface {
	ResolveForTenant(ctx context.Context, tenantID uuid.UUID) (map[string]bool, error)
}

// SettingResolver resolves settings for a tenant.
type SettingResolver interface {
	ResolveForTenant(ctx context.Context, tenantID uuid.UUID) (map[string]string, error)
}

// Service handles API key lifecycle.
type Service struct {
	db       *pgxpool.Pool
	redis    *redis.Client
	flags    FlagResolver
	settings SettingResolver
	log      *slog.Logger
}

// New creates an API key Service.
func New(db *pgxpool.Pool, r *redis.Client, flags FlagResolver, settings SettingResolver, log *slog.Logger) *Service {
	return &Service{db: db, redis: r, flags: flags, settings: settings, log: log}
}

// Create generates a new API key. Returns the domain.APIKey and the raw token.
// The raw token is only returned once — it is never stored.
func (s *Service) Create(ctx context.Context, p domain.CreateAPIKeyParams) (*domain.APIKey, string, error) {
	raw, hash, err := generateAPIKey()
	if err != nil {
		return nil, "", fmt.Errorf("apikey: generate: %w", err)
	}

	key := &domain.APIKey{
		ID:        uuid.New(),
		TenantID:  p.TenantID,
		Name:      p.Name,
		KeyHash:   hash,
		Scopes:    p.Scopes,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: time.Now(),
	}
	createdBy := p.CreatedBy
	key.CreatedBy = &createdBy

	_, err = s.db.Exec(ctx, `
		INSERT INTO api_keys (id, tenant_id, name, key_hash, scopes, created_by, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		key.ID, key.TenantID, key.Name, key.KeyHash, key.Scopes, key.CreatedBy, key.ExpiresAt)
	if err != nil {
		return nil, "", fmt.Errorf("apikey: persist: %w", err)
	}
	return key, raw, nil
}

// Validate resolves a raw API key to a ResolvedSession.
// Returns (nil, nil) if the key does not exist or is invalid.
func (s *Service) Validate(ctx context.Context, rawToken string) (*domain.ResolvedSession, error) {
	hash := sha256Hex(rawToken)
	cacheKey := apiKeyPrefix + hash

	// Cache-aside: check Redis first (5-min TTL)
	if data, err := s.redis.Get(ctx, cacheKey); err == nil {
		var resolved domain.ResolvedSession
		if json.Unmarshal([]byte(data), &resolved) == nil {
			return &resolved, nil
		}
	}

	// DB lookup
	row := s.db.QueryRow(ctx, `
		SELECT id, tenant_id, name, scopes, created_by, expires_at, revoked_at, created_at
		FROM api_keys
		WHERE key_hash = $1`, hash)

	var key domain.APIKey
	if err := row.Scan(&key.ID, &key.TenantID, &key.Name, &key.Scopes,
		&key.CreatedBy, &key.ExpiresAt, &key.RevokedAt, &key.CreatedAt); err != nil {
		return nil, nil // not found
	}

	if !key.IsValid() {
		return nil, nil
	}

	// Build ResolvedSession
	flags, _ := s.flags.ResolveForTenant(ctx, key.TenantID)
	settings, _ := s.settings.ResolveForTenant(ctx, key.TenantID)

	var userID uuid.UUID
	if key.CreatedBy != nil {
		userID = *key.CreatedBy
	}
	resolved := &domain.ResolvedSession{
		UserID:    userID,
		UserType:  domain.UserTypeAPI,
		TenantID:  key.TenantID,
		EntityScope: domain.EntityScope{Type: domain.EntityScopeAll},
		Configuration: domain.Configuration{
			Flags:    orEmptyFlags(flags),
			Settings: orEmptySettings(settings),
			Prefs:    map[string]string{},
		},
	}

	// Cache with 5-min TTL (revocation gap accepted by design)
	if data, err := json.Marshal(resolved); err == nil {
		_ = s.redis.SetEX(ctx, cacheKey, string(data), apiKeyCacheTTL)
	}

	return resolved, nil
}

// Revoke marks the key revoked in DB and evicts from Redis.
func (s *Service) Revoke(ctx context.Context, keyID uuid.UUID) error {
	var hash string
	row := s.db.QueryRow(ctx, `
		UPDATE api_keys SET revoked_at = NOW()
		WHERE id = $1 RETURNING key_hash`, keyID)
	if err := row.Scan(&hash); err != nil {
		return fmt.Errorf("apikey: revoke: %w", err)
	}
	_ = s.redis.Del(ctx, apiKeyPrefix+hash)
	return nil
}

// List returns all API keys for a tenant (no key_hash exposed).
func (s *Service) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.APIKey, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, tenant_id, name, scopes, created_by, expires_at, revoked_at, created_at
		FROM api_keys
		WHERE tenant_id = $1
		ORDER BY created_at DESC`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []*domain.APIKey
	for rows.Next() {
		var k domain.APIKey
		if err := rows.Scan(&k.ID, &k.TenantID, &k.Name, &k.Scopes,
			&k.CreatedBy, &k.ExpiresAt, &k.RevokedAt, &k.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, &k)
	}
	return keys, rows.Err()
}

// --- helpers ---

func generateAPIKey() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	raw = rawKeyPrefix + hex.EncodeToString(b) // "eak_" + 64 hex chars
	hash = sha256Hex(raw)
	return
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func orEmptyFlags(m map[string]bool) map[string]bool {
	if m != nil {
		return m
	}
	return map[string]bool{}
}

func orEmptySettings(m map[string]string) map[string]string {
	if m != nil {
		return m
	}
	return map[string]string{}
}

