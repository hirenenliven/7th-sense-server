package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerProductRoutes(router *http.ServeMux, prefix string, productHandler *handler.ProductHandler) {
	router.HandleFunc("GET "+prefix+"/products", productHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/products/{id}", productHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/products", productHandler.Create)
	router.HandleFunc("PUT "+prefix+"/products/{id}", productHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/products/{id}", productHandler.Delete)
}
