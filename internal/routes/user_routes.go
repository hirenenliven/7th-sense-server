package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserRoutes(mux *http.ServeMux, prefix string, h *handler.UserHandler) {
	mux.HandleFunc("GET "+prefix+"/users", h.GetAll)
	mux.HandleFunc("GET "+prefix+"/users/{id}", h.GetByID)
	mux.HandleFunc("POST "+prefix+"/users", h.Create)
	mux.HandleFunc("PUT "+prefix+"/users/{id}", h.Update)
	mux.HandleFunc("DELETE "+prefix+"/users/{id}", h.Delete)
}
