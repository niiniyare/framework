// Package forecourt is the Awo built-in Forecourt module.
// Provides: Site, Tank, Pump, Nozzle, MeterReading, DipReading, ShiftClose, FleetCard.
package forecourt

import (
	"awo.so/framework/internal/core"
	"awo.so/framework/modules/forecourt/entities"
)

// Register adds all Forecourt module EntityDefinitions to the system registry.
func Register(reg *core.EntityRegistry) {
	reg.MustRegister(entities.SiteDefinition)
	reg.MustRegister(entities.TankDefinition)
	reg.MustRegister(entities.PumpDefinition)
	reg.MustRegister(entities.NozzleDefinition)
	reg.MustRegister(entities.MeterReadingDefinition)
	reg.MustRegister(entities.DipReadingDefinition)
	reg.MustRegister(entities.ShiftCloseDefinition)
	reg.MustRegister(entities.FleetCardDefinition)
}
