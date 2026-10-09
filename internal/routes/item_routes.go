package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerItemRoutes(mux *http.ServeMux, h *handler.ItemHandler) {
	mux.HandleFunc("GET "+apiPrefix+"/items", h.GetAll)
	mux.HandleFunc("GET "+apiPrefix+"/items/{id}", h.GetByID)
	mux.HandleFunc("POST "+apiPrefix+"/items", h.Create)
	mux.HandleFunc("PUT "+apiPrefix+"/items/{id}", h.Update)
	mux.HandleFunc("DELETE "+apiPrefix+"/items/{id}", h.Delete)
}
