// Command zones regenerates the zone outlines and country points from Natural Earth.
// Run it from the back module root: `go run ./cmd/zones`.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"frise-chrono/back/internal/app/zones"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := zones.Run(ctx); err != nil {
		slog.Error("zone generation failed", "error", err)
		os.Exit(1)
	}
}
