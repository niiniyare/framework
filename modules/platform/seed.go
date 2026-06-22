package platform

import (
	"context"
	"fmt"
	"time"

	"awo.so/framework/internal/core"
	"awo.so/framework/internal/store"
	"awo.so/framework/pkg/filter"
)

// SeedSystemRoles creates the system roles in a freshly provisioned tenant schema.
// Safe to call multiple times — skips existing roles.
func SeedSystemRoles(ctx context.Context, repo store.EntityRepository) error {
	systemRoles := []map[string]any{
		{"name": "admin", "label": "Administrator", "is_system": true,
			"description": "Full access to all entities and settings"},
		{"name": "viewer", "label": "Viewer", "is_system": true,
			"description": "Read-only access to all entities"},
		{"name": "finance_manager", "label": "Finance Manager", "is_system": true,
			"description": "Full access to finance module"},
		{"name": "accountant", "label": "Accountant", "is_system": true,
			"description": "Create and submit journal entries, view accounts"},
		{"name": "hr_manager", "label": "HR Manager", "is_system": true,
			"description": "Full access to HR module"},
		{"name": "inventory_clerk", "label": "Inventory Clerk", "is_system": true,
			"description": "Create and manage stock entries"},
		{"name": "cashier", "label": "Cashier", "is_system": true,
			"description": "Operate forecourt shifts and record transactions"},
	}

	for _, role := range systemRoles {
		f := filter.New().Where("name", role["name"]).Build()
		exists, err := repo.Exists(ctx, "Role", f)
		if err != nil {
			return fmt.Errorf("seed roles: check exists %s: %w", role["name"], err)
		}
		if exists {
			continue
		}
		if _, err := repo.Create(ctx, "Role", role); err != nil {
			return fmt.Errorf("seed roles: create %s: %w", role["name"], err)
		}
	}
	return nil
}

// SeedAdminUser creates the initial superuser for a new tenant.
// passwordHash must already be bcrypt-hashed by the caller.
// Returns nil record (no error) if the user already exists.
func SeedAdminUser(ctx context.Context, repo store.EntityRepository, email, fullName, passwordHash string) (*core.EntityRecord, error) {
	f := filter.New().Where("email", email).Build()
	exists, err := repo.Exists(ctx, "User", f)
	if err != nil {
		return nil, fmt.Errorf("seed admin: check exists: %w", err)
	}
	if exists {
		return nil, nil
	}

	user, err := repo.Create(ctx, "User", map[string]any{
		"email":         email,
		"full_name":     fullName,
		"password_hash": passwordHash,
		"is_active":     true,
		"is_super":      true,
	})
	if err != nil {
		return nil, fmt.Errorf("seed admin: create user: %w", err)
	}
	return user, nil
}

// SeedSettings creates the singleton Settings record if absent.
func SeedSettings(ctx context.Context, repo store.EntityRepository) error {
	count, err := repo.Count(ctx, "Settings", filter.New().Build())
	if err != nil {
		return fmt.Errorf("seed settings: count: %w", err)
	}
	if count > 0 {
		return nil
	}
	_, err = repo.Create(ctx, "Settings", map[string]any{
		"timezone":                  "Africa/Nairobi",
		"date_format":               "DD/MM/YYYY",
		"language":                  "en",
		"currency":                  "KES",
		"enable_finance":            true,
		"session_timeout_hours":     24,
		"max_concurrent_sessions":   5,
		"password_min_length":       8,
		"etims_environment":         "sandbox",
		"api_rate_limit_per_minute": 300,
		"primary_color":             "#1890ff",
		"created_at":                time.Now().UTC(),
	})
	return err
}
