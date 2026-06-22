package core

import "time"

// DocStatus is the submission state of a document entity.
type DocStatus string

const (
	DocStatusDraft     DocStatus = "Draft"
	DocStatusSubmitted DocStatus = "Submitted"
	DocStatusCancelled DocStatus = "Cancelled"
)

// EntityRecord is the uniform value type for all entity instances in the framework.
// Both system (SQL-backed) and custom (JSONB-backed) entities are surfaced as EntityRecord
// above the store layer. Hooks, permissions, and SDUI page builders all operate on this type.
type EntityRecord struct {
	ID         string
	EntityName string
	TenantID   string

	// Data holds field values. For system entities the keys are Go snake_case field names.
	// For custom entities the keys match the CustomFieldDef.key values.
	Data map[string]any

	DocStatus DocStatus

	CreatedAt time.Time
	UpdatedAt time.Time

	// IsSystem is true when this record is backed by a typed SQL table (system entity).
	IsSystem bool
}

// Get returns the value for a field by name. Returns nil if absent.
func (r *EntityRecord) Get(field string) any {
	if r.Data == nil {
		return nil
	}
	return r.Data[field]
}

// Set sets a field value. Initialises Data if nil.
func (r *EntityRecord) Set(field string, value any) {
	if r.Data == nil {
		r.Data = make(map[string]any)
	}
	r.Data[field] = value
}

// Clone returns a shallow copy of the record.
func (r *EntityRecord) Clone() *EntityRecord {
	if r == nil {
		return nil
	}
	c := *r
	c.Data = make(map[string]any, len(r.Data))
	for k, v := range r.Data {
		c.Data[k] = v
	}
	return &c
}
