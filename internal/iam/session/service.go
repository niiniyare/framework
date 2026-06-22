// Package session manages the creation, retrieval and invalidation of user sessions.
package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"awo.so/framework/internal/iam/domain"
	"awo.so/framework/internal/redis"
)

const (
	sessionKeyPrefix = "session:"
	mfaPendingPrefix = "mfa:login:pending:"
	userSessionsKey  = "user_sessions:"
	mfaPendingTTL    = 5 * time.Minute
)

// FlagResolver resolves all feature flags for a tenant.
type FlagResolver interface {
	ResolveForTenant(ctx context.Context, tenantID uuid.UUID) (map[string]bool, error)
}

// SettingResolver resolves all settings for a tenant.
type SettingResolver interface {
	ResolveForTenant(ctx context.Context, tenantID uuid.UUID) (map[string]string, error)
}

// CreateParams holds everything needed to build a new session.
type CreateParams struct {
	User      *domain.User
	IP        string
	UserAgent string
}

// Service manages session lifecycle.
type Service struct {
	db       *pgxpool.Pool
	redis    *redis.Client
	flags    FlagResolver
	settings SettingResolver
	log      *slog.Logger
}

// New creates a session Service.
// flags may be nil initially — call SetFlags before serving traffic.
func New(db *pgxpool.Pool, r *redis.Client, flags FlagResolver, settings SettingResolver, log *slog.Logger) *Service {
	return &Service{db: db, redis: r, flags: flags, settings: settings, log: log}
}

// SetFlags wires the flag resolver after both services are constructed.
// Must be called before the first Login call.
func (s *Service) SetFlags(flags FlagResolver) { s.flags = flags }

// Create builds and persists a new session for the given user.
// Returns the ResolvedSession and the raw (un-hashed) token.
func (s *Service) Create(ctx context.Context, p CreateParams) (*domain.ResolvedSession, string, error) {
	// Generate raw token and compute hash
	raw, hash, err := generateToken()
	if err != nil {
		return nil, "", fmt.Errorf("session: generate token: %w", err)
	}

	// Resolve tenant flags and settings for embedding in session
	var flags map[string]bool
	if s.flags != nil {
		if f, err := s.flags.ResolveForTenant(ctx, p.User.TenantID); err != nil {
			s.log.Warn("session: resolve flags failed, continuing with empty", slog.Any("err", err))
			flags = map[string]bool{}
		} else {
			flags = f
		}
	} else {
		flags = map[string]bool{}
	}
	settings, err := s.settings.ResolveForTenant(ctx, p.User.TenantID)
	if err != nil {
		s.log.Warn("session: resolve settings failed, continuing with empty", slog.Any("err", err))
		settings = map[string]string{}
	}

	// Determine TTL from settings
	ttlHours := 8
	if v, ok := settings["iam.session_ttl_hours"]; ok {
		if h := parseInt(v); h > 0 {
			ttlHours = h
		}
	}
	expiresAt := time.Now().Add(time.Duration(ttlHours) * time.Hour)

	scope := domain.EntityScope{Type: domain.EntityScopeAll}
	cfg := domain.Configuration{
		Flags:    flags,
		Settings: settings,
		Prefs:    map[string]string{},
	}

	resolved := &domain.ResolvedSession{
		UserID:        p.User.ID,
		UserType:      coalesce(p.User.UserType, domain.UserTypeInternal),
		TenantID:      p.User.TenantID,
		DisplayName:   p.User.FullName,
		EntityScope:   scope,
		Configuration: cfg,
		ExpiresAt:     expiresAt,
		SessionToken:  hash,
	}

	// Persist to DB
	scopeJSON, _ := json.Marshal(scope)
	cfgJSON, _ := json.Marshal(cfg)
	_, err = s.db.Exec(ctx, `
		INSERT INTO user_sessions
			(tenant_id, user_id, user_type, session_token,
			 entity_scope, configuration, ip_address, user_agent,
			 is_active, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE, $9)`,
		p.User.TenantID, p.User.ID, resolved.UserType, hash,
		scopeJSON, cfgJSON, p.IP, p.UserAgent, expiresAt)
	if err != nil {
		return nil, "", fmt.Errorf("session: persist: %w", err)
	}

	// Cache in Redis
	if err := s.cacheSession(ctx, hash, resolved, expiresAt); err != nil {
		s.log.Warn("session: cache write failed", slog.Any("err", err))
	}

	return resolved, raw, nil
}

// Get retrieves a session by raw token. Returns nil (no error) if not found.
func (s *Service) Get(ctx context.Context, rawToken string) (*domain.ResolvedSession, error) {
	hash := sha256Hex(rawToken)
	cacheKey := sessionKeyPrefix + hash

	// Redis cache-aside
	if data, err := s.redis.Get(ctx, cacheKey); err == nil {
		var resolved domain.ResolvedSession
		if json.Unmarshal([]byte(data), &resolved) == nil {
			return &resolved, nil
		}
	}

	// DB fallback
	row := s.db.QueryRow(ctx, `
		SELECT u.id, s.user_type, s.tenant_id, u.full_name,
		       s.entity_scope, s.configuration, s.expires_at, s.session_token
		FROM user_sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.session_token = $1
		  AND s.is_active = TRUE
		  AND s.expires_at > NOW()`, hash)

	var (
		resolved  domain.ResolvedSession
		scopeJSON []byte
		cfgJSON   []byte
	)
	if err := row.Scan(
		&resolved.UserID, &resolved.UserType, &resolved.TenantID, &resolved.DisplayName,
		&scopeJSON, &cfgJSON, &resolved.ExpiresAt, &resolved.SessionToken,
	); err != nil {
		return nil, nil // not found or expired
	}
	if err := json.Unmarshal(scopeJSON, &resolved.EntityScope); err != nil {
		return nil, fmt.Errorf("session: parse entity scope: %w", err)
	}
	if err := json.Unmarshal(cfgJSON, &resolved.Configuration); err != nil {
		return nil, fmt.Errorf("session: parse configuration: %w", err)
	}

	// Re-cache
	if err := s.cacheSession(ctx, hash, &resolved, resolved.ExpiresAt); err != nil {
		s.log.Warn("session: re-cache failed", slog.Any("err", err))
	}

	return &resolved, nil
}

// Invalidate marks a session inactive and evicts it from Redis.
func (s *Service) Invalidate(ctx context.Context, rawToken string) error {
	hash := sha256Hex(rawToken)
	_, err := s.db.Exec(ctx, `
		UPDATE user_sessions SET is_active = FALSE, updated_at = NOW()
		WHERE session_token = $1`, hash)
	if err != nil {
		return fmt.Errorf("session: invalidate db: %w", err)
	}
	_ = s.redis.Del(ctx, sessionKeyPrefix+hash)
	return nil
}

// InvalidateByUser evicts all active sessions for a user from Redis.
// DB rows are left for audit; they expire naturally.
func (s *Service) InvalidateByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE user_sessions SET is_active = FALSE
		WHERE user_id = $1 AND is_active = TRUE`, userID)
	_ = s.redis.Del(ctx, userSessionsKey+userID.String())
	return err
}

// InvalidateByTenant evicts all active sessions for all users of a tenant.
func (s *Service) InvalidateByTenant(ctx context.Context, tenantID uuid.UUID) error {
	_, err := s.db.Exec(ctx, `
		UPDATE user_sessions SET is_active = FALSE
		WHERE tenant_id = $1 AND is_active = TRUE`, tenantID)
	return err
}

// StoreMFAPending stores a pending MFA token in Redis with 5-min TTL.
func (s *Service) StoreMFAPending(ctx context.Context, pendingToken string, userID uuid.UUID) error {
	return s.redis.SetEX(ctx, mfaPendingPrefix+pendingToken, userID.String(), mfaPendingTTL)
}

// ConsumeMFAPending atomically reads and deletes a pending MFA token.
// Returns "" if not found.
func (s *Service) ConsumeMFAPending(ctx context.Context, pendingToken string) (string, error) {
	key := mfaPendingPrefix + pendingToken
	val, err := s.redis.Get(ctx, key)
	if errors.Is(err, redis.ErrKeyNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	_ = s.redis.Del(ctx, key)
	return val, nil
}

// --- helpers ---

func (s *Service) cacheSession(ctx context.Context, hash string, resolved *domain.ResolvedSession, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}
	data, err := json.Marshal(resolved)
	if err != nil {
		return err
	}
	return s.redis.SetEX(ctx, sessionKeyPrefix+hash, string(data), ttl)
}

func generateToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	raw = hex.EncodeToString(b) // 64 hex chars
	hash = sha256Hex(raw)
	return
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return n
}
