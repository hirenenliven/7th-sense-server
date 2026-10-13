package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachineHandler handles HTTP requests for machines.
type MachineHandler struct {
	service *service.MachineService
}

// NewMachineHandler creates a MachineHandler that uses machineService.
func NewMachineHandler(machineService *service.MachineService) *MachineHandler {
	return &MachineHandler{service: machineService}
}

// GetAll handles GET /api/v1/machines.
func (machineHandler *MachineHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	machines, err := machineHandler.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machines)
}

// GetByID handles GET /api/v1/machines/{id}.
func (machineHandler *MachineHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machine, err := machineHandler.service.GetByID(request.Context(), id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machine)
}

// Create handles POST /api/v1/machines.
func (machineHandler *MachineHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var machineRequest dto.CreateMachineRequest
	if err := decodeAndValidate(writer, request, &machineRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machine, err := machineHandler.service.Create(request.Context(), machineRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, machine)
}

// Update handles PUT /api/v1/machines/{id}.
func (machineHandler *MachineHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var machineRequest dto.UpdateMachineRequest
	if err := decodeAndValidate(writer, request, &machineRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machine, err := machineHandler.service.Update(request.Context(), id, machineRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machine)
}

// Delete handles DELETE /api/v1/machines/{id}.
func (machineHandler *MachineHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := machineHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Machine deleted")
}
