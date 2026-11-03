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
//
// @Summary      Get idle reasons
// @Description  Returns the idle reasons the user may read. A non-nil customerID keeps only that customer's reasons.
// @Description
// @Description  Used by pages: Downtime / idle reasons
// @Tags         Reasons
// @ID           getIdleReasons
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.IdleReason}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/idle-reasons [get]
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
//
// @Summary      Create idle reason
// @Description  Adds a idle reason. A customer or admin (the logged-in user) adds it to their own company, whatever customer the request names; the superadmin to the customer named.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Downtime / idle reasons
// @Tags         Reasons
// @ID           createIdleReason
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateIdleReasonRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.IdleReason}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/idle-reasons [post]
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
//
// @Summary      Update idle reason
// @Description  Replaces the idle reason. A blank status keeps the saved one.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid status.
// @Description
// @Description  Used by pages: Downtime / idle reasons
// @Tags         Reasons
// @ID           updateIdleReason
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateIdleReasonRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.IdleReason}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/idle-reasons/{id} [put]
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
//
// @Summary      Delete idle reason
// @Description  Soft-deletes the idle reason. A customer or admin (the logged-in user) deletes only their company's reasons ("Idle reason not found" for others).
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Downtime / idle reasons
// @Tags         Reasons
// @ID           deleteIdleReason
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/idle-reasons/{id} [delete]
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
