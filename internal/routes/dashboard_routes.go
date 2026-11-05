package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
)

// Page-wise chart routes: /productivity/<chart> and /machine-health/<chart>. Everyone logged in
// reads the machines they may access; the data is what the CNC collector recorded (read only).
func registerDashboardRoutes(router *http.ServeMux, prefix string, dashboardHandler *handler.DashboardHandler) {
	router.HandleFunc("GET "+prefix+"/productivity/mode-timeline", dashboardHandler.GetModeTimeline)
	router.HandleFunc("GET "+prefix+"/productivity/mode-totals", dashboardHandler.GetModeTotals)
	router.HandleFunc("GET "+prefix+"/productivity/oee", dashboardHandler.GetOEE)
	router.HandleFunc("GET "+prefix+"/productivity/cycle-counts", dashboardHandler.GetCycleCounts)
	router.HandleFunc("GET "+prefix+"/productivity/finish-cycles", dashboardHandler.GetFinishCycles)
	router.HandleFunc("GET "+prefix+"/productivity/machine-summaries", dashboardHandler.GetMachineSummaries)
	router.HandleFunc("GET "+prefix+"/machine-health/machine-time-by-day", dashboardHandler.GetMachineTimeByDay)
}
