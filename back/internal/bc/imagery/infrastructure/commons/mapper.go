package commons

import (
	"html"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/internal/core/timespan"
	"frise-chrono/back/internal/core/zoning"
)

const (
	sourceName = "wikimedia_commons"
	authority  = "Wikimedia Commons"
)

var (
	htmlTags   = regexp.MustCompile(`<[^>]*>`)
	whitespace = regexp.MustCompile(`\s+`)
)

// institutions are original holders recognized in the Commons "Credit" field.
var institutions = []string{
	"Nationaal Archief",
	"Bibliothèque nationale de France",
	"Library of Congress",
	"Bundesarchiv",
	"ETH-Bibliothek",
	"Médiathèque de l'architecture et du patrimoine",
	"National Archives and Records Administration",
	"State Library of Victoria",
	"State Library of New South Wales",
	"National Library of Australia",
	"National Library of Ireland",
	"Rijksmuseum",
	"Smithsonian Institution",
	"Imperial War Museums",
	"Archives nationales",
	"Deutsche Fotothek",
}

// toImage maps a Commons file page to an image; it reports false when the page has no
// usable file.
func toImage(p page, country string, decade timespan.Decade) (domain.Image, bool) {
	zone, known := zoning.ZoneOfCountry(country)
	if p.Missing || len(p.ImageInfo) == 0 || !known {
		return domain.Image{}, false
	}
	info := p.ImageInfo[0]
	if info.URL == "" {
		return domain.Image{}, false
	}

	location := toLocation(p, zone, country)

	return domain.Image{
		Title: title(p.Title),
		Kind:  domain.KindPhotograph,
		File: domain.File{
			ThumbnailURL: firstNonEmpty(info.ThumbURL, info.URL),
			DisplayURL:   firstNonEmpty(info.ResponsiveURLs["2"], info.ResponsiveURLs["1.5"], info.ThumbURL, info.URL),
			OriginalURL:  info.URL,
			Width:        info.Width,
			Height:       info.Height,
		},
		Time:     timespan.FromText(info.ExtMetadata["DateTimeOriginal"].Value, decade),
		Location: location,
		Certification: domain.Certification{
			Level:          domain.CertificationSource,
			Authority:      authority,
			OriginalHolder: originalHolder(info.ExtMetadata["Credit"].Value),
		},
		Provenance: domain.Provenance{
			Source:      sourceName,
			SourceID:    p.Title,
			SourceURL:   info.DescriptionURL,
			License:     plainText(info.ExtMetadata["LicenseShortName"].Value),
			Author:      plainText(info.ExtMetadata["Artist"].Value),
			RetrievedAt: time.Now().UTC(),
		},
		LearnMore: []domain.Link{historyArticle(country)},
	}, true
}

// toLocation places the image at its coordinates when Commons has them, otherwise at a
// point inside its country, when that country has one.
func toLocation(p page, zone zoning.ID, country string) domain.Location {
	location := domain.Location{Zone: zone, Country: country, Granularity: domain.GranularityCountry}
	if len(p.Coordinates) > 0 {
		location.Granularity = domain.GranularityPoint
		location.Point = &domain.Point{Latitude: p.Coordinates[0].Lat, Longitude: p.Coordinates[0].Lon}
		return location
	}
	if point, found := zoning.PointOfCountry(country); found {
		location.Point = &domain.Point{Latitude: point[1], Longitude: point[0]}
	}
	return location
}

// title turns "File:Market day in Berlin.jpg" into "Market day in Berlin".
func title(pageTitle string) string {
	name := strings.TrimPrefix(pageTitle, "File:")
	return strings.TrimSuffix(name, path.Ext(name))
}

func originalHolder(credit string) string {
	text := plainText(credit)
	for _, institution := range institutions {
		if strings.Contains(text, institution) {
			return institution
		}
	}
	return ""
}

// historyArticle links to the English Wikipedia article on the history of the country.
func historyArticle(country string) domain.Link {
	label := "History of " + country
	return domain.Link{
		Label: label,
		URL:   "https://en.wikipedia.org/wiki/" + url.PathEscape(strings.ReplaceAll(label, " ", "_")),
	}
}

func plainText(value string) string {
	text := html.UnescapeString(htmlTags.ReplaceAllString(value, " "))
	return strings.TrimSpace(whitespace.ReplaceAllString(text, " "))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
