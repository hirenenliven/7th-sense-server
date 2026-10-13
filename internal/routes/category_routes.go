package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerCategoryRoutes(router *http.ServeMux, prefix string, categoryHandler *handler.CategoryHandler) {
	router.HandleFunc("GET "+prefix+"/categories", categoryHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/categories/{id}", categoryHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/categories", categoryHandler.Create)
	router.HandleFunc("PUT "+prefix+"/categories/{id}", categoryHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/categories/{id}", categoryHandler.Delete)
}
