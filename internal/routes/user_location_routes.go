package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Every user reads their company's rows; only managers change them (see targetUserID).
func registerUserLocationRoutes(router *http.ServeMux, prefix string, userLocationHandler *handler.UserLocationHandler) {
	router.HandleFunc("GET "+prefix+"/user-location", userLocationHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-location/{id}", userLocationHandler.GetByID)
	router.Handle("POST "+prefix+"/user-location", managersOnly(userLocationHandler.Create))
	router.Handle("PUT "+prefix+"/user-location/{id}", managersOnly(userLocationHandler.Update))
	router.Handle("DELETE "+prefix+"/user-location/{id}", managersOnly(userLocationHandler.Delete))
}
