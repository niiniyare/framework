// Package tenant manages the tenant lifecycle and per-tenant resource registry.
package tenant

import "time"

// Plan is a tenant's subscription tier.
type Plan string

const (
	PlanStarter    Plan = "starter"
	PlanPro        Plan = "pro"
	PlanEnterprise Plan = "enterprise"
)

// Status is the operational state of a tenant.
type Status string

const (
	StatusActive        Status = "active"
	StatusSuspended     Status = "suspended"
	StatusTrial         Status = "trial"
	StatusDecommissioned Status = "decommissioned"
)

// Tenant holds the resolved, in-memory representation of a tenant.
// It is loaded from the database and cached in Redis.
type Tenant struct {
	ID           string
	Slug         string
	Name         string
	Plan         Plan
	Status       Status
	SchemaName   string // PostgreSQL schema: "tenant_<slug>"
	FeatureFlags map[string]bool
	CreatedAt    time.Time
	SuspendedAt  *time.Time
}

// IsActive returns true if the tenant can serve requests.
func (t *Tenant) IsActive() bool {
	return t.Status == StatusActive || t.Status == StatusTrial
}

// MaxDBConns returns the connection pool ceiling for this tenant's plan.
func (t *Tenant) MaxDBConns() int32 {
	switch t.Plan {
	case PlanEnterprise:
		return 50
	case PlanPro:
		return 20
	default:
		return 5
	}
}
