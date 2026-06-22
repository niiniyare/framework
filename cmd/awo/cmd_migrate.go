package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	var (
		targetTenant string
		dryRun       bool
		verify       bool
		rollback     bool
		steps        int
		diffName     string
	)

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration operations",
	}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply pending migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "migrations/system"
			if targetTenant != "" {
				dir = fmt.Sprintf("migrations/tenant_template") // TODO: resolve per-tenant schema
			}
			return runAtlas("migrate", "apply", "--dir", "file://"+dir)
		},
	}
	applyCmd.Flags().StringVar(&targetTenant, "tenant", "", "Apply to a specific tenant schema only")

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runAtlas("migrate", "status", "--dir", "file://migrations/system")
		},
	}

	diffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Generate a new versioned migration file",
		RunE: func(_ *cobra.Command, _ []string) error {
			if diffName == "" {
				diffName = "change"
			}
			return runAtlas("migrate", "diff", diffName, "--dir", "file://migrations/system")
		},
	}
	diffCmd.Flags().StringVar(&diffName, "name", "", "Migration description slug")

	_ = dryRun
	_ = verify
	_ = rollback
	_ = steps

	cmd.AddCommand(applyCmd, statusCmd, diffCmd)
	return cmd
}

func runAtlas(args ...string) error {
	path, err := exec.LookPath("atlas")
	if err != nil {
		return fmt.Errorf("atlas binary not found in PATH — install from https://atlasgo.io/getting-started")
	}
	c := exec.Command(path, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
