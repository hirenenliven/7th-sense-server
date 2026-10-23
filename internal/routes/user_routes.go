package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// User management is for the superadmin, customers and admins, each limited to
// the users they may manage (service.UserService); /me is every user's own
// record, and every user may save their own theme and language.
func registerUserRoutes(router *http.ServeMux, prefix string, userHandler *handler.UserHandler) {
	router.Handle("GET "+prefix+"/users", managersOnly(userHandler.GetAll))
	router.Handle("GET "+prefix+"/users/{id}", managersOnly(userHandler.GetByID))
	router.Handle("POST "+prefix+"/users", managersOnly(userHandler.Create))
	router.Handle("PUT "+prefix+"/users/{id}", managersOnly(userHandler.Update))
	router.Handle("DELETE "+prefix+"/users/{id}", managersOnly(userHandler.Delete))
	router.HandleFunc("GET "+prefix+"/me", userHandler.Me)
	router.HandleFunc("PUT "+prefix+"/me/theme", userHandler.UpdateTheme)
	router.HandleFunc("PUT "+prefix+"/me/language", userHandler.UpdateLanguage)
}
