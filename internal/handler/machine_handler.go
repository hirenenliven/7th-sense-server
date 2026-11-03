package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachineHandler handles HTTP requests for machines. The superadmin sees and
// changes every customer's machines; a customer and its admins add and edit their
// company's; other users only read the machines they may access
// (service.MachineService decides which).
type MachineHandler struct {
	service *service.MachineService
}

// NewMachineHandler creates a MachineHandler that uses machineService.
func NewMachineHandler(machineService *service.MachineService) *MachineHandler {
	return &MachineHandler{service: machineService}
}

// GetAll handles GET /api/v1/machines[?userId=<customer id>].
//
// @Summary      Get machines
// @Description  Returns the machines the user may access (see machineScope). A non-nil customerID keeps only that customer's machines.
// @Description
// @Description  Used by pages: Add / edit customer or user, Customer / user detail, All machines, Machine health, Productivity, Add / edit part, Add / edit production schedule
// @Tags         Machines
// @ID           getMachines
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.Machine}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines [get]
func (machineHandler *MachineHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	customerID, err := queryID(request, "userId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machines, err := machineHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machines)
}

// GetByID handles GET /api/v1/machines/{id}[?userId=<customer id>].
//
// @Summary      Get machine by id
// @Description  Returns the machine when the user may access it (see machineScope), else "Machine not found". A non-nil customerID finds it only when that customer owns it.
// @Description
// @Description  Used by pages: Add / edit machine
// @Tags         Machines
// @ID           getMachineById
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=model.Machine}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id} [get]
func (machineHandler *MachineHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	customerID, err := queryID(request, "userId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machine, err := machineHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machine)
}

// Create handles POST /api/v1/machines.
//
// @Summary      Create machine
// @Description  Adds a machine. A customer or admin (the logged-in user) adds it to their own company, whatever customer the request names; the superadmin to the customer named.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid password.
// @Description
// @Description  Used by pages: Add / edit machine
// @Tags         Machines
// @ID           createMachine
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateMachineRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.Machine}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines [post]
func (machineHandler *MachineHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var machineRequest dto.CreateMachineRequest
	if err := decodeAndValidate(writer, request, &machineRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machine, err := machineHandler.service.Create(request.Context(), actor, machineRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, machine)
}

// Update handles PUT /api/v1/machines/{id}.
//
// @Summary      Update machine
// @Description  Replaces the form fields. A blank password or status keeps the saved one.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid password; Please enter valid status.
// @Description
// @Description  Used by pages: Add / edit machine
// @Tags         Machines
// @ID           updateMachine
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateMachineRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.Machine}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id} [put]
func (machineHandler *MachineHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var machineRequest dto.UpdateMachineRequest
	if err := decodeAndValidate(writer, request, &machineRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	machine, err := machineHandler.service.Update(request.Context(), actor, id, machineRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, machine)
}

// Delete handles DELETE /api/v1/machines/{id}.
//
// @Summary      Delete machine
// @Description  Delete handles DELETE /api/v1/machines/{id}.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: All machines
// @Tags         Machines
// @ID           deleteMachine
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id} [delete]
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
