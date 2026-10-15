package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// HolidayHandler handles HTTP requests for holidays and week offs. The superadmin sees
// and changes every customer's; everyone else reads their company's and managers change
// them (service.HolidayService decides which).
type HolidayHandler struct {
	service *service.HolidayService
}

// NewHolidayHandler creates a HolidayHandler that uses holidayService.
func NewHolidayHandler(holidayService *service.HolidayService) *HolidayHandler {
	return &HolidayHandler{service: holidayService}
}

// GetAll handles GET /api/v1/holidays[?userId=<customer id>].
func (holidayHandler *HolidayHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	holidays, err := holidayHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, holidays)
}

// GetByID handles GET /api/v1/holidays/{id}[?userId=<customer id>].
func (holidayHandler *HolidayHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	holiday, err := holidayHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, holiday)
}

// Create handles POST /api/v1/holidays.
func (holidayHandler *HolidayHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var holidayRequest dto.CreateHolidayRequest
	if err := decodeAndValidate(writer, request, &holidayRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	holiday, err := holidayHandler.service.Create(request.Context(), actor, holidayRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, holiday)
}

// Update handles PUT /api/v1/holidays/{id}.
func (holidayHandler *HolidayHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var holidayRequest dto.UpdateHolidayRequest
	if err := decodeAndValidate(writer, request, &holidayRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	holiday, err := holidayHandler.service.Update(request.Context(), actor, id, holidayRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, holiday)
}

// Delete handles DELETE /api/v1/holidays/{id}.
func (holidayHandler *HolidayHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := holidayHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Holiday deleted")
}
