package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachineHandler handles HTTP requests for machines. Admins see and change every
// customer's machines; other users only read their own (see machineCustomerFilter).
type MachineHandler struct {
	service *service.MachineService
}

// NewMachineHandler creates a MachineHandler that uses machineService.
func NewMachineHandler(machineService *service.MachineService) *MachineHandler {
	return &MachineHandler{service: machineService}
}

// GetAll handles GET /api/v1/machines[?userId=<customer id>].
func (machineHandler *MachineHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	customerID, err := machineCustomerFilter(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machines, err := machineHandler.service.GetAll(request.Context(), customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machines)
}

// GetByID handles GET /api/v1/machines/{id}.
func (machineHandler *MachineHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	customerID, err := machineCustomerFilter(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machine, err := machineHandler.service.GetByID(request.Context(), customerID, id)
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

// machineCustomerFilter returns whose machines the request may read: for an admin the
// customer named by ?userId=<id>, or nil for every customer; for any other user always
// their own id. A non-admin naming another user gets Forbidden.
func machineCustomerFilter(request *http.Request) (*int64, error) {
	currentUser, ok := middleware.CurrentUser(request.Context())
	if !ok {
		return nil, response.Unauthorized
	}
	userID, err := queryID(request, "userId")
	if err != nil {
		return nil, err
	}
	if model.IsAdmin(currentUser.Role) {
		return userID, nil
	}
	if userID != nil && *userID != currentUser.ID {
		return nil, response.Forbidden
	}
	return &currentUser.ID, nil
}
