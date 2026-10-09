package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerItemRoutes(router *http.ServeMux, prefix string, itemHandler *handler.ItemHandler) {
	router.HandleFunc("GET "+prefix+"/items", itemHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/items/{id}", itemHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/items", itemHandler.Create)
	router.HandleFunc("PUT "+prefix+"/items/{id}", itemHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/items/{id}", itemHandler.Delete)
}
