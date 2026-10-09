package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserRoutes(router *http.ServeMux, prefix string, userHandler *handler.UserHandler) {
	router.HandleFunc("GET "+prefix+"/users", userHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/users/{id}", userHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/users", userHandler.Create)
	router.HandleFunc("PUT "+prefix+"/users/{id}", userHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/users/{id}", userHandler.Delete)
}
