package presentation

import (
	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/internal/shared/apigen"
)

func toAPIImage(image domain.Image) apigen.Image {
	links := make([]apigen.Link, 0, len(image.LearnMore))
	for _, link := range image.LearnMore {
		links = append(links, apigen.Link{Label: link.Label, Url: link.URL})
	}
	return apigen.Image{
		Id:    image.ID(),
		Title: image.Title,
		Kind:  apigen.ImageKind(image.Kind),
		File: apigen.ImageFile{
			ThumbnailUrl: image.File.ThumbnailURL,
			OriginalUrl:  image.File.OriginalURL,
			Width:        image.File.Width,
			Height:       image.File.Height,
		},
		Time: apigen.TimeSpan{Edtf: image.Time.EDTF, Label: image.Time.Label},
		Location: apigen.Location{
			Zone:        apigen.ZoneId(image.Location.Zone),
			Granularity: apigen.LocationGranularity(image.Location.Granularity),
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
		LearnMore: links,
	}
}

func optional(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
