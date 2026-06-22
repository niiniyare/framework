package core

import (
	"fmt"
	"sync"
)

// EntityRegistry holds all EntityDefinitions known to a single tenant.
// System entities are registered once at process start (global registry).
// Custom entities are loaded per-tenant at tenant boot.
//
// The registry is safe for concurrent reads. Writes (Register) are expected
// only during startup or tenant boot, not during request serving.
type EntityRegistry struct {
	mu          sync.RWMutex
	definitions map[string]*EntityDefinition
}

// NewEntityRegistry returns an empty EntityRegistry.
func NewEntityRegistry() *EntityRegistry {
	return &EntityRegistry{definitions: make(map[string]*EntityDefinition)}
}

// Register adds an EntityDefinition. Returns an error if the name is already
// registered. Names are case-sensitive.
func (r *EntityRegistry) Register(def *EntityDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.definitions[def.Name]; exists {
		return fmt.Errorf("entityregistry: %q already registered", def.Name)
	}
	// Freeze: pre-build the field index so concurrent reads are safe.
	def.buildIndex()
	r.definitions[def.Name] = def
	return nil
}

// MustRegister is like Register but panics on error. Intended for init-time
// registration of system entities where a conflict is a programming error.
func (r *EntityRegistry) MustRegister(def *EntityDefinition) {
	if err := r.Register(def); err != nil {
		panic(err)
	}
}

// Get returns the EntityDefinition for the given name, or an error if not found.
func (r *EntityRegistry) Get(name string) (*EntityDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.definitions[name]
	if !ok {
		return nil, fmt.Errorf("entityregistry: %q not found", name)
	}
	return def, nil
}

// Has returns true if the named entity is registered.
func (r *EntityRegistry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.definitions[name]
	return ok
}

// All returns a snapshot of all registered EntityDefinitions.
func (r *EntityRegistry) All() []*EntityDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*EntityDefinition, 0, len(r.definitions))
	for _, d := range r.definitions {
		out = append(out, d)
	}
	return out
}

// Len returns the number of registered definitions.
func (r *EntityRegistry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.definitions)
}
