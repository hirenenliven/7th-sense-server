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
//
// @Summary      Get models
// @Description  Returns the models. A non-nil controllerID keeps only that controller's models.
// @Description
// @Description  Used by pages: Add / edit machine, Machine models
// @Tags         Machine catalog
// @ID           getModels
// @Produce      json
// @Security     cookieAuth
// @Param        controllerId query int false "Controller id"
// @Success      200 {object} response.Envelope{data=[]model.MachineModel}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/models [get]
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
//
// @Summary      Create model
// @Description  Create handles POST /api/v1/models.
// @Description  Roles: superadmin only (others get 403).
// @Description  Validation: Please select controller; Please enter valid name.
// @Description
// @Description  Used by pages: Machine models
// @Tags         Machine catalog
// @ID           createModel
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.MachineModelRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.MachineModel}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/models [post]
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
//
// @Summary      Update model
// @Description  Replaces the model's controller and name. Its controller cannot change while machines use it, because each machine stores the controller of its model.
// @Description  Roles: superadmin only (others get 403).
// @Description  Validation: Please select controller; Please enter valid name.
// @Description
// @Description  Used by pages: Machine models
// @Tags         Machine catalog
// @ID           updateModel
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.MachineModelRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.MachineModel}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/models/{id} [put]
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
//
// @Summary      Delete model
// @Description  Soft-deletes the model. A model used by machines is kept.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Machine models
// @Tags         Machine catalog
// @ID           deleteModel
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/models/{id} [delete]
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
