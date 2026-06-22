// Package store defines the EntityRepository interface — the persistence contract
// that all store implementations must satisfy. Framework code above this layer
// depends only on this interface, never on ent, pgx, or any other ORM directly.
package store

import (
	"context"
	"errors"
	"fmt"

	"awo.so/framework/internal/core"
	"awo.so/framework/pkg/filter"
)

// EntityRepository is the persistence abstraction for all entity types.
// Implementations exist for:
//   - ent (system entities, SQL-backed)
//   - jsonb (custom entities, JSONB-backed)
//   - fake (in-memory, used in tests via pkg/testkit)
type EntityRepository interface {
	// FindByID returns a single record by primary key.
	// Returns ErrNotFound if no record exists with that ID.
	FindByID(ctx context.Context, entityName, id string) (*core.EntityRecord, error)

	// FindMany returns all records matching the filter along with the total count
	// (pre-pagination, for cursor pagination metadata).
	FindMany(ctx context.Context, entityName string, f filter.Filter) ([]*core.EntityRecord, int64, error)

	// Exists returns true if at least one record matches the filter.
	Exists(ctx context.Context, entityName string, f filter.Filter) (bool, error)

	// Count returns the number of records matching the filter (no data fetched).
	Count(ctx context.Context, entityName string, f filter.Filter) (int64, error)

	// Create inserts a new record and returns it with framework-assigned fields
	// (ID, TenantID, CreatedAt, UpdatedAt, DocStatus) populated.
	Create(ctx context.Context, entityName string, data map[string]any) (*core.EntityRecord, error)

	// Update applies the provided data patch to the record with the given ID.
	// Only keys present in data are updated; absent keys are left unchanged.
	// Returns ErrNotFound if the record does not exist.
	Update(ctx context.Context, entityName, id string, data map[string]any) (*core.EntityRecord, error)

	// Delete permanently removes the record with the given ID.
	// Returns ErrNotFound if the record does not exist.
	Delete(ctx context.Context, entityName, id string) error

	// BulkCreate inserts multiple records in a single round-trip.
	// Returns all created records and a slice of per-record errors (nil entries for successes).
	BulkCreate(ctx context.Context, entityName string, batch []map[string]any) ([]*core.EntityRecord, []error)

	// WithTx executes fn inside a database transaction.
	// The EntityRepository passed to fn is scoped to that transaction.
	// If fn returns an error, the transaction is rolled back.
	// Nested calls create savepoints (implementation-dependent).
	WithTx(ctx context.Context, fn func(ctx context.Context, tx EntityRepository) error) error
}

// --- Sentinel errors --------------------------------------------------------

// ErrNotFound is returned when FindByID or Delete targets a non-existent record.
var ErrNotFound = errors.New("record not found")

// ErrConflict is returned on unique constraint violations.
var ErrConflict = errors.New("record already exists")

// ErrInvalidFilter is returned when a filter cannot be translated to the store's query language.
var ErrInvalidFilter = errors.New("invalid filter")

// NotFoundError wraps ErrNotFound with entity context.
type NotFoundError struct {
	EntityName string
	ID         string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s record %q not found", e.EntityName, e.ID)
}
func (e *NotFoundError) Is(target error) bool { return target == ErrNotFound }

// ConflictError wraps ErrConflict with field context.
type ConflictError struct {
	EntityName string
	Field      string
	Value      any
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s: %q already exists with %s=%v", e.EntityName, e.Value, e.Field, e.Value)
}
func (e *ConflictError) Is(target error) bool { return target == ErrConflict }
