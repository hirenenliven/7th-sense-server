package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// A state is only added and removed, so there is no PUT.
// Every user reads their company's rows; only managers change them (see targetUserID).
func registerUserStateRoutes(router *http.ServeMux, prefix string, userStateHandler *handler.UserStateHandler) {
	router.HandleFunc("GET "+prefix+"/user-state", userStateHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-state/{id}", userStateHandler.GetByID)
	router.Handle("POST "+prefix+"/user-state", managersOnly(userStateHandler.Create))
	router.Handle("DELETE "+prefix+"/user-state/{id}", managersOnly(userStateHandler.Delete))
}
