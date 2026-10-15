// Package domain holds the geometry model of the geography context.
package domain

import "frise-chrono/back/internal/core/zoning"

// Position is a [longitude, latitude] pair in WGS 84 degrees.
type Position [2]float64

// Ring is a closed line of positions; the first ring of a polygon is its outer border.
type Ring []Position

// Polygon is an outer ring followed by its holes.
type Polygon []Ring

// Outline is the shape of a zone, made of one or more polygons.
type Outline struct {
	Zone     zoning.Zone
	Polygons []Polygon
}
