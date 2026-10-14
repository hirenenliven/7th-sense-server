package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads the machines they may access; only the superadmin changes them.
func registerMachineRoutes(router *http.ServeMux, prefix string, machineHandler *handler.MachineHandler) {
	router.HandleFunc("GET "+prefix+"/machines", machineHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/machines/{id}", machineHandler.GetByID)
	router.Handle("POST "+prefix+"/machines", superAdminOnly(machineHandler.Create))
	router.Handle("PUT "+prefix+"/machines/{id}", superAdminOnly(machineHandler.Update))
	router.Handle("DELETE "+prefix+"/machines/{id}", superAdminOnly(machineHandler.Delete))
}
