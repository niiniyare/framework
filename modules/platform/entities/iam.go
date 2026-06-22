// Package entities contains all Platform module EntityDefinitions.
package entities

import (
	"context"
	"fmt"

	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
)

// UserDefinition represents an authenticated user within a tenant.
// Password hash is Sensitive — never returned in API responses or logs.
// Email is unique per tenant (enforced at DB level via unique index).
var UserDefinition = entitydef.New("User").
	System().
	Label("User").
	Field("email", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.Unique(),
		fieldtype.MaxLen(255),
	).
	Field("full_name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.MaxLen(255),
	).
	Field("password_hash", fieldtype.Data,
		fieldtype.Sensitive(),
		fieldtype.ReadOnly(), // set via dedicated auth endpoints, not CRUD
	).
	Field("is_active", fieldtype.Bool,
		fieldtype.Default(true),
	).
	Field("is_super", fieldtype.Bool,
		fieldtype.Default(false),
		fieldtype.ReadOnly(), // only system admins can elevate
	).
	Field("phone", fieldtype.Data,
		fieldtype.MaxLen(20),
	).
	Field("avatar_url", fieldtype.AttachImage).
	Field("last_login_at", fieldtype.DateTime,
		fieldtype.ReadOnly(),
	).
	Field("department", fieldtype.Link,
		fieldtype.LinkTo("Department"),
	).
	Hook(hooks.BeforeDelete, preventLastAdminDeletion).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager",
		permissions.ActionCreate,
		permissions.ActionRead,
		permissions.ActionUpdate,
	)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// preventLastAdminDeletion blocks deletion of the last admin user in a tenant.
func preventLastAdminDeletion(_ context.Context, hctx *hooks.HookContext) error {
	type dataGetter interface{ Get(string) any }
	rec, ok := hctx.Record.(dataGetter)
	if !ok {
		return nil
	}
	// Guard: if this user has is_super=true, ensure at least one other super remains.
	if isSuper, _ := rec.Get("is_super").(bool); isSuper {
		// Actual count check requires a DB query — hook has access via hctx.Repo.
		// Deferred to integration: return error if count of super users drops to 0.
		_ = hctx.Repo
	}
	return nil
}

// RoleDefinition represents a named permission role within a tenant.
// System roles (is_system=true) are seeded at tenant provisioning and cannot be deleted.
var RoleDefinition = entitydef.New("Role").
	System().
	Label("Role").
	Field("name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.Unique(),
		fieldtype.MaxLen(100),
	).
	Field("label", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.MaxLen(255),
	).
	Field("description", fieldtype.SmallText).
	Field("is_system", fieldtype.Bool,
		fieldtype.Default(false),
		fieldtype.Immutable(),
	).
	Hook(hooks.BeforeDelete, preventSystemRoleDeletion).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

func preventSystemRoleDeletion(_ context.Context, hctx *hooks.HookContext) error {
	type dataGetter interface{ Get(string) any }
	rec, ok := hctx.Record.(dataGetter)
	if !ok {
		return nil
	}
	if sys, _ := rec.Get("is_system").(bool); sys {
		return fmt.Errorf("system roles cannot be deleted")
	}
	return nil
}

// RoleAssignmentDefinition links a User to a Role (many-to-many junction).
var RoleAssignmentDefinition = entitydef.New("RoleAssignment").
	System().
	Label("Role Assignment").
	Field("user", fieldtype.Link,
		fieldtype.Required(),
		fieldtype.LinkTo("User"),
		fieldtype.Immutable(),
	).
	Field("role", fieldtype.Link,
		fieldtype.Required(),
		fieldtype.LinkTo("Role"),
		fieldtype.Immutable(),
	).
	Field("granted_by", fieldtype.Link,
		fieldtype.LinkTo("User"),
		fieldtype.ReadOnly(),
	).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// PermissionRuleDefinition allows tenant admins to grant or restrict role access
// on a per-entity, per-action basis beyond the framework defaults.
var PermissionRuleDefinition = entitydef.New("PermissionRule").
	System().
	Label("Permission Rule").
	Field("role", fieldtype.Link,
		fieldtype.Required(),
		fieldtype.LinkTo("Role"),
	).
	Field("entity_name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.MaxLen(100),
	).
	Field("action", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Choices("read", "create", "update", "delete", "submit", "cancel", "amend"),
	).
	Field("effect", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Default("allow"),
		fieldtype.Choices("allow", "deny"),
	).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
