package core

import (
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
	"awo.so/framework/pkg/sdui"
	"awo.so/framework/pkg/workflow"
)

// EntityType distinguishes system entities from custom (tenant-defined) entities.
type EntityType string

const (
	EntityTypeSystem EntityType = "system" // Go-typed, SQL-backed via ent
	EntityTypeCustom EntityType = "custom" // tenant-defined, JSONB-backed
)

// EntityDefinition is the central primitive of the Awo framework.
// It is immutable after registration — the builder constructs it and
// EntityRegistry.Register freezes it. All framework subsystems (store,
// hooks, permissions, SDUI, workflow) read from this struct.
type EntityDefinition struct {
	// Identity
	Name  string
	Label string // human-readable, used in UI

	Type EntityType

	// Fields declared on this entity.
	Fields []Field

	// fieldIndex is a pre-built name→Field map for O(1) lookup.
	fieldIndex map[string]*Field

	// Lifecycle hooks — multiple hooks per event are chained in declaration order.
	Hooks map[hooks.HookEvent][]hooks.HookFn

	// RBAC permissions attached to this entity.
	Permissions []permissions.Permission

	// Row-level privacy policies.
	PrivacyPolicies []permissions.PrivacyPolicy

	// Temporal workflow bindings.
	Workflows []workflow.WorkflowBinding

	// PageBuilder overrides the default SDUI page generation for this entity.
	// If nil, the framework generates a standard CRUD page.
	PageBuilder sdui.PageBuilderFn

	// Behaviour flags
	IsSubmittable bool // enables on_submit / on_cancel hooks and status transitions
	IsSingleton   bool // only one record per tenant (e.g. Settings)
}

// FieldByName returns the named field, or nil if not found.
// O(1) after the first call (index is built lazily on first lookup).
func (d *EntityDefinition) FieldByName(name string) *Field {
	if d.fieldIndex == nil {
		d.buildIndex()
	}
	return d.fieldIndex[name]
}

// buildIndex constructs the fieldIndex from the Fields slice.
func (d *EntityDefinition) buildIndex() {
	d.fieldIndex = make(map[string]*Field, len(d.Fields))
	for i := range d.Fields {
		d.fieldIndex[d.Fields[i].Name] = &d.Fields[i]
	}
}

// HooksFor returns the hook chain for an event, or nil if none registered.
func (d *EntityDefinition) HooksFor(event hooks.HookEvent) []hooks.HookFn {
	if d.Hooks == nil {
		return nil
	}
	return d.Hooks[event]
}
