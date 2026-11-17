package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// registerOperatorTaskRoutes registers the maintenance task master list and the operator tasks
// (server ADR 0015). Everyone logged in reads (operator tasks: their company's); only the
// superadmin adds, edits and deletes the tasks of the list and the operator tasks.
func registerOperatorTaskRoutes(router *http.ServeMux, prefix string, taskHandler *handler.MaintenanceTaskHandler, operatorTaskHandler *handler.OperatorTaskHandler) {
	router.HandleFunc("GET "+prefix+"/maintenance-tasks", taskHandler.GetAll)
	router.Handle("POST "+prefix+"/maintenance-tasks", superAdminOnly(taskHandler.Create))
	router.Handle("PUT "+prefix+"/maintenance-tasks/{id}", superAdminOnly(taskHandler.Update))
	router.Handle("DELETE "+prefix+"/maintenance-tasks/{id}", superAdminOnly(taskHandler.Delete))
	router.HandleFunc("GET "+prefix+"/operator-tasks", operatorTaskHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/operator-tasks/{id}", operatorTaskHandler.GetByID)
	router.Handle("POST "+prefix+"/operator-tasks", superAdminOnly(operatorTaskHandler.Create))
	router.Handle("PUT "+prefix+"/operator-tasks/{id}", superAdminOnly(operatorTaskHandler.Update))
	router.Handle("DELETE "+prefix+"/operator-tasks/{id}", superAdminOnly(operatorTaskHandler.Delete))
}
