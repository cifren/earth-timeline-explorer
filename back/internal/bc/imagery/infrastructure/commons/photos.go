package commons

import (
	"context"
	"fmt"
	"net/url"

	"frise-chrono/back/internal/bc/imagery/domain"
	"frise-chrono/back/internal/core/timespan"
)

const thumbnailWidth = "1280"

// PhotoSource draws random photographs from the Commons category tree of a country
// and a decade.
type PhotoSource struct {
	client *Client
}

// NewPhotoSource returns a photo source using the given client.
func NewPhotoSource(client *Client) *PhotoSource {
	return &PhotoSource{client: client}
}

type searchResult struct {
	Query struct {
		Pages []page `json:"pages"`
	} `json:"query"`
}

type page struct {
	Title       string `json:"title"`
	Missing     bool   `json:"missing"`
	Coordinates []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coordinates"`
	ImageInfo []imageInfo `json:"imageinfo"`
}

type imageInfo struct {
	URL            string                   `json:"url"`
	DescriptionURL string                   `json:"descriptionurl"`
	ThumbURL       string                   `json:"thumburl"`
	Width          int                      `json:"width"`
	Height         int                      `json:"height"`
	ExtMetadata    map[string]metadataField `json:"extmetadata"`
}

type metadataField struct {
	Value string `json:"value"`
}

// RandomPhoto draws one photograph. It reports false when the draw returned no usable
// file (empty category tree, deleted file).
func (s *PhotoSource) RandomPhoto(ctx context.Context, country string, decade timespan.Decade) (domain.Image, bool, error) {
	params := url.Values{
		"generator":           {"search"},
		"gsrnamespace":        {"6"},
		"gsrsearch":           {fmt.Sprintf(`filetype:bitmap deepcat:"%ds photographs of %s"`, decade, country)},
		"gsrsort":             {"random"},
		"gsrlimit":            {"1"},
		"prop":                {"imageinfo|coordinates"},
		"iiprop":              {"url|size|extmetadata"},
		"iiurlwidth":          {thumbnailWidth},
		"iiextmetadatafilter": {"DateTimeOriginal|LicenseShortName|Artist|Credit"},
	}
	var result searchResult
	if err := s.client.query(ctx, params, &result); err != nil {
		return domain.Image{}, false, err
	}
	if len(result.Query.Pages) == 0 {
		return domain.Image{}, false, nil
	}
	image, usable := toImage(result.Query.Pages[0], country, decade)
	return image, usable, nil
}
