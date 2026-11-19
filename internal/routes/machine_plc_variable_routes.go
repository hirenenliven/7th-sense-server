package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Only the superadmin sets a machine's PLC node ids: they are installation details of the
// controller, like the connection checked by the diagnostics.
func registerMachinePlcVariableRoutes(router *http.ServeMux, prefix string, plcVariableHandler *handler.MachinePlcVariableHandler) {
	router.Handle("GET "+prefix+"/machines/{id}/plc-variables", superAdminOnly(plcVariableHandler.GetAll))
	router.Handle("DELETE "+prefix+"/machines/{id}/plc-variables", superAdminOnly(plcVariableHandler.ResetAll))
	router.Handle("PUT "+prefix+"/machines/{id}/plc-variables/{variableId}", superAdminOnly(plcVariableHandler.Update))
	router.Handle("DELETE "+prefix+"/machines/{id}/plc-variables/{variableId}", superAdminOnly(plcVariableHandler.Reset))
}
