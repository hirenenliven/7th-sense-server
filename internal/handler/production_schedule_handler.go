package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// ProductionScheduleHandler handles HTTP requests for production schedules. The superadmin
// sees and changes every
// customer's schedules; everyone else reads their company's and managers change them
// (service.ProductionScheduleService decides which).
type ProductionScheduleHandler struct {
	service *service.ProductionScheduleService
}

// NewProductionScheduleHandler creates a ProductionScheduleHandler that uses scheduleService.
func NewProductionScheduleHandler(scheduleService *service.ProductionScheduleService) *ProductionScheduleHandler {
	return &ProductionScheduleHandler{service: scheduleService}
}

// GetAll handles GET /api/v1/production-schedules[?userId=<customer id>].
//
// @Summary      Get production schedules
// @Description  Returns the schedules the user may read (see companyOwner), without their items. A non-nil customerID keeps only that customer's schedules.
// @Description
// @Description  Used by pages: Production schedules
// @Tags         Production schedules
// @ID           getProductionSchedules
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.ProductionSchedule}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/production-schedules [get]
func (scheduleHandler *ProductionScheduleHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	schedules, err := scheduleHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, schedules)
}

// GetByID handles GET /api/v1/production-schedules/{id}[?userId=<customer id>].
//
// @Summary      Get production schedule by id
// @Description  Returns the schedule with its items when the user may read it, else "Production schedule not found".
// @Description
// @Description  Used by pages: Add / edit production schedule
// @Tags         Production schedules
// @ID           getProductionScheduleById
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=model.ProductionSchedule}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/production-schedules/{id} [get]
func (scheduleHandler *ProductionScheduleHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	schedule, err := scheduleHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, schedule)
}

// Create handles POST /api/v1/production-schedules.
//
// @Summary      Create production schedule
// @Description  Adds a schedule with its items. A customer or admin (the logged-in user) adds it to their own company, whatever customer the request names; the superadmin to the customer named.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Add / edit production schedule
// @Tags         Production schedules
// @ID           createProductionSchedule
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateProductionScheduleRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.ProductionSchedule}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/production-schedules [post]
func (scheduleHandler *ProductionScheduleHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var scheduleRequest dto.CreateProductionScheduleRequest
	if err := decodeAndValidate(writer, request, &scheduleRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	schedule, err := scheduleHandler.service.Create(request.Context(), actor, scheduleRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, schedule)
}

// Update handles PUT /api/v1/production-schedules/{id}.
//
// @Summary      Update production schedule
// @Description  Replaces the schedule and its items. Items with an id are updated (they must be the schedule's), items without one are added and saved items missing from the request are deleted.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid status.
// @Description
// @Description  Used by pages: Add / edit production schedule
// @Tags         Production schedules
// @ID           updateProductionSchedule
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateProductionScheduleRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.ProductionSchedule}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/production-schedules/{id} [put]
func (scheduleHandler *ProductionScheduleHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var scheduleRequest dto.UpdateProductionScheduleRequest
	if err := decodeAndValidate(writer, request, &scheduleRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	schedule, err := scheduleHandler.service.Update(request.Context(), actor, id, scheduleRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, schedule)
}

// Delete handles DELETE /api/v1/production-schedules/{id}.
//
// @Summary      Delete production schedule
// @Description  Soft-deletes the schedule. A customer or admin (the logged-in user) deletes only their company's schedules ("Production schedule not found" for others).
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Production schedules
// @Tags         Production schedules
// @ID           deleteProductionSchedule
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/production-schedules/{id} [delete]
func (scheduleHandler *ProductionScheduleHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := scheduleHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Production schedule deleted")
}
