package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/response"
)

const apiPrefix = "/api"

// Handlers groups every handler the router needs.
type Handlers struct {
	User *handler.UserHandler
	Item *handler.ItemHandler
}

// New builds the main router with all routes and global middleware.
func New(handlers Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.SendResponse(w, r, response.Success, nil)
	})

	// Register API versions
	registerV1Routes(mux, handlers)

	return middleware.Chain(mux, middleware.Logger, middleware.Recoverer)
}

// registerV1Routes registers all v1 API routes.
func registerV1Routes(mux *http.ServeMux, handlers Handlers) {
	v1Prefix := apiPrefix + "/v1"

	if handlers.User != nil {
		registerUserRoutes(mux, v1Prefix, handlers.User)
	}
	if handlers.Item != nil {
		registerItemRoutes(mux, v1Prefix, handlers.Item)
	}
}
