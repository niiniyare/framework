// Package tenant manages the tenant lifecycle and per-tenant resource registry.
package tenant

import (
	"time"

	"github.com/google/uuid"
)

// Plan is a tenant's subscription tier.
type Plan string

const (
	PlanBasic        Plan = "Basic"
	PlanProfessional Plan = "Professional"
	PlanEnterprise   Plan = "Enterprise"
)

// Status is the operational state of a tenant.
type Status string

const (
	StatusPending    Status = "PENDING"
	StatusActive     Status = "ACTIVE"
	StatusSuspended  Status = "SUSPENDED"
	StatusArchived   Status = "ARCHIVED"
)

// Tenant holds the resolved, in-memory representation of a tenant.
type Tenant struct {
	ID           uuid.UUID
	Slug         string
	Name         string
	Email        string
	Subdomain    string
	Status       Status
	Plan         Plan
	SchemaName   string // PostgreSQL schema: "tenant_<slug>"
	CurrencyCode string
	Timezone     string
	FeatureFlags map[string]bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// IsActive returns true if the tenant can serve requests.
func (t *Tenant) IsActive() bool {
	return t.Status == StatusActive
}

// MaxDBConns returns the connection pool ceiling for this tenant's plan.
func (t *Tenant) MaxDBConns() int32 {
	switch t.Plan {
	case PlanEnterprise:
		return 50
	case PlanProfessional:
		return 20
	default:
		return 5
	}
}
