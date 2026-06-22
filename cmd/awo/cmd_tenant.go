package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func tenantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tenant",
		Short: "Tenant lifecycle management",
	}

	var (
		tenantName string
		tenantSlug string
		tenantPlan string
	)

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Provision a new tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tenantSlug == "" || tenantName == "" {
				return fmt.Errorf("--slug and --name are required")
			}
			// TODO: connect to DB and call TenantLifecycle.Provision
			fmt.Printf("tenant.create: slug=%s name=%s plan=%s (DB wiring pending)\n", tenantSlug, tenantName, tenantPlan)
			return nil
		},
	}
	createCmd.Flags().StringVar(&tenantName, "name", "", "Tenant display name")
	createCmd.Flags().StringVar(&tenantSlug, "slug", "", "Tenant slug (URL-safe identifier)")
	createCmd.Flags().StringVar(&tenantPlan, "plan", "starter", "Subscription plan (starter|pro|enterprise)")

	suspendCmd := &cobra.Command{
		Use:   "suspend",
		Short: "Suspend a tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tenantSlug == "" {
				return fmt.Errorf("--slug is required")
			}
			fmt.Printf("tenant.suspend: slug=%s (DB wiring pending)\n", tenantSlug)
			return nil
		},
	}
	suspendCmd.Flags().StringVar(&tenantSlug, "slug", "", "Tenant slug")

	activateCmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate or reactivate a tenant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tenantSlug == "" {
				return fmt.Errorf("--slug is required")
			}
			fmt.Printf("tenant.activate: slug=%s (DB wiring pending)\n", tenantSlug)
			return nil
		},
	}
	activateCmd.Flags().StringVar(&tenantSlug, "slug", "", "Tenant slug")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tenants",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("tenant.list: (DB wiring pending)")
			return nil
		},
	}

	cmd.AddCommand(createCmd, suspendCmd, activateCmd, listCmd)
	return cmd
}
