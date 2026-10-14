package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Everyone logged in reads their company's parts (the superadmin every customer's) and
// their drawings. Managers add, edit and delete their company's parts and upload drawings.
func registerPartRoutes(router *http.ServeMux, prefix string, partHandler *handler.PartHandler) {
	router.HandleFunc("GET "+prefix+"/parts", partHandler.GetAll)
	router.HandleFunc("GET "+prefix+"/parts/{id}", partHandler.GetByID)
	router.HandleFunc("GET "+prefix+"/parts/{id}/processes/{processId}/drawing", partHandler.Drawing)
	router.Handle("POST "+prefix+"/parts", managersOnly(partHandler.Create))
	router.Handle("POST "+prefix+"/parts/drawings", managersOnly(partHandler.UploadDrawing))
	router.Handle("PUT "+prefix+"/parts/{id}", managersOnly(partHandler.Update))
	router.Handle("DELETE "+prefix+"/parts/{id}", managersOnly(partHandler.Delete))
}
