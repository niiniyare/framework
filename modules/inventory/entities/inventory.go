// Package entities contains all Inventory module EntityDefinitions.
package entities

import (
	"context"
	"fmt"

	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
)

// ItemGroupDefinition is a hierarchical category for items.
var ItemGroupDefinition = entitydef.New("ItemGroup").
	System().
	Label("Item Group").
	Field("name", fieldtype.Data, fieldtype.Required(), fieldtype.Unique(), fieldtype.MaxLen(255)).
	Field("parent_group", fieldtype.Link, fieldtype.LinkTo("ItemGroup")).
	Field("description", fieldtype.SmallText).
	Field("is_group", fieldtype.Bool, fieldtype.Default(false)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("inventory_clerk",
		permissions.ActionRead, permissions.ActionCreate, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// ItemDefinition represents a stockable product or service.
var ItemDefinition = entitydef.New("Item").
	System().
	Label("Item").
	Field("code", fieldtype.Data,
		fieldtype.Required(), fieldtype.Unique(), fieldtype.Immutable(), fieldtype.MaxLen(50)).
	Field("name", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("description", fieldtype.LongText).
	Field("item_group", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("ItemGroup")).
	Field("unit_of_measure", fieldtype.Select,
		fieldtype.Required(), fieldtype.Default("Unit"),
		fieldtype.Choices("Unit", "Kg", "Litre", "Metre", "Box", "Carton", "Piece", "Set")).
	Field("valuation_method", fieldtype.Select,
		fieldtype.Default("FIFO"),
		fieldtype.Choices("FIFO", "WeightedAverage")).
	Field("standard_rate", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("is_stock_item", fieldtype.Bool, fieldtype.Default(true)).
	Field("is_service_item", fieldtype.Bool, fieldtype.Default(false)).
	Field("reorder_level", fieldtype.Float, fieldtype.Default(float64(0))).
	Field("reorder_qty", fieldtype.Float, fieldtype.Default(float64(0))).
	Field("is_disabled", fieldtype.Bool, fieldtype.Default(false)).
	Field("barcode", fieldtype.Data, fieldtype.MaxLen(50)).
	Field("image", fieldtype.AttachImage).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("inventory_clerk",
		permissions.ActionRead, permissions.ActionCreate, permissions.ActionUpdate)).
	Allow(permissions.Grant("finance_manager", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// WarehouseDefinition is a physical storage location.
var WarehouseDefinition = entitydef.New("Warehouse").
	System().
	Label("Warehouse").
	Field("name", fieldtype.Data, fieldtype.Required(), fieldtype.Unique(), fieldtype.MaxLen(255)).
	Field("code", fieldtype.Data, fieldtype.Unique(), fieldtype.MaxLen(20)).
	Field("warehouse_type", fieldtype.Select,
		fieldtype.Default("Stores"),
		fieldtype.Choices("Stores", "Transit", "WIP", "Scrap")).
	Field("address", fieldtype.SmallText).
	Field("is_group", fieldtype.Bool, fieldtype.Default(false)).
	Field("parent_warehouse", fieldtype.Link, fieldtype.LinkTo("Warehouse")).
	Field("is_disabled", fieldtype.Bool, fieldtype.Default(false)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("inventory_clerk", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// BinDefinition tracks current stock quantity of an Item in a Warehouse.
// Updated automatically by StockLedgerEntry — not directly editable.
var BinDefinition = entitydef.New("Bin").
	System().
	Label("Bin").
	Field("item", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Item"), fieldtype.Immutable()).
	Field("warehouse", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Warehouse"), fieldtype.Immutable()).
	Field("actual_qty", fieldtype.Float, fieldtype.ReadOnly(), fieldtype.Default(float64(0))).
	Field("reserved_qty", fieldtype.Float, fieldtype.ReadOnly(), fieldtype.Default(float64(0))).
	Field("ordered_qty", fieldtype.Float, fieldtype.ReadOnly(), fieldtype.Default(float64(0))).
	Field("valuation_rate", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("stock_value", fieldtype.Currency, fieldtype.ReadOnly()).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("inventory_clerk", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// StockEntryDefinition is the header for a stock movement document.
// Entry types: Receipt, Issue, Transfer, Adjustment.
var StockEntryDefinition = entitydef.New("StockEntry").
	System().
	Label("Stock Entry").
	Field("entry_number", fieldtype.Data,
		fieldtype.ReadOnly(), fieldtype.WithNamingSeries("SE-{YYYY}-{MM}-{SEQ:5}")).
	Field("entry_type", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Choices("Receipt", "Issue", "Transfer", "Adjustment")).
	Field("posting_date", fieldtype.Date, fieldtype.Required()).
	Field("source_warehouse", fieldtype.Link, fieldtype.LinkTo("Warehouse")).
	Field("target_warehouse", fieldtype.Link, fieldtype.LinkTo("Warehouse")).
	Field("remarks", fieldtype.SmallText).
	Field("total_value", fieldtype.Currency, fieldtype.ReadOnly()).
	Submittable().
	Hook(hooks.BeforeSave, validateWarehouseForEntryType).
	Hook(hooks.OnSubmit, validateStockLinesNotEmpty).
	Allow(permissions.Grant("inventory_clerk",
		permissions.ActionCreate, permissions.ActionRead,
		permissions.ActionUpdate, permissions.ActionSubmit)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

func validateWarehouseForEntryType(_ context.Context, hctx *hooks.HookContext) error {
	type dataGetter interface{ Get(string) any }
	rec, ok := hctx.Record.(dataGetter)
	if !ok {
		data, ok2 := hctx.Data["entry_type"]
		if !ok2 {
			return nil
		}
		_ = data
		return nil
	}
	entryType, _ := rec.Get("entry_type").(string)
	switch entryType {
	case "Transfer":
		if rec.Get("source_warehouse") == nil || rec.Get("target_warehouse") == nil {
			return fmt.Errorf("transfer entries require both source and target warehouse")
		}
	case "Receipt":
		if rec.Get("target_warehouse") == nil {
			return fmt.Errorf("receipt entries require a target warehouse")
		}
	case "Issue":
		if rec.Get("source_warehouse") == nil {
			return fmt.Errorf("issue entries require a source warehouse")
		}
	}
	return nil
}

func validateStockLinesNotEmpty(_ context.Context, _ *hooks.HookContext) error {
	// Full validation (line count check) requires a DB query via hctx.Repo.
	// Deferred to integration — hook stub ensures the binding is registered.
	return nil
}

// StockEntryLineDefinition is a single item line within a StockEntry.
var StockEntryLineDefinition = entitydef.New("StockEntryLine").
	System().
	Label("Stock Entry Line").
	Field("stock_entry", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("StockEntry"), fieldtype.Immutable()).
	Field("item", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Item")).
	Field("qty", fieldtype.Float,
		fieldtype.Required(), fieldtype.MinVal(0)).
	Field("uom", fieldtype.Select,
		fieldtype.Choices("Unit", "Kg", "Litre", "Metre", "Box", "Carton", "Piece", "Set")).
	Field("basic_rate", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("amount", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("serial_no", fieldtype.SmallText).
	Field("batch_no", fieldtype.Data, fieldtype.MaxLen(50)).
	Allow(permissions.Grant("inventory_clerk",
		permissions.ActionCreate, permissions.ActionRead, permissions.ActionUpdate)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// StockLedgerEntryDefinition is the immutable audit trail of every stock movement.
// Created automatically by the StockEntry submit workflow — never created directly via API.
var StockLedgerEntryDefinition = entitydef.New("StockLedgerEntry").
	System().
	Label("Stock Ledger Entry").
	Field("posting_date", fieldtype.Date, fieldtype.Required(), fieldtype.Immutable()).
	Field("item", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Item"), fieldtype.Immutable()).
	Field("warehouse", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Warehouse"), fieldtype.Immutable()).
	Field("actual_qty", fieldtype.Float, fieldtype.Required(), fieldtype.Immutable()).
	Field("qty_after_transaction", fieldtype.Float, fieldtype.Required(), fieldtype.Immutable()).
	Field("valuation_rate", fieldtype.Currency, fieldtype.Immutable()).
	Field("stock_value", fieldtype.Currency, fieldtype.Immutable()).
	Field("stock_value_difference", fieldtype.Currency, fieldtype.Immutable()).
	Field("voucher_type", fieldtype.Data, fieldtype.Immutable(), fieldtype.MaxLen(50)).
	Field("voucher_no", fieldtype.Data, fieldtype.Immutable(), fieldtype.MaxLen(50)).
	Allow(permissions.Grant("admin", permissions.ActionRead)).
	Allow(permissions.Grant("inventory_clerk", permissions.ActionRead)).
	Allow(permissions.Grant("finance_manager", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
