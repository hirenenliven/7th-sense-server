package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerMachineRoutes(router *http.ServeMux, prefix string, machineHandler *handler.MachineHandler) {
	router.HandleFunc("GET "+prefix+"/machines", machineHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/machines/{id}", machineHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/machines", machineHandler.Create)
	router.HandleFunc("PUT "+prefix+"/machines/{id}", machineHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/machines/{id}", machineHandler.Delete)
}
