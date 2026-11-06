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
//
// @Summary      Get downtime reasons
// @Description  Returns the downtime reasons the user may read. A non-nil customerID keeps only that customer's reasons.
// @Description
// @Description  Used by pages: Downtime / idle reasons, Operations Center, Productivity, Plant Intelligence
// @Tags         Reasons
// @ID           getDowntimeReasons
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.DowntimeReason}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/downtime-reasons [get]
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
//
// @Summary      Create downtime reason
// @Description  Adds a downtime reason. A customer or admin (the logged-in user) adds it to their own company, whatever customer the request names; the superadmin to the customer named.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Downtime / idle reasons
// @Tags         Reasons
// @ID           createDowntimeReason
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateDowntimeReasonRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.DowntimeReason}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/downtime-reasons [post]
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
//
// @Summary      Update downtime reason
// @Description  Replaces the downtime reason. A blank status keeps the saved one.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid status.
// @Description
// @Description  Used by pages: Downtime / idle reasons
// @Tags         Reasons
// @ID           updateDowntimeReason
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateDowntimeReasonRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.DowntimeReason}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/downtime-reasons/{id} [put]
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
//
// @Summary      Delete downtime reason
// @Description  Soft-deletes the downtime reason. A customer or admin (the logged-in user) deletes only their company's reasons ("Downtime reason not found" for others).
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Downtime / idle reasons
// @Tags         Reasons
// @ID           deleteDowntimeReason
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/downtime-reasons/{id} [delete]
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
