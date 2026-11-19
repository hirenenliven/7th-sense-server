package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachineAccuracyHandler handles HTTP requests for the machine accuracy list (the Machine Accuracy
// options of the machine and part forms).
type MachineAccuracyHandler struct {
	service *service.MachineAccuracyService
}

// NewMachineAccuracyHandler creates a MachineAccuracyHandler that uses accuracyService.
func NewMachineAccuracyHandler(accuracyService *service.MachineAccuracyService) *MachineAccuracyHandler {
	return &MachineAccuracyHandler{service: accuracyService}
}

// GetAll handles GET /api/v1/machine-accuracies.
//
// @Summary      Get machine accuracies
// @Description  Returns the machine accuracy list (shared by every customer), ordered by id, each with usageCount (the machines and part processes that use it). Every logged-in user reads it.
// @Description
// @Description  Used by pages: Machine accuracy, Add / edit machine, Add / edit part
// @Tags         Machine accuracies
// @ID           getMachineAccuracies
// @Produce      json
// @Security     cookieAuth
// @Success      200 {object} response.Envelope{data=[]model.MachineAccuracy}
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-accuracies [get]
func (accuracyHandler *MachineAccuracyHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	accuracies, err := accuracyHandler.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, accuracies)
}

// Create handles POST /api/v1/machine-accuracies: a new accuracy for the list.
//
// @Summary      Add machine accuracy
// @Description  Adds an accuracy to the list; every customer's machine and part forms offer it at once. Spaces are trimmed and joined; the name is unique, ignoring case.
// @Description  Roles: superadmin (others get 403).
// @Description  Validation: Please enter machine accuracy; Machine accuracy must be at most 50 characters; This machine accuracy already exists.
// @Description
// @Description  Used by pages: Machine accuracy
// @Tags         Machine accuracies
// @ID           createMachineAccuracy
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.MachineAccuracyRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.MachineAccuracy}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-accuracies [post]
func (accuracyHandler *MachineAccuracyHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var accuracyRequest dto.MachineAccuracyRequest
	if err := decodeAndValidate(writer, request, &accuracyRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	accuracy, err := accuracyHandler.service.Create(request.Context(), accuracyRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, accuracy)
}

// Update handles PUT /api/v1/machine-accuracies/{id}.
//
// @Summary      Update machine accuracy
// @Description  Renames the accuracy. Machines and part processes that use it show the new name.
// @Description  Roles: superadmin (others get 403).
// @Description  Validation: Please enter machine accuracy; Machine accuracy must be at most 50 characters; This machine accuracy already exists; Machine accuracy not found.
// @Description
// @Description  Used by pages: Machine accuracy
// @Tags         Machine accuracies
// @ID           updateMachineAccuracy
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Machine accuracy id"
// @Param        request body dto.MachineAccuracyRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.MachineAccuracy}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Machine accuracy not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-accuracies/{id} [put]
func (accuracyHandler *MachineAccuracyHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	var accuracyRequest dto.MachineAccuracyRequest
	if err := decodeAndValidate(writer, request, &accuracyRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	accuracy, err := accuracyHandler.service.Update(request.Context(), id, accuracyRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, accuracy)
}

// Delete handles DELETE /api/v1/machine-accuracies/{id}.
//
// @Summary      Delete machine accuracy
// @Description  Soft-deletes the accuracy. One used by machines or part processes is refused (400): change it there first.
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Machine accuracy
// @Tags         Machine accuracies
// @ID           deleteMachineAccuracy
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Machine accuracy id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Accuracy in use"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Machine accuracy not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-accuracies/{id} [delete]
func (accuracyHandler *MachineAccuracyHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	if err := accuracyHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Machine accuracy deleted")
}
