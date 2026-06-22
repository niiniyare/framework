package entities

import (
	"context"
	"fmt"

	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
)

// JournalEntryDefinition is the header record for a double-entry GL posting.
// It is submittable — posting only occurs on submit, not on save.
var JournalEntryDefinition = entitydef.New("JournalEntry").
	System().
	Label("Journal Entry").
	Field(
		"entry_number", fieldtype.Data,
		fieldtype.ReadOnly(),
		fieldtype.WithNamingSeries("JE-{YYYY}-{MM}-{SEQ:5}"),
	).
	Field(
		"posting_date", fieldtype.Date,
		fieldtype.Required(),
	).
	Field(
		"reference", fieldtype.Data,
		fieldtype.MaxLen(100),
	).
	Field(
		"narration", fieldtype.LongText,
		fieldtype.Required(),
	).
	Field(
		"total_debit", fieldtype.Currency,
		fieldtype.ReadOnly(),
	).
	Field(
		"total_credit", fieldtype.Currency,
		fieldtype.ReadOnly(),
	).
	Field(
		"fiscal_period", fieldtype.Link,
		fieldtype.LinkTo("FiscalPeriod"),
		fieldtype.Required(),
	).
	Submittable().
	Hook(hooks.OnSubmit, validateDoubleEntry).
	Allow(permissions.Grant("accountant", permissions.ActionCreate, permissions.ActionRead, permissions.ActionUpdate, permissions.ActionSubmit)).
	Allow(permissions.Grant("finance_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// validateDoubleEntry enforces that debits == credits before the entry is posted.
func validateDoubleEntry(_ context.Context, hctx *hooks.HookContext) error {
	// Record is *core.EntityRecord at runtime; access via Data map through the any interface.
	type dataGetter interface{ Get(string) any }
	rec, ok := hctx.Record.(dataGetter)
	if !ok {
		return nil
	}
	debit := toFloat(rec.Get("total_debit"))
	credit := toFloat(rec.Get("total_credit"))
	if debit != credit {
		return fmt.Errorf("journal entry: debits (%.4f) must equal credits (%.4f)", debit, credit)
	}
	if debit == 0 {
		return fmt.Errorf("journal entry: entry has no lines")
	}
	return nil
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	}
	return 0
}

// JournalEntryLineDefinition is a single debit or credit line within a JournalEntry.
var JournalEntryLineDefinition = entitydef.New("JournalEntryLine").
	System().
	Label("Journal Entry Line").
	Field(
		"journal_entry", fieldtype.Link,
		fieldtype.Required(),
		fieldtype.LinkTo("JournalEntry"),
		fieldtype.Immutable(),
	).
	Field(
		"account", fieldtype.Link,
		fieldtype.Required(),
		fieldtype.LinkTo("Account"),
	).
	Field(
		"debit", fieldtype.Currency,
		fieldtype.Default(float64(0)),
		fieldtype.MinVal(0),
	).
	Field(
		"credit", fieldtype.Currency,
		fieldtype.Default(float64(0)),
		fieldtype.MinVal(0),
	).
	Field(
		"cost_centre", fieldtype.Data,
		fieldtype.MaxLen(100),
	).
	Field("description", fieldtype.SmallText).
	Allow(permissions.Grant("accountant", permissions.ActionCreate, permissions.ActionRead, permissions.ActionUpdate)).
	Allow(permissions.Grant("finance_manager", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
