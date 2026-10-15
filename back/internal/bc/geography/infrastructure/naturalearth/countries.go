// Package naturalearth downloads country borders from Natural Earth (public domain).
package naturalearth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CountriesURL is the Admin 0 countries layer, 1:50m scale, in GeoJSON.
const CountriesURL = "https://raw.githubusercontent.com/nvkelso/natural-earth-vector/master/geojson/ne_50m_admin_0_countries.geojson"

const downloadTimeout = 2 * time.Minute

// Country is a country border with its UN M49 region and sub-region.
type Country struct {
	Name string
	// Names are every English name of the country (administrative, short, long, formal).
	Names     []string
	Region    string
	Subregion string
	// LabelPoint is where a map places the country label, inside it: [longitude, latitude].
	LabelPoint [2]float64
	Geometry   json.RawMessage
}

type featureCollection struct {
	Features []struct {
		Properties struct {
			Name      string  `json:"ADMIN"`
			ShortName string  `json:"NAME"`
			LongName  string  `json:"NAME_LONG"`
			English   string  `json:"NAME_EN"`
			Formal    string  `json:"FORMAL_EN"`
			Region    string  `json:"REGION_UN"`
			Subregion string  `json:"SUBREGION"`
			LabelX    float64 `json:"LABEL_X"`
			LabelY    float64 `json:"LABEL_Y"`
		} `json:"properties"`
		Geometry json.RawMessage `json:"geometry"`
	} `json:"features"`
}

// FetchCountries downloads every country border from a Natural Earth countries layer.
func FetchCountries(ctx context.Context, url string) ([]Country, error) {
	ctx, cancel := context.WithTimeout(ctx, downloadTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("download Natural Earth countries: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download Natural Earth countries: HTTP %d", response.StatusCode)
	}

	var collection featureCollection
	if err := json.NewDecoder(response.Body).Decode(&collection); err != nil {
		return nil, fmt.Errorf("decode Natural Earth countries: %w", err)
	}
	countries := make([]Country, 0, len(collection.Features))
	for _, feature := range collection.Features {
		properties := feature.Properties
		countries = append(countries, Country{
			Name:       properties.Name,
			Names:      []string{properties.Name, properties.ShortName, properties.LongName, properties.English, properties.Formal},
			Region:     properties.Region,
			Subregion:  properties.Subregion,
			LabelPoint: [2]float64{properties.LabelX, properties.LabelY},
			Geometry:   feature.Geometry,
		})
	}
	return countries, nil
}

// Index finds countries by any of their names, ignoring case, a leading "the" and a
// trailing "(country)", so that "the Netherlands" finds "Netherlands" and "Georgia
// (country)" finds "Georgia".
type Index map[string]Country

// NewIndex indexes the countries by every one of their names.
func NewIndex(countries []Country) Index {
	index := Index{}
	for _, country := range countries {
		for _, name := range country.Names {
			key := indexKey(name)
			if _, taken := index[key]; key == "" || taken {
				continue
			}
			index[key] = country
		}
	}
	return index
}

// Find returns the country having the given name.
func (i Index) Find(name string) (Country, bool) {
	country, found := i[indexKey(name)]
	return country, found
}

func indexKey(name string) string {
	key := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(name)), "the ")
	return strings.TrimSuffix(key, " (country)")
}
