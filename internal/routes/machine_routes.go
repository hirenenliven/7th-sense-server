package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in can read; only admins can change.
func registerMachineRoutes(router *http.ServeMux, prefix string, machineHandler *handler.MachineHandler) {
	router.HandleFunc("GET "+prefix+"/machines", machineHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/machines/{id}", machineHandler.GetByID)
	router.Handle("POST "+prefix+"/machines", adminOnly(machineHandler.Create))
	router.Handle("PUT "+prefix+"/machines/{id}", adminOnly(machineHandler.Update))
	router.Handle("DELETE "+prefix+"/machines/{id}", adminOnly(machineHandler.Delete))
}
