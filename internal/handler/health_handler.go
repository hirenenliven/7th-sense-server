package handler

import (
	"net/http"

	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// HealthHandler serves the verification endpoints of CNC data collection.
type HealthHandler struct {
	service *service.HealthService
}

// NewHealthHandler creates a HealthHandler that uses healthService.
func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{service: healthService}
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
