package entities

import (
	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/permissions"
)

// FiscalYearDefinition represents a financial year (e.g. 2024-01-01 to 2024-12-31).
var FiscalYearDefinition = entitydef.New("FiscalYear").
	System().
	Label("Fiscal Year").
	Field("name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.Unique(),
		fieldtype.MaxLen(50),
	).
	Field("start_date", fieldtype.Date,
		fieldtype.Required(),
	).
	Field("end_date", fieldtype.Date,
		fieldtype.Required(),
	).
	Field("is_closed", fieldtype.Bool,
		fieldtype.Default(false),
		fieldtype.ReadOnly(),
	).
	Submittable(). // submit = close the year
	Allow(permissions.Grant("finance_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// FiscalPeriodDefinition is a sub-period within a FiscalYear (monthly by convention).
var FiscalPeriodDefinition = entitydef.New("FiscalPeriod").
	System().
	Label("Fiscal Period").
	Field("fiscal_year", fieldtype.Link,
		fieldtype.Required(),
		fieldtype.LinkTo("FiscalYear"),
		fieldtype.Immutable(),
	).
	Field("name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.MaxLen(50),
	).
	Field("start_date", fieldtype.Date,
		fieldtype.Required(),
	).
	Field("end_date", fieldtype.Date,
		fieldtype.Required(),
	).
	Field("is_closed", fieldtype.Bool,
		fieldtype.Default(false),
		fieldtype.ReadOnly(),
	).
	Allow(permissions.Grant("finance_manager", permissions.AllActions...)).
	Allow(permissions.Grant("accountant", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
