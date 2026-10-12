package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerBranchRoutes(router *http.ServeMux, prefix string, branchHandler *handler.BranchHandler) {
	router.HandleFunc("GET "+prefix+"/branches", branchHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/branches/{id}", branchHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/branches", branchHandler.Create)
	router.HandleFunc("PUT "+prefix+"/branches/{id}", branchHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/branches/{id}", branchHandler.Delete)
}
