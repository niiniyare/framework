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
	"awo.so/framework/internal/sdui"
	"awo.so/framework/internal/server"
	"awo.so/framework/internal/settings"
	"awo.so/framework/internal/store"
	"awo.so/framework/internal/store/tenantpool"
	"awo.so/framework/internal/tenant"
	"awo.so/framework/internal/workflow"
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
		port         int
		dbURL        string
		redisAddr    string
		temporalAddr string
		webDir       string
		apiOnly      bool
		workerOnly   bool
		noUI         bool
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
				authzSvc       *authz.Service
			)

			if db != nil && redisClient != nil {
				var err error
				authzSvc, err = authz.New(db, log, 5*time.Minute)
				if err != nil {
					return err
				}

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

			// ── Temporal dispatcher (best-effort) ────────────────────────────
			var dispatcher *workflow.Dispatcher
			if !apiOnly {
				if temporalAddr == "" {
					temporalAddr = os.Getenv("TEMPORAL_HOST_PORT")
				}
				if temporalAddr == "" {
					temporalAddr = workflow.DefaultConfig().HostPort
				}
				if wfClient, err := workflow.New(workflow.Config{
					HostPort:  temporalAddr,
					Namespace: "default",
				}); err != nil {
					log.Warn("Temporal unavailable — workflow dispatch disabled", slog.Any("err", err))
				} else {
					defer wfClient.Close()
					dispatcher = workflow.NewDispatcher(wfClient, log)
					log.Info("temporal connected", slog.String("addr", temporalAddr))
				}
			}

			// ── SDUI cached page builder ──────────────────────────────────────
			var pageBuilder *sdui.CachedBuilder
			if redisClient != nil {
				pageBuilder = sdui.NewCachedBuilder(redisClient)
			}

			deps := &server.Deps{
				SystemRegistry: systemRegistry,
				TenantRegistry: tenantReg,
				HookExecutor:   internalHooks.New(),
				Evaluator:      internalPerms.New(),
				Authz:          authzSvc,
				Dispatcher:     dispatcher,
				PageBuilder:    pageBuilder,
				RepoFor:        repoFor,
				Log:            log,
			}

			cfg := server.DefaultConfig()
			cfg.Port = port
			if noUI {
				cfg.WebDir = ""
			} else {
				if webDir != "" {
					cfg.WebDir = webDir
				}
				// Dev fallback: if dist/index.html doesn't exist but web/ source does,
				// serve index.html from web/ and public assets from web/public/.
				if _, err := os.Stat(cfg.WebDir + "/index.html"); os.IsNotExist(err) {
					if _, srcErr := os.Stat("./web/index.html"); srcErr == nil {
						cfg.WebDir = "./web"
						cfg.WebPublicDir = "./web/public"
						log.Info("web/dist not built — serving from web/ source (dev mode)")
					}
				}
			}

			app := server.New(cfg, deps, tenantResolver, sessionStore, log)

			if authDeps != nil {
				server.MountAuthRoutes(app, authDeps)
			}

			if authzSvc != nil && db != nil {
				server.MountIAMRoutes(app, &server.IAMDeps{
					Authz: authzSvc,
					DB:    db,
					Log:   log,
				})
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
	cmd.Flags().StringVar(&temporalAddr, "temporal", "", "Temporal host:port (overrides TEMPORAL_HOST_PORT)")
	cmd.Flags().StringVar(&webDir, "web-dir", "", "Path to built web assets (default: ./web/dist)")
	cmd.Flags().BoolVar(&noUI, "no-ui", false, "Disable static web UI serving (API only)")
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
