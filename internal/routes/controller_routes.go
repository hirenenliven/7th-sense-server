package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerControllerRoutes(router *http.ServeMux, prefix string, controllerHandler *handler.ControllerHandler) {
	router.HandleFunc("GET "+prefix+"/controllers", controllerHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/controllers/{id}", controllerHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/controllers", controllerHandler.Create)
	router.HandleFunc("PUT "+prefix+"/controllers/{id}", controllerHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/controllers/{id}", controllerHandler.Delete)
}
