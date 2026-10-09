package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserRoutes(mux *http.ServeMux, h *handler.UserHandler) {
	mux.HandleFunc("GET "+apiPrefix+"/users", h.GetAll)
	mux.HandleFunc("GET "+apiPrefix+"/users/{id}", h.GetByID)
	mux.HandleFunc("POST "+apiPrefix+"/users", h.Create)
	mux.HandleFunc("PUT "+apiPrefix+"/users/{id}", h.Update)
	mux.HandleFunc("DELETE "+apiPrefix+"/users/{id}", h.Delete)
}
