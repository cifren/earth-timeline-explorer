// Package zones is the composition root of the zones command. It is separate from the
// HTTP service's composition root because it links GEOS, which requires cgo.
package zones

import (
	"context"
	"log/slog"

	"frise-chrono/back/internal/bc/geography/application"
	"frise-chrono/back/internal/bc/geography/infrastructure"
)

// Run regenerates the zone outlines file embedded in the HTTP service.
func Run(ctx context.Context) error {
	if err := application.GenerateOutlines(ctx, infrastructure.OutlinesFile); err != nil {
		return err
	}
	slog.InfoContext(ctx, "zone outlines written", "path", infrastructure.OutlinesFile)
	return nil
}
