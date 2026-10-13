package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

func registerUserBranchRoutes(router *http.ServeMux, prefix string, userBranchHandler *handler.UserBranchHandler) {
	router.HandleFunc("GET "+prefix+"/user-branch", userBranchHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-branch/{id}", userBranchHandler.GetByID)
	router.HandleFunc("POST "+prefix+"/user-branch", userBranchHandler.Create)
	router.HandleFunc("PUT "+prefix+"/user-branch/{id}", userBranchHandler.Update)
	router.HandleFunc("DELETE "+prefix+"/user-branch/{id}", userBranchHandler.Delete)
}
