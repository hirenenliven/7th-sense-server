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
//
// @Summary      Get shifts
// @Description  Returns the shifts the user may read, each with its breaks. A non-nil customerID keeps only that customer's shifts.
// @Description
// @Description  Used by pages: Productivity, Add / edit production schedule, Shifts, Breaks, Operations Center, Plant Intelligence
// @Tags         Shifts
// @ID           getShifts
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.Shift}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/shifts [get]
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
//
// @Summary      Create shift
// @Description  Adds a shift with its breaks. A customer or admin (the logged-in user) adds it to their own company, whatever customer the request names; the superadmin to the customer named.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Shifts
// @Tags         Shifts
// @ID           createShift
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateShiftRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.Shift}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/shifts [post]
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
//
// @Summary      Update shift
// @Description  Replaces the shift and its breaks: breaks left out are released (no shift). A blank status keeps the saved one.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid status.
// @Description
// @Description  Used by pages: Shifts
// @Tags         Shifts
// @ID           updateShift
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateShiftRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.Shift}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/shifts/{id} [put]
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
//
// @Summary      Delete shift
// @Description  Soft-deletes the shift and releases its breaks (they stay in the customer's list without a shift). A customer or admin (the logged-in user) deletes only their company's shifts ("Shift not found" for others).
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Shifts
// @Tags         Shifts
// @ID           deleteShift
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/shifts/{id} [delete]
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
