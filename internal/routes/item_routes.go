package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerItemRoutes(mux *http.ServeMux, prefix string, h *handler.ItemHandler) {
	mux.HandleFunc("GET "+prefix+"/items", h.GetAll)
	mux.HandleFunc("GET "+prefix+"/items/{id}", h.GetByID)
	mux.HandleFunc("POST "+prefix+"/items", h.Create)
	mux.HandleFunc("PUT "+prefix+"/items/{id}", h.Update)
	mux.HandleFunc("DELETE "+prefix+"/items/{id}", h.Delete)
}
