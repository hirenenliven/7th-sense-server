package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
)

// registerAuthRoutes registers login and logout, which work without a session,
// and the current-user and log-in-as routes, which need one (log in as: superadmin only).
func registerAuthRoutes(router *http.ServeMux, prefix string, authHandler *handler.AuthHandler, authenticate middleware.Middleware) {
	router.HandleFunc("POST "+prefix+"/auth/login", authHandler.Login)
	router.HandleFunc("POST "+prefix+"/auth/logout", authHandler.Logout)
	router.Handle("GET "+prefix+"/auth/me", authenticate(http.HandlerFunc(authHandler.Me)))
	router.Handle("POST "+prefix+"/auth/login-as", authenticate(superAdminOnly(authHandler.LoginAs)))
}
