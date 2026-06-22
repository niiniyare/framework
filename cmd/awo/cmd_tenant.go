package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	"awo.so/framework/internal/audit"
	"awo.so/framework/internal/authz"
	iamauth "awo.so/framework/internal/iam/auth"
	"awo.so/framework/internal/tenant"
)

func tenantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenant",
		Short: "Tenant lifecycle management",
	}
	cmd.AddCommand(tenantCreateCmd(), tenantSuspendCmd(), tenantActivateCmd(), tenantListCmd())
	return cmd
}

// openDB opens a pgxpool from --db flag or DATABASE_URL env.
func openDB(dbURL string) (*pgxpool.Pool, error) {
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		return nil, fmt.Errorf("database URL required: set --db or DATABASE_URL")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return pgxpool.New(ctx, dbURL)
}

func buildLifecycle(db *pgxpool.Pool) (*tenant.Lifecycle, error) {
	log := newLogger()
	az, err := authz.New(db, log, 5*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("authz: %w", err)
	}
	au := audit.New(db, log)
	return tenant.NewLifecycle(db, az, au, log), nil
}

// tenantCreateCmd provisions a new tenant interactively or via flags.
func tenantCreateCmd() *cobra.Command {
	var (
		dbURL        string
		name         string
		plan         string
		email        string
		adminName    string
		adminEmail   string
		adminPass    string
		currencyCode string
		timezone     string
		modules      []string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Provision a new tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("── New Tenant ──────────────────────────────")

			var err error
			if name == "" {
				name, err = askRequired("Tenant name")
				if err != nil {
					return err
				}
			}
			if email == "" {
				email, err = ask("Contact email", "")
				if err != nil {
					return err
				}
			}
			if plan == "" {
				plan, err = askSelect("Plan", []string{"Basic", "Professional", "Enterprise"}, "Basic")
				if err != nil {
					return err
				}
			}
			if currencyCode == "" {
				currencyCode, err = ask("Currency code", "KES")
				if err != nil {
					return err
				}
			}
			if timezone == "" {
				timezone, err = ask("Timezone", "Africa/Nairobi")
				if err != nil {
					return err
				}
			}

			fmt.Println("── Admin User ──────────────────────────────")

			if adminName == "" {
				adminName, err = ask("Admin full name", "Admin")
				if err != nil {
					return err
				}
			}
			if adminEmail == "" {
				adminEmail, err = askRequired("Admin email")
				if err != nil {
					return err
				}
			}
			if adminPass == "" {
				for {
					adminPass, err = askPassword("Admin password")
					if err != nil {
						return err
					}
					if len(adminPass) >= 8 {
						break
					}
					fmt.Println("  (password must be at least 8 characters)")
				}
			}

			if len(modules) == 0 {
				raw, err := ask("Modules (comma-separated)", "finance")
				if err != nil {
					return err
				}
				for _, m := range strings.Split(raw, ",") {
					m = strings.TrimSpace(m)
					if m != "" {
						modules = append(modules, m)
					}
				}
			}

			fmt.Println("────────────────────────────────────────────")

			db, err := openDB(dbURL)
			if err != nil {
				return err
			}
			defer db.Close()

			lc, err := buildLifecycle(db)
			if err != nil {
				return err
			}

			passwordHash, err := iamauth.HashPassword(adminPass)
			if err != nil {
				return fmt.Errorf("hash password: %w", err)
			}

			fmt.Print("Provisioning... ")
			t, err := lc.Provision(context.Background(), tenant.ProvisionParams{
				Name:         name,
				Email:        email,
				CurrencyCode: currencyCode,
				Timezone:     timezone,
				AdminName:    adminName,
				AdminEmail:   adminEmail,
				PasswordHash: passwordHash,
				Plan:         normalizePlan(plan),
				Modules:      modules,
			})
			if err != nil {
				fmt.Println("failed.")
				return err
			}

			fmt.Printf("done.\n\n")
			fmt.Printf("  Tenant ID : %s\n", t.ID)
			fmt.Printf("  Slug      : %s\n", t.Slug)
			fmt.Printf("  Status    : %s\n", t.Status)
			fmt.Printf("  Plan      : %s\n", plan)
			fmt.Printf("  Admin     : %s <%s>\n", adminName, adminEmail)
			return nil
		},
	}

	cmd.Flags().StringVar(&dbURL, "db", "", "PostgreSQL connection URL (overrides DATABASE_URL)")
	cmd.Flags().StringVar(&name, "name", "", "Tenant display name")
	cmd.Flags().StringVar(&plan, "plan", "", "Plan: Basic|Professional|Enterprise")
	cmd.Flags().StringVar(&email, "email", "", "Tenant contact email")
	cmd.Flags().StringVar(&currencyCode, "currency", "", "ISO 4217 currency code")
	cmd.Flags().StringVar(&timezone, "timezone", "", "IANA timezone")
	cmd.Flags().StringVar(&adminName, "admin-name", "", "Admin user full name")
	cmd.Flags().StringVar(&adminEmail, "admin-email", "", "Admin user email")
	cmd.Flags().StringVar(&adminPass, "admin-pass", "", "Admin user initial password")
	cmd.Flags().StringSliceVar(&modules, "modules", nil, "Modules to enable (e.g. finance,hr)")
	return cmd
}

func tenantSuspendCmd() *cobra.Command {
	var (
		dbURL  string
		slug   string
		reason string
		by     string
	)

	cmd := &cobra.Command{
		Use:   "suspend",
		Short: "Suspend an active tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if slug == "" {
				slug, err = askRequired("Tenant slug")
				if err != nil {
					return err
				}
			}
			if reason == "" {
				reason, err = ask("Reason", "")
				if err != nil {
					return err
				}
			}
			if by == "" {
				by, err = ask("Suspended by", "cli")
				if err != nil {
					return err
				}
			}

			db, err := openDB(dbURL)
			if err != nil {
				return err
			}
			defer db.Close()

			t, err := tenant.NewDBResolver(db).Resolve(slug)
			if err != nil {
				return err
			}

			lc, err := buildLifecycle(db)
			if err != nil {
				return err
			}

			if err := lc.Suspend(context.Background(), t.ID, reason, by); err != nil {
				return err
			}

			fmt.Printf("Suspended: %s (%s)\n", slug, t.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&dbURL, "db", "", "PostgreSQL connection URL (overrides DATABASE_URL)")
	cmd.Flags().StringVar(&slug, "slug", "", "Tenant slug")
	cmd.Flags().StringVar(&reason, "reason", "", "Suspension reason")
	cmd.Flags().StringVar(&by, "by", "", "Actor performing the suspension")
	return cmd
}

func tenantActivateCmd() *cobra.Command {
	var (
		dbURL  string
		slug   string
		reason string
		by     string
	)

	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Reactivate a suspended tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if slug == "" {
				slug, err = askRequired("Tenant slug")
				if err != nil {
					return err
				}
			}
			if reason == "" {
				reason, err = ask("Reason", "")
				if err != nil {
					return err
				}
			}
			if by == "" {
				by, err = ask("Reactivated by", "cli")
				if err != nil {
					return err
				}
			}

			db, err := openDB(dbURL)
			if err != nil {
				return err
			}
			defer db.Close()

			t, err := tenant.NewDBResolver(db).Resolve(slug)
			if err != nil {
				return err
			}

			lc, err := buildLifecycle(db)
			if err != nil {
				return err
			}

			if err := lc.Reactivate(context.Background(), t.ID, reason, by); err != nil {
				return err
			}

			fmt.Printf("Activated: %s (%s)\n", slug, t.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&dbURL, "db", "", "PostgreSQL connection URL (overrides DATABASE_URL)")
	cmd.Flags().StringVar(&slug, "slug", "", "Tenant slug")
	cmd.Flags().StringVar(&reason, "reason", "", "Reactivation reason")
	cmd.Flags().StringVar(&by, "by", "", "Actor performing the reactivation")
	return cmd
}

func tenantListCmd() *cobra.Command {
	var dbURL string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tenants",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := openDB(dbURL)
			if err != nil {
				return err
			}
			defer db.Close()

			rows, err := db.Query(context.Background(), `
				SELECT id, slug, name, status, plan, created_at
				FROM tenants
				WHERE deleted_at IS NULL
				ORDER BY created_at DESC`)
			if err != nil {
				return err
			}
			defer rows.Close()

			fmt.Printf("%-36s  %-20s  %-30s  %-11s  %-12s  %s\n",
				"ID", "SLUG", "NAME", "STATUS", "PLAN", "CREATED")
			fmt.Println(strings.Repeat("─", 122))

			count := 0
			for rows.Next() {
				var (
					id                         uuid.UUID
					slug, name, status, plan   string
					createdAt                  time.Time
				)
				if err := rows.Scan(&id, &slug, &name, &status, &plan, &createdAt); err != nil {
					return err
				}
				fmt.Printf("%-36s  %-20s  %-30s  %-11s  %-12s  %s\n",
					id, slug, name, status, plan, createdAt.Format("2006-01-02"))
				count++
			}
			if err := rows.Err(); err != nil {
				return err
			}
			fmt.Printf("\n%d tenant(s)\n", count)
			return nil
		},
	}

	cmd.Flags().StringVar(&dbURL, "db", "", "PostgreSQL connection URL (overrides DATABASE_URL)")
	return cmd
}

// normalizePlan normalises user-supplied plan string to Title case.
func normalizePlan(s string) string {
	switch strings.ToLower(s) {
	case "professional", "pro":
		return "Professional"
	case "enterprise":
		return "Enterprise"
	default:
		return "Basic"
	}
}
