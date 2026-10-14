package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// User management is for the superadmin, customers and admins, each limited to
// the users they may manage (service.UserService); /me is every user's own
// record, and every user may save their own theme.
func registerUserRoutes(router *http.ServeMux, prefix string, userHandler *handler.UserHandler) {
	router.Handle("GET "+prefix+"/users", userManagersOnly(userHandler.GetAll))
	router.Handle("GET "+prefix+"/users/{id}", userManagersOnly(userHandler.GetByID))
	router.Handle("POST "+prefix+"/users", userManagersOnly(userHandler.Create))
	router.Handle("PUT "+prefix+"/users/{id}", userManagersOnly(userHandler.Update))
	router.Handle("DELETE "+prefix+"/users/{id}", userManagersOnly(userHandler.Delete))
	router.HandleFunc("GET "+prefix+"/me", userHandler.Me)
	router.HandleFunc("PUT "+prefix+"/me/theme", userHandler.UpdateTheme)
}
