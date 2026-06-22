package tenant

import (
	"fmt"
	"sync"

	"awo.so/framework/internal/core"
)

// Entry holds all per-tenant runtime resources.
type Entry struct {
	Tenant   *Tenant
	Registry *core.EntityRegistry // custom entity definitions for this tenant
	// pgxPool and entClient are added in Phase 3 (store wiring).
}

// Registry is the in-process store of active tenants and their resources.
// It is populated lazily as tenants make requests.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]*Entry // keyed by tenant ID
}

// NewRegistry returns an empty tenant Registry.
func NewRegistry() *Registry {
	return &Registry{entries: make(map[string]*Entry)}
}

// GetOrLoad returns the Entry for the given tenant ID.
// If the tenant is not loaded yet, the provided loader function is called.
// The loader is called without the lock held; concurrent callers for the same
// tenant may both call the loader — the second result is discarded.
func (r *Registry) GetOrLoad(tenantID string, loader func() (*Entry, error)) (*Entry, error) {
	r.mu.RLock()
	if e, ok := r.entries[tenantID]; ok {
		r.mu.RUnlock()
		return e, nil
	}
	r.mu.RUnlock()

	entry, err := loader()
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	// Double-check after acquiring write lock.
	if existing, ok := r.entries[tenantID]; ok {
		r.mu.Unlock()
		return existing, nil
	}
	r.entries[tenantID] = entry
	r.mu.Unlock()
	return entry, nil
}

// Evict removes the entry for a tenant (e.g. on suspension or schema change).
func (r *Registry) Evict(tenantID string) {
	r.mu.Lock()
	delete(r.entries, tenantID)
	r.mu.Unlock()
}

// All returns a snapshot of all loaded tenant IDs.
func (r *Registry) All() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.entries))
	for id := range r.entries {
		ids = append(ids, id)
	}
	return ids
}

// ErrTenantNotFound is returned when a tenant slug cannot be resolved.
type ErrTenantNotFound struct{ Slug string }

func (e *ErrTenantNotFound) Error() string { return fmt.Sprintf("tenant %q not found", e.Slug) }
