// Package commons is the gateway to the Wikimedia Commons API: it fetches records and
// maps them to the imagery domain (anti-corruption layer).
package commons

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"frise-chrono/back/internal/bc/imagery/domain"
)

const (
	endpoint    = "https://commons.wikimedia.org/w/api.php"
	callTimeout = 3 * time.Second
	maxLag      = "5"
)

// Client calls the Commons API following the Wikimedia API etiquette: identifying
// User-Agent, maxlag, one request at a time, bounded duration.
type Client struct {
	httpClient *http.Client
	userAgent  string
	serial     chan struct{}
}

// NewClient returns a client sending the given User-Agent.
func NewClient(userAgent string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: callTimeout},
		userAgent:  userAgent,
		serial:     make(chan struct{}, 1),
	}
}

type apiError struct {
	Error *struct {
		Code string `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}

// query runs an API query and decodes its JSON answer into out.
func (c *Client) query(ctx context.Context, params url.Values, out any) error {
	select {
	case c.serial <- struct{}{}:
		defer func() { <-c.serial }()
	case <-ctx.Done():
		return fmt.Errorf("%w: %w", domain.ErrSourceUnavailable, ctx.Err())
	}

	ctx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("formatversion", "2")
	params.Set("maxlag", maxLag)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", c.userAgent)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrSourceUnavailable, err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: HTTP %d", domain.ErrSourceUnavailable, response.StatusCode)
	}

	var body json.RawMessage
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return fmt.Errorf("%w: %w", domain.ErrSourceUnavailable, err)
	}
	var failure apiError
	if err := json.Unmarshal(body, &failure); err == nil && failure.Error != nil {
		return fmt.Errorf("%w: %s: %s", domain.ErrSourceUnavailable, failure.Error.Code, failure.Error.Info)
	}
	return json.Unmarshal(body, out)
}
