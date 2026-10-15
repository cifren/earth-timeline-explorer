// Package presentation exposes the geography context over HTTP.
package presentation

import (
	"context"

	"frise-chrono/back/internal/bc/geography/infrastructure"
	"frise-chrono/back/internal/core/zoning"
	"frise-chrono/back/internal/shared/apigen"
)

const (
	zonesCacheControl    = "public, max-age=86400"
	geometryCacheControl = "public, max-age=2592000"
)

// Handlers implements the geography operations of the API.
type Handlers struct {
	geometry apigen.ZoneFeatureCollection
	etag     string
}

// NewHandlers prepares the geometry response once from the loaded outlines.
func NewHandlers(outlines infrastructure.Outlines) *Handlers {
	return &Handlers{geometry: toFeatureCollection(outlines), etag: outlines.ETag}
}

// ListZones implements GET /zones.
func (h *Handlers) ListZones(_ context.Context, _ apigen.ListZonesRequestObject) (apigen.ListZonesResponseObject, error) {
	zones := zoning.All()
	items := make([]apigen.Zone, 0, len(zones))
	for _, zone := range zones {
		items = append(items, toAPIZone(zone))
	}
	cacheControl := zonesCacheControl
	return apigen.ListZones200JSONResponse{
		Body:    apigen.ZoneList{Items: items},
		Headers: apigen.ListZones200ResponseHeaders{CacheControl: &cacheControl},
	}, nil
}

// GetZoneGeometry implements GET /zones/geometry, answering 304 when the client copy
// is current.
func (h *Handlers) GetZoneGeometry(_ context.Context, request apigen.GetZoneGeometryRequestObject) (apigen.GetZoneGeometryResponseObject, error) {
	if request.Params.IfNoneMatch != nil && *request.Params.IfNoneMatch == h.etag {
		return apigen.GetZoneGeometry304Response{}, nil
	}
	cacheControl := geometryCacheControl
	etag := h.etag
	return apigen.GetZoneGeometry200ApplicationGeoPlusJSONResponse{
		Body:    h.geometry,
		Headers: apigen.GetZoneGeometry200ResponseHeaders{CacheControl: &cacheControl, ETag: &etag},
	}, nil
}

func toAPIZone(zone zoning.Zone) apigen.Zone {
	return apigen.Zone{Id: apigen.ZoneId(zone.ID), Name: zone.Name}
}

func toFeatureCollection(outlines infrastructure.Outlines) apigen.ZoneFeatureCollection {
	features := make([]apigen.ZoneFeature, 0, len(outlines.Zones))
	for _, outline := range outlines.Zones {
		polygons := make([][][][]float64, 0, len(outline.Polygons))
		for _, polygon := range outline.Polygons {
			rings := make([][][]float64, 0, len(polygon))
			for _, ring := range polygon {
				positions := make([][]float64, 0, len(ring))
				for _, position := range ring {
					positions = append(positions, []float64{position[0], position[1]})
				}
				rings = append(rings, positions)
			}
			polygons = append(polygons, rings)
		}
		features = append(features, apigen.ZoneFeature{
			Type:       "Feature",
			Properties: toAPIZone(outline.Zone),
			Geometry:   apigen.MultiPolygon{Type: "MultiPolygon", Coordinates: polygons},
		})
	}
	return apigen.ZoneFeatureCollection{Type: "FeatureCollection", Features: features}
}
