package commons

import (
	"context"
	"fmt"
	"net/url"
	"strings"

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
	ResponsiveURLs map[string]string        `json:"responsiveUrls"`
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
	params := fileParams()
	params.Set("generator", "search")
	params.Set("gsrnamespace", "6")
	params.Set("gsrsearch", fmt.Sprintf(`filetype:bitmap deepcat:"%ds photographs of %s"`, decade, country))
	params.Set("gsrsort", "random")
	params.Set("gsrlimit", "1")
	return s.firstPhoto(ctx, params, country, decade)
}

// Photo fetches the photograph a reference points to, drawn earlier for the decade. It
// reports false when the reference is not a Commons file or the file no longer exists.
func (s *PhotoSource) Photo(ctx context.Context, reference domain.Reference, decade timespan.Decade) (domain.Image, bool, error) {
	if reference.Source != sourceName || !strings.HasPrefix(reference.SourceID, "File:") {
		return domain.Image{}, false, nil
	}
	params := fileParams()
	params.Set("titles", reference.SourceID)
	return s.firstPhoto(ctx, params, reference.Country, decade)
}

// fileParams asks for what an image needs from a file page.
func fileParams() url.Values {
	return url.Values{
		"prop":                {"imageinfo|coordinates"},
		"iiprop":              {"url|size|extmetadata"},
		"iiurlwidth":          {thumbnailWidth},
		"iiextmetadatafilter": {"DateTimeOriginal|LicenseShortName|Artist|Credit"},
	}
}

func (s *PhotoSource) firstPhoto(ctx context.Context, params url.Values, country string, decade timespan.Decade) (domain.Image, bool, error) {
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
