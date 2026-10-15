package problem

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrite(t *testing.T) {
	recorder := httptest.NewRecorder()

	Write(recorder, New(http.StatusBadRequest, "The decade must be a multiple of 10."))

	if recorder.Code != http.StatusBadRequest || recorder.Header().Get("Content-Type") != ContentType {
		t.Fatalf("status %d, content type %q", recorder.Code, recorder.Header().Get("Content-Type"))
	}
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"type": "about:blank", "title": "Bad Request", "status": 400.0, "detail": "The decade must be a multiple of 10."}
	for key, value := range want {
		if body[key] != value {
			t.Errorf("%s = %v, want %v", key, body[key], value)
		}
	}
}
