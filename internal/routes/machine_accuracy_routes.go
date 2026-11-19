package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// registerMachineAccuracyRoutes registers the machine accuracy list (server ADR 0017): everyone
// logged in reads (machine and part forms); only the superadmin adds, edits and deletes.
func registerMachineAccuracyRoutes(router *http.ServeMux, prefix string, accuracyHandler *handler.MachineAccuracyHandler) {
	router.HandleFunc("GET "+prefix+"/machine-accuracies", accuracyHandler.GetAll)
	router.Handle("POST "+prefix+"/machine-accuracies", superAdminOnly(accuracyHandler.Create))
	router.Handle("PUT "+prefix+"/machine-accuracies/{id}", superAdminOnly(accuracyHandler.Update))
	router.Handle("DELETE "+prefix+"/machine-accuracies/{id}", superAdminOnly(accuracyHandler.Delete))
}
