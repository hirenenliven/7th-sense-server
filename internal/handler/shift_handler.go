package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// ShiftHandler handles HTTP requests for shifts and their breaks. The superadmin sees
// and changes every customer's shifts; everyone else reads their company's and
// managers change them (service.ShiftService decides which).
type ShiftHandler struct {
	service *service.ShiftService
}

// NewShiftHandler creates a ShiftHandler that uses shiftService.
func NewShiftHandler(shiftService *service.ShiftService) *ShiftHandler {
	return &ShiftHandler{service: shiftService}
}

// GetAll handles GET /api/v1/shifts[?userId=<customer id>].
func (shiftHandler *ShiftHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	shifts, err := shiftHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, shifts)
}

// GetByID handles GET /api/v1/shifts/{id}[?userId=<customer id>].
func (shiftHandler *ShiftHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	shift, err := shiftHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, shift)
}

// Create handles POST /api/v1/shifts.
func (shiftHandler *ShiftHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var shiftRequest dto.CreateShiftRequest
	if err := decodeAndValidate(writer, request, &shiftRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	shift, err := shiftHandler.service.Create(request.Context(), actor, shiftRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, shift)
}

// Update handles PUT /api/v1/shifts/{id}.
func (shiftHandler *ShiftHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var shiftRequest dto.UpdateShiftRequest
	if err := decodeAndValidate(writer, request, &shiftRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	shift, err := shiftHandler.service.Update(request.Context(), actor, id, shiftRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, shift)
}

// Delete handles DELETE /api/v1/shifts/{id}.
func (shiftHandler *ShiftHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := shiftHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Shift deleted")
}
