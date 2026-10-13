package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerBrandRoutes(router *http.ServeMux, prefix string, brandHandler *handler.BrandHandler) {
	router.HandleFunc("GET "+prefix+"/brands", brandHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/brands/{id}", brandHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/brands", brandHandler.Create)
	router.HandleFunc("PUT "+prefix+"/brands/{id}", brandHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/brands/{id}", brandHandler.Delete)
}
