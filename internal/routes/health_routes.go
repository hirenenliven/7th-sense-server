package routes

import (
	"io/fs"
	"net/http"

	"7th-sense-server/internal/handler"
	"7th-sense-server/web"
)

// registerHealthRoutes registers the machine list with current modes, the machine mode and
// cycle history and their chart page. All are public (no session) so the chart page can be
// opened on its own.
func registerHealthRoutes(router *http.ServeMux, prefix string, healthHandler *handler.HealthHandler) {
	router.HandleFunc("GET "+prefix+"/health/machines", healthHandler.GetMachineModes)
	router.HandleFunc("GET "+prefix+"/health/machines/{id}/modes", healthHandler.GetMachineModeHistory)
	router.HandleFunc("GET "+prefix+"/health/machines/{id}/cycles", healthHandler.GetMachineCycleHistory)

	modeChart, err := fs.Sub(web.ModeChart, "mode-chart")
	if err != nil {
		panic(err) // the directory is embedded at build time, so this cannot fail at run time
	}
	router.Handle("GET /mode-chart/", http.StripPrefix("/mode-chart/", http.FileServerFS(modeChart)))
}
