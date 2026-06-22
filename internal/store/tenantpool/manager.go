// Package tenantpool manages per-tenant pgxpool connections.
// Each tenant gets a pool with search_path set to their schema.
package tenantpool

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"awo.so/framework/internal/store"
	"awo.so/framework/internal/store/jsonb"
	"awo.so/framework/internal/tenant"
)

// Manager creates and caches per-tenant connection pools.
type Manager struct {
	baseURL string
	log     *slog.Logger
	mu      sync.RWMutex
	pools   map[string]*pgxpool.Pool // keyed by schema name
}

// New creates a Manager with the given base PostgreSQL connection URL.
func New(baseURL string, log *slog.Logger) *Manager {
	return &Manager{
		baseURL: baseURL,
		log:     log,
		pools:   make(map[string]*pgxpool.Pool),
	}
}

// RepoFor returns a func(t *tenant.Tenant) store.EntityRepository for wiring
// into server.Deps.RepoFor. Uses a background context for pool creation.
func (m *Manager) RepoFor(ctx context.Context) func(t *tenant.Tenant) store.EntityRepository {
	return func(t *tenant.Tenant) store.EntityRepository {
		pool, err := m.poolFor(ctx, t.SchemaName)
		if err != nil {
			m.log.Error("tenantpool: failed to get pool",
				slog.String("schema", t.SchemaName), slog.Any("err", err))
			return nil
		}
		return jsonb.New(pool, t.ID.String())
	}
}

// Evict closes and removes the pool for a schema.
func (m *Manager) Evict(schemaName string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p, ok := m.pools[schemaName]; ok {
		p.Close()
		delete(m.pools, schemaName)
	}
}

// Close closes all managed pools.
func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, p := range m.pools {
		p.Close()
	}
	m.pools = make(map[string]*pgxpool.Pool)
}

func (m *Manager) poolFor(ctx context.Context, schemaName string) (*pgxpool.Pool, error) {
	m.mu.RLock()
	if p, ok := m.pools[schemaName]; ok {
		m.mu.RUnlock()
		return p, nil
	}
	m.mu.RUnlock()

	sep := "?"
	for _, c := range m.baseURL {
		if c == '?' {
			sep = "&"
			break
		}
	}
	dsn := m.baseURL + sep + "search_path=" + schemaName + ",public"

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse config for %q: %w", schemaName, err)
	}
	cfg.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect for %q: %w", schemaName, err)
	}

	m.mu.Lock()
	if existing, ok := m.pools[schemaName]; ok {
		m.mu.Unlock()
		pool.Close()
		return existing, nil
	}
	m.pools[schemaName] = pool
	m.mu.Unlock()
	return pool, nil
}
