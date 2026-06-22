// Package crm is the Awo built-in CRM module.
// Provides: Customer, Contact, Address, Lead, Opportunity.
package crm

import (
	"awo.so/framework/internal/core"
	"awo.so/framework/modules/crm/entities"
)

// Register adds all CRM module EntityDefinitions to the system registry.
func Register(reg *core.EntityRegistry) {
	reg.MustRegister(entities.CustomerDefinition)
	reg.MustRegister(entities.ContactDefinition)
	reg.MustRegister(entities.AddressDefinition)
	reg.MustRegister(entities.LeadDefinition)
	reg.MustRegister(entities.OpportunityDefinition)
}
