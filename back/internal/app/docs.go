package app

import (
	"encoding/json"
	"net/http"

	"frise-chrono/back/internal/shared/apigen"
	"frise-chrono/back/internal/shared/problem"
)

// documentationPage renders the OpenAPI document with Redoc.
const documentationPage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Frise Chrono API</title>
</head>
<body>
  <redoc spec-url="/api/openapi.json"></redoc>
  <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>`

// openAPIDocument serves the contract embedded in the generated code.
func openAPIDocument(w http.ResponseWriter, _ *http.Request) {
	spec, err := apigen.GetSwagger()
	if err != nil {
		problem.Write(w, problem.New(http.StatusInternalServerError, "The API document is unavailable."))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(spec)
}

func documentation(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(documentationPage))
}
