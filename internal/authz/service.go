package authz

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	casbinmodel "github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service implements RBAC enforcement via Casbin SyncedEnforcer.
type Service struct {
	enforcer *casbin.SyncedEnforcer
	db       *pgxpool.Pool
	log      *slog.Logger
}

// New creates a Service, loads all rules from DB, and starts auto-refresh.
// refreshInterval: how often to reload policies from DB (0 = no auto-refresh).
func New(db *pgxpool.Pool, log *slog.Logger, refreshInterval time.Duration) (*Service, error) {
	m, err := casbinmodel.NewModelFromString(CasbinModel)
	if err != nil {
		return nil, fmt.Errorf("authz: parse model: %w", err)
	}

	adapter := NewPgxAdapter(db)
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("authz: new enforcer: %w", err)
	}

	if refreshInterval > 0 {
		enforcer.StartAutoLoadPolicy(refreshInterval)
	}

	return &Service{enforcer: enforcer, db: db, log: log}, nil
}

// Enforce checks if subject can perform action on object in domain.
// Never hits DB on the hot path — SyncedEnforcer holds in-memory model.
func (s *Service) Enforce(ctx context.Context, req Request) (bool, error) {
	ok, err := s.enforcer.Enforce(req.Subject, req.Domain, req.Object, req.Action)
	if err != nil {
		return false, fmt.Errorf("authz enforce: %w", err)
	}
	return ok, nil
}

// EnforceBatch checks multiple requests in one call.
func (s *Service) EnforceBatch(ctx context.Context, reqs []Request) ([]bool, error) {
	casbinReqs := make([][]any, len(reqs))
	for i, r := range reqs {
		casbinReqs[i] = []any{r.Subject, r.Domain, r.Object, r.Action}
	}
	results, err := s.enforcer.BatchEnforce(casbinReqs)
	if err != nil {
		return nil, fmt.Errorf("authz batch enforce: %w", err)
	}
	return results, nil
}

// AssignRole adds a g-rule (user → role in domain) to Casbin and inserts
// a row into role_assignments for audit/expiry tracking.
func (s *Service) AssignRole(ctx context.Context, subject, role, domain string, opts ...AssignOpt) error {
	o := applyOpts(opts)

	// Insert into Casbin (g-rule: subject inherits role in domain)
	if _, err := s.enforcer.AddRoleForUserInDomain(subject, role, domain); err != nil {
		return fmt.Errorf("authz assign role: %w", err)
	}

	// Persist assignment metadata for audit and expiry
	_, err := s.db.Exec(ctx, `
		INSERT INTO role_assignments
			(subject, role_name, domain, assigned_by, delegated_by, expires_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, TRUE)
		ON CONFLICT (subject, role_name, domain) DO UPDATE
			SET assigned_by  = EXCLUDED.assigned_by,
			    expires_at   = EXCLUDED.expires_at,
			    is_active    = TRUE`,
		subject, role, domain, o.AssignedBy, o.DelegatedBy, o.ExpiresAt)
	if err != nil {
		return fmt.Errorf("authz assign role: persist: %w", err)
	}
	return nil
}

// RevokeRole removes the g-rule from Casbin and marks assignment inactive.
func (s *Service) RevokeRole(ctx context.Context, subject, role, domain string) error {
	if _, err := s.enforcer.DeleteRoleForUserInDomain(subject, role, domain); err != nil {
		return fmt.Errorf("authz revoke role: %w", err)
	}
	_, err := s.db.Exec(ctx, `
		UPDATE role_assignments SET is_active = FALSE
		WHERE subject=$1 AND role_name=$2 AND domain=$3`,
		subject, role, domain)
	return err
}

// GetRoles returns direct roles for subject in domain.
func (s *Service) GetRoles(_ context.Context, subject, domain string) ([]string, error) {
	return s.enforcer.GetRolesForUserInDomain(subject, domain), nil
}

// GetImplicitRoles returns all inherited roles (including via role hierarchy).
func (s *Service) GetImplicitRoles(_ context.Context, subject, domain string) ([]string, error) {
	return s.enforcer.GetImplicitRolesForUser(subject, domain)
}

// HasRole checks whether subject has role in domain.
func (s *Service) HasRole(_ context.Context, subject, role, domain string) (bool, error) {
	roles := s.enforcer.GetRolesForUserInDomain(subject, domain)
	for _, r := range roles {
		if r == role {
			return true, nil
		}
	}
	return false, nil
}

// AddPolicy adds a p-rule to Casbin.
func (s *Service) AddPolicy(_ context.Context, p Policy) error {
	effect := p.Effect
	if effect == "" {
		effect = EffectAllow
	}
	_, err := s.enforcer.AddPolicy(p.Subject, p.Domain, p.Object, p.Action, effect)
	return err
}

// RemovePolicy removes a p-rule from Casbin.
func (s *Service) RemovePolicy(_ context.Context, p Policy) error {
	effect := p.Effect
	if effect == "" {
		effect = EffectAllow
	}
	_, err := s.enforcer.RemovePolicy(p.Subject, p.Domain, p.Object, p.Action, effect)
	return err
}

// GetPolicies returns all p-rules for a domain.
func (s *Service) GetPolicies(_ context.Context, domain string) ([]Policy, error) {
	rules, err := s.enforcer.GetFilteredPolicy(1, domain)
	if err != nil {
		return nil, fmt.Errorf("authz get policies: %w", err)
	}
	policies := make([]Policy, 0, len(rules))
	for _, r := range rules {
		p := Policy{
			Subject: strAt(r, 0),
			Domain:  strAt(r, 1),
			Object:  strAt(r, 2),
			Action:  strAt(r, 3),
			Effect:  strAt(r, 4),
		}
		policies = append(policies, p)
	}
	return policies, nil
}

// InvalidateCache forces a full policy reload from DB.
func (s *Service) InvalidateCache(_ context.Context) error {
	return s.enforcer.LoadPolicy()
}

// BootstrapTenantAdmin adds a wildcard allow policy and assigns tenant_admin role
// to the given user. Idempotent.
func (s *Service) BootstrapTenantAdmin(ctx context.Context, userID, tenantID uuid.UUID) error {
	subject := TenantSubject(userID.String())
	domain := TenantDomain(tenantID.String())

	if err := s.AddPolicy(ctx, Policy{
		Subject: "tenant_admin",
		Domain:  domain,
		Object:  "*",
		Action:  "*",
		Effect:  EffectAllow,
	}); err != nil {
		return fmt.Errorf("authz bootstrap: add policy: %w", err)
	}

	return s.AssignRole(ctx, subject, "tenant_admin", domain,
		WithAssignedBy(PlatformSubject("system")))
}

// SeedDefaultRoles adds the standard Casbin policies for built-in roles in a tenant domain.
func (s *Service) SeedDefaultRoles(ctx context.Context, tenantID uuid.UUID) error {
	domain := TenantDomain(tenantID.String())

	defaultPolicies := []Policy{
		{Subject: "admin", Domain: domain, Object: "*", Action: "*", Effect: EffectAllow},
		{Subject: "viewer", Domain: domain, Object: "*", Action: "read", Effect: EffectAllow},
		{Subject: "finance_manager", Domain: domain, Object: "Account", Action: "*", Effect: EffectAllow},
		{Subject: "finance_manager", Domain: domain, Object: "JournalEntry", Action: "*", Effect: EffectAllow},
		{Subject: "finance_manager", Domain: domain, Object: "FiscalYear", Action: "*", Effect: EffectAllow},
		{Subject: "finance_manager", Domain: domain, Object: "FiscalPeriod", Action: "*", Effect: EffectAllow},
		{Subject: "accountant", Domain: domain, Object: "JournalEntry", Action: "create", Effect: EffectAllow},
		{Subject: "accountant", Domain: domain, Object: "JournalEntry", Action: "read", Effect: EffectAllow},
		{Subject: "accountant", Domain: domain, Object: "JournalEntry", Action: "submit", Effect: EffectAllow},
		{Subject: "accountant", Domain: domain, Object: "Account", Action: "read", Effect: EffectAllow},
		{Subject: "inventory_clerk", Domain: domain, Object: "StockEntry", Action: "create", Effect: EffectAllow},
		{Subject: "inventory_clerk", Domain: domain, Object: "StockEntry", Action: "read", Effect: EffectAllow},
		{Subject: "inventory_clerk", Domain: domain, Object: "StockEntry", Action: "update", Effect: EffectAllow},
		{Subject: "inventory_clerk", Domain: domain, Object: "StockEntry", Action: "submit", Effect: EffectAllow},
		{Subject: "inventory_clerk", Domain: domain, Object: "Item", Action: "read", Effect: EffectAllow},
		{Subject: "hr_manager", Domain: domain, Object: "Employee", Action: "*", Effect: EffectAllow},
		{Subject: "hr_manager", Domain: domain, Object: "LeaveRequest", Action: "*", Effect: EffectAllow},
		{Subject: "hr_manager", Domain: domain, Object: "PayrollRun", Action: "*", Effect: EffectAllow},
		{Subject: "cashier", Domain: domain, Object: "ShiftClose", Action: "create", Effect: EffectAllow},
		{Subject: "cashier", Domain: domain, Object: "ShiftClose", Action: "read", Effect: EffectAllow},
		{Subject: "cashier", Domain: domain, Object: "ShiftClose", Action: "update", Effect: EffectAllow},
		{Subject: "cashier", Domain: domain, Object: "ShiftClose", Action: "submit", Effect: EffectAllow},
	}

	for _, p := range defaultPolicies {
		if err := s.AddPolicy(ctx, p); err != nil && !isDuplicate(err) {
			return fmt.Errorf("authz seed roles: add policy %s/%s/%s: %w",
				p.Subject, p.Object, p.Action, err)
		}
	}
	return nil
}

func isDuplicate(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate")
}
