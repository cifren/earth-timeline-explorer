package application

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/internal/bc/imagery/infrastructure/commons"
	"frise-chrono/back/internal/bc/imagery/infrastructure/commons/commonstest"
	"frise-chrono/back/internal/core/zoning"
	"frise-chrono/back/internal/shared/cache"
)

func newDrawer(t *testing.T, fake *commonstest.Fake) *Drawer {
	t.Helper()
	client := commons.NewClient(fake.Start(t), "FriseChrono/test")
	store := cache.New(redis.NewClient(&redis.Options{Addr: miniredis.RunT(t).Addr()}))
	return NewDrawer(commons.NewCountryCatalog(client, store), commons.NewPhotoSource(client))
}

func files(country string, count int) []commonstest.File {
	result := make([]commonstest.File, 0, count)
	for i := range count {
		result = append(result, commonstest.File{
			Title: fmt.Sprintf("File:%s %d.jpg", country, i),
			URL:   fmt.Sprintf("https://upload.wikimedia.org/%s-%d.jpg", country, i),
		})
	}
	return result
}

func TestZonesWithImages(t *testing.T) {
	drawer := newDrawer(t, &commonstest.Fake{Countries: map[string][]string{
		"1950s": {"France", "Japan", "Atlantis", "Brazil"},
	}})

	zones, err := drawer.ZonesWithImages(context.Background(), 1950)

	want := []zoning.ID{zoning.LatinAmericaCaribbean, zoning.Europe, zoning.EastAndSouthEastAsia}
	if err != nil || !slices.Equal(zones, want) {
		t.Fatalf("ZonesWithImages = %v, %v; want %v", zones, err, want)
	}
}

func TestDraw(t *testing.T) {
	countries := map[string][]string{"1950s": {"France", "Japan"}}
	tests := []struct {
		name      string
		files     map[string][]commonstest.File
		zone      zoning.ID
		limit     int
		wantCount int
	}{
		{name: "one image of the zone", files: map[string][]commonstest.File{"France": files("France", 3)}, zone: zoning.Europe, limit: 1, wantCount: 1},
		{name: "several distinct images", files: map[string][]commonstest.File{"France": files("France", 5)}, zone: zoning.Europe, limit: 5, wantCount: 5},
		{name: "fewer files than requested", files: map[string][]commonstest.File{"France": files("France", 2)}, zone: zoning.Europe, limit: 5, wantCount: 2},
		{name: "empty categories", files: nil, zone: zoning.Europe, limit: 1, wantCount: 0},
		{name: "zone without country", files: map[string][]commonstest.File{"France": files("France", 3)}, zone: zoning.Africa, limit: 1, wantCount: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			drawer := newDrawer(t, &commonstest.Fake{Countries: countries, Files: test.files})

			images, err := drawer.Draw(context.Background(), test.zone, 1950, test.limit)

			if err != nil || len(images) != test.wantCount {
				t.Fatalf("Draw returned %d images, %v; want %d", len(images), err, test.wantCount)
			}
			ids := map[string]bool{}
			for _, image := range images {
				if ids[image.ID()] {
					t.Fatalf("image %s drawn twice", image.ID())
				}
				ids[image.ID()] = true
				if image.Location.Zone != test.zone {
					t.Fatalf("image of zone %s, want %s", image.Location.Zone, test.zone)
				}
			}
		})
	}
}

func TestDrawReportsUnavailableSource(t *testing.T) {
	drawer := newDrawer(t, &commonstest.Fake{Failing: true})

	_, err := drawer.Draw(context.Background(), zoning.Europe, 1950, 1)

	if !errors.Is(err, domain.ErrSourceUnavailable) {
		t.Fatalf("error = %v, want ErrSourceUnavailable", err)
	}
}
