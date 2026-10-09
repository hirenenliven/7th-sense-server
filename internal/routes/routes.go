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
	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(writer http.ResponseWriter, request *http.Request) {
		response.SendResponse(writer, request, response.Success, nil)
	})

	// Register API versions
	registerV1Routes(router, handlers)

	return middleware.Chain(router, middleware.Logger, middleware.Recoverer)
}

// registerV1Routes registers all v1 API routes.
func registerV1Routes(router *http.ServeMux, handlers Handlers) {
	v1Prefix := apiPrefix + "/v1"

	if handlers.User != nil {
		registerUserRoutes(router, v1Prefix, handlers.User)
	}
	if handlers.Item != nil {
		registerItemRoutes(router, v1Prefix, handlers.Item)
	}
}
