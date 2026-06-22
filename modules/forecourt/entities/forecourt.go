// Package entities contains all Forecourt module EntityDefinitions.
package entities

import (
	"context"
	"fmt"

	"awo.so/framework/pkg/entitydef"
	"awo.so/framework/pkg/fieldtype"
	"awo.so/framework/pkg/hooks"
	"awo.so/framework/pkg/permissions"
)

// SiteDefinition is a fuel retail site (petrol station).
var SiteDefinition = entitydef.New("Site").
	System().
	Label("Site").
	Field("code", fieldtype.Data,
		fieldtype.Required(), fieldtype.Unique(), fieldtype.Immutable(), fieldtype.MaxLen(20)).
	Field("name", fieldtype.Data, fieldtype.Required(), fieldtype.MaxLen(255)).
	Field("address", fieldtype.SmallText).
	Field("county", fieldtype.Data, fieldtype.MaxLen(100)).
	Field("epra_licence_no", fieldtype.Data, fieldtype.MaxLen(50)).
	Field("nema_permit_no", fieldtype.Data, fieldtype.MaxLen(50)).
	Field("manager", fieldtype.Link, fieldtype.LinkTo("Employee")).
	Field("is_active", fieldtype.Bool, fieldtype.Default(true)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("cashier", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// TankDefinition is an underground fuel storage tank at a Site.
var TankDefinition = entitydef.New("Tank").
	System().
	Label("Tank").
	Field("tank_number", fieldtype.Data,
		fieldtype.Required(), fieldtype.MaxLen(20)).
	Field("site", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Site"), fieldtype.Immutable()).
	Field("product", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Choices("PMS", "AGO", "IK", "LPG", "AdBlue")).
	Field("capacity_litres", fieldtype.Float, fieldtype.Required(), fieldtype.MinVal(0)).
	Field("dead_stock_litres", fieldtype.Float, fieldtype.Default(float64(0))).
	Field("current_volume_litres", fieldtype.Float, fieldtype.ReadOnly(), fieldtype.Default(float64(0))).
	Field("is_active", fieldtype.Bool, fieldtype.Default(true)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("cashier", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// PumpDefinition is a dispenser unit on the forecourt.
var PumpDefinition = entitydef.New("Pump").
	System().
	Label("Pump").
	Field("pump_number", fieldtype.Data,
		fieldtype.Required(), fieldtype.MaxLen(20)).
	Field("site", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Site"), fieldtype.Immutable()).
	Field("tank", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Tank")).
	Field("is_active", fieldtype.Bool, fieldtype.Default(true)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("cashier", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// NozzleDefinition is a nozzle on a Pump (one pump may have multiple nozzles).
var NozzleDefinition = entitydef.New("Nozzle").
	System().
	Label("Nozzle").
	Field("nozzle_number", fieldtype.Data,
		fieldtype.Required(), fieldtype.MaxLen(10)).
	Field("pump", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Pump"), fieldtype.Immutable()).
	Field("is_active", fieldtype.Bool, fieldtype.Default(true)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("cashier", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// MeterReadingDefinition records totaliser meter readings at shift open/close.
var MeterReadingDefinition = entitydef.New("MeterReading").
	System().
	Label("Meter Reading").
	Field("shift_close", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("ShiftClose"), fieldtype.Immutable()).
	Field("nozzle", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Nozzle"), fieldtype.Immutable()).
	Field("opening_meter", fieldtype.Float, fieldtype.Required(), fieldtype.Immutable()).
	Field("closing_meter", fieldtype.Float, fieldtype.Required()).
	Field("volume_sold", fieldtype.Float, fieldtype.ReadOnly()).
	Field("unit_price", fieldtype.Currency, fieldtype.Required()).
	Field("sales_amount", fieldtype.Currency, fieldtype.ReadOnly()).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("cashier",
		permissions.ActionCreate, permissions.ActionRead, permissions.ActionUpdate)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// DipReadingDefinition records manual tank dip measurements.
var DipReadingDefinition = entitydef.New("DipReading").
	System().
	Label("Dip Reading").
	Field("tank", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Tank"), fieldtype.Immutable()).
	Field("reading_date", fieldtype.Date, fieldtype.Required(), fieldtype.Immutable()).
	Field("reading_time", fieldtype.Time, fieldtype.Immutable()).
	Field("dip_mm", fieldtype.Float, fieldtype.Required()).
	Field("volume_litres", fieldtype.Float, fieldtype.Required()).
	Field("water_dip_mm", fieldtype.Float, fieldtype.Default(float64(0))).
	Field("shift_close", fieldtype.Link, fieldtype.LinkTo("ShiftClose")).
	Field("recorded_by", fieldtype.Link, fieldtype.LinkTo("Employee")).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("cashier",
		permissions.ActionCreate, permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

// ShiftCloseDefinition is the daily/shift reconciliation document for a Site.
var ShiftCloseDefinition = entitydef.New("ShiftClose").
	System().
	Label("Shift Close").
	Submittable().
	Field("shift_number", fieldtype.Data,
		fieldtype.ReadOnly(), fieldtype.WithNamingSeries("SC-{YYYY}-{MM}-{DD}-{SEQ:4}")).
	Field("site", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Site"), fieldtype.Immutable()).
	Field("shift_date", fieldtype.Date, fieldtype.Required(), fieldtype.Immutable()).
	Field("shift_type", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Choices("Morning", "Afternoon", "Night", "Full")).
	Field("opened_by", fieldtype.Link, fieldtype.LinkTo("Employee")).
	Field("closed_by", fieldtype.Link, fieldtype.LinkTo("Employee")).
	Field("open_time", fieldtype.Time).
	Field("close_time", fieldtype.Time).
	Field("total_fuel_sales", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("total_cash_collected", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("total_mpesa", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("total_fleet_card", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("total_credit", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("variance", fieldtype.Currency, fieldtype.ReadOnly()).
	Field("remarks", fieldtype.SmallText).
	Hook(hooks.BeforeSave, validateShiftCloseDates).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("cashier",
		permissions.ActionCreate, permissions.ActionRead,
		permissions.ActionUpdate, permissions.ActionSubmit)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()

func validateShiftCloseDates(_ context.Context, hctx *hooks.HookContext) error {
	type dataGetter interface{ Get(string) any }
	rec, ok := hctx.Record.(dataGetter)
	if !ok {
		return nil
	}
	open, _ := rec.Get("open_time").(string)
	close, _ := rec.Get("close_time").(string)
	if open != "" && close != "" && close < open {
		return fmt.Errorf("close_time must be after open_time")
	}
	return nil
}

// FleetCardDefinition represents a pre-paid/credit fuel card issued to a fleet customer.
var FleetCardDefinition = entitydef.New("FleetCard").
	System().
	Label("Fleet Card").
	Field("card_number", fieldtype.Data,
		fieldtype.Required(), fieldtype.Unique(), fieldtype.Immutable(), fieldtype.MaxLen(50)).
	Field("customer", fieldtype.Link,
		fieldtype.Required(), fieldtype.LinkTo("Customer")).
	Field("vehicle_registration", fieldtype.Data, fieldtype.MaxLen(20)).
	Field("card_type", fieldtype.Select,
		fieldtype.Required(),
		fieldtype.Default("Prepaid"),
		fieldtype.Choices("Prepaid", "Credit")).
	Field("credit_limit", fieldtype.Currency, fieldtype.Default(float64(0))).
	Field("balance", fieldtype.Currency, fieldtype.ReadOnly(), fieldtype.Default(float64(0))).
	Field("allowed_products", fieldtype.MultiSelect,
		fieldtype.Choices("PMS", "AGO", "IK", "LPG", "AdBlue")).
	Field("is_active", fieldtype.Bool, fieldtype.Default(true)).
	Allow(permissions.Grant("admin", permissions.AllActions...)).
	Allow(permissions.Grant("cashier", permissions.ActionRead)).
	Allow(permissions.Grant("viewer", permissions.ActionRead)).
	MustBuild()
