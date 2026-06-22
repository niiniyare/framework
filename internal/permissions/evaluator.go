// Package permissions implements RBAC and privacy policy evaluation.
package permissions

import (
	"context"
	"errors"
	"fmt"

	"awo.so/framework/internal/core"
	"awo.so/framework/pkg/filter"
	"awo.so/framework/pkg/permissions"
)

// ErrForbidden is returned when a permission check fails.
var ErrForbidden = errors.New("forbidden")

// ForbiddenError wraps ErrForbidden with operation context.
type ForbiddenError struct {
	EntityName string
	Action     permissions.Action
	Reason     string
}

func (e *ForbiddenError) Error() string {
	return fmt.Sprintf("forbidden: %s/%s — %s", e.EntityName, e.Action, e.Reason)
}
func (e *ForbiddenError) Is(target error) bool { return target == ErrForbidden }

// Evaluator checks permissions and applies privacy policies for an entity operation.
type Evaluator struct{}

// New returns a new Evaluator.
func New() *Evaluator { return &Evaluator{} }

// Check verifies that the principal may perform action on the given record.
// record may be nil for collection-level checks (create, list).
// Returns ErrForbidden if any Permission denies the operation.
// Superusers bypass RBAC but NOT privacy policies.
func (e *Evaluator) Check(
	ctx context.Context,
	def *core.EntityDefinition,
	principal *permissions.Principal,
	action permissions.Action,
	record *core.EntityRecord,
) error {
	if len(def.Permissions) == 0 {
		// No permissions declared → deny by default (fail-closed).
		// Exception: superusers are always permitted.
		if principal.IsSuper {
			return nil
		}
		return &ForbiddenError{
			EntityName: def.Name,
			Action:     action,
			Reason:     "no permissions declared on entity",
		}
	}

	// Any single permission granting access is sufficient (OR logic across permissions).
	for _, p := range def.Permissions {
		if p.Allowed(ctx, principal, action, record) {
			return nil
		}
	}

	return &ForbiddenError{
		EntityName: def.Name,
		Action:     action,
		Reason:     fmt.Sprintf("principal roles %v not granted %s", principal.Roles, action),
	}
}

// ApplyPrivacy modifies f in-place so that only rows visible to the principal are returned.
// All privacy policies are applied (AND logic — every policy must pass).
func (e *Evaluator) ApplyPrivacy(
	ctx context.Context,
	def *core.EntityDefinition,
	principal *permissions.Principal,
	f *filter.Filter,
) error {
	for _, policy := range def.PrivacyPolicies {
		if err := policy.FilterQuery(ctx, principal, f); err != nil {
			return fmt.Errorf("privacy policy on %s: %w", def.Name, err)
		}
	}
	return nil
}
