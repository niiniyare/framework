package tenant

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBResolver resolves tenants by slug from the system DB.
// Implements the middleware.TenantResolver interface.
type DBResolver struct {
	db *pgxpool.Pool
}

// NewDBResolver creates a tenant resolver backed by the system database.
func NewDBResolver(db *pgxpool.Pool) *DBResolver {
	return &DBResolver{db: db}
}

// Resolve looks up a tenant by slug. Returns ErrTenantNotFound if absent.
func (r *DBResolver) Resolve(slug string) (*Tenant, error) {
	ctx := context.Background()
	row := r.db.QueryRow(ctx, `
		SELECT id, slug, name, email,
		       COALESCE(subdomain, ''),
		       status, plan,
		       COALESCE(currency_code, 'KES'),
		       COALESCE(timezone, 'Africa/Nairobi'),
		       schema_name,
		       created_at, updated_at
		FROM tenants
		WHERE slug = $1 AND deleted_at IS NULL`, slug)

	var t Tenant
	var statusStr, planStr string
	if err := row.Scan(
		&t.ID, &t.Slug, &t.Name, &t.Email, &t.Subdomain,
		&statusStr, &planStr,
		&t.CurrencyCode, &t.Timezone, &t.SchemaName,
		&t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &ErrTenantNotFound{Slug: slug}
		}
		return nil, fmt.Errorf("resolve tenant %q: %w", slug, err)
	}

	t.Status = Status(statusStr)
	t.Plan = Plan(planStr)
	return &t, nil
}
