package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads their company's shifts (the superadmin every customer's).
// Managers add, edit and delete their company's; breaks are saved with their shift.
func registerShiftRoutes(router *http.ServeMux, prefix string, shiftHandler *handler.ShiftHandler) {
	router.HandleFunc("GET "+prefix+"/shifts", shiftHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/shifts/{id}", shiftHandler.GetByID)
	router.Handle("POST "+prefix+"/shifts", managersOnly(shiftHandler.Create))
	router.Handle("PUT "+prefix+"/shifts/{id}", managersOnly(shiftHandler.Update))
	router.Handle("DELETE "+prefix+"/shifts/{id}", managersOnly(shiftHandler.Delete))
}
