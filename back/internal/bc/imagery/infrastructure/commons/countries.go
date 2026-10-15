package commons

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"frise-chrono/back/internal/core/timespan"
	"frise-chrono/back/internal/shared/cache"
)

var countriesCachePolicy = cache.Policy{FreshFor: 7 * 24 * time.Hour, KeepFor: 90 * 24 * time.Hour}

// CountryCatalog lists the countries having a photograph category for a decade.
type CountryCatalog struct {
	client *Client
	cache  *cache.Store
}

// NewCountryCatalog returns a catalog reading Commons through the cache.
func NewCountryCatalog(client *Client, store *cache.Store) *CountryCatalog {
	return &CountryCatalog{client: client, cache: store}
}

type categoryMembers struct {
	Continue struct {
		CmContinue string `json:"cmcontinue"`
	} `json:"continue"`
	Query struct {
		CategoryMembers []struct {
			Title string `json:"title"`
		} `json:"categorymembers"`
	} `json:"query"`
}

// Countries returns the country names of "Category:<decade>s photographs of <country>".
func (c *CountryCatalog) Countries(ctx context.Context, decade timespan.Decade) ([]string, error) {
	key := fmt.Sprintf("imagery:commons:countries:%d", decade)
	return cache.GetOrRefresh(ctx, c.cache, key, countriesCachePolicy, func(ctx context.Context) ([]string, error) {
		return c.fetch(ctx, decade)
	})
}

func (c *CountryCatalog) fetch(ctx context.Context, decade timespan.Decade) ([]string, error) {
	prefix := fmt.Sprintf("Category:%ds photographs of ", decade)
	params := url.Values{
		"list":    {"categorymembers"},
		"cmtitle": {fmt.Sprintf("Category:%ds photographs by country", decade)},
		"cmtype":  {"subcat"},
		"cmlimit": {"max"},
	}
	var countries []string
	for {
		var page categoryMembers
		if err := c.client.query(ctx, params, &page); err != nil {
			return nil, err
		}
		for _, member := range page.Query.CategoryMembers {
			country, isCountry := strings.CutPrefix(member.Title, prefix)
			if isCountry {
				countries = append(countries, country)
			}
		}
		if page.Continue.CmContinue == "" {
			return countries, nil
		}
		params.Set("cmcontinue", page.Continue.CmContinue)
	}
}
