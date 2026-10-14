package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Branches are shared by every company: everyone logged in reads them; only the superadmin changes them.
func registerBranchRoutes(router *http.ServeMux, prefix string, branchHandler *handler.BranchHandler) {
	router.HandleFunc("GET "+prefix+"/branches", branchHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/branches/{id}", branchHandler.GetByID)
	router.Handle("POST "+prefix+"/branches", superAdminOnly(branchHandler.Create))
	router.Handle("PUT "+prefix+"/branches/{id}", superAdminOnly(branchHandler.Update))
	router.Handle("DELETE "+prefix+"/branches/{id}", superAdminOnly(branchHandler.Delete))
}
