package handler

import (
	"net/http"

	"7th-sense-server/internal/response"
)

// GetVirtualScreen handles GET /api/v1/productivity/virtual-screen?machineId=: what the machine is
// doing now (mode, program, live CNC values, last alarm, running cycle and its job).
//
// @Summary      Get virtual screen
// @Description  GetVirtualScreen handles GET /api/v1/productivity/virtual-screen?machineId=: what the machine is doing now: mode and program (open mode period), operation mode, the latest CNC values (spindle load / speed / override, feed rate / override, X / Y / Z axis load, tool number, E-stop; null when never read), the last alarm raised (active while no alarm_off followed), the running cycle (elapsed and planned seconds) and its job (the running cycle's, else the one scheduled now: quantity, completed cycles, times, part, process, clientele). Read only; the page reloads it every few seconds.
// @Description
// @Description  Used by pages: Monitoring (Virtual Screen tab)
// @Tags         Productivity
// @ID           getVirtualScreen
// @Produce      json
// @Security     cookieAuth
// @Param        machineId query int true "Machine id"
// @Success      200 {object} response.Envelope{data=dto.VirtualScreenResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Machine not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/productivity/virtual-screen [get]
func (dashboardHandler *DashboardHandler) GetVirtualScreen(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	machineID, err := queryID(request, "machineId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	if machineID == nil {
		response.SendResponse(writer, request, response.NewError(http.StatusBadRequest, "Please select a machine"), nil)
		return
	}

	screen, err := dashboardHandler.service.GetVirtualScreen(request.Context(), actor, *machineID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, screen)
}
