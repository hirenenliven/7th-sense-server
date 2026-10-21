package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads their company's downtime reasons (the superadmin every customer's).
// Managers add, edit and delete their company's.
func registerDowntimeReasonRoutes(router *http.ServeMux, prefix string, downtimeReasonHandler *handler.DowntimeReasonHandler) {
	router.HandleFunc("GET "+prefix+"/downtime-reasons", downtimeReasonHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/downtime-reasons/{id}", downtimeReasonHandler.GetByID)
	router.Handle("POST "+prefix+"/downtime-reasons", managersOnly(downtimeReasonHandler.Create))
	router.Handle("PUT "+prefix+"/downtime-reasons/{id}", managersOnly(downtimeReasonHandler.Update))
	router.Handle("DELETE "+prefix+"/downtime-reasons/{id}", managersOnly(downtimeReasonHandler.Delete))
}
