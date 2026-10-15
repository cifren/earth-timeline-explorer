// Package application holds the use cases of the imagery context.
package application

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"slices"

	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/internal/bc/imagery/infrastructure/commons"
	"frise-chrono/back/internal/core/timespan"
	"frise-chrono/back/internal/core/zoning"
)

// attemptsPerImage bounds the draws spent on one image when files are missing.
const attemptsPerImage = 3

// Drawer draws random images of a zone and a decade.
type Drawer struct {
	countries *commons.CountryCatalog
	photos    *commons.PhotoSource
}

// NewDrawer returns a drawer reading the given catalog and photo source.
func NewDrawer(countries *commons.CountryCatalog, photos *commons.PhotoSource) *Drawer {
	return &Drawer{countries: countries, photos: photos}
}

// ZonesWithImages returns the zones having at least one country with photographs for
// the decade, in display order.
func (d *Drawer) ZonesWithImages(ctx context.Context, decade timespan.Decade) ([]zoning.ID, error) {
	countries, err := d.countries.Countries(ctx, decade)
	if err != nil {
		return nil, err
	}
	var zones []zoning.ID
	for _, zone := range zoning.All() {
		if len(countriesIn(ctx, zone.ID, countries)) > 0 {
			zones = append(zones, zone.ID)
		}
	}
	return zones, nil
}

// Draw returns up to limit distinct random images. Each image comes from a country of
// the zone drawn uniformly; a missing file triggers a new draw, within a bounded number
// of attempts. Images already drawn are kept when the source fails midway.
func (d *Drawer) Draw(ctx context.Context, zone zoning.ID, decade timespan.Decade, limit int) ([]domain.Image, error) {
	all, err := d.countries.Countries(ctx, decade)
	if err != nil {
		return nil, err
	}
	countries := countriesIn(ctx, zone, all)
	if len(countries) == 0 {
		return []domain.Image{}, nil
	}

	images := make([]domain.Image, 0, limit)
	for attempt := 0; attempt < limit*attemptsPerImage && len(images) < limit; attempt++ {
		country := countries[rand.IntN(len(countries))]
		image, found, err := d.photos.RandomPhoto(ctx, country, decade)
		if err != nil && len(images) > 0 {
			return images, nil
		}
		if err != nil {
			return nil, err
		}
		if !found || slices.ContainsFunc(images, func(drawn domain.Image) bool { return drawn.ID() == image.ID() }) {
			continue
		}
		images = append(images, image)
	}
	return images, nil
}

// Image fetches again an image drawn for the decade, from its identifier. It reports false
// when the identifier is malformed or the image no longer exists at its source.
func (d *Drawer) Image(ctx context.Context, id string, decade timespan.Decade) (domain.Image, bool, error) {
	reference, valid := domain.ParseID(id)
	if !valid {
		return domain.Image{}, false, nil
	}
	return d.photos.Photo(ctx, reference, decade)
}

func countriesIn(ctx context.Context, zone zoning.ID, countries []string) []string {
	var inZone []string
	for _, country := range countries {
		countryZone, known := zoning.ZoneOfCountry(country)
		if !known {
			slog.DebugContext(ctx, "country without zone", "country", country)
			continue
		}
		if countryZone == zone {
			inZone = append(inZone, country)
		}
	}
	return inZone
}
