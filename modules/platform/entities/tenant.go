package entities

import (
	"context"

	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/permissions"
)

// TenantEntityDefinition exposes tenant metadata as a queryable entity.
// Only system superusers can read or modify tenant records via the API.
// Tenant provisioning/suspension goes through tenant.Lifecycle.
var TenantEntityDefinition = entitydef.New("Tenant").
	System().
	Label("Tenant").
	Field("slug", fieldtype.Data,
		fieldtype.Required(), fieldtype.Unique(), fieldtype.Immutable(), fieldtype.MaxLen(63)).
	Field("name", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("email", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("subdomain", fieldtype.Data, fieldtype.Unique(), fieldtype.MaxLen(63)).
	Field("plan", fieldtype.Select,
		fieldtype.Required(), fieldtype.Default("Basic"),
		fieldtype.Choices("Basic", "Professional", "Enterprise")).
	Field("status", fieldtype.Select,
		fieldtype.Required(), fieldtype.Default("PENDING"), fieldtype.ReadOnly(),
		fieldtype.Choices("PENDING", "ACTIVE", "SUSPENDED", "ARCHIVED")).
	Field("industry", fieldtype.Data, fieldtype.MaxLen(100)).
	Field("company_size", fieldtype.Data, fieldtype.MaxLen(50)).
	Field("currency_code", fieldtype.Data, fieldtype.Default("KES"), fieldtype.MaxLen(3)).
	Field("timezone", fieldtype.Data, fieldtype.Default("Africa/Nairobi"), fieldtype.MaxLen(50)).
	Field("schema_name", fieldtype.Data, fieldtype.ReadOnly(), fieldtype.MaxLen(70)).
	Allow(&superOnlyPermission{}).
	MustBuild()

// superOnlyPermission grants access only to system superusers.
type superOnlyPermission struct{}

func (s *superOnlyPermission) Allowed(_ context.Context, p *permissions.Principal, _ permissions.Action, _ any) bool {
	return p.IsSuper
}
