//go:build integration

package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"

	"frise-chrono/back/internal/bc/imagery/infrastructure/commons/commonstest"
)

// testRedisDB keeps the integration tests apart from the development data.
const testRedisDB = 15

func startAPI(t *testing.T, fake *commonstest.Fake) *httptest.Server {
	t.Helper()
	address := os.Getenv("REDIS_ADDR")
	if address == "" {
		t.Fatal("REDIS_ADDR must point to the Redis of the stack")
	}
	redisClient := redis.NewClient(&redis.Options{Addr: address, DB: testRedisDB})
	t.Cleanup(func() { _ = redisClient.Close() })
	if err := redisClient.FlushDB(context.Background()).Err(); err != nil {
		t.Fatalf("Redis unavailable: %v", err)
	}

	config := Config{CommonsUserAgent: "FriseChrono/test", CORSAllowedOrigins: []string{"http://localhost:5173"}}
	handler, err := newHandler(config, fake.Start(t), redisClient)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func get(t *testing.T, server *httptest.Server, path string, headers map[string]string) (*http.Response, map[string]any) {
	t.Helper()
	request, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	for name, value := range headers {
		request.Header.Set(name, value)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = response.Body.Close() }()
	var body map[string]any
	_ = json.NewDecoder(response.Body).Decode(&body)
	return response, body
}

func europe1950() *commonstest.Fake {
	return &commonstest.Fake{
		Countries: map[string][]string{"1950s": {"France", "the Netherlands", "Japan"}},
		Files: map[string][]commonstest.File{
			"France":          {{Title: "File:Paris 1953.jpg", URL: "https://upload.wikimedia.org/paris.jpg", Date: "1953-04-03", License: "CC0"}},
			"the Netherlands": {{Title: "File:Rotterdam 1958.jpg", URL: "https://upload.wikimedia.org/rotterdam.jpg", Date: "circa 1958", License: "CC BY-SA 4.0"}},
		},
	}
}

func TestZonesAndGeometry(t *testing.T) {
	server := startAPI(t, &commonstest.Fake{})

	response, body := get(t, server, "/api/v1/zones", nil)
	if response.StatusCode != http.StatusOK || len(body["items"].([]any)) != 9 {
		t.Fatalf("GET /zones: %d %v", response.StatusCode, body)
	}

	response, body = get(t, server, "/api/v1/zones/geometry", nil)
	etag := response.Header.Get("ETag")
	if response.StatusCode != http.StatusOK || etag == "" || len(body["features"].([]any)) != 9 {
		t.Fatalf("GET /zones/geometry: %d, ETag %q", response.StatusCode, etag)
	}

	response, _ = get(t, server, "/api/v1/zones/geometry", map[string]string{"If-None-Match": etag})
	if response.StatusCode != http.StatusNotModified {
		t.Fatalf("conditional GET /zones/geometry: %d, want 304", response.StatusCode)
	}
}

func TestZonesWithImagesAreCachedInRedis(t *testing.T) {
	fake := europe1950()
	server := startAPI(t, fake)

	for range 2 {
		response, body := get(t, server, "/api/v1/decades/1950/zones", nil)
		zones := body["zones"].([]any)
		if response.StatusCode != http.StatusOK || len(zones) != 2 || zones[0] != "europe" || zones[1] != "east-southeast-asia" {
			t.Fatalf("GET /decades/1950/zones: %d %v", response.StatusCode, body)
		}
	}
	if fake.Requests() != 1 {
		t.Fatalf("Commons received %d requests, want 1 thanks to the cache", fake.Requests())
	}
}

func TestDrawImages(t *testing.T) {
	server := startAPI(t, europe1950())

	response, body := get(t, server, "/api/v1/images?zone=europe&decade=1950&limit=2", map[string]string{"Origin": "http://localhost:5173"})

	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /images: %d %v", response.StatusCode, body)
	}
	if response.Header.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatal("missing CORS header for the web origin")
	}
	items := body["items"].([]any)
	if len(items) != 2 {
		t.Fatalf("got %d images, want 2", len(items))
	}
	labels := map[string]bool{}
	for _, item := range items {
		image := item.(map[string]any)
		labels[image["time"].(map[string]any)["label"].(string)] = true
		if image["location"].(map[string]any)["zone"] != "europe" {
			t.Fatalf("image outside Europe: %v", image)
		}
	}
	if !labels["3 April 1953"] || !labels["around 1958"] {
		t.Fatalf("date labels = %v", labels)
	}
}

func TestDrawImagesErrors(t *testing.T) {
	tests := []struct {
		name       string
		fake       *commonstest.Fake
		path       string
		wantStatus int
	}{
		{name: "invalid decade", fake: europe1950(), path: "/api/v1/images?zone=europe&decade=1955", wantStatus: http.StatusBadRequest},
		{name: "missing zone", fake: europe1950(), path: "/api/v1/images?decade=1950", wantStatus: http.StatusBadRequest},
		{name: "source unavailable", fake: &commonstest.Fake{Failing: true}, path: "/api/v1/images?zone=europe&decade=1950", wantStatus: http.StatusServiceUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, body := get(t, startAPI(t, test.fake), test.path, nil)
			if response.StatusCode != test.wantStatus || response.Header.Get("Content-Type") != "application/problem+json" {
				t.Fatalf("status %d, content type %q, body %v", response.StatusCode, response.Header.Get("Content-Type"), body)
			}
			if body["status"] != float64(test.wantStatus) || body["title"] == "" {
				t.Fatalf("problem body = %v", body)
			}
		})
	}
}
