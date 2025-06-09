// Package app is the root package that composes all application components
// into a runnable service.
package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
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
)

// App represents the core application layer.
//
// The struct combines all main components and manages their lifecycle:
// - HTTP router (chi)
// - Background workers
// - Configuration
// - Storage
//
// Should be created once during application startup using New()
// and managed as a single unit.
type App struct {
	router  *chi.Mux
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

	worker := worker.NewDeletionProcessor(deleteBatchService)

	shortenHandler := handlers.NewShortenHandler(shortenerService, authService, cfg.ShortBaseAddr)
	apiShortenHandler := handlers.NewAPIShortenHandler(shortenerService, authService, cfg.ShortBaseAddr)
	retrieveHandler := handlers.NewRetrieveHandler(shortenerService)
	shortenBatchHandler := handlers.NewShortenBatchHandler(batchShortenerService, authService, cfg.ShortBaseAddr)
	retrieveBatchHandler := handlers.NewRetrieveBatchHandler(batchShortenerService, authService, cfg.ShortBaseAddr)
	pingHandler := handlers.NewPingHandler(pingService)
	deleteBatchHandler := handlers.NewDeleteBatchHandler(worker, authService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Timeout(cfg.Timeout))
		r.Use(middleware.Compress)

		r.Route("/api", func(r chi.Router) {
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

	return &App{
		router:  r,
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
func (app *App) Run() (err error) {
	defer func() {
		if errStrgClose := app.storage.Close(); errStrgClose != nil {
			err = fmt.Errorf("%v; storage close failed: %w", err, errStrgClose)
		}
	}()
	defer func() {
		if errLogSync := logger.Log.Sync(); errLogSync != nil {
			err = fmt.Errorf("%v; log sync failed: %w", err, errLogSync)
		}
	}()

	app.worker.Run(tickerPeriod, app.cfg.Timeout)
	defer app.worker.Shutdown()

	logger.Log.Info(fmt.Sprintf("Server started successfully! Address: %s Storage Type: %s", app.cfg.ServerAddr, app.cfg.StorageType))
	printBuildInfo()

	if app.cfg.EnableHTTPS {
		err = http.ListenAndServeTLS(app.cfg.ServerAddr, "cert.pem", "key.pem", app.router)
		if err != nil {
			return err
		}
	} else {
		err = http.ListenAndServe(app.cfg.ServerAddr, app.router)
		if err != nil {
			return err
		}
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
