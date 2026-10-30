package handler

import (
	"net/http"

	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// HealthHandler serves the verification endpoints of CNC data collection.
type HealthHandler struct {
	service      *service.HealthService
	cycleService *service.MachineCycleService
}

// NewHealthHandler creates a HealthHandler that uses healthService, and cycleService for the
// machine cycles.
func NewHealthHandler(healthService *service.HealthService, cycleService *service.MachineCycleService) *HealthHandler {
	return &HealthHandler{service: healthService, cycleService: cycleService}
}

// GetMachineModes handles GET /api/v1/health/machines: every machine with the mode it is in
// now. Public: no session needed.
func (healthHandler *HealthHandler) GetMachineModes(writer http.ResponseWriter, request *http.Request) {
	machines, err := healthHandler.service.GetMachineModes(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machines)
}

// GetMachineModeHistory handles GET /api/v1/health/machines/{id}/modes: the machine's mode
// periods of the last 10 days and the time per mode. Public: no session needed.
func (healthHandler *HealthHandler) GetMachineModeHistory(writer http.ResponseWriter, request *http.Request) {
	machineID, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	history, err := healthHandler.service.GetMachineModeHistory(request.Context(), machineID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, history)
}

// GetMachineCycleHistory handles GET /api/v1/health/machines/{id}/cycles: the machine's cycles
// of the last 24 hours with their totals. Public: no session needed.
func (healthHandler *HealthHandler) GetMachineCycleHistory(writer http.ResponseWriter, request *http.Request) {
	machineID, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	cycles, err := healthHandler.cycleService.GetMachineCycleHistory(request.Context(), machineID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, cycles)
}
