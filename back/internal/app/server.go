package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	geographyinfra "frise-chrono/back/internal/bc/geography/infrastructure"
	geography "frise-chrono/back/internal/bc/geography/presentation"
	imageryapp "frise-chrono/back/internal/bc/imagery/application"
	"frise-chrono/back/internal/bc/imagery/infrastructure/commons"
	imagery "frise-chrono/back/internal/bc/imagery/presentation"
	"frise-chrono/back/internal/shared/cache"
	"frise-chrono/back/internal/shared/problem"
	"frise-chrono/back/internal/shared/telemetry"
	"frise-chrono/back/var/apigen"
)

const (
	apiBaseURL      = "/api/v1"
	shutdownTimeout = 10 * time.Second
)

// Aliases give each context's handlers a distinct embedded field name.
type (
	geographyAPI = geography.Handlers
	imageryAPI   = imagery.Handlers
)

// api assembles the operations of every context into the generated interface.
type api struct {
	*geographyAPI
	*imageryAPI
}

// Run serves the API until ctx is cancelled, then shuts down gracefully.
func Run(ctx context.Context, config Config) error {
	if config.TracesEnabled {
		stopTracing, err := telemetry.StartTracing(ctx)
		if err != nil {
			return err
		}
		defer func() { _ = stopTracing(context.Background()) }()
	}

	redisClient := redis.NewClient(&redis.Options{Addr: config.RedisAddr})
	defer func() { _ = redisClient.Close() }()
	if err := redisotel.InstrumentTracing(redisClient); err != nil {
		return fmt.Errorf("redis tracing: %w", err)
	}

	handler, err := newHandler(config, commons.Endpoint, redisClient)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              config.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server listening", "addr", config.HTTPAddr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server stopped: %w", err)
	case <-ctx.Done():
	}

	slog.Info("server shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}

// newHandler wires the contexts into the HTTP handler of the API, calling Wikimedia
// Commons at the given endpoint.
func newHandler(config Config, commonsEndpoint string, redisClient *redis.Client) (http.Handler, error) {
	outlines, err := geographyinfra.LoadOutlines()
	if err != nil {
		return nil, err
	}

	commonsClient := commons.NewClient(commonsEndpoint, config.CommonsUserAgent)
	store := cache.New(redisClient)
	drawer := imageryapp.NewDrawer(commons.NewCountryCatalog(commonsClient, store), commons.NewPhotoSource(commonsClient))

	strict := apigen.NewStrictHandlerWithOptions(
		api{geography.NewHandlers(outlines), imagery.NewHandlers(drawer)},
		nil,
		apigen.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  badRequest,
			ResponseErrorHandlerFunc: internalError,
		},
	)

	mux := http.NewServeMux()
	apigen.HandlerWithOptions(strict, apigen.StdHTTPServerOptions{
		BaseURL:          apiBaseURL,
		BaseRouter:       mux,
		ErrorHandlerFunc: badRequest,
	})
	mux.HandleFunc("GET /api/health", health)
	mux.HandleFunc("GET /api/openapi.json", openAPIDocument)
	mux.HandleFunc("GET /api/docs", documentation)

	return otelhttp.NewHandler(withRequestLog(withCORS(config.CORSAllowedOrigins, withRecovery(mux))), "api"), nil
}

func badRequest(w http.ResponseWriter, _ *http.Request, err error) {
	problem.Write(w, problem.New(http.StatusBadRequest, err.Error()))
}

func internalError(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), "request failed", "error", err)
	problem.Write(w, problem.New(http.StatusInternalServerError, "An unexpected error occurred."))
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
