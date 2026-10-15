package presentation

import (
	"context"
	"testing"

	"frise-chrono/back/internal/bc/geography/domain"
	"frise-chrono/back/internal/bc/geography/infrastructure"
	"frise-chrono/back/internal/core/zoning"
	"frise-chrono/back/var/apigen"
)

func newHandlers() *Handlers {
	europe, _ := zoning.Find(zoning.Europe)
	square := domain.Polygon{domain.Ring{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}}
	return NewHandlers(infrastructure.Outlines{
		Zones: []domain.Outline{{Zone: europe, Polygons: []domain.Polygon{square}}},
		ETag:  `"abc"`,
	})
}

func TestListZones(t *testing.T) {
	response, err := newHandlers().ListZones(context.Background(), apigen.ListZonesRequestObject{})
	if err != nil {
		t.Fatal(err)
	}
	zones := response.(apigen.ListZones200JSONResponse)
	if len(zones.Body.Items) != 9 || zones.Body.Items[0].Id != "northern-america" {
		t.Fatalf("zones = %+v", zones.Body.Items)
	}
	if zones.Headers.CacheControl == nil {
		t.Fatal("missing Cache-Control")
	}
}

func TestGetZoneGeometry(t *testing.T) {
	current := `"abc"`
	outdated := `"old"`
	tests := []struct {
		name        string
		ifNoneMatch *string
		wantStatus  int
	}{
		{name: "first download", ifNoneMatch: nil, wantStatus: 200},
		{name: "outdated copy", ifNoneMatch: &outdated, wantStatus: 200},
		{name: "current copy", ifNoneMatch: &current, wantStatus: 304},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := newHandlers().GetZoneGeometry(context.Background(), apigen.GetZoneGeometryRequestObject{
				Params: apigen.GetZoneGeometryParams{IfNoneMatch: test.ifNoneMatch},
			})
			if err != nil {
				t.Fatal(err)
			}
			switch geometry := response.(type) {
			case apigen.GetZoneGeometry304Response:
				if test.wantStatus != 304 {
					t.Fatalf("got 304, want %d", test.wantStatus)
				}
			case apigen.GetZoneGeometry200ApplicationGeoPlusJSONResponse:
				if test.wantStatus != 200 || len(geometry.Body.Features) != 1 || *geometry.Headers.ETag != current {
					t.Fatalf("got 200 %+v, want %d", geometry.Headers, test.wantStatus)
				}
			default:
				t.Fatalf("unexpected response %T", response)
			}
		})
	}
}
