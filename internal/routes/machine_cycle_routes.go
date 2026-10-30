package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads the cycles of the machines they may access. Only the CNC collector
// writes cycles, so there is no POST, PUT or DELETE.
func registerMachineCycleRoutes(router *http.ServeMux, prefix string, cycleHandler *handler.MachineCycleHandler) {
	router.HandleFunc("GET "+prefix+"/machine-cycles", cycleHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/machine-cycles/{id}", cycleHandler.GetByID)
}
