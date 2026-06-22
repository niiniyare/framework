// Package inventory is the Awo built-in Inventory module.
// Provides: Item, ItemGroup, Warehouse, Bin, StockEntry, StockEntryLine, StockLedger.
package inventory

import (
	"awo.so/framework/internal/core"
	"awo.so/framework/modules/inventory/entities"
)

// Register adds all Inventory module EntityDefinitions to the system registry.
func Register(reg *core.EntityRegistry) {
	reg.MustRegister(entities.ItemGroupDefinition)
	reg.MustRegister(entities.ItemDefinition)
	reg.MustRegister(entities.WarehouseDefinition)
	reg.MustRegister(entities.BinDefinition)
	reg.MustRegister(entities.StockEntryDefinition)
	reg.MustRegister(entities.StockEntryLineDefinition)
	reg.MustRegister(entities.StockLedgerEntryDefinition)
}
