package presentation

import (
	"context"
	"testing"
	"time"

	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/internal/core/timespan"
	"frise-chrono/back/internal/core/zoning"
	"frise-chrono/back/var/apigen"
)

func TestDrawImagesRejectsInvalidRequests(t *testing.T) {
	two := 2
	tooMany := 6
	tests := []struct {
		name   string
		params apigen.DrawImagesParams
	}{
		{name: "decade not a multiple of 10", params: apigen.DrawImagesParams{Zone: "europe", Decade: 1955}},
		{name: "decade before the timeline", params: apigen.DrawImagesParams{Zone: "europe", Decade: 1700}},
		{name: "unknown zone", params: apigen.DrawImagesParams{Zone: "atlantis", Decade: 1950, Limit: &two}},
		{name: "limit too high", params: apigen.DrawImagesParams{Zone: "europe", Decade: 1950, Limit: &tooMany}},
	}
	handlers := NewHandlers(nil)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := handlers.DrawImages(context.Background(), apigen.DrawImagesRequestObject{Params: test.params})
			if err != nil {
				t.Fatal(err)
			}
			problem, isBadRequest := response.(apigen.DrawImages400ApplicationProblemPlusJSONResponse)
			if !isBadRequest || problem.Status != 400 || problem.Detail == nil {
				t.Fatalf("response = %#v, want a 400 problem with a detail", response)
			}
		})
	}
}

func TestListZonesWithImagesRejectsInvalidDecade(t *testing.T) {
	response, err := NewHandlers(nil).ListZonesWithImages(context.Background(), apigen.ListZonesWithImagesRequestObject{Decade: 1955})
	if err != nil {
		t.Fatal(err)
	}
	if _, isBadRequest := response.(apigen.ListZonesWithImages400ApplicationProblemPlusJSONResponse); !isBadRequest {
		t.Fatalf("response = %#v, want a 400 problem", response)
	}
}

func TestToAPIImage(t *testing.T) {
	image := domain.Image{
		Title:         "Harbour",
		Kind:          domain.KindPhotograph,
		Time:          timespan.TimeSpan{EDTF: "195X", Label: "1950s"},
		Location:      domain.Location{Zone: zoning.Europe, Country: "France", Granularity: domain.GranularityCountry},
		Certification: domain.Certification{Level: domain.CertificationSource, Authority: "Wikimedia Commons"},
		Provenance: domain.Provenance{
			Source: "wikimedia_commons", SourceID: "File:Harbour.jpg", License: "Public domain",
			RetrievedAt: time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	got := toAPIImage(image)

	if got.Id != "wikimedia_commons:France:File:Harbour.jpg" || got.Location.Zone != "europe" || got.Time.Label != "1950s" {
		t.Fatalf("toAPIImage = %+v", got)
	}
	if got.Provenance.Author != nil || got.Certification.OriginalHolder != nil {
		t.Fatal("empty optional fields must be omitted")
	}
	if got.LearnMore == nil {
		t.Fatal("learnMore must be an empty list, not null")
	}
	if got.Location.MapLinks == nil || len(got.Location.MapLinks) != 0 {
		t.Fatalf("mapLinks = %+v, want an empty list without a point", got.Location.MapLinks)
	}
	if got.Location.Point != nil {
		t.Fatal("a country-level location has no point")
	}

	image.Location = domain.Location{Zone: zoning.Europe, Country: "France", Granularity: domain.GranularityPoint, Point: &domain.Point{Latitude: 48.85, Longitude: 2.35}}
	point := toAPIImage(image).Location.Point
	if point == nil || point.Type != apigen.PointTypePoint || point.Coordinates[0] != 2.35 || point.Coordinates[1] != 48.85 {
		t.Fatalf("point = %+v, want a GeoJSON point [longitude, latitude]", point)
	}
}
