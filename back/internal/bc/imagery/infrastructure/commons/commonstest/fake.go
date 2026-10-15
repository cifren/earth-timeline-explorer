// Package commonstest serves a fake Wikimedia Commons API for tests: category listings
// and random file searches answered from in-memory data.
package commonstest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"
)

// File is a file page served by the fake.
type File struct {
	Title          string
	URL            string
	Date           string
	License        string
	Artist         string
	Credit         string
	HasCoordinates bool
}

// Fake answers the Commons queries used by the back.
type Fake struct {
	// Countries lists, per decade ("1950s"), the countries having a photograph category.
	Countries map[string][]string
	// Files lists, per country, the files returned in turn by random searches.
	Files map[string][]File
	// Failing makes every request answer an API error.
	Failing bool

	mu        sync.Mutex
	drawn     map[string]int
	userAgent string
	requests  int
}

var deepcat = regexp.MustCompile(`deepcat:"(\d{4}s) photographs of (.+)"`)

// Start serves the fake until the test ends and returns its endpoint.
func (f *Fake) Start(t *testing.T) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(server.Close)
	return server.URL
}

// UserAgent returns the User-Agent of the last request.
func (f *Fake) UserAgent() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.userAgent
}

// Requests returns the number of requests received.
func (f *Fake) Requests() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.requests
}

func (f *Fake) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.requests++
	f.userAgent = r.UserAgent()
	query := r.URL.Query()
	w.Header().Set("Content-Type", "application/json")

	if f.Failing {
		_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]string{"code": "maxlag", "info": "Waiting for a database server"}})
		return
	}
	if query.Get("list") == "categorymembers" {
		f.serveCategories(w, query.Get("cmtitle"))
		return
	}
	if query.Has("titles") {
		f.serveTitle(w, query.Get("titles"))
		return
	}
	f.serveSearch(w, query.Get("gsrsearch"))
}

func (f *Fake) serveTitle(w http.ResponseWriter, title string) {
	for _, files := range f.Files {
		for _, file := range files {
			if file.Title == title {
				_ = json.NewEncoder(w).Encode(map[string]any{"query": map[string]any{"pages": []any{page(file)}}})
				return
			}
		}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"query": map[string]any{"pages": []any{map[string]any{"title": title, "missing": true}}}})
}

func (f *Fake) serveCategories(w http.ResponseWriter, title string) {
	decade := strings.TrimSuffix(strings.TrimPrefix(title, "Category:"), " photographs by country")
	members := []map[string]string{{"title": "Category:Black and white photographs in the " + decade + " by country"}}
	for _, country := range f.Countries[decade] {
		members = append(members, map[string]string{"title": "Category:" + decade + " photographs of " + country})
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"query": map[string]any{"categorymembers": members}})
}

func (f *Fake) serveSearch(w http.ResponseWriter, search string) {
	match := deepcat.FindStringSubmatch(search)
	if match == nil || len(f.Files[match[2]]) == 0 {
		_ = json.NewEncoder(w).Encode(map[string]any{"batchcomplete": true})
		return
	}
	if f.drawn == nil {
		f.drawn = map[string]int{}
	}
	files := f.Files[match[2]]
	file := files[f.drawn[match[2]]%len(files)]
	f.drawn[match[2]]++
	_ = json.NewEncoder(w).Encode(map[string]any{"query": map[string]any{"pages": []any{page(file)}}})
}

func page(file File) map[string]any {
	result := map[string]any{
		"title": file.Title,
		"imageinfo": []map[string]any{{
			"url":            file.URL,
			"thumburl":       file.URL + "?width=1280",
			"responsiveUrls": map[string]string{"1.5": file.URL + "?width=1920", "2": file.URL + "?width=2560"},
			"descriptionurl": "https://commons.wikimedia.org/wiki/" + file.Title,
			"width":          4000,
			"height":         3000,
			"extmetadata": map[string]map[string]string{
				"DateTimeOriginal": {"value": file.Date},
				"LicenseShortName": {"value": file.License},
				"Artist":           {"value": file.Artist},
				"Credit":           {"value": file.Credit},
			},
		}},
	}
	if file.HasCoordinates {
		result["coordinates"] = []map[string]float64{{"lat": 48.85, "lon": 2.35}}
	}
	return result
}
