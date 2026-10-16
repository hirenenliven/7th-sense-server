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
