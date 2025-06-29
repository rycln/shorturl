// Package app is the root package that composes all application components
// into a runnable service.
package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rycln/shorturl/internal/config"
	"github.com/rycln/shorturl/internal/contextkeys"
	"github.com/rycln/shorturl/internal/handlers"
	"github.com/rycln/shorturl/internal/logger"
	"github.com/rycln/shorturl/internal/middleware"
	"github.com/rycln/shorturl/internal/services"
	"github.com/rycln/shorturl/internal/storage"
	"github.com/rycln/shorturl/internal/worker"
)

// buildInfo holds application build metadata that can be set during compilation.
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// Package-level constants defining core application parameters.
const (
	// lengthOfShortURL defines the character length of generated short URLs.
	lengthOfShortURL = 7

	// jwtExpires sets the lifetime duration for JWT authentication tokens.
	// Used in auth service when generating new tokens.
	jwtExpires = time.Duration(2) * time.Hour

	// tickerPeriod specifies the interval for batch operations processing.
	tickerPeriod = time.Duration(10) * time.Second

	// shutdownTimeout defines timeout for graceful shutdown
	shutdownTimeout = 5 * time.Second
)

// App represents the core application layer.
//
// The struct combines all main components and manages their lifecycle:
// - HTTP server
// - Background workers
// - Configuration
// - Storage
//
// Should be created once during application startup using New()
// and managed as a single unit.
type App struct {
	server  *http.Server
	storage storage.Storage
	worker  *worker.DeletionProcessor
	cfg     *config.Cfg
}

// New constructs and initializes the complete application.
//
// Steps performed:
// 1. Creates all storage layers
// 2. Initializes services
// 3. Configures HTTP routing
// 4. Prepares background workers
//
// Returns error if any component fails to initialize.
func New() (*App, error) {
	cfg, err := config.NewConfigBuilder().
		WithConfigFile().
		WithFlagParsing().
		WithEnvParsing().
		WithDefaultJWTKey().
		Build()
	if err != nil {
		return nil, fmt.Errorf("can't initialize config: %v", err)
	}

	err = logger.LogInit(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("can't initialize logger: %v", err)
	}

	scfg := storage.NewStorageConfig(
		storage.WithDatabaseDsn(cfg.DatabaseDsn),
		storage.WithFilePath(cfg.StorageFilePath),
		storage.WithStorageType(cfg.StorageType),
	)
	strg, err := storage.Factory(scfg)
	if err != nil {
		return nil, fmt.Errorf("can't initialize storage: %v", err)
	}

	hashService := services.NewHashGen(lengthOfShortURL)
	shortenerService := services.NewShortener(strg, hashService)
	batchShortenerService := services.NewBatchShortener(strg, hashService)
	pingService := services.NewPing(strg)
	authService := services.NewAuth(cfg.Key, jwtExpires)
	deleteBatchService := services.NewBatchDeleter(strg)
	statsService := services.NewStatsCollector(strg)

	worker := worker.NewDeletionProcessor(deleteBatchService)

	shortenHandler := handlers.NewShortenHandler(shortenerService, authService, cfg.ShortBaseAddr)
	apiShortenHandler := handlers.NewAPIShortenHandler(shortenerService, authService, cfg.ShortBaseAddr)
	retrieveHandler := handlers.NewRetrieveHandler(shortenerService)
	shortenBatchHandler := handlers.NewShortenBatchHandler(batchShortenerService, authService, cfg.ShortBaseAddr)
	retrieveBatchHandler := handlers.NewRetrieveBatchHandler(batchShortenerService, authService, cfg.ShortBaseAddr)
	pingHandler := handlers.NewPingHandler(pingService)
	deleteBatchHandler := handlers.NewDeleteBatchHandler(worker, authService)
	statsHandler := handlers.NewStatsHandler(statsService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Timeout(cfg.Timeout))
		r.Use(middleware.Compress)

		r.Route("/api", func(r chi.Router) {
			r.Route("/internal", func(r chi.Router) {
				r.Use(middleware.TrustedSubnet(cfg.TrustedSubnet))
				r.Get("/stats", statsHandler.ServeHTTP)
			})
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.JWT)
				r.Route("/shorten", func(r chi.Router) {
					r.Post("/batch", shortenBatchHandler.ServeHTTP)
					r.Post("/", apiShortenHandler.ServeHTTP)
				})
				r.Route("/user/urls", func(r chi.Router) {
					r.Get("/", retrieveBatchHandler.ServeHTTP)
					r.Delete("/", deleteBatchHandler.ServeHTTP)
				})
			})
		})
		r.With(authMiddleware.JWT).Post("/", shortenHandler.ServeHTTP)

		r.Get("/ping", pingHandler.ServeHTTP)
		r.Get("/{short}", func(res http.ResponseWriter, req *http.Request) {
			ctx := context.WithValue(req.Context(), contextkeys.ShortURL, chi.URLParam(req, "short"))
			retrieveHandler.ServeHTTP(res, req.WithContext(ctx))
		})
	})

	r.Group(func(r chi.Router) {
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)

		r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		r.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
		r.Handle("/debug/pprof/block", pprof.Handler("block"))
		r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		r.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	})

	s := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}

	return &App{
		server:  s,
		storage: strg,
		worker:  worker,
		cfg:     cfg,
	}, nil
}

// Run starts the application services.
//
// Launches:
// - HTTP server (blocking call)
// - Background deletion processor
func (app *App) Run() error {
	doneCh := app.worker.Run(tickerPeriod, app.cfg.Timeout)

	go func() {
		if app.cfg.EnableHTTPS {
			err := app.server.ListenAndServeTLS("cert.pem", "key.pem")
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		} else {
			err := app.server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}
	}()

	logger.Log.Info(fmt.Sprintf("Server started successfully! Address: %s Storage Type: %s", app.cfg.ServerAddr, app.cfg.StorageType))
	printBuildInfo()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-shutdown

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := app.shutdown(shutdownCtx, doneCh)
	if err != nil {
		return fmt.Errorf("shutdown error: %v", err)
	}

	err = app.cleanup()
	if err != nil {
		return fmt.Errorf("cleanup error: %v", err)
	}

	log.Println(strings.TrimPrefix(os.Args[0], "./") + " shutted down gracefully")

	return nil
}

// shutdown gracefully shuts down the application components.
// It performs the following steps in order:
//  1. Shuts down the HTTP server with the given context
//  2. Shuts down the worker component
//  3. Waits for either worker completion (doneCh) or context timeout
func (app *App) shutdown(ctx context.Context, doneCh <-chan struct{}) error {
	if err := app.server.Shutdown(ctx); err != nil {
		return err
	}

	app.worker.Shutdown()

	select {
	case <-ctx.Done():
		return fmt.Errorf("worker shutdown timeout: %w", ctx.Err())
	case <-doneCh:
	}

	return nil
}

// cleanup performs resource cleanup operations for the application.
// It handles:
//   - Closing storage connections
//   - Syncing logger buffers (ignoring EINVAL errors for non-buffered logger)
func (app *App) cleanup() error {
	if err := app.storage.Close(); err != nil {
		return fmt.Errorf("storage close failed: %w", err)
	}

	if err := logger.Log.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
		return fmt.Errorf("log sync failed: %w", err)
	}

	return nil
}

// printBuildInfo displays the build metadata in a standardized format.
func printBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
