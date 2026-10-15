package presentation

import (
	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/var/apigen"
)

func toAPIImage(image domain.Image) apigen.Image {
	return apigen.Image{
		Id:    image.ID(),
		Title: image.Title,
		Kind:  apigen.ImageKind(image.Kind),
		File: apigen.ImageFile{
			ThumbnailUrl: image.File.ThumbnailURL,
			DisplayUrl:   image.File.DisplayURL,
			OriginalUrl:  image.File.OriginalURL,
			Width:        image.File.Width,
			Height:       image.File.Height,
		},
		Time: apigen.TimeSpan{Edtf: image.Time.EDTF, Label: image.Time.Label},
		Location: apigen.Location{
			Zone:        apigen.ZoneId(image.Location.Zone),
			Granularity: apigen.LocationGranularity(image.Location.Granularity),
			Point:       toAPIPoint(image.Location.Point),
			MapLinks:    toAPILinks(image.Location.MapLinks()),
		},
		Certification: apigen.Certification{
			Level:          apigen.CertificationLevel(image.Certification.Level),
			Authority:      image.Certification.Authority,
			OriginalHolder: optional(image.Certification.OriginalHolder),
		},
		Provenance: apigen.Provenance{
			Source:      apigen.ProvenanceSource(image.Provenance.Source),
			SourceId:    image.Provenance.SourceID,
			SourceUrl:   image.Provenance.SourceURL,
			License:     image.Provenance.License,
			Author:      optional(image.Provenance.Author),
			RetrievedAt: image.Provenance.RetrievedAt,
		},
		LearnMore: toAPILinks(image.LearnMore),
	}
}

func toAPIPoint(point *domain.Point) *apigen.Point {
	if point == nil {
		return nil
	}
	return &apigen.Point{Type: apigen.PointTypePoint, Coordinates: []float64{point.Longitude, point.Latitude}}
}

func toAPILinks(links []domain.Link) []apigen.Link {
	result := make([]apigen.Link, 0, len(links))
	for _, link := range links {
		result = append(result, apigen.Link{Label: link.Label, Url: link.URL})
	}
	return result
}

func optional(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
