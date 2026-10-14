package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// User management is admin-only; /me is every user's own record, and every
// user may save their own theme.
func registerUserRoutes(router *http.ServeMux, prefix string, userHandler *handler.UserHandler) {
	router.Handle("GET "+prefix+"/users", adminOnly(userHandler.GetAll))
	router.Handle("GET "+prefix+"/users/{id}", adminOnly(userHandler.GetByID))
	router.Handle("POST "+prefix+"/users", adminOnly(userHandler.Create))
	router.Handle("PUT "+prefix+"/users/{id}", adminOnly(userHandler.Update))
	router.Handle("DELETE "+prefix+"/users/{id}", adminOnly(userHandler.Delete))
	router.HandleFunc("GET "+prefix+"/me", userHandler.Me)
	router.HandleFunc("PUT "+prefix+"/me/theme", userHandler.UpdateTheme)
}
