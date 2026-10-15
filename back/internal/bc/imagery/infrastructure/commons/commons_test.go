package commons

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/internal/bc/imagery/infrastructure/commons/commonstest"
	"frise-chrono/back/internal/core/timespan"
	"frise-chrono/back/internal/core/zoning"
	"frise-chrono/back/internal/shared/cache"
)

const userAgent = "FriseChrono/test (https://example.org)"

func newStore(t *testing.T) *cache.Store {
	t.Helper()
	server := miniredis.RunT(t)
	return cache.New(redis.NewClient(&redis.Options{Addr: server.Addr()}))
}

func TestClientSendsUserAgentAndReportsAPIErrors(t *testing.T) {
	fake := &commonstest.Fake{Failing: true}
	client := NewClient(fake.Start(t), userAgent)

	_, err := NewCountryCatalog(client, newStore(t)).Countries(context.Background(), 1950)

	if !errors.Is(err, domain.ErrSourceUnavailable) {
		t.Fatalf("error = %v, want ErrSourceUnavailable", err)
	}
	if fake.UserAgent() != userAgent {
		t.Fatalf("User-Agent = %q, want %q", fake.UserAgent(), userAgent)
	}
}

func TestClientReportsUnreachableSource(t *testing.T) {
	client := NewClient("http://127.0.0.1:1", userAgent)

	_, _, err := NewPhotoSource(client).RandomPhoto(context.Background(), "France", 1950)

	if !errors.Is(err, domain.ErrSourceUnavailable) {
		t.Fatalf("error = %v, want ErrSourceUnavailable", err)
	}
}

func TestCountriesKeepsOnlyCountryCategories(t *testing.T) {
	fake := &commonstest.Fake{Countries: map[string][]string{"1950s": {"France", "the Netherlands"}}}
	catalog := NewCountryCatalog(NewClient(fake.Start(t), userAgent), newStore(t))

	countries, err := catalog.Countries(context.Background(), 1950)

	if err != nil || !slices.Equal(countries, []string{"France", "the Netherlands"}) {
		t.Fatalf("Countries = %v, %v", countries, err)
	}
}

func TestCountriesAreCached(t *testing.T) {
	fake := &commonstest.Fake{Countries: map[string][]string{"1950s": {"France"}}}
	catalog := NewCountryCatalog(NewClient(fake.Start(t), userAgent), newStore(t))

	for range 3 {
		if _, err := catalog.Countries(context.Background(), 1950); err != nil {
			t.Fatal(err)
		}
	}
	if fake.Requests() != 1 {
		t.Fatalf("Commons received %d requests, want 1", fake.Requests())
	}
}

func TestRandomPhotoMapsTheFile(t *testing.T) {
	fake := &commonstest.Fake{Files: map[string][]commonstest.File{"France": {{
		Title:          "File:Market day in Paris.jpg",
		URL:            "https://upload.wikimedia.org/market.jpg",
		Date:           "03/04/1953",
		License:        "CC BY-SA 4.0",
		Artist:         `<a href="https://example.org">J. Dupont</a>`,
		Credit:         "Collection of the Nationaal Archief, The Hague",
		HasCoordinates: true,
	}}}}
	source := NewPhotoSource(NewClient(fake.Start(t), userAgent))

	image, found, err := source.RandomPhoto(context.Background(), "France", 1950)

	if err != nil || !found {
		t.Fatalf("RandomPhoto = %v, %v", found, err)
	}
	checks := map[string][2]string{
		"title":          {image.Title, "Market day in Paris"},
		"id":             {image.ID(), "wikimedia_commons:File:Market day in Paris.jpg"},
		"date":           {image.Time.Label, "3 April 1953"},
		"zone":           {string(image.Location.Zone), string(zoning.Europe)},
		"granularity":    {string(image.Location.Granularity), string(domain.GranularityPoint)},
		"author":         {image.Provenance.Author, "J. Dupont"},
		"license":        {image.Provenance.License, "CC BY-SA 4.0"},
		"authority":      {image.Certification.Authority, "Wikimedia Commons"},
		"originalHolder": {image.Certification.OriginalHolder, "Nationaal Archief"},
		"thumbnail":      {image.File.ThumbnailURL, "https://upload.wikimedia.org/market.jpg?width=1280"},
		"learnMore":      {image.LearnMore[0].URL, "https://en.wikipedia.org/wiki/History_of_France"},
	}
	for field, values := range checks {
		if values[0] != values[1] {
			t.Errorf("%s = %q, want %q", field, values[0], values[1])
		}
	}
}

func TestRandomPhotoReportsEmptyCategory(t *testing.T) {
	source := NewPhotoSource(NewClient((&commonstest.Fake{}).Start(t), userAgent))

	_, found, err := source.RandomPhoto(context.Background(), "France", 1950)

	if err != nil || found {
		t.Fatalf("RandomPhoto = %v, %v; want not found, no error", found, err)
	}
}

func TestToImage(t *testing.T) {
	usable := page{Title: "File:Harbour.jpg", ImageInfo: []imageInfo{{URL: "https://upload.wikimedia.org/h.jpg"}}}
	tests := []struct {
		name    string
		page    page
		country string
		want    bool
	}{
		{name: "usable file", page: usable, country: "France", want: true},
		{name: "missing file", page: page{Title: "File:Gone.jpg", Missing: true}, country: "France", want: false},
		{name: "no image information", page: page{Title: "File:Empty.jpg"}, country: "France", want: false},
		{name: "country without zone", page: usable, country: "Atlantis", want: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			image, ok := toImage(test.page, test.country, timespan.Decade(1950))
			if ok != test.want {
				t.Fatalf("toImage usable = %v, want %v", ok, test.want)
			}
			if ok && (image.Location.Granularity != domain.GranularityCountry || image.Time.Label != "1950s") {
				t.Fatalf("toImage = %+v", image)
			}
		})
	}
}
