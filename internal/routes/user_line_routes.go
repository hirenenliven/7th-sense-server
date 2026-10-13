package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserLineRoutes(router *http.ServeMux, prefix string, userLineHandler *handler.UserLineHandler) {
	router.HandleFunc("GET "+prefix+"/user-line", userLineHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-line/{id}", userLineHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-line", userLineHandler.Create)
	router.HandleFunc("PUT "+prefix+"/user-line/{id}", userLineHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/user-line/{id}", userLineHandler.Delete)
}
