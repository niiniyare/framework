package entities

import (
	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/permissions"
)

// AccountDefinition declares the Chart of Accounts entity.
// Accounts form a hierarchy via the parent_account Link field.
// Account types follow standard double-entry: Asset, Liability, Equity, Income, Expense.
var AccountDefinition = entitydef.New("Account").
	System().
	Label("Account").
	Field("code", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.Unique(),
		fieldtype.Immutable(),
		fieldtype.MaxLen(20),
	).
	Field("name", fieldtype.Data,
		fieldtype.Required(),
		fieldtype.MaxLen(255),
	).
	Field("account_type", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Choices("Asset", "Liability", "Equity", "Income", "Expense"),
	).
	Field("parent_account", fieldtype.Link,
		fieldtype.LinkTo("Account"),
	).
	Field("is_group", fieldtype.Bool,
		fieldtype.Default(false),
	).
	Field("currency", fieldtype.Select,
		fieldtype.Default("KES"),
		fieldtype.Choices("KES", "USD", "EUR", "GBP"),
	).
	Field("description", fieldtype.SmallText).
	Field("is_disabled", fieldtype.Bool,
		fieldtype.Default(false),
	).
	Allow(permissions.Grant("accountant", permissions.AllActions...)).
	Allow(permissions.Grant("finance_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
