package domain

import (
	"slices"
	"testing"
)

func TestMapLinks(t *testing.T) {
	tests := []struct {
		name     string
		location Location
		want     []Link
	}{
		{
			name:     "known point",
			location: Location{Granularity: GranularityPoint, Point: &Point{Latitude: 48.8566, Longitude: 2.3522}},
			want: []Link{
				{Label: "GeoHack", URL: "https://geohack.toolforge.org/geohack.php?params=48.8566_N_2.3522_E"},
				{Label: "OpenStreetMap", URL: "https://www.openstreetmap.org/?mlat=48.8566&mlon=2.3522#map=16/48.8566/2.3522"},
				{Label: "Google Maps", URL: "https://www.google.com/maps/search/?api=1&query=48.8566%2C2.3522"},
			},
		},
		{
			name:     "southern and western hemispheres",
			location: Location{Granularity: GranularityPoint, Point: &Point{Latitude: -22.9068, Longitude: -43.1729}},
			want: []Link{
				{Label: "GeoHack", URL: "https://geohack.toolforge.org/geohack.php?params=22.9068_S_43.1729_W"},
				{Label: "OpenStreetMap", URL: "https://www.openstreetmap.org/?mlat=-22.9068&mlon=-43.1729#map=16/-22.9068/-43.1729"},
				{Label: "Google Maps", URL: "https://www.google.com/maps/search/?api=1&query=-22.9068%2C-43.1729"},
			},
		},
		{
			name:     "country only",
			location: Location{Granularity: GranularityCountry},
			want:     []Link{},
		},
		{
			name:     "point of the country only",
			location: Location{Granularity: GranularityCountry, Point: &Point{Latitude: 55.97, Longitude: 9.02}},
			want:     []Link{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.location.MapLinks(); !slices.Equal(got, test.want) {
				t.Fatalf("MapLinks() =\n%v\nwant\n%v", got, test.want)
			}
		})
	}
}
