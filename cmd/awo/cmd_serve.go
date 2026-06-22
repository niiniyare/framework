package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	"awo.so/framework/internal/audit"
	"awo.so/framework/internal/authz"
	"awo.so/framework/internal/core"
	"awo.so/framework/internal/flagsvc"
	internalHooks "awo.so/framework/internal/hooks"
	iamauth "awo.so/framework/internal/iam/auth"
	"awo.so/framework/internal/iam/apikey"
	"awo.so/framework/internal/iam/session"
	"awo.so/framework/internal/middleware"
	internalPerms "awo.so/framework/internal/permissions"
	internalredis "awo.so/framework/internal/redis"
	"awo.so/framework/internal/server"
	"awo.so/framework/internal/settings"
	"awo.so/framework/internal/store"
	"awo.so/framework/internal/store/tenantpool"
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
		port      int
		dbURL     string
		redisAddr string
		apiOnly   bool
		workerOnly bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the Awo HTTP server (and Temporal worker by default)",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := newLogger()
			log.Info("awo serve starting", slog.String("env", flagEnv), slog.Int("port", port))

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			// ── Database ──────────────────────────────────────────────────────
			if dbURL == "" {
				dbURL = os.Getenv("DATABASE_URL")
			}
			var db *pgxpool.Pool
			if dbURL != "" {
				pool, err := pgxpool.New(ctx, dbURL)
				if err != nil {
					return err
				}
				defer pool.Close()
				db = pool
				log.Info("database connected")
			} else {
				log.Warn("DATABASE_URL not set — DB-backed services disabled")
			}

			// ── Redis ─────────────────────────────────────────────────────────
			if redisAddr == "" {
				redisAddr = os.Getenv("REDIS_URL")
			}
			if redisAddr == "" {
				redisAddr = "127.0.0.1:6379"
			}
			var redisClient *internalredis.Client
			if rc, err := internalredis.New(internalredis.Config{Addr: redisAddr}); err != nil {
				log.Warn("Redis unavailable — caching disabled", slog.Any("err", err))
			} else if err := rc.Ping(ctx); err != nil {
				log.Warn("Redis ping failed", slog.Any("err", err))
			} else {
				redisClient = rc
				log.Info("redis connected", slog.String("addr", redisAddr))
			}

			// ── IAM services ──────────────────────────────────────────────────
			var (
				tenantResolver middleware.TenantResolver = &noopTenantResolver{}
				sessionStore   middleware.SessionStore   = &noopSessionStore{}
				authDeps       *server.AuthDeps
			)

			if db != nil && redisClient != nil {
				authzSvc, err := authz.New(db, log, 5*time.Minute)
				if err != nil {
					return err
				}
				_ = authzSvc

				auditSvc := audit.New(db, log)
				_ = auditSvc

				settingsSvc := settings.New(db, redisClient, log)

				// Break circular dep: session ← flags ← session
				sessionSvc := session.New(db, redisClient, nil, settingsSvc, log)
				flagsSvc := flagsvc.New(db, redisClient, sessionSvc, log)
				sessionSvc.SetFlags(flagsSvc)

				authSvc := iamauth.New(db, sessionSvc, log)
				apiKeySvc := apikey.New(db, redisClient, flagsSvc, settingsSvc, log)

				tenantResolver = tenant.NewDBResolver(db)
				sessionStore = &resolvedSessionBridge{sessions: sessionSvc}

				authDeps = &server.AuthDeps{
					Auth:    authSvc,
					APIKeys: apiKeySvc,
					Log:     log,
				}
			} else {
				log.Warn("IAM services disabled (DB or Redis not available)")
			}

			// ── Entity registry ───────────────────────────────────────────────
			systemRegistry := core.NewEntityRegistry()
			platform.Register(systemRegistry)
			finance.Register(systemRegistry)
			inventory.Register(systemRegistry)
			hr.Register(systemRegistry)
			crm.Register(systemRegistry)
			forecourt.Register(systemRegistry)

			tenantReg := tenant.NewRegistry()

			// Per-tenant repo factory (nil when no DB URL configured)
			var repoFor func(t *tenant.Tenant) store.EntityRepository
			if dbURL != "" {
				poolMgr := tenantpool.New(dbURL, log)
				defer poolMgr.Close()
				repoFor = poolMgr.RepoFor(ctx)
			}

			deps := &server.Deps{
				SystemRegistry: systemRegistry,
				TenantRegistry: tenantReg,
				HookExecutor:   internalHooks.New(),
				Evaluator:      internalPerms.New(),
				RepoFor:        repoFor,
				Log:            log,
			}

			cfg := server.DefaultConfig()
			cfg.Port = port

			app := server.New(cfg, deps, tenantResolver, sessionStore, log)

			if authDeps != nil {
				server.MountAuthRoutes(app, authDeps)
			}

			// ── Run ───────────────────────────────────────────────────────────
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
	cmd.Flags().StringVar(&dbURL, "db", "", "PostgreSQL connection URL (overrides DATABASE_URL)")
	cmd.Flags().StringVar(&redisAddr, "redis", "", "Redis address (overrides REDIS_URL)")
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

// resolvedSessionBridge converts session.ResolvedSession → permissions.Principal
// for the existing auth middleware interface.
type resolvedSessionBridge struct {
	sessions *session.Service
}

func (b *resolvedSessionBridge) GetSession(ctx context.Context, _, rawToken string) (*permissions.Principal, error) {
	resolved, err := b.sessions.Get(ctx, rawToken)
	if err != nil || resolved == nil {
		return nil, iamauth.ErrInvalidSession
	}
	return &permissions.Principal{
		UserID:   resolved.UserID.String(),
		TenantID: resolved.TenantID.String(),
		Roles:    []string{},
		IsSuper:  resolved.IsPlatform(),
	}, nil
}

// noopTenantResolver returns a synthetic active tenant — allows server to start
// without a database for development/smoke testing.
type noopTenantResolver struct{}

func (n *noopTenantResolver) Resolve(slug string) (*tenant.Tenant, error) {
	return &tenant.Tenant{
		Slug:   slug,
		Name:   slug,
		Status: tenant.StatusActive,
		Plan:   tenant.PlanBasic,
	}, nil
}

// noopSessionStore rejects every session — forces 401 on all authenticated routes.
type noopSessionStore struct{}

func (n *noopSessionStore) GetSession(_ context.Context, _, _ string) (*permissions.Principal, error) {
	return nil, iamauth.ErrInvalidSession
}
