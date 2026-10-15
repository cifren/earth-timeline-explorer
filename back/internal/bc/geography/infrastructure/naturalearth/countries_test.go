package naturalearth

import "testing"

func TestIndexFindsACountryByAnyName(t *testing.T) {
	netherlands := Country{Name: "Netherlands", Names: []string{"Netherlands", "Netherlands", "", "Netherlands", "Kingdom of the Netherlands"}}
	index := NewIndex([]Country{netherlands})

	for _, name := range []string{"the Netherlands", "NETHERLANDS", "Kingdom of the Netherlands"} {
		if country, found := index.Find(name); !found || country.Name != "Netherlands" {
			t.Errorf("Find(%q) = %q, %v", name, country.Name, found)
		}
	}
	georgia := NewIndex([]Country{{Name: "Georgia", Names: []string{"Georgia"}}})
	if _, found := georgia.Find("Georgia (country)"); !found {
		t.Error("Georgia (country) did not match Georgia")
	}
	if _, found := index.Find(""); found {
		t.Error("an empty name matched a country")
	}
}
