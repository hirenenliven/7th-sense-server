package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Managers read their company's licence keys (the superadmin every customer's). Only the
// superadmin generates, saves, edits and deletes them.
func registerLicenceKeyRoutes(router *http.ServeMux, prefix string, licenceKeyHandler *handler.LicenceKeyHandler) {
	router.Handle("GET "+prefix+"/licence-keys", managersOnly(licenceKeyHandler.GetAll))
	router.Handle("GET "+prefix+"/licence-keys/{id}", managersOnly(licenceKeyHandler.GetByID))
	router.Handle("POST "+prefix+"/licence-keys/generate", superAdminOnly(licenceKeyHandler.Generate))
	router.Handle("POST "+prefix+"/licence-keys", superAdminOnly(licenceKeyHandler.Create))
	router.Handle("PUT "+prefix+"/licence-keys/{id}", superAdminOnly(licenceKeyHandler.Update))
	router.Handle("DELETE "+prefix+"/licence-keys/{id}", superAdminOnly(licenceKeyHandler.Delete))
}
