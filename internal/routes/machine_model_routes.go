package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in can read (the machine form lists them); only the superadmin can change.
func registerMachineModelRoutes(router *http.ServeMux, prefix string, machineModelHandler *handler.MachineModelHandler) {
	router.HandleFunc("GET "+prefix+"/models", machineModelHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/models/{id}", machineModelHandler.GetByID)
	router.Handle("POST "+prefix+"/models", superAdminOnly(machineModelHandler.Create))
	router.Handle("PUT "+prefix+"/models/{id}", superAdminOnly(machineModelHandler.Update))
	router.Handle("DELETE "+prefix+"/models/{id}", superAdminOnly(machineModelHandler.Delete))
}
