package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// BreakHandler handles HTTP requests for breaks: a customer's breaks its shifts take.
// The superadmin sees and changes every customer's breaks; everyone else reads their
// company's and managers change them (service.BreakService decides which).
type BreakHandler struct {
	service *service.BreakService
}

// NewBreakHandler creates a BreakHandler that uses breakService.
func NewBreakHandler(breakService *service.BreakService) *BreakHandler {
	return &BreakHandler{service: breakService}
}

// GetAll handles GET /api/v1/breaks[?userId=<customer id>].
//
// @Summary      Get breaks
// @Description  Returns the breaks the user may read. A non-nil customerID keeps only that customer's breaks.
// @Description
// @Description  Used by pages: Shifts, Breaks
// @Tags         Breaks
// @ID           getBreaks
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.Break}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/breaks [get]
func (breakHandler *BreakHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	breaks, err := breakHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, breaks)
}

// GetByID handles GET /api/v1/breaks/{id}[?userId=<customer id>].
func (breakHandler *BreakHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	shiftBreak, err := breakHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, shiftBreak)
}

// Create handles POST /api/v1/breaks.
//
// @Summary      Create break
// @Description  Adds a break, on the shift named or on none. A customer or admin (the logged-in user) adds it to their own company, whatever customer the request names; the superadmin to the customer named.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Breaks
// @Tags         Breaks
// @ID           createBreak
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateBreakRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.Break}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/breaks [post]
func (breakHandler *BreakHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var breakRequest dto.CreateBreakRequest
	if err := decodeAndValidate(writer, request, &breakRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	shiftBreak, err := breakHandler.service.Create(request.Context(), actor, breakRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, shiftBreak)
}

// Update handles PUT /api/v1/breaks/{id}.
//
// @Summary      Update break
// @Description  Replaces the break, its shift included (nil takes it off its shift). A blank status keeps the saved one.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid status.
// @Description
// @Description  Used by pages: Breaks
// @Tags         Breaks
// @ID           updateBreak
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateBreakRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.Break}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/breaks/{id} [put]
func (breakHandler *BreakHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var breakRequest dto.UpdateBreakRequest
	if err := decodeAndValidate(writer, request, &breakRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	shiftBreak, err := breakHandler.service.Update(request.Context(), actor, id, breakRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, shiftBreak)
}

// Delete handles DELETE /api/v1/breaks/{id}.
//
// @Summary      Delete break
// @Description  Soft-deletes the break (its shift no longer takes it). A customer or admin (the logged-in user) deletes only their company's breaks ("Break not found" for others).
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Breaks
// @Tags         Breaks
// @ID           deleteBreak
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/breaks/{id} [delete]
func (breakHandler *BreakHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := breakHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Break deleted")
}
