package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
)

// registerAuthRoutes registers login and logout, which work without a session,
// and the current-user route, which needs one.
func registerAuthRoutes(router *http.ServeMux, prefix string, authHandler *handler.AuthHandler, authenticate middleware.Middleware) {
	router.HandleFunc("POST "+prefix+"/auth/login", authHandler.Login)
	router.HandleFunc("POST "+prefix+"/auth/logout", authHandler.Logout)
	router.Handle("GET "+prefix+"/auth/me", authenticate(http.HandlerFunc(authHandler.Me)))
}
