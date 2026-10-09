package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/response"
)

const apiPrefix = "/api/v1"

// Handlers groups every handler the router needs.
type Handlers struct {
	User *handler.UserHandler
	Item *handler.ItemHandler
}

// New builds the main router with all routes and global middleware.
func New(h Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.Send(w, http.StatusOK, "OK")
	})

	registerUserRoutes(mux, h.User)
	registerItemRoutes(mux, h.Item)

	return middleware.Chain(mux, middleware.Logger, middleware.Recoverer)
}
