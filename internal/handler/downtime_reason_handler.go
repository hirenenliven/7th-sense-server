package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// DowntimeReasonHandler handles HTTP requests for downtime reasons. The superadmin sees and
// changes every customer's reasons; everyone else reads their company's and managers
// change them (service.DowntimeReasonService decides which).
type DowntimeReasonHandler struct {
	service *service.DowntimeReasonService
}

// NewDowntimeReasonHandler creates a DowntimeReasonHandler that uses downtimeReasonService.
func NewDowntimeReasonHandler(downtimeReasonService *service.DowntimeReasonService) *DowntimeReasonHandler {
	return &DowntimeReasonHandler{service: downtimeReasonService}
}

// GetAll handles GET /api/v1/downtime-reasons[?userId=<customer id>].
func (downtimeReasonHandler *DowntimeReasonHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	reasons, err := downtimeReasonHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, reasons)
}

// GetByID handles GET /api/v1/downtime-reasons/{id}[?userId=<customer id>].
func (downtimeReasonHandler *DowntimeReasonHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	reason, err := downtimeReasonHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, reason)
}

// Create handles POST /api/v1/downtime-reasons.
func (downtimeReasonHandler *DowntimeReasonHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var reasonRequest dto.CreateDowntimeReasonRequest
	if err := decodeAndValidate(writer, request, &reasonRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	reason, err := downtimeReasonHandler.service.Create(request.Context(), actor, reasonRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, reason)
}

// Update handles PUT /api/v1/downtime-reasons/{id}.
func (downtimeReasonHandler *DowntimeReasonHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var reasonRequest dto.UpdateDowntimeReasonRequest
	if err := decodeAndValidate(writer, request, &reasonRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	reason, err := downtimeReasonHandler.service.Update(request.Context(), actor, id, reasonRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, reason)
}

// Delete handles DELETE /api/v1/downtime-reasons/{id}.
func (downtimeReasonHandler *DowntimeReasonHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := downtimeReasonHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Downtime reason deleted")
}
