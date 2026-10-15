// Command server runs the Frise Chrono HTTP API.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"frise-chrono/back/internal/app"
	"frise-chrono/back/internal/shared/telemetry"
)

func main() {
	slog.SetDefault(slog.New(telemetry.NewLogHandler(slog.NewJSONHandler(os.Stdout, nil))))

	config, err := app.LoadConfig()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, config); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
