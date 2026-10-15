// Package naturalearth downloads country borders from Natural Earth (public domain).
package naturalearth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CountriesURL is the Admin 0 countries layer, 1:50m scale, in GeoJSON.
const CountriesURL = "https://raw.githubusercontent.com/nvkelso/natural-earth-vector/master/geojson/ne_50m_admin_0_countries.geojson"

const downloadTimeout = 2 * time.Minute

// Country is a country border with its UN M49 region and sub-region.
type Country struct {
	Name      string
	Region    string
	Subregion string
	Geometry  json.RawMessage
}

type featureCollection struct {
	Features []struct {
		Properties struct {
			Name      string `json:"ADMIN"`
			Region    string `json:"REGION_UN"`
			Subregion string `json:"SUBREGION"`
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
		countries = append(countries, Country{
			Name:      feature.Properties.Name,
			Region:    feature.Properties.Region,
			Subregion: feature.Properties.Subregion,
			Geometry:  feature.Geometry,
		})
	}
	return countries, nil
}
