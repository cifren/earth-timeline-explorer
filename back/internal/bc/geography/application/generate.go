// Package application holds the use cases of the geography context.
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"frise-chrono/back/internal/bc/geography/domain"
	"frise-chrono/back/internal/bc/geography/infrastructure"
	"frise-chrono/back/internal/bc/geography/infrastructure/geometry"
	"frise-chrono/back/internal/bc/geography/infrastructure/naturalearth"
	"frise-chrono/back/internal/core/zoning"
)

// GenerateOutlines downloads the country borders, merges them zone by zone and writes
// the outlines file served by the API. Each polygon of a country joins the country's
// zone, except overseas territories, which join the zone they lie in.
func GenerateOutlines(ctx context.Context, path string) error {
	countries, err := naturalearth.FetchCountries(ctx, naturalearth.CountriesURL)
	if err != nil {
		return err
	}

	bordersByZone := map[zoning.ID][]json.RawMessage{}
	for _, country := range countries {
		zone, known := zoning.ZoneOfRegion(country.Region, country.Subregion)
		if !known {
			slog.InfoContext(ctx, "country outside every zone", "country", country.Name, "region", country.Region, "subregion", country.Subregion)
			continue
		}
		parts, err := geometry.Parts(country.Geometry)
		if err != nil {
			return fmt.Errorf("split %s: %w", country.Name, err)
		}
		for _, part := range parts {
			partZone := zone
			if territoryZone, overseas := zoning.ZoneOfTerritory(part.Longitude, part.Latitude); overseas {
				partZone = territoryZone
			}
			bordersByZone[partZone] = append(bordersByZone[partZone], part.Geometry)
		}
	}

	outlines := make([]domain.Outline, 0, len(bordersByZone))
	for _, zone := range zoning.All() {
		borders := bordersByZone[zone.ID]
		if len(borders) == 0 {
			return fmt.Errorf("no country found for zone %s", zone.ID)
		}
		polygons, err := geometry.Merge(borders)
		if err != nil {
			return fmt.Errorf("merge zone %s: %w", zone.ID, err)
		}
		slog.InfoContext(ctx, "zone merged", "zone", zone.ID, "countries", len(borders), "polygons", len(polygons))
		outlines = append(outlines, domain.Outline{Zone: zone, Polygons: polygons})
	}
	return infrastructure.WriteOutlines(path, outlines)
}
