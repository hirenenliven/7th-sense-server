// Package apidocs serves the API documentation: the OpenAPI 3 spec (openapi.yaml, generated from
// the swag annotations on the handlers and converted to OpenAPI 3 by tools/apidocs) and Swagger UI
// at /api-docs/ when API_DOCS=true.
//
// Regenerate the spec and the Postman collection after changing an endpoint:
//
//	go generate ./docs/api
package apidocs

//go:generate go tool swag init --quiet -g cmd/server/main.go -d ../.. -o .build --ot json --parseInternal
//go:generate go run ../../tools/apidocs

import (
	_ "embed"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:embed openapi.yaml
var spec []byte

// Path is where Swagger UI is served.
const Path = "/api-docs/"

// Handler serves openapi.yaml and Swagger UI under Path.
func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+Path+"openapi.yaml", func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/yaml")
		_, _ = writer.Write(spec)
	})
	mux.Handle("GET "+Path, httpSwagger.Handler(httpSwagger.URL(Path+"openapi.yaml")))
	return mux
}
