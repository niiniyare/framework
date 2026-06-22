// Package platform is the Awo built-in Platform module.
// It provides the foundational entities every tenant requires before any
// business module can function: IAM, organisation structure, settings, and feature flags.
//
// Registration order within this module matters — IAM entities must be registered
// before Organisation, Settings must be last (it references everything).
//
// Call platform.Register(systemRegistry) as the first module during startup.
package platform

import (
	"awo.so/framework/internal/core"
	"awo.so/framework/modules/platform/entities"
)

// Register adds all Platform module EntityDefinitions to the system registry.
func Register(reg *core.EntityRegistry) {
	// IAM — must come first; all other entities reference User and Role.
	reg.MustRegister(entities.UserDefinition)
	reg.MustRegister(entities.RoleDefinition)
	reg.MustRegister(entities.RoleAssignmentDefinition)
	reg.MustRegister(entities.PermissionRuleDefinition)

	// Organisation
	reg.MustRegister(entities.OrganisationDefinition)
	reg.MustRegister(entities.DepartmentDefinition)

	// Tenant management (admin-only entities)
	reg.MustRegister(entities.TenantEntityDefinition)

	// Feature flags
	reg.MustRegister(entities.FeatureFlagDefinition)
	reg.MustRegister(entities.TenantFlagOverrideDefinition)

	// Settings — singleton, registered last
	reg.MustRegister(entities.SettingsDefinition)
}
