package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads their company's holidays (the superadmin every customer's).
// Managers add, edit and delete their company's.
func registerHolidayRoutes(router *http.ServeMux, prefix string, holidayHandler *handler.HolidayHandler) {
	router.HandleFunc("GET "+prefix+"/holidays", holidayHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/holidays/{id}", holidayHandler.GetByID)
	router.Handle("POST "+prefix+"/holidays", managersOnly(holidayHandler.Create))
	router.Handle("PUT "+prefix+"/holidays/{id}", managersOnly(holidayHandler.Update))
	router.Handle("DELETE "+prefix+"/holidays/{id}", managersOnly(holidayHandler.Delete))
}
