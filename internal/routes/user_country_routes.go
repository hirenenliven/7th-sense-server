package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// A country is only added and removed, so there is no PUT.
func registerUserCountryRoutes(router *http.ServeMux, prefix string, userCountryHandler *handler.UserCountryHandler) {
	router.HandleFunc("GET "+prefix+"/user-country", userCountryHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-country/{id}", userCountryHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-country", userCountryHandler.Create)
	router.HandleFunc("DELETE "+prefix+"/user-country/{id}", userCountryHandler.Delete)
}
