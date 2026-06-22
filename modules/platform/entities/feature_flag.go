package entities

import (
	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/permissions"
)

// FeatureFlagDefinition is the system-level definition of a feature flag.
// Flags are defined by framework/module developers and cannot be created by tenants.
// Tenants override flag values via TenantFlagOverrideDefinition.
var FeatureFlagDefinition = entitydef.New("FeatureFlag").
	System().
	Label("Feature Flag").
	Field("name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.Unique(),
		fieldtype.Immutable(),
		fieldtype.MaxLen(100),
	).
	Field("label", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.MaxLen(255),
	).
	Field("description", fieldtype.LongText).
	Field("flag_type", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Default("boolean"),
		fieldtype.Choices("boolean", "string", "percentage"),
	).
	Field("default_value", fieldtype.Bool,
		fieldtype.Default(false),
	).
	Field("default_string_value", fieldtype.Data,
		fieldtype.MaxLen(255),
	).
	Field("default_percentage", fieldtype.Int,
		fieldtype.Default(0),
		fieldtype.MinVal(0),
		fieldtype.MaxVal(100),
	).
	Field("stage", fieldtype.Select,
		fieldtype.Default("draft"),
		fieldtype.Choices("draft", "active", "deprecated", "removed"),
	).
	Field("depends_on", fieldtype.Data,
		fieldtype.MaxLen(100), // name of a flag that must be active first
	).
	// Only superusers can create/modify flag definitions.
	Allow(&superOnlyPermission{}).
	MustBuild()

// TenantFlagOverrideDefinition allows a tenant to override a system feature flag value.
// Evaluated after the system default: system → tenant override → user override.
var TenantFlagOverrideDefinition = entitydef.New("TenantFlagOverride").
	System().
	Label("Feature Flag Override").
	Field("flag_name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.Immutable(),
		fieldtype.MaxLen(100),
	).
	Field("boolean_value", fieldtype.Bool,
		fieldtype.Default(false),
	).
	Field("string_value", fieldtype.Data,
		fieldtype.MaxLen(255),
	).
	Field("percentage_value", fieldtype.Int,
		fieldtype.Default(0),
		fieldtype.MinVal(0),
		fieldtype.MaxVal(100),
	).
	Field("reason", fieldtype.SmallText). // why this tenant has a non-default value
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
