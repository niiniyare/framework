// Package finance is the Awo built-in Finance module.
// It registers the core accounting entities: Account, JournalEntry, JournalEntryLine,
// FiscalYear, and FiscalPeriod.
//
// Register finance by calling finance.Register(systemRegistry) during process startup
// before the HTTP server begins accepting connections.
package finance

import (
	"awo.so/framework/internal/core"
	"awo.so/framework/modules/finance/entities"
)

// Register adds all Finance module EntityDefinitions to the system registry.
func Register(reg *core.EntityRegistry) {
	reg.MustRegister(entities.AccountDefinition)
	reg.MustRegister(entities.JournalEntryDefinition)
	reg.MustRegister(entities.JournalEntryLineDefinition)
	reg.MustRegister(entities.FiscalYearDefinition)
	reg.MustRegister(entities.FiscalPeriodDefinition)
}
