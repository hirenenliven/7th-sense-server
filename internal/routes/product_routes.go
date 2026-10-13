package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in can read; only admins can change.
func registerProductRoutes(router *http.ServeMux, prefix string, productHandler *handler.ProductHandler) {
	router.HandleFunc("GET "+prefix+"/products", productHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/products/{id}", productHandler.GetByID)
	router.Handle("POST "+prefix+"/products", adminOnly(productHandler.Create))
	router.Handle("PUT "+prefix+"/products/{id}", adminOnly(productHandler.Update))
	router.Handle("DELETE "+prefix+"/products/{id}", adminOnly(productHandler.Delete))
}
