package core

import "fmt"

// EntityResolver resolves entity names to definitions.
// It checks a tenant-scoped registry first, then falls back to the system
// (global) registry. This implements the Custom → System fallback transparently.
type EntityResolver struct {
	system *EntityRegistry // global; shared across all tenants; system entities only
	tenant *EntityRegistry // per-tenant; custom entities + tenant overrides
}

// NewEntityResolver creates a resolver that checks tenant before system.
// tenant may be nil — in that case only system entities are resolvable.
func NewEntityResolver(system, tenant *EntityRegistry) *EntityResolver {
	return &EntityResolver{system: system, tenant: tenant}
}

// Resolve returns the EntityDefinition for the given name.
// Tenant registry takes precedence over system registry.
// Returns ErrEntityNotFound if not found in either.
func (r *EntityResolver) Resolve(name string) (*EntityDefinition, error) {
	if r.tenant != nil {
		if def, err := r.tenant.Get(name); err == nil {
			return def, nil
		}
	}
	if r.system != nil {
		if def, err := r.system.Get(name); err == nil {
			return def, nil
		}
	}
	return nil, &ErrEntityNotFound{Name: name}
}

// MustResolve panics if the entity is not found. Use only in init paths.
func (r *EntityResolver) MustResolve(name string) *EntityDefinition {
	def, err := r.Resolve(name)
	if err != nil {
		panic(err)
	}
	return def
}

// ErrEntityNotFound is returned when an entity name cannot be resolved.
type ErrEntityNotFound struct {
	Name string
}

func (e *ErrEntityNotFound) Error() string {
	return fmt.Sprintf("entity %q not found", e.Name)
}

// IsEntityNotFound returns true if the error is an ErrEntityNotFound.
func IsEntityNotFound(err error) bool {
	_, ok := err.(*ErrEntityNotFound)
	return ok
}
