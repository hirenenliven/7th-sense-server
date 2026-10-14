package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachineModelHandler handles HTTP requests for machine models (models table).
type MachineModelHandler struct {
	service *service.MachineModelService
}

// NewMachineModelHandler creates a MachineModelHandler that uses machineModelService.
func NewMachineModelHandler(machineModelService *service.MachineModelService) *MachineModelHandler {
	return &MachineModelHandler{service: machineModelService}
}

// GetAll handles GET /api/v1/models. ?controllerId= keeps only that controller's models.
func (machineModelHandler *MachineModelHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	controllerID, err := queryID(request, "controllerId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	models, err := machineModelHandler.service.GetAll(request.Context(), controllerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, models)
}

// GetByID handles GET /api/v1/models/{id}.
func (machineModelHandler *MachineModelHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machineModel, err := machineModelHandler.service.GetByID(request.Context(), id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machineModel)
}

// Create handles POST /api/v1/models.
func (machineModelHandler *MachineModelHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var modelRequest dto.MachineModelRequest
	if err := decodeAndValidate(writer, request, &modelRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machineModel, err := machineModelHandler.service.Create(request.Context(), modelRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, machineModel)
}

// Update handles PUT /api/v1/models/{id}.
func (machineModelHandler *MachineModelHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var modelRequest dto.MachineModelRequest
	if err := decodeAndValidate(writer, request, &modelRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machineModel, err := machineModelHandler.service.Update(request.Context(), id, modelRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machineModel)
}

// Delete handles DELETE /api/v1/models/{id}.
func (machineModelHandler *MachineModelHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := machineModelHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Model deleted")
}
