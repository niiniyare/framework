// Package testkit provides testing utilities for Awo framework consumers.
package testkit

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"awo.so/framework/internal/core"
	"awo.so/framework/internal/store"
	"awo.so/framework/pkg/filter"
)

// FakeStore is an in-memory implementation of store.EntityRepository.
// It is thread-safe and suitable for unit tests that exercise hooks,
// permissions, and handler logic without a real database.
//
// Limitations:
//   - FindMany applies only equality (Eq) conditions; other operators return all records.
//   - WithTx executes the function immediately with no real transaction isolation.
type FakeStore struct {
	mu       sync.RWMutex
	records  map[string]map[string]*core.EntityRecord // [entityName][id]
	tenantID string
}

// NewFakeStore creates a FakeStore scoped to the given tenant ID.
func NewFakeStore(tenantID string) *FakeStore {
	return &FakeStore{
		tenantID: tenantID,
		records:  make(map[string]map[string]*core.EntityRecord),
	}
}

func (f *FakeStore) bucket(entityName string) map[string]*core.EntityRecord {
	if f.records[entityName] == nil {
		f.records[entityName] = make(map[string]*core.EntityRecord)
	}
	return f.records[entityName]
}

// FindByID implements store.EntityRepository.
func (f *FakeStore) FindByID(_ context.Context, entityName, id string) (*core.EntityRecord, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	rec, ok := f.bucket(entityName)[id]
	if !ok {
		return nil, &store.NotFoundError{EntityName: entityName, ID: id}
	}
	return rec.Clone(), nil
}

// FindMany implements store.EntityRepository with basic equality filtering.
func (f *FakeStore) FindMany(_ context.Context, entityName string, fil filter.Filter) ([]*core.EntityRecord, int64, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	bucket := f.bucket(entityName)
	var results []*core.EntityRecord
	for _, rec := range bucket {
		if matchesFilter(rec, fil) {
			results = append(results, rec.Clone())
		}
	}

	total := int64(len(results))

	// Apply limit
	limit := fil.Limit
	if limit <= 0 {
		limit = 20
	}
	if len(results) > limit {
		results = results[:limit]
	}

	return results, total, nil
}

// Exists implements store.EntityRepository.
func (f *FakeStore) Exists(ctx context.Context, entityName string, fil filter.Filter) (bool, error) {
	_, total, err := f.FindMany(ctx, entityName, filter.Filter{Conditions: fil.Conditions, Limit: 1})
	return total > 0, err
}

// Count implements store.EntityRepository.
func (f *FakeStore) Count(ctx context.Context, entityName string, fil filter.Filter) (int64, error) {
	_, total, err := f.FindMany(ctx, entityName, filter.Filter{Conditions: fil.Conditions})
	return total, err
}

// Create implements store.EntityRepository.
func (f *FakeStore) Create(_ context.Context, entityName string, data map[string]any) (*core.EntityRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	id := uuid.New().String()
	now := time.Now().UTC()
	rec := &core.EntityRecord{
		ID:         id,
		EntityName: entityName,
		TenantID:   f.tenantID,
		Data:       copyMap(data),
		DocStatus:  core.DocStatusDraft,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	f.bucket(entityName)[id] = rec
	return rec.Clone(), nil
}

// Update implements store.EntityRepository.
func (f *FakeStore) Update(_ context.Context, entityName, id string, data map[string]any) (*core.EntityRecord, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	rec, ok := f.bucket(entityName)[id]
	if !ok {
		return nil, &store.NotFoundError{EntityName: entityName, ID: id}
	}
	for k, v := range data {
		rec.Data[k] = v
	}
	rec.UpdatedAt = time.Now().UTC()
	return rec.Clone(), nil
}

// Delete implements store.EntityRepository.
func (f *FakeStore) Delete(_ context.Context, entityName, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	bucket := f.bucket(entityName)
	if _, ok := bucket[id]; !ok {
		return &store.NotFoundError{EntityName: entityName, ID: id}
	}
	delete(bucket, id)
	return nil
}

// BulkCreate implements store.EntityRepository.
func (f *FakeStore) BulkCreate(ctx context.Context, entityName string, batch []map[string]any) ([]*core.EntityRecord, []error) {
	results := make([]*core.EntityRecord, len(batch))
	errs := make([]error, len(batch))
	for i, data := range batch {
		rec, err := f.Create(ctx, entityName, data)
		results[i] = rec
		errs[i] = err
	}
	return results, errs
}

// WithTx implements store.EntityRepository. No actual transaction isolation.
func (f *FakeStore) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.EntityRepository) error) error {
	return fn(ctx, f)
}

// Ensure FakeStore satisfies the interface at compile time.
var _ store.EntityRepository = (*FakeStore)(nil)

// --- helpers ----------------------------------------------------------------

func copyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	c := make(map[string]any, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}

// matchesFilter applies equality checks only.
func matchesFilter(rec *core.EntityRecord, f filter.Filter) bool {
	for _, cond := range f.Conditions {
		if cond.Op != filter.Eq {
			continue // non-equality operators not supported in fake
		}
		val, ok := rec.Data[cond.Field]
		if !ok {
			return false
		}
		if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", cond.Value) {
			return false
		}
	}
	return true
}

// notFoundSentinel is used to silence the unused import error for errors pkg.
var _ = errors.New
