package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apidocs "7th-sense-server/docs/api"
	"7th-sense-server/internal/routes"
)

// TestAPIDocs checks Swagger UI and the embedded OpenAPI spec are served at /api-docs/ when
// enabled, and that /api-docs/ does not exist when they are off.
func TestAPIDocs(t *testing.T) {
	router := routes.New(routes.Handlers{APIDocs: apidocs.Handler()}, asTestUser)
	for path, want := range map[string]string{
		"/api-docs/index.html":   "swagger-ui",
		"/api-docs/openapi.yaml": "openapi: 3.0.3",
	} {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), want) {
			t.Errorf("GET %s: status %d, want 200 with %q", path, recorder.Code, want)
		}
	}

	off := routes.New(routes.Handlers{}, asTestUser)
	recorder := httptest.NewRecorder()
	off.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api-docs/index.html", nil))
	if recorder.Code != http.StatusNotFound {
		t.Errorf("API docs off: status %d, want 404", recorder.Code)
	}
}
