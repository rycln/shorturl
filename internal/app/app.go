package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

const (
	lengthOfShortURL = 7
	jwtExpires       = time.Duration(2) * time.Hour
	tickerPeriod     = time.Duration(10) * time.Second
)

type App struct {
	router  *chi.Mux
	storage storage.Storage
	worker  *worker.DeletionProcessor
	cfg     *config.Cfg
}

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

	shortenHandler := handlers.NewShortenHandler(shortenerService, cfg.ShortBaseAddr)
	apiShortenHandler := handlers.NewAPIShortenHandler(shortenerService, cfg.ShortBaseAddr)
	retrieveHandler := handlers.NewRetrieveHandler(shortenerService)
	shortenBatchHandler := handlers.NewShortenBatchHandler(batchShortenerService, cfg.ShortBaseAddr)
	retrieveBatchHandler := handlers.NewRetrieveBatchHandler(batchShortenerService, cfg.ShortBaseAddr)
	pingHandler := handlers.NewPingHandler(pingService)
	deleteBatchHandler := handlers.NewDeleteBatchHandler(worker)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Compress)
	r.Use(authMiddleware.JWT)

	r.Post("/api/shorten/batch", shortenBatchHandler.HandleHTTP)
	r.Post("/api/shorten", apiShortenHandler.HandleHTTP)
	r.Get("/api/user/urls", retrieveBatchHandler.HandleHTTP)
	r.Delete("/api/user/urls", deleteBatchHandler.HandleHTTP)
	r.Get("/ping", pingHandler.HandleHTTP)
	r.Get("/{short}", func(res http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), contextkeys.ShortURL, chi.URLParam(req, "short"))
		retrieveHandler.HandleHTTP(res, req.WithContext(ctx))
	})
	r.Post("/", shortenHandler.ServeHTTP)

	return &App{
		router:  r,
		storage: strg,
		worker:  worker,
		cfg:     cfg,
	}, nil
}

func (app *App) Run() error {
	defer app.storage.Close()
	defer logger.Log.Sync()

	app.worker.Run(tickerPeriod, app.cfg.Timeout)
	defer app.worker.Shutdown()

	logger.Log.Info(fmt.Sprintf("Server started successfully! Address: %s Storage Type: %s", app.cfg.ServerAddr, app.cfg.StorageType))

	err := http.ListenAndServe(app.cfg.ServerAddr, app.router)
	if err != nil {
		return err
	}
	return nil
}
