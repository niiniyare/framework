// Package hooks defines the types for EntityDefinition lifecycle hooks.
package hooks

import "context"

// HookEvent identifies which point in the lifecycle a hook fires.
type HookEvent string

const (
	BeforeValidate HookEvent = "before_validate"
	BeforeSave     HookEvent = "before_save"
	AfterSave      HookEvent = "after_save"
	BeforeDelete   HookEvent = "before_delete"
	OnSubmit       HookEvent = "on_submit"
	OnCancel       HookEvent = "on_cancel"
)

// HookContext carries all data available to a hook at execution time.
// The concrete types for Record, PrevRecord, and Principal are defined in
// internal/core to avoid circular imports; hooks receive them as any and
// type-assert when needed. Framework consumers use the helper methods instead.
type HookContext struct {
	Event      HookEvent
	EntityName string

	// Record is the current (post-validation) record being saved.
	// Nil on create before the first save.
	Record any

	// PrevRecord is the record state before this operation.
	// Nil on create.
	PrevRecord any

	// Data is the raw incoming write payload (after field coercion, before hooks).
	Data map[string]any

	// Principal is the authenticated user performing the operation.
	Principal any

	// Repo is the EntityRepository for the current tenant/transaction.
	// Hooks may call Repo.FindMany / Repo.FindByID for cross-entity validation.
	// Hooks MUST NOT call Repo.Create/Update/Delete — use AfterSave for side effects.
	Repo any
}

// HookFn is the function signature for all lifecycle hooks.
// Return a non-nil error to abort the operation and roll back the transaction.
type HookFn func(ctx context.Context, hctx *HookContext) error
