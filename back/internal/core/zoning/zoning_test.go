package zoning

import "testing"

func TestAllKeepsDisplayOrder(t *testing.T) {
	want := []ID{NorthernAmerica, LatinAmericaCaribbean, Europe, Africa, MiddleEast,
		SouthAndCentralAsia, EastAndSouthEastAsia, Oceania, Antarctica}
	zones := All()
	if len(zones) != len(want) {
		t.Fatalf("All() returned %d zones, want %d", len(zones), len(want))
	}
	for i, zone := range zones {
		if zone.ID != want[i] {
			t.Errorf("All()[%d] = %s, want %s", i, zone.ID, want[i])
		}
	}
}

func TestAllReturnsACopy(t *testing.T) {
	All()[0].Name = "changed"
	if All()[0].Name == "changed" {
		t.Fatal("All() exposes its internal slice")
	}
}

func TestFind(t *testing.T) {
	zone, found := Find(Europe)
	if !found || zone.Name != "Europe" {
		t.Fatalf("Find(Europe) = %+v, %v", zone, found)
	}
	if _, found := Find("atlantis"); found {
		t.Fatal("Find(atlantis) found a zone")
	}
}

func TestZoneOfCountry(t *testing.T) {
	tests := []struct {
		country string
		want    ID
		found   bool
	}{
		{"France", Europe, true},
		{"the United States", NorthernAmerica, true},
		{"Mexico", LatinAmericaCaribbean, true},
		{"Russia", Europe, true},
		{"Turkey", MiddleEast, true},
		{"Iran", SouthAndCentralAsia, true},
		{"Georgia (country)", MiddleEast, true},
		{"the Soviet Union", Europe, true},
		{"Belgian Congo", Africa, true},
		{"Korea", EastAndSouthEastAsia, true},
		{"Antarctica", Antarctica, true},
		{"Atlantis", "", false},
	}
	for _, test := range tests {
		t.Run(test.country, func(t *testing.T) {
			got, found := ZoneOfCountry(test.country)
			if got != test.want || found != test.found {
				t.Fatalf("ZoneOfCountry(%q) = %q, %v; want %q, %v", test.country, got, found, test.want, test.found)
			}
		})
	}
}

func TestZoneOfRegion(t *testing.T) {
	tests := []struct {
		region, subregion string
		want              ID
		found             bool
	}{
		{"Americas", "Northern America", NorthernAmerica, true},
		{"Americas", "South America", LatinAmericaCaribbean, true},
		{"Americas", "Caribbean", LatinAmericaCaribbean, true},
		{"Europe", "Eastern Europe", Europe, true},
		{"Africa", "Northern Africa", Africa, true},
		{"Asia", "Western Asia", MiddleEast, true},
		{"Asia", "Central Asia", SouthAndCentralAsia, true},
		{"Asia", "South-Eastern Asia", EastAndSouthEastAsia, true},
		{"Oceania", "Melanesia", Oceania, true},
		{"Antarctica", "Antarctica", Antarctica, true},
		{"Americas", "Seven seas (open ocean)", "", false},
	}
	for _, test := range tests {
		t.Run(test.subregion, func(t *testing.T) {
			got, found := ZoneOfRegion(test.region, test.subregion)
			if got != test.want || found != test.found {
				t.Fatalf("ZoneOfRegion(%q, %q) = %q, %v; want %q, %v", test.region, test.subregion, got, found, test.want, test.found)
			}
		})
	}
}

func TestZoneOfTerritory(t *testing.T) {
	tests := []struct {
		name                string
		longitude, latitude float64
		want                ID
		found               bool
	}{
		{"French Guiana", -53.1, 3.9, LatinAmericaCaribbean, true},
		{"Martinique", -61.0, 14.6, LatinAmericaCaribbean, true},
		{"Réunion", 55.5, -21.1, Africa, true},
		{"Paris is not a territory", 2.35, 48.85, "", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, found := ZoneOfTerritory(test.longitude, test.latitude)
			if got != test.want || found != test.found {
				t.Fatalf("ZoneOfTerritory(%v, %v) = %q, %v; want %q, %v", test.longitude, test.latitude, got, found, test.want, test.found)
			}
		})
	}
}
