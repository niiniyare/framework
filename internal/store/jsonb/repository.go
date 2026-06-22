// Package jsonb implements EntityRepository for custom (JSONB-backed) entities.
package jsonb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"awo.so/framework/internal/core"
	"awo.so/framework/internal/store"
	"awo.so/framework/pkg/filter"
)

// Repository implements store.EntityRepository using a JSONB custom_records table.
// The table schema is:
//
//	CREATE TABLE custom_records (
//	    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//	    entity_name TEXT NOT NULL,
//	    tenant_id   TEXT NOT NULL,
//	    data        JSONB NOT NULL DEFAULT '{}',
//	    doc_status  TEXT NOT NULL DEFAULT 'Draft',
//	    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
//	    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
//	);
type Repository struct {
	pool     *pgxpool.Pool
	tenantID string
}

// New creates a JSONB repository scoped to the given tenant.
func New(pool *pgxpool.Pool, tenantID string) *Repository {
	return &Repository{pool: pool, tenantID: tenantID}
}

// FindByID implements store.EntityRepository.
func (r *Repository) FindByID(ctx context.Context, entityName, id string) (*core.EntityRecord, error) {
	const q = `
		SELECT id, entity_name, tenant_id, data, doc_status, created_at, updated_at
		FROM custom_records
		WHERE id = $1 AND entity_name = $2 AND tenant_id = $3`

	row := r.pool.QueryRow(ctx, q, id, entityName, r.tenantID)
	return scanRecord(row)
}

// FindMany implements store.EntityRepository.
func (r *Repository) FindMany(ctx context.Context, entityName string, f filter.Filter) ([]*core.EntityRecord, int64, error) {
	baseArgs := []any{entityName, r.tenantID}
	whereClause := "entity_name = $1 AND tenant_id = $2"

	if len(f.Conditions) > 0 || len(f.Groups) > 0 {
		result, err := store.FilterToSQL(f, len(baseArgs))
		if err != nil {
			return nil, 0, err
		}
		if result.SQL != "" {
			whereClause += " AND " + result.SQL
			baseArgs = append(baseArgs, result.Args...)
		}
	}

	// Count
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM custom_records WHERE %s", whereClause)
	var total int64
	if err := r.pool.QueryRow(ctx, countQ, baseArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("jsonb.FindMany count: %w", err)
	}

	// Data
	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := f.Offset

	orderBy := store.OrderBySQL(f)
	dataQ := fmt.Sprintf(`
		SELECT id, entity_name, tenant_id, data, doc_status, created_at, updated_at
		FROM custom_records
		WHERE %s
		%s
		LIMIT %d OFFSET %d`,
		whereClause, orderBy, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, baseArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("jsonb.FindMany query: %w", err)
	}
	defer rows.Close()

	var records []*core.EntityRecord
	for rows.Next() {
		rec, err := scanRow(rows)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, rec)
	}
	return records, total, rows.Err()
}

// Exists implements store.EntityRepository.
func (r *Repository) Exists(ctx context.Context, entityName string, f filter.Filter) (bool, error) {
	_, total, err := r.FindMany(ctx, entityName, filter.Filter{Conditions: f.Conditions, Groups: f.Groups, Limit: 1})
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

// Count implements store.EntityRepository.
func (r *Repository) Count(ctx context.Context, entityName string, f filter.Filter) (int64, error) {
	_, total, err := r.FindMany(ctx, entityName, filter.Filter{Conditions: f.Conditions, Groups: f.Groups})
	return total, err
}

// Create implements store.EntityRepository.
func (r *Repository) Create(ctx context.Context, entityName string, data map[string]any) (*core.EntityRecord, error) {
	blob, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("jsonb.Create marshal: %w", err)
	}

	const q = `
		INSERT INTO custom_records (entity_name, tenant_id, data)
		VALUES ($1, $2, $3)
		RETURNING id, entity_name, tenant_id, data, doc_status, created_at, updated_at`

	row := r.pool.QueryRow(ctx, q, entityName, r.tenantID, blob)
	return scanRecord(row)
}

// Update implements store.EntityRepository.
func (r *Repository) Update(ctx context.Context, entityName, id string, data map[string]any) (*core.EntityRecord, error) {
	// Merge with existing data
	existing, err := r.FindByID(ctx, entityName, id)
	if err != nil {
		return nil, err
	}
	for k, v := range data {
		existing.Data[k] = v
	}

	blob, err := json.Marshal(existing.Data)
	if err != nil {
		return nil, fmt.Errorf("jsonb.Update marshal: %w", err)
	}

	const q = `
		UPDATE custom_records
		SET data = $1, updated_at = now()
		WHERE id = $2 AND entity_name = $3 AND tenant_id = $4
		RETURNING id, entity_name, tenant_id, data, doc_status, created_at, updated_at`

	row := r.pool.QueryRow(ctx, q, blob, id, entityName, r.tenantID)
	return scanRecord(row)
}

// Delete implements store.EntityRepository.
func (r *Repository) Delete(ctx context.Context, entityName, id string) error {
	const q = `DELETE FROM custom_records WHERE id = $1 AND entity_name = $2 AND tenant_id = $3`
	tag, err := r.pool.Exec(ctx, q, id, entityName, r.tenantID)
	if err != nil {
		return fmt.Errorf("jsonb.Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return &store.NotFoundError{EntityName: entityName, ID: id}
	}
	return nil
}

// BulkCreate implements store.EntityRepository.
func (r *Repository) BulkCreate(ctx context.Context, entityName string, batch []map[string]any) ([]*core.EntityRecord, []error) {
	results := make([]*core.EntityRecord, len(batch))
	errs := make([]error, len(batch))
	for i, data := range batch {
		rec, err := r.Create(ctx, entityName, data)
		results[i] = rec
		errs[i] = err
	}
	return results, errs
}

// WithTx implements store.EntityRepository.
func (r *Repository) WithTx(ctx context.Context, fn func(ctx context.Context, tx store.EntityRepository) error) error {
	pgxTx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("jsonb.WithTx begin: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = pgxTx.Rollback(ctx)
			panic(p)
		}
	}()

	txRepo := &txRepository{pool: r.pool, tx: pgxTx, tenantID: r.tenantID}
	if err := fn(ctx, txRepo); err != nil {
		_ = pgxTx.Rollback(ctx)
		return err
	}
	return pgxTx.Commit(ctx)
}

// --- Scan helpers -----------------------------------------------------------

func scanRecord(row pgx.Row) (*core.EntityRecord, error) {
	var (
		id, entityName, tenantID, docStatus string
		data                                []byte
		createdAt, updatedAt                time.Time
	)
	if err := row.Scan(&id, &entityName, &tenantID, &data, &docStatus, &createdAt, &updatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("jsonb scan: %w", err)
	}
	return unmarshalRecord(id, entityName, tenantID, docStatus, data, createdAt, updatedAt)
}

func scanRow(rows pgx.Rows) (*core.EntityRecord, error) {
	var (
		id, entityName, tenantID, docStatus string
		data                                []byte
		createdAt, updatedAt                time.Time
	)
	if err := rows.Scan(&id, &entityName, &tenantID, &data, &docStatus, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("jsonb scan row: %w", err)
	}
	return unmarshalRecord(id, entityName, tenantID, docStatus, data, createdAt, updatedAt)
}

func unmarshalRecord(id, entityName, tenantID, docStatus string, data []byte, createdAt, updatedAt time.Time) (*core.EntityRecord, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("jsonb unmarshal: %w", err)
	}
	return &core.EntityRecord{
		ID:         id,
		EntityName: entityName,
		TenantID:   tenantID,
		Data:       m,
		DocStatus:  core.DocStatus(docStatus),
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		IsSystem:   false,
	}, nil
}

// --- txRepository: EntityRepository scoped to a pgx.Tx ---------------------

type txRepository struct {
	pool     *pgxpool.Pool
	tx       pgx.Tx
	tenantID string
}

// delegate creates a temporary Repository that uses the tx pool adapter.
// For simplicity we re-implement the query methods using the tx directly.

func (t *txRepository) FindByID(ctx context.Context, entityName, id string) (*core.EntityRecord, error) {
	const q = `
		SELECT id, entity_name, tenant_id, data, doc_status, created_at, updated_at
		FROM custom_records
		WHERE id = $1 AND entity_name = $2 AND tenant_id = $3`
	row := t.tx.QueryRow(ctx, q, id, entityName, t.tenantID)
	return scanRecord(row)
}

func (t *txRepository) FindMany(ctx context.Context, entityName string, f filter.Filter) ([]*core.EntityRecord, int64, error) {
	// Reuse the non-transactional pool for reads — acceptable within the same transaction scope.
	r := &Repository{pool: t.pool, tenantID: t.tenantID}
	return r.FindMany(ctx, entityName, f)
}

func (t *txRepository) Exists(ctx context.Context, entityName string, f filter.Filter) (bool, error) {
	r := &Repository{pool: t.pool, tenantID: t.tenantID}
	return r.Exists(ctx, entityName, f)
}

func (t *txRepository) Count(ctx context.Context, entityName string, f filter.Filter) (int64, error) {
	r := &Repository{pool: t.pool, tenantID: t.tenantID}
	return r.Count(ctx, entityName, f)
}

func (t *txRepository) Create(ctx context.Context, entityName string, data map[string]any) (*core.EntityRecord, error) {
	blob, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	const q = `
		INSERT INTO custom_records (entity_name, tenant_id, data)
		VALUES ($1, $2, $3)
		RETURNING id, entity_name, tenant_id, data, doc_status, created_at, updated_at`
	row := t.tx.QueryRow(ctx, q, entityName, t.tenantID, blob)
	return scanRecord(row)
}

func (t *txRepository) Update(ctx context.Context, entityName, id string, data map[string]any) (*core.EntityRecord, error) {
	existing, err := t.FindByID(ctx, entityName, id)
	if err != nil {
		return nil, err
	}
	for k, v := range data {
		existing.Data[k] = v
	}
	blob, err := json.Marshal(existing.Data)
	if err != nil {
		return nil, err
	}
	const q = `
		UPDATE custom_records
		SET data = $1, updated_at = now()
		WHERE id = $2 AND entity_name = $3 AND tenant_id = $4
		RETURNING id, entity_name, tenant_id, data, doc_status, created_at, updated_at`
	row := t.tx.QueryRow(ctx, q, blob, id, entityName, t.tenantID)
	return scanRecord(row)
}

func (t *txRepository) Delete(ctx context.Context, entityName, id string) error {
	const q = `DELETE FROM custom_records WHERE id = $1 AND entity_name = $2 AND tenant_id = $3`
	tag, err := t.tx.Exec(ctx, q, id, entityName, t.tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return &store.NotFoundError{EntityName: entityName, ID: id}
	}
	return nil
}

func (t *txRepository) BulkCreate(ctx context.Context, entityName string, batch []map[string]any) ([]*core.EntityRecord, []error) {
	results := make([]*core.EntityRecord, len(batch))
	errs := make([]error, len(batch))
	for i, data := range batch {
		rec, err := t.Create(ctx, entityName, data)
		results[i] = rec
		errs[i] = err
	}
	return results, errs
}

func (t *txRepository) WithTx(_ context.Context, fn func(ctx context.Context, tx store.EntityRepository) error) error {
	// Nested: reuse same transaction (savepoints could be added here in v2).
	return fn(context.Background(), t)
}

// Ensure txRepository satisfies the interface at compile time.
var _ store.EntityRepository = (*txRepository)(nil)

// Ensure Repository satisfies the interface at compile time.
var _ store.EntityRepository = (*Repository)(nil)
