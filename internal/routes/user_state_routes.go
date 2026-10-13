package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// A state is only added and removed, so there is no PUT.
func registerUserStateRoutes(router *http.ServeMux, prefix string, userStateHandler *handler.UserStateHandler) {
	router.HandleFunc("GET "+prefix+"/user-state", userStateHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-state/{id}", userStateHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-state", userStateHandler.Create)
	router.HandleFunc("DELETE "+prefix+"/user-state/{id}", userStateHandler.Delete)
}
