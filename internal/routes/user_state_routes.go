package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// A state is only added and removed, so there is no PUT.
// Countries and states are assigned customer-wise by the superadmin only; every user
// reads their company's rows (see targetUserID).
func registerUserStateRoutes(router *http.ServeMux, prefix string, userStateHandler *handler.UserStateHandler) {
	router.HandleFunc("GET "+prefix+"/user-state", userStateHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-state/{id}", userStateHandler.GetByID)
	router.Handle("POST "+prefix+"/user-state", superAdminOnly(userStateHandler.Create))
	router.Handle("DELETE "+prefix+"/user-state/{id}", superAdminOnly(userStateHandler.Delete))
}
