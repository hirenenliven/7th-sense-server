package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in can read; only admins can change.
func registerCategoryRoutes(router *http.ServeMux, prefix string, categoryHandler *handler.CategoryHandler) {
	router.HandleFunc("GET "+prefix+"/categories", categoryHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/categories/{id}", categoryHandler.GetByID)
	router.Handle("POST "+prefix+"/categories", superAdminOnly(categoryHandler.Create))
	router.Handle("PUT "+prefix+"/categories/{id}", superAdminOnly(categoryHandler.Update))
	router.Handle("DELETE "+prefix+"/categories/{id}", superAdminOnly(categoryHandler.Delete))
}
