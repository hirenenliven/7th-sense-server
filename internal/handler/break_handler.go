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
