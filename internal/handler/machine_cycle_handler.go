package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachineCycleHandler serves the machine cycles the CNC collector records. They are read only:
// every logged-in user reads the cycles of the machines they may access
// (service.MachineCycleService decides which).
type MachineCycleHandler struct {
	service *service.MachineCycleService
}

// NewMachineCycleHandler creates a MachineCycleHandler that uses cycleService.
func NewMachineCycleHandler(cycleService *service.MachineCycleService) *MachineCycleHandler {
	return &MachineCycleHandler{service: cycleService}
}

// GetAll handles GET /api/v1/machine-cycles?machineId=<id>[&from=YYYY-MM-DD][&to=YYYY-MM-DD]:
// the machine's cycles that started from the from day to the to day (both included; to
// defaults to today, from to the to day), with their totals.
func (cycleHandler *MachineCycleHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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
	listRequest := dto.MachineCycleListRequest{From: request.URL.Query().Get("from"), To: request.URL.Query().Get("to")}
	if machineID != nil {
		listRequest.MachineID = *machineID
	}
	if err := listRequest.Validate(); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	cycles, err := cycleHandler.service.GetAll(request.Context(), actor, listRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, cycles)
}

// GetByID handles GET /api/v1/machine-cycles/{id}.
func (cycleHandler *MachineCycleHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	cycle, err := cycleHandler.service.GetByID(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, cycle)
}
