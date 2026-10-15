package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads their company's breaks (the superadmin every customer's).
// Managers add, edit and delete their company's.
func registerBreakRoutes(router *http.ServeMux, prefix string, breakHandler *handler.BreakHandler) {
	router.HandleFunc("GET "+prefix+"/breaks", breakHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/breaks/{id}", breakHandler.GetByID)
	router.Handle("POST "+prefix+"/breaks", managersOnly(breakHandler.Create))
	router.Handle("PUT "+prefix+"/breaks/{id}", managersOnly(breakHandler.Update))
	router.Handle("DELETE "+prefix+"/breaks/{id}", managersOnly(breakHandler.Delete))
}
