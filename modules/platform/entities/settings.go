package entities

import (
	"context"
	"fmt"

	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
)

// SettingsDefinition is a tenant-scoped singleton that stores all tenant configuration.
// One record exists per tenant. Created automatically at tenant provisioning.
// Organised into logical sections: locale, modules, integrations, limits, branding.
var SettingsDefinition = entitydef.New("Settings").
	System().
	Label("Settings").
	Singleton().

	// --- Locale & Region ---
	Field("timezone", fieldtype.Select,
		fieldtype.Default("Africa/Nairobi"),
		fieldtype.Choices(
			"Africa/Nairobi", "Africa/Lagos", "Africa/Johannesburg",
			"Africa/Dar_es_Salaam", "Africa/Kampala", "UTC",
		),
	).
	Field("date_format", fieldtype.Select,
		fieldtype.Default("DD/MM/YYYY"),
		fieldtype.Choices("DD/MM/YYYY", "MM/DD/YYYY", "YYYY-MM-DD"),
	).
	Field("time_format", fieldtype.Select,
		fieldtype.Default("HH:mm"),
		fieldtype.Choices("HH:mm", "hh:mm A"),
	).
	Field("language", fieldtype.Select,
		fieldtype.Default("en"),
		fieldtype.Choices("en", "sw"), // English, Swahili
	).
	Field("currency", fieldtype.Select,
		fieldtype.Default("KES"),
		fieldtype.Choices("KES", "USD", "EUR", "GBP", "UGX", "TZS"),
	).
	Field("number_format", fieldtype.Select,
		fieldtype.Default("1,234.56"),
		fieldtype.Choices("1,234.56", "1.234,56", "1 234.56"),
	).

	// --- Module Enablement ---
	Field("enable_finance", fieldtype.Bool, fieldtype.Default(true)).
	Field("enable_inventory", fieldtype.Bool, fieldtype.Default(false)).
	Field("enable_hr", fieldtype.Bool, fieldtype.Default(false)).
	Field("enable_crm", fieldtype.Bool, fieldtype.Default(false)).
	Field("enable_forecourt", fieldtype.Bool, fieldtype.Default(false)).

	// --- Security & Session ---
	Field("session_timeout_hours", fieldtype.Int,
		fieldtype.Default(24),
		fieldtype.MinVal(1),
		fieldtype.MaxVal(720), // max 30 days
	).
	Field("max_concurrent_sessions", fieldtype.Int,
		fieldtype.Default(5),
		fieldtype.MinVal(1),
		fieldtype.MaxVal(50),
	).
	Field("require_mfa", fieldtype.Bool, fieldtype.Default(false)).
	Field("password_min_length", fieldtype.Int,
		fieldtype.Default(8),
		fieldtype.MinVal(6),
		fieldtype.MaxVal(128),
	).
	Field("password_requires_uppercase", fieldtype.Bool, fieldtype.Default(false)).
	Field("password_requires_number", fieldtype.Bool, fieldtype.Default(true)).

	// --- KRA eTIMS Integration (Kenya-specific) ---
	Field("etims_enabled", fieldtype.Bool, fieldtype.Default(false)).
	Field("etims_taxpayer_pin", fieldtype.Data,
		fieldtype.MaxLen(20),
		fieldtype.Sensitive(),
	).
	Field("etims_device_serial", fieldtype.Data,
		fieldtype.MaxLen(50),
	).
	Field("etims_environment", fieldtype.Select,
		fieldtype.Default("sandbox"),
		fieldtype.Choices("sandbox", "production"),
	).

	// --- NEMA Compliance (forecourt) ---
	Field("nema_tank_variance_threshold_pct", fieldtype.Float,
		fieldtype.Default(float64(0.5)), // 0.5% daily variance limit
	).

	// --- Rate Limits ---
	Field("api_rate_limit_per_minute", fieldtype.Int,
		fieldtype.Default(300),
		fieldtype.MinVal(10),
	).

	// --- Branding ---
	Field("primary_color", fieldtype.Data,
		fieldtype.Default("#1890ff"),
		fieldtype.MaxLen(7), // #RRGGBB
	).
	Field("logo_url", fieldtype.AttachImage).

	// --- Hooks ---
	Hook(hooks.BeforeDelete, preventSettingsDeletion).

	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("finance_manager",
		permissions.ActionRead,
		permissions.ActionUpdate,
	)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// preventSettingsDeletion ensures the Settings singleton is never deleted.
func preventSettingsDeletion(_ context.Context, _ *hooks.HookContext) error {
	return fmt.Errorf("settings record cannot be deleted")
}
