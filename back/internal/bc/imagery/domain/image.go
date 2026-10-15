// Package domain holds the image model of the imagery context.
package domain

import (
	"errors"
	"slices"
	"strings"
	"time"

	"frise-chrono/back/internal/core/timespan"
	"frise-chrono/back/internal/core/zoning"
)

// ErrSourceUnavailable is returned when an image source cannot be reached in time.
var ErrSourceUnavailable = errors.New("image source unavailable")

// Kind is the nature of an image.
type Kind string

// KindPhotograph is a photograph of the period it shows.
const KindPhotograph Kind = "photograph"

// Granularity tells how precisely the location of an image is known.
type Granularity string

// Location granularities.
const (
	GranularityPoint   Granularity = "point"
	GranularityCountry Granularity = "country"
)

// CertificationLevel is the trust level shown to the user.
type CertificationLevel string

// CertificationSource means the image is implicitly certified by its source.
const CertificationSource CertificationLevel = "source"

// Image is an image anchored in a zone and a period.
type Image struct {
	Title         string
	Kind          Kind
	File          File
	Time          timespan.TimeSpan
	Location      Location
	Certification Certification
	Provenance    Provenance
	LearnMore     []Link
}

// ID identifies an image across sources and holds what is needed to fetch it again: its
// source, its country and its identifier at the source, e.g.
// "wikimedia_commons:France:File:Market day.jpg".
func (i Image) ID() string {
	return strings.Join([]string{i.Provenance.Source, i.Location.Country, i.Provenance.SourceID}, ":")
}

// Reference is an image identifier split into its parts.
type Reference struct {
	Source   string
	Country  string
	SourceID string
}

// ParseID splits an image identifier built by Image.ID.
func ParseID(id string) (Reference, bool) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 || slices.Contains(parts, "") {
		return Reference{}, false
	}
	return Reference{Source: parts[0], Country: parts[1], SourceID: parts[2]}, true
}

// File locates the image files at their source.
type File struct {
	ThumbnailURL string
	// DisplayURL is a large rendition, in a format browsers show, for full-screen viewing.
	DisplayURL  string
	OriginalURL string
	Width       int
	Height      int
}

// Location is where the image was taken. The country is internal: clients see the zone.
type Location struct {
	Zone        zoning.ID
	Country     string
	Granularity Granularity
	// Point is where the image was taken with GranularityPoint, a point inside the
	// country with GranularityCountry; it may be unknown for the latter.
	Point *Point
}

// Point is a WGS 84 position, in decimal degrees.
type Point struct {
	Latitude  float64
	Longitude float64
}

// Certification tells who vouches for the date and location of the image.
type Certification struct {
	Level          CertificationLevel
	Authority      string
	OriginalHolder string
}

// Provenance tells where the image was fetched from.
type Provenance struct {
	Source      string
	SourceID    string
	SourceURL   string
	License     string
	Author      string
	RetrievedAt time.Time
}

// Link points to further reading.
type Link struct {
	Label string
	URL   string
}
