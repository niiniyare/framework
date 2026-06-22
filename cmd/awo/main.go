// Command awo is the Awo Framework CLI.
// It provides subcommands for serving, scaffolding, migrating, and tenant management.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "awo",
	Short: "Awo Framework CLI",
	Long:  "CLI for managing Awo Framework applications: serve, scaffold, migrate, and tenant operations.",
}

// Global flags
var (
	flagEnv      string
	flagConfig   string
	flagLogLevel string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&flagEnv, "env", "development", "Target environment (development|staging|production)")
	rootCmd.PersistentFlags().StringVar(&flagConfig, "config", "", "Path to config file (defaults to awo.yaml in CWD)")
	rootCmd.PersistentFlags().StringVar(&flagLogLevel, "log-level", "info", "Log level (debug|info|warn|error)")

	rootCmd.AddCommand(serveCmd())
	rootCmd.AddCommand(scaffoldCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(tenantCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
