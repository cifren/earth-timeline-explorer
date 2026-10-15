// Package problem builds and writes Problem Details responses (RFC 9457).
package problem

import (
	"encoding/json"
	"net/http"

	"frise-chrono/back/internal/shared/apigen"
)

// ContentType is the media type of Problem Details.
const ContentType = "application/problem+json"

// New returns a problem without a specific type, titled after the HTTP status.
func New(status int, detail string) apigen.Problem {
	return apigen.Problem{
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Detail: &detail,
	}
}

// Write sends a problem as the HTTP response.
func Write(w http.ResponseWriter, p apigen.Problem) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(p.Status)
	_ = json.NewEncoder(w).Encode(p)
}
