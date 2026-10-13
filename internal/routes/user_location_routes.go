package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserLocationRoutes(router *http.ServeMux, prefix string, userLocationHandler *handler.UserLocationHandler) {
	router.HandleFunc("GET "+prefix+"/user-location", userLocationHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-location/{id}", userLocationHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-location", userLocationHandler.Create)
	router.HandleFunc("PUT "+prefix+"/user-location/{id}", userLocationHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/user-location/{id}", userLocationHandler.Delete)
}
