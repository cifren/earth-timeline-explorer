package domain

import "testing"

func TestImageIDRoundTrip(t *testing.T) {
	image := Image{
		Location:   Location{Country: "France"},
		Provenance: Provenance{Source: "wikimedia_commons", SourceID: "File:Market day: Paris.jpg"},
	}

	reference, valid := ParseID(image.ID())

	want := Reference{Source: "wikimedia_commons", Country: "France", SourceID: "File:Market day: Paris.jpg"}
	if !valid || reference != want {
		t.Fatalf("ParseID(%q) = %+v, %v, want %+v", image.ID(), reference, valid, want)
	}
	for _, malformed := range []string{"", "wikimedia_commons", "wikimedia_commons:France", "wikimedia_commons::File:x.jpg"} {
		if _, valid := ParseID(malformed); valid {
			t.Errorf("ParseID(%q) is valid", malformed)
		}
	}
}
