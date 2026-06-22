package entities

import (
	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/permissions"
)

// OrganisationDefinition is the top-level legal entity of a tenant.
// Each tenant has exactly one Organisation record (Singleton).
// It stores legal name, registration details, and branding.
var OrganisationDefinition = entitydef.New("Organisation").
	System().
	Label("Organisation").
	Singleton(). // one record per tenant
	Field("legal_name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.MaxLen(255),
	).
	Field("trading_name", fieldtype.Data,
		fieldtype.MaxLen(255),
	).
	Field("registration_number", fieldtype.Data,
		fieldtype.MaxLen(50),
	).
	Field("tax_pin", fieldtype.Data,
		fieldtype.MaxLen(50),
		fieldtype.Sensitive(), // KRA PIN — sensitive
	).
	Field("vat_number", fieldtype.Data,
		fieldtype.MaxLen(50),
	).
	Field("address_line1", fieldtype.Data,
		fieldtype.MaxLen(255),
	).
	Field("address_line2", fieldtype.Data,
		fieldtype.MaxLen(255),
	).
	Field("city", fieldtype.Data,
		fieldtype.MaxLen(100),
	).
	Field("country", fieldtype.Select,
		fieldtype.Default("KE"),
		fieldtype.Choices("KE", "UG", "TZ", "RW", "ET", "NG", "GH", "ZA"),
	).
	Field("phone", fieldtype.Data,
		fieldtype.MaxLen(20),
	).
	Field("email", fieldtype.Data,
		fieldtype.MaxLen(255),
	).
	Field("website", fieldtype.Data,
		fieldtype.MaxLen(255),
	).
	Field("logo", fieldtype.AttachImage).
	Field("currency", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Default("KES"),
		fieldtype.Choices("KES", "USD", "EUR", "GBP", "UGX", "TZS", "RWF"),
	).
	Field("fiscal_year_start", fieldtype.Select,
		fieldtype.Default("January"),
		fieldtype.Choices(
			"January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December",
		),
	).
	Field("timezone", fieldtype.Select,
		fieldtype.Default("Africa/Nairobi"),
		fieldtype.Choices("Africa/Nairobi", "Africa/Lagos", "Africa/Johannesburg", "UTC"),
	).
	Field("date_format", fieldtype.Select,
		fieldtype.Default("DD/MM/YYYY"),
		fieldtype.Choices("DD/MM/YYYY", "MM/DD/YYYY", "YYYY-MM-DD"),
	).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("finance_manager", permissions.ActionRead, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// DepartmentDefinition models the organisational unit hierarchy.
// Departments form a tree via parent_department (materialised path or adjacency list).
var DepartmentDefinition = entitydef.New("Department").
	System().
	Label("Department").
	Field("name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.MaxLen(255),
	).
	Field("code", fieldtype.Data,
		fieldtype.Unique(),
		fieldtype.MaxLen(20),
	).
	Field("parent_department", fieldtype.Link,
		fieldtype.LinkTo("Department"),
	).
	Field("head_of_department", fieldtype.Link,
		fieldtype.LinkTo("User"),
	).
	Field("cost_centre_code", fieldtype.Data,
		fieldtype.MaxLen(20),
	).
	Field("is_active", fieldtype.Bool,
		fieldtype.Default(true),
	).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("hr_manager", permissions.ActionCreate, permissions.ActionRead, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
