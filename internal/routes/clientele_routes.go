package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads their company's clienteles (the superadmin every customer's).
// Managers add, edit and delete their company's.
func registerClienteleRoutes(router *http.ServeMux, prefix string, clienteleHandler *handler.ClienteleHandler) {
	router.HandleFunc("GET "+prefix+"/clienteles", clienteleHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/clienteles/{id}", clienteleHandler.GetByID)
	router.Handle("POST "+prefix+"/clienteles", managersOnly(clienteleHandler.Create))
	router.Handle("PUT "+prefix+"/clienteles/{id}", managersOnly(clienteleHandler.Update))
	router.Handle("DELETE "+prefix+"/clienteles/{id}", managersOnly(clienteleHandler.Delete))
}
