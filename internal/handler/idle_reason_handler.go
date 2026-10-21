package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// IdleReasonHandler handles HTTP requests for idle reasons. The superadmin sees and
// changes every customer's reasons; everyone else reads their company's and managers
// change them (service.IdleReasonService decides which).
type IdleReasonHandler struct {
	service *service.IdleReasonService
}

// NewIdleReasonHandler creates a IdleReasonHandler that uses idleReasonService.
func NewIdleReasonHandler(idleReasonService *service.IdleReasonService) *IdleReasonHandler {
	return &IdleReasonHandler{service: idleReasonService}
}

// GetAll handles GET /api/v1/idle-reasons[?userId=<customer id>].
func (idleReasonHandler *IdleReasonHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	reasons, err := idleReasonHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, reasons)
}

// GetByID handles GET /api/v1/idle-reasons/{id}[?userId=<customer id>].
func (idleReasonHandler *IdleReasonHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	reason, err := idleReasonHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, reason)
}

// Create handles POST /api/v1/idle-reasons.
func (idleReasonHandler *IdleReasonHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var reasonRequest dto.CreateIdleReasonRequest
	if err := decodeAndValidate(writer, request, &reasonRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	reason, err := idleReasonHandler.service.Create(request.Context(), actor, reasonRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, reason)
}

// Update handles PUT /api/v1/idle-reasons/{id}.
func (idleReasonHandler *IdleReasonHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var reasonRequest dto.UpdateIdleReasonRequest
	if err := decodeAndValidate(writer, request, &reasonRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	reason, err := idleReasonHandler.service.Update(request.Context(), actor, id, reasonRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, reason)
}

// Delete handles DELETE /api/v1/idle-reasons/{id}.
func (idleReasonHandler *IdleReasonHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := idleReasonHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Idle reason deleted")
}
