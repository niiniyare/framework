// Package entitydef provides the public builder API for declaring EntityDefinitions.
// This is the primary API surface for module developers and application developers.
//
// Example:
//
//	def, err := entitydef.New("SalesOrder").
//	    Label("Sales Order").
//	    Field("customer", fieldtype.Link, fieldtype.Required(), fieldtype.LinkTo("Customer")).
//	    Field("total", fieldtype.Currency, fieldtype.ReadOnly()).
//	    Field("status", fieldtype.Select, fieldtype.Choices("Draft","Submitted","Cancelled")).
//	    Submittable().
//	    Hook(hooks.OnSubmit, validateStockHook).
//	    Hook(hooks.AfterSave, triggerFulfillmentWorkflow).
//	    Allow(permissions.Grant("sales_manager", permissions.AllActions...)).
//	    Build()
package entitydef

import (
	"errors"
	"fmt"

	"awo.so/framework/internal/core"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
	"awo.so/framework/pkg/sdui"
	"awo.so/framework/pkg/workflow"
)

// Builder accumulates options for an EntityDefinition.
// Obtain one via New().
type Builder struct {
	name   string
	label  string
	etype  core.EntityType
	fields []core.Field
	hks    map[hooks.HookEvent][]hooks.HookFn
	perms  []permissions.Permission
	priv   []permissions.PrivacyPolicy
	wfs    []workflow.WorkflowBinding
	page   sdui.PageBuilderFn
	submit bool
	single bool
	errs   []error
}

// New starts building an EntityDefinition with the given name.
// The name must be PascalCase and unique within its registry.
// By default the entity type is Custom (JSONB-backed).
// Call .System() to declare a system (SQL-backed) entity.
func New(name string) *Builder {
	if name == "" {
		b := &Builder{}
		b.errs = append(b.errs, errors.New("entitydef: name must not be empty"))
		return b
	}
	return &Builder{
		name:  name,
		label: name,
		etype: core.EntityTypeCustom,
		hks:   make(map[hooks.HookEvent][]hooks.HookFn),
	}
}

// Label sets the human-readable label shown in the UI.
func (b *Builder) Label(label string) *Builder {
	b.label = label
	return b
}

// System declares this as a system entity (Go-typed, SQL-backed via ent).
// Most framework consumers will not call this; it is used by the framework's
// built-in module definitions.
func (b *Builder) System() *Builder {
	b.etype = core.EntityTypeSystem
	return b
}

// Field adds a field to the entity.
func (b *Builder) Field(name string, ft fieldtype.FieldType, opts ...fieldtype.FieldOption) *Builder {
	if name == "" {
		b.errs = append(b.errs, fmt.Errorf("entitydef %q: field name must not be empty", b.name))
		return b
	}
	// Check for duplicate
	for _, f := range b.fields {
		if f.Name == name {
			b.errs = append(b.errs, fmt.Errorf("entitydef %q: duplicate field %q", b.name, name))
			return b
		}
	}

	o := fieldtype.Apply(opts)
	b.fields = append(b.fields, core.Field{
		Name:         name,
		Label:        name, // overridable via UI metadata
		Type:         ft,
		Required:     o.Required,
		Unique:       o.Unique,
		Immutable:    o.Immutable,
		Sensitive:    o.Sensitive,
		Translatable: o.Translatable,
		ReadOnly:     o.ReadOnly,
		Default:      o.Default,
		MaxLen:       o.MaxLen,
		Min:          o.Min,
		Max:          o.Max,
		Choices:      o.Choices,
		LinkTarget:   o.LinkTarget,
		NamingSeries: o.NamingSeries,
	})
	return b
}

// Hook registers a lifecycle hook for the given event.
// Multiple hooks on the same event are chained in registration order.
func (b *Builder) Hook(event hooks.HookEvent, fn hooks.HookFn) *Builder {
	if fn == nil {
		b.errs = append(b.errs, fmt.Errorf("entitydef %q: nil hook for event %q", b.name, event))
		return b
	}
	b.hks[event] = append(b.hks[event], fn)
	return b
}

// Allow attaches an RBAC permission rule.
func (b *Builder) Allow(p permissions.Permission) *Builder {
	b.perms = append(b.perms, p)
	return b
}

// Privacy attaches a row-level privacy policy.
func (b *Builder) Privacy(p permissions.PrivacyPolicy) *Builder {
	b.priv = append(b.priv, p)
	return b
}

// Workflow attaches a Temporal workflow binding.
func (b *Builder) Workflow(w workflow.WorkflowBinding) *Builder {
	if w.WorkflowName == "" {
		b.errs = append(b.errs, fmt.Errorf("entitydef %q: workflow binding has empty WorkflowName", b.name))
		return b
	}
	b.wfs = append(b.wfs, w)
	return b
}

// PageBuilder sets a custom SDUI page builder, overriding the framework default.
func (b *Builder) PageBuilder(fn sdui.PageBuilderFn) *Builder {
	b.page = fn
	return b
}

// Submittable enables the on_submit / on_cancel hooks and DocStatus transitions.
func (b *Builder) Submittable() *Builder {
	b.submit = true
	return b
}

// Singleton marks the entity as having at most one record per tenant (e.g. Settings).
func (b *Builder) Singleton() *Builder {
	b.single = true
	return b
}

// Build validates the accumulated options and returns an immutable EntityDefinition.
// Returns an error if any invalid options were supplied.
func (b *Builder) Build() (*core.EntityDefinition, error) {
	if len(b.errs) > 0 {
		return nil, errors.Join(b.errs...)
	}
	if b.name == "" {
		return nil, errors.New("entitydef: name must not be empty")
	}

	def := &core.EntityDefinition{
		Name:            b.name,
		Label:           b.label,
		Type:            b.etype,
		Fields:          b.fields,
		Hooks:           b.hks,
		Permissions:     b.perms,
		PrivacyPolicies: b.priv,
		Workflows:       b.wfs,
		PageBuilder:     b.page,
		IsSubmittable:   b.submit,
		IsSingleton:     b.single,
	}
	// Pre-build the field index so it's ready before the definition is frozen.
	// EntityRegistry.Register will call buildIndex again, which is idempotent.
	return def, nil
}

// MustBuild is like Build but panics on error. Use only in init/var blocks.
func (b *Builder) MustBuild() *core.EntityDefinition {
	def, err := b.Build()
	if err != nil {
		panic(err)
	}
	return def
}
