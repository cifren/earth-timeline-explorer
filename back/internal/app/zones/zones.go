// Package zones is the composition root of the zones command. It is separate from the
// HTTP service's composition root because it links GEOS, which requires cgo.
package zones

import (
	"context"
	"log/slog"

	"frise-chrono/back/internal/bc/geography/application"
	"frise-chrono/back/internal/bc/geography/infrastructure"
)

// Run regenerates the zone outlines and country points embedded in the HTTP service.
func Run(ctx context.Context) error {
	if err := application.Generate(ctx, infrastructure.OutlinesFile, infrastructure.CountryPointsFile); err != nil {
		return err
	}
	slog.InfoContext(ctx, "geographic data written", "outlines", infrastructure.OutlinesFile, "countryPoints", infrastructure.CountryPointsFile)
	return nil
}
