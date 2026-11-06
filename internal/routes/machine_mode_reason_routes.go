package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Whoever may read a machine records why it was down or idle, operators included, so no role guard:
// service.MachineModeReasonService checks the machine on every request.
func registerMachineModeReasonRoutes(router *http.ServeMux, prefix string, reasonHandler *handler.MachineModeReasonHandler) {
	router.HandleFunc("GET "+prefix+"/machine-mode-reasons", reasonHandler.GetByPeriod)
	router.HandleFunc("POST "+prefix+"/machine-mode-reasons", reasonHandler.Create)
	router.HandleFunc("PUT "+prefix+"/machine-mode-reasons/{id}", reasonHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/machine-mode-reasons/{id}", reasonHandler.Delete)
}
