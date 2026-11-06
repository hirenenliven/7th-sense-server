package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachineModeReasonHandler handles HTTP requests for the reasons recorded against the Downtime and
// Idle periods of a machine (service.MachineModeReasonService decides what is allowed).
type MachineModeReasonHandler struct {
	service *service.MachineModeReasonService
}

// NewMachineModeReasonHandler creates a MachineModeReasonHandler that uses reasonService.
func NewMachineModeReasonHandler(reasonService *service.MachineModeReasonService) *MachineModeReasonHandler {
	return &MachineModeReasonHandler{service: reasonService}
}

// GetByPeriod handles GET /api/v1/machine-mode-reasons?periodId=<machine_mode_log id>.
//
// @Summary      Get machine mode reasons
// @Description  Returns the whole Downtime or Idle period (not clipped to a chart window), which master list its reasons come from, the machine's customer and the reasons recorded so far in time order. A period of another mode, or one still running, answers 400 with the message to show.
// @Description  Roles: anyone who may read the machine (others get 404).
// @Description
// @Description  Used by pages: Productivity
// @Tags         Reasons
// @ID           getMachineModeReasons
// @Produce      json
// @Security     cookieAuth
// @Param        periodId query int true "machine_mode_log id (a Gantt bar's id)"
// @Success      200 {object} response.Envelope{data=dto.MachineModeReasonsResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-mode-reasons [get]
func (reasonHandler *MachineModeReasonHandler) GetByPeriod(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	periodID, err := queryID(request, "periodId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	if periodID == nil {
		response.SendResponse(writer, request, response.NewError(http.StatusBadRequest, "Please select machine mode period"), nil)
		return
	}

	reasons, err := reasonHandler.service.GetByPeriod(request.Context(), actor, *periodID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, reasons)
}

// Create handles POST /api/v1/machine-mode-reasons.
//
// @Summary      Create machine mode reason
// @Description  Records a reason for a part of a finished Downtime or Idle period.
// @Description  Roles: anyone who may read the machine (others get 404).
// @Description  Validation: the time lies inside the period and overlaps no other reason of it; the reason is active, of the machine's customer, in the list the period's mode uses.
// @Description
// @Description  Used by pages: Productivity
// @Tags         Reasons
// @ID           createMachineModeReason
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateMachineModeReasonRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.MachineModeReason}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-mode-reasons [post]
func (reasonHandler *MachineModeReasonHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var reasonRequest dto.CreateMachineModeReasonRequest
	if err := decodeAndValidate(writer, request, &reasonRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	reason, err := reasonHandler.service.Create(request.Context(), actor, reasonRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, reason)
}

// Update handles PUT /api/v1/machine-mode-reasons/{id}.
//
// @Summary      Update machine mode reason
// @Description  Replaces the reason and the time of a recorded row, by the rules of create. The reason the row already has stays allowed when it has since been made inactive or deleted from its list.
// @Description  Roles: anyone who may read the machine (others get 404).
// @Description
// @Description  Used by pages: Productivity
// @Tags         Reasons
// @ID           updateMachineModeReason
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateMachineModeReasonRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.MachineModeReason}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-mode-reasons/{id} [put]
func (reasonHandler *MachineModeReasonHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var reasonRequest dto.UpdateMachineModeReasonRequest
	if err := decodeAndValidate(writer, request, &reasonRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	reason, err := reasonHandler.service.Update(request.Context(), actor, id, reasonRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, reason)
}

// Delete handles DELETE /api/v1/machine-mode-reasons/{id}.
//
// @Summary      Delete machine mode reason
// @Description  Soft-deletes a recorded row.
// @Description  Roles: anyone who may read the machine (others get 404).
// @Description
// @Description  Used by pages: Productivity
// @Tags         Reasons
// @ID           deleteMachineModeReason
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-mode-reasons/{id} [delete]
func (reasonHandler *MachineModeReasonHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := reasonHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Machine mode reason deleted")
}
