package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// A country is only added and removed, so there is no PUT.
// Countries and states are assigned customer-wise by the superadmin only; every user
// reads their company's rows (see targetUserID).
func registerUserCountryRoutes(router *http.ServeMux, prefix string, userCountryHandler *handler.UserCountryHandler) {
	router.HandleFunc("GET "+prefix+"/user-country", userCountryHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/user-country/{id}", userCountryHandler.GetByID)
	router.Handle("POST "+prefix+"/user-country", superAdminOnly(userCountryHandler.Create))
	router.Handle("DELETE "+prefix+"/user-country/{id}", superAdminOnly(userCountryHandler.Delete))
}
