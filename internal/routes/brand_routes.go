package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in can read; only admins can change.
func registerBrandRoutes(router *http.ServeMux, prefix string, brandHandler *handler.BrandHandler) {
	router.HandleFunc("GET "+prefix+"/brands", brandHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/brands/{id}", brandHandler.GetByID)
	router.Handle("POST "+prefix+"/brands", adminOnly(brandHandler.Create))
	router.Handle("PUT "+prefix+"/brands/{id}", adminOnly(brandHandler.Update))
	router.Handle("DELETE "+prefix+"/brands/{id}", adminOnly(brandHandler.Delete))
}
