package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"awo.so/framework/internal/core"
	internalHooks "awo.so/framework/internal/hooks"
	internalPerms "awo.so/framework/internal/permissions"
	"awo.so/framework/internal/server"
	"awo.so/framework/internal/tenant"
	"awo.so/framework/modules/crm"
	"awo.so/framework/modules/finance"
	"awo.so/framework/modules/forecourt"
	"awo.so/framework/modules/hr"
	"awo.so/framework/modules/inventory"
	"awo.so/framework/modules/platform"
	"awo.so/framework/pkg/permissions"
)

func serveCmd() *cobra.Command {
	var (
		port       int
		apiOnly    bool
		workerOnly bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Awo HTTP server (and Temporal worker by default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := newLogger()
			log.Info("awo serve starting", slog.String("env", flagEnv), slog.Int("port", port))

			// System entity registry — populated by module registrations.
			systemRegistry := core.NewEntityRegistry()

			// Register modules — order matters: platform first, then domain modules.
			platform.Register(systemRegistry)
			finance.Register(systemRegistry)
			inventory.Register(systemRegistry)
			hr.Register(systemRegistry)
			crm.Register(systemRegistry)
			forecourt.Register(systemRegistry)

			// Tenant registry — populated lazily on first request per tenant.
			tenantReg := tenant.NewRegistry()

			deps := &server.Deps{
				SystemRegistry: systemRegistry,
				TenantRegistry: tenantReg,
				HookExecutor:   internalHooks.New(),
				Evaluator:      internalPerms.New(),
				RepoFor:        nil, // TODO: wire real repo in Phase 3
				Log:            log,
			}

			cfg := server.DefaultConfig()
			cfg.Port = port

			// TODO: wire real tenant resolver and session store from Redis/DB.
			// For now, use a no-op resolver so the binary compiles and runs.
			app := server.New(cfg, deps, &noopTenantResolver{}, &noopSessionStore{}, log)

			// Graceful shutdown
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			errCh := make(chan error, 1)
			go func() {
				addr := server.Addr(cfg)
				log.Info("listening", slog.String("addr", addr))
				errCh <- app.Listen(addr)
			}()

			select {
			case <-ctx.Done():
				log.Info("shutting down")
				return app.Shutdown()
			case err := <-errCh:
				return err
			}
		},
	}

	cmd.Flags().IntVar(&port, "port", 8080, "Port to listen on")
	cmd.Flags().BoolVar(&apiOnly, "api-only", false, "Start Fiber only (no Temporal worker)")
	cmd.Flags().BoolVar(&workerOnly, "worker-only", false, "Start Temporal worker only (no Fiber)")

	return cmd
}

func newLogger() *slog.Logger {
	level := slog.LevelInfo
	switch flagLogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if flagEnv == "development" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

// --- Temporary no-op implementations until DB/Redis wiring is complete ------

type noopTenantResolver struct{}

func (n *noopTenantResolver) Resolve(slug string) (*tenant.Tenant, error) {
	return nil, fmt.Errorf("tenant resolver not configured")
}

type noopSessionStore struct{}

func (n *noopSessionStore) GetSession(_ context.Context, tenantID, sessionID string) (*permissions.Principal, error) {
	return nil, fmt.Errorf("session store not configured")
}
