package handler

import (
	"net/http"

	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachineDiagnosticsHandler handles the superadmin's checks of one machine's data collection
// (service.MachineDiagnosticsService).
type MachineDiagnosticsHandler struct {
	service *service.MachineDiagnosticsService
}

// NewMachineDiagnosticsHandler creates a MachineDiagnosticsHandler that uses diagnosticsService.
func NewMachineDiagnosticsHandler(diagnosticsService *service.MachineDiagnosticsService) *MachineDiagnosticsHandler {
	return &MachineDiagnosticsHandler{service: diagnosticsService}
}

// GetDiagnostics handles GET /api/v1/machines/{id}/diagnostics.
//
// @Summary      Get machine diagnostics
// @Description  What the server knows about one machine's data collection now: set-up, collector status, current mode, every variable (readable or why not, latest value, last stored number), recent machine events, and the checks worked out from them with the reason and what to do for each problem. Sends nothing to the machine.
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Machine diagnostics
// @Tags         Machines
// @ID           getMachineDiagnostics
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Machine id"
// @Success      200 {object} response.Envelope{data=dto.MachineDiagnosticsResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id}/diagnostics [get]
func (diagnosticsHandler *MachineDiagnosticsHandler) GetDiagnostics(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	diagnostics, err := diagnosticsHandler.service.GetDiagnostics(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, diagnostics)
}

// CheckNetwork handles GET /api/v1/machines/{id}/diagnostics/network.
//
// @Summary      Check machine network
// @Description  Tests the machine's stored live IP and port from this server, where the collector connects from: one ping and one TCP connect that sends no data. Takes up to a few seconds, so the page runs it on opening and when asked, not on a timer.
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Machine diagnostics
// @Tags         Machines
// @ID           checkMachineNetwork
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Machine id"
// @Success      200 {object} response.Envelope{data=dto.MachineNetworkResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id}/diagnostics/network [get]
func (diagnosticsHandler *MachineDiagnosticsHandler) CheckNetwork(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	network, err := diagnosticsHandler.service.CheckNetwork(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, network)
}
