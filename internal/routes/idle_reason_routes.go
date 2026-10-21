package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads their company's idle reasons (the superadmin every customer's).
// Managers add, edit and delete their company's.
func registerIdleReasonRoutes(router *http.ServeMux, prefix string, idleReasonHandler *handler.IdleReasonHandler) {
	router.HandleFunc("GET "+prefix+"/idle-reasons", idleReasonHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/idle-reasons/{id}", idleReasonHandler.GetByID)
	router.Handle("POST "+prefix+"/idle-reasons", managersOnly(idleReasonHandler.Create))
	router.Handle("PUT "+prefix+"/idle-reasons/{id}", managersOnly(idleReasonHandler.Update))
	router.Handle("DELETE "+prefix+"/idle-reasons/{id}", managersOnly(idleReasonHandler.Delete))
}
