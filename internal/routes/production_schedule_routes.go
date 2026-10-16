package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads their company's production schedules (the superadmin every
// customer's). Managers add, edit and delete their company's schedules.
func registerProductionScheduleRoutes(router *http.ServeMux, prefix string, scheduleHandler *handler.ProductionScheduleHandler) {
	router.HandleFunc("GET "+prefix+"/production-schedules", scheduleHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/production-schedules/{id}", scheduleHandler.GetByID)
	router.Handle("POST "+prefix+"/production-schedules", managersOnly(scheduleHandler.Create))
	router.Handle("PUT "+prefix+"/production-schedules/{id}", managersOnly(scheduleHandler.Update))
	router.Handle("DELETE "+prefix+"/production-schedules/{id}", managersOnly(scheduleHandler.Delete))
}
