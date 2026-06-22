// Package permissions defines the RBAC and privacy policy interfaces.
package permissions

import (
	"context"

	"awo.so/framework/pkg/filter"
)

// Action is an operation a principal can attempt on an entity.
type Action string

const (
	ActionRead   Action = "read"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionSubmit Action = "submit"
	ActionCancel Action = "cancel"
	ActionAmend  Action = "amend"
)

// AllActions is the full set of actions, useful for granting full access.
var AllActions = []Action{
	ActionRead, ActionCreate, ActionUpdate, ActionDelete,
	ActionSubmit, ActionCancel, ActionAmend,
}

// Principal identifies an authenticated user and their roles within a tenant.
type Principal struct {
	UserID   string
	TenantID string
	Roles    []string
	IsSuper  bool // bypasses RBAC (not privacy policies)
}

// HasRole returns true if the principal holds the given role.
func (p *Principal) HasRole(role string) bool {
	for _, r := range p.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// RoleGrant declares that a role may perform a set of actions on an entity.
type RoleGrant struct {
	Role    string
	Actions []Action
}

// Permission is an RBAC rule attached to an EntityDefinition.
// Allowed is called before every operation. Record may be nil for collection checks.
type Permission interface {
	Allowed(ctx context.Context, principal *Principal, action Action, record any) bool
}

// PrivacyPolicy restricts which rows a principal can see or modify.
// FilterQuery modifies f in-place, adding predicates that enforce row-level visibility.
type PrivacyPolicy interface {
	FilterQuery(ctx context.Context, principal *Principal, f *filter.Filter) error
}

// RolePermission is the standard implementation of Permission backed by RoleGrant list.
type RolePermission struct {
	Grants []RoleGrant
}

// Allowed returns true if any of the principal's roles grants the action.
// Superusers bypass the check entirely.
func (rp *RolePermission) Allowed(_ context.Context, p *Principal, action Action, _ any) bool {
	if p.IsSuper {
		return true
	}
	for _, g := range rp.Grants {
		if p.HasRole(g.Role) {
			for _, a := range g.Actions {
				if a == action {
					return true
				}
			}
		}
	}
	return false
}

// Grant is a shorthand constructor for a RolePermission with a single RoleGrant.
func Grant(role string, actions ...Action) *RolePermission {
	return &RolePermission{Grants: []RoleGrant{{Role: role, Actions: actions}}}
}

// AllowAll is a permission that grants every action to every authenticated principal.
// Use only for system maintenance entities or during development.
type AllowAll struct{}

func (AllowAll) Allowed(_ context.Context, _ *Principal, _ Action, _ any) bool { return true }
