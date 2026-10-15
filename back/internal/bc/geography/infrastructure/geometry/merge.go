// Package geometry merges and simplifies borders with GEOS. It requires cgo and is used
// only by the zones command, never by the HTTP service.
package geometry

import (
	"encoding/json"
	"fmt"

	"github.com/twpayne/go-geos"

	"frise-chrono/back/internal/bc/geography/domain"
)

// Tuning in degrees: gaps between neighbouring countries narrower than closeGaps are
// filled so that no border remains inside a zone; simplify keeps the file light;
// precision rounds coordinates.
const (
	closeGaps = 0.05
	simplify  = 0.05
	precision = 0.001
)

// Merge unions GeoJSON geometries into one simplified outline.
func Merge(geometries []json.RawMessage) ([]domain.Polygon, error) {
	parts := make([]*geos.Geom, 0, len(geometries))
	for _, raw := range geometries {
		part, err := geos.NewGeomFromGeoJSON(string(raw))
		if err != nil {
			return nil, fmt.Errorf("read geometry: %w", err)
		}
		parts = append(parts, part.MakeValid())
	}

	merged := geos.NewCollection(geos.TypeIDGeometryCollection, parts).
		UnaryUnion().
		Buffer(closeGaps, 4).
		Buffer(-closeGaps, 4).
		TopologyPreserveSimplify(simplify).
		SetPrecision(precision, geos.PrecisionRuleNone)

	return toPolygons(merged.ToGeoJSON(0))
}

// Part is one polygon of a geometry, with a point lying inside it.
type Part struct {
	Geometry  json.RawMessage
	Longitude float64
	Latitude  float64
}

// Parts splits a GeoJSON Polygon or MultiPolygon into its polygons.
func Parts(raw json.RawMessage) ([]Part, error) {
	geom, err := geos.NewGeomFromGeoJSON(string(raw))
	if err != nil {
		return nil, fmt.Errorf("read geometry: %w", err)
	}
	parts := make([]Part, 0, geom.NumGeometries())
	for i := range geom.NumGeometries() {
		polygon := geom.Geometry(i)
		inside := polygon.PointOnSurface()
		parts = append(parts, Part{
			Geometry:  json.RawMessage(polygon.ToGeoJSON(0)),
			Longitude: inside.X(),
			Latitude:  inside.Y(),
		})
	}
	return parts, nil
}

type geoJSONGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

// toPolygons reads a GeoJSON Polygon or MultiPolygon as a list of polygons.
func toPolygons(geoJSON string) ([]domain.Polygon, error) {
	var geometry geoJSONGeometry
	if err := json.Unmarshal([]byte(geoJSON), &geometry); err != nil {
		return nil, err
	}
	switch geometry.Type {
	case "Polygon":
		var polygon domain.Polygon
		err := json.Unmarshal(geometry.Coordinates, &polygon)
		return []domain.Polygon{polygon}, err
	case "MultiPolygon":
		var polygons []domain.Polygon
		err := json.Unmarshal(geometry.Coordinates, &polygons)
		return polygons, err
	default:
		return nil, fmt.Errorf("unexpected merged geometry type %q", geometry.Type)
	}
}
