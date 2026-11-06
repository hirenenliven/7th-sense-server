package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Only the superadmin checks a machine's data collection: it shows every customer's connection
// details and makes the server probe the machine's address.
func registerMachineDiagnosticsRoutes(router *http.ServeMux, prefix string, diagnosticsHandler *handler.MachineDiagnosticsHandler) {
	router.Handle("GET "+prefix+"/machines/{id}/diagnostics", superAdminOnly(diagnosticsHandler.GetDiagnostics))
	router.Handle("GET "+prefix+"/machines/{id}/diagnostics/network", superAdminOnly(diagnosticsHandler.CheckNetwork))
}
