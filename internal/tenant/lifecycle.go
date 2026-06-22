package tenant

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"awo.so/framework/internal/audit"
	"awo.so/framework/internal/authz"
)

// ProvisionParams holds everything needed to create and activate a new tenant.
type ProvisionParams struct {
	Name         string
	Email        string
	Subdomain    string
	Industry     string
	CompanySize  string
	CurrencyCode string
	Timezone     string
	AdminName    string
	AdminEmail   string
	PasswordHash string // bcrypt-hashed by caller
	Plan         string // "Basic" | "Professional" | "Enterprise"
	Modules      []string
}

// Lifecycle manages the tenant provisioning, suspension and reactivation flows.
type Lifecycle struct {
	db    *pgxpool.Pool
	authz *authz.Service
	audit *audit.Service
	log   *slog.Logger
}

// NewLifecycle creates a Lifecycle manager.
func NewLifecycle(db *pgxpool.Pool, az *authz.Service, au *audit.Service, log *slog.Logger) *Lifecycle {
	return &Lifecycle{db: db, authz: az, audit: au, log: log}
}

// Provision creates a tenant, admin user, seeds roles, and activates.
// The entire flow runs in a single DB transaction — failure rolls back everything.
func (l *Lifecycle) Provision(ctx context.Context, p ProvisionParams) (*Tenant, error) {
	var (
		tenant    Tenant
		tenantID  = uuid.New()
		adminUser uuid.UUID
	)

	err := withTx(ctx, l.db, func(ctx context.Context, tx pgx.Tx) error {
		// Step 1: Create tenant record (PENDING)
		slug := generateSlug(p.Name)
		_, err := tx.Exec(ctx, `
			INSERT INTO tenants
			    (id, slug, name, email, subdomain, status,
			     industry, company_size, currency_code, timezone)
			VALUES ($1, $2, $3, $4, $5, 'PENDING', $6, $7, $8, $9)`,
			tenantID, slug, p.Name, p.Email, nullStr(p.Subdomain),
			nullStr(p.Industry), nullStr(p.CompanySize),
			coalesce(p.CurrencyCode, "KES"), coalesce(p.Timezone, "Africa/Nairobi"))
		if err != nil {
			return fmt.Errorf("provision: create tenant: %w", err)
		}

		// Step 2: Resource limit configuration
		maxUsers, maxEntities, maxTx, storageMB := planLimits(p.Plan)
		mods := p.Modules
		if len(mods) == 0 {
			mods = []string{"finance"}
		}
		modsJSON, _ := json.Marshal(mods)
		_, err = tx.Exec(ctx, `
			INSERT INTO tenant_configurations
			    (tenant_id, max_users, max_entities, max_transactions_month,
			     storage_quota_mb, allowed_modules)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (tenant_id) DO NOTHING`,
			tenantID, maxUsers, maxEntities, maxTx, storageMB, modsJSON)
		if err != nil {
			return fmt.Errorf("provision: create config: %w", err)
		}

		// Step 3: Create admin user
		adminUser = uuid.New()
		_, err = tx.Exec(ctx, `
			INSERT INTO users
			    (id, tenant_id, email, full_name, password_hash,
			     user_type, is_active, is_super)
			VALUES ($1, $2, $3, $4, $5, 'INTERNAL', TRUE, TRUE)`,
			adminUser, tenantID, p.AdminEmail, p.AdminName, p.PasswordHash)
		if err != nil {
			return fmt.Errorf("provision: create admin user: %w", err)
		}

		// Step 4: Activate
		_, err = tx.Exec(ctx, `
			UPDATE tenants SET status = 'ACTIVE', updated_at = NOW() WHERE id = $1`, tenantID)
		if err != nil {
			return fmt.Errorf("provision: activate: %w", err)
		}

		tenant = Tenant{
			ID:        tenantID,
			Slug:      slug,
			Name:      p.Name,
			Email:     p.Email,
			Status:    StatusActive,
			CreatedAt: time.Now(),
		}

		// Step 5: Emit domain event (same tx → at-least-once delivery)
		_ = l.audit.Emit(ctx, audit.DomainEvent{
			EventType: audit.EventTenantProvisioned,
			TenantID:  &tenantID,
			Payload: map[string]any{
				"tenant_id":   tenantID.String(),
				"tenant_name": p.Name,
				"admin_email": p.AdminEmail,
			},
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Step 6: Bootstrap Casbin roles (best-effort — separate store)
	if err := l.authz.BootstrapTenantAdmin(ctx, adminUser, tenantID); err != nil {
		l.log.Error("provision: bootstrap admin failed",
			slog.String("tenant_id", tenantID.String()), slog.Any("err", err))
	}
	if err := l.authz.SeedDefaultRoles(ctx, tenantID); err != nil {
		l.log.Error("provision: seed default roles failed",
			slog.String("tenant_id", tenantID.String()), slog.Any("err", err))
	}

	return &tenant, nil
}

// Suspend transitions a tenant from ACTIVE to SUSPENDED.
func (l *Lifecycle) Suspend(ctx context.Context, tenantID uuid.UUID, reason, suspendedBy string) error {
	tag, err := l.db.Exec(ctx, `
		UPDATE tenants
		SET status     = 'SUSPENDED',
		    updated_at = NOW(),
		    metadata   = COALESCE(metadata, '{}'::jsonb) || jsonb_build_object(
		        'suspension_reason', $2::text,
		        'suspended_at',      NOW()::text,
		        'suspended_by',      $3::text
		    )
		WHERE id = $1 AND status = 'ACTIVE'`, tenantID, reason, suspendedBy)
	if err != nil {
		return fmt.Errorf("suspend tenant: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("suspend tenant: not found or not ACTIVE")
	}
	l.audit.Log(ctx, audit.Event{
		TenantID:     &tenantID,
		EventType:    audit.EventTenantSuspended,
		ActorSubject: suspendedBy,
		ActorDomain:  authz.DomainPlatform,
		Metadata:     map[string]any{"reason": reason},
	})
	return nil
}

// Reactivate transitions a tenant from SUSPENDED back to ACTIVE.
func (l *Lifecycle) Reactivate(ctx context.Context, tenantID uuid.UUID, reason, reactivatedBy string) error {
	tag, err := l.db.Exec(ctx, `
		UPDATE tenants
		SET status     = 'ACTIVE',
		    updated_at = NOW(),
		    metadata   = COALESCE(metadata, '{}'::jsonb) || jsonb_build_object(
		        'reactivation_reason', $2::text,
		        'reactivated_at',      NOW()::text,
		        'reactivated_by',      $3::text
		    )
		WHERE id = $1 AND status = 'SUSPENDED'`, tenantID, reason, reactivatedBy)
	if err != nil {
		return fmt.Errorf("reactivate tenant: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reactivate tenant: not found or not SUSPENDED")
	}
	l.audit.Log(ctx, audit.Event{
		TenantID:     &tenantID,
		EventType:    audit.EventTenantReactivated,
		ActorSubject: reactivatedBy,
		ActorDomain:  authz.DomainPlatform,
		Metadata:     map[string]any{"reason": reason},
	})
	return nil
}

// --- helpers ---

func withTx(ctx context.Context, db *pgxpool.Pool, fn func(context.Context, pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	if err := fn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func generateSlug(name string) string {
	s := strings.ToLower(name)
	s = nonAlnum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return uuid.New().String()[:8]
	}
	return s
}

func planLimits(plan string) (maxUsers, maxEntities, maxTx, storageMB int) {
	switch plan {
	case "Enterprise":
		return 500, 200, 100000, 102400
	case "Professional":
		return 50, 20, 10000, 10240
	default: // Basic
		return 10, 5, 1000, 1024
	}
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
