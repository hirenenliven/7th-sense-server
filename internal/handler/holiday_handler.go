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
//
// @Summary      Get holidays
// @Description  Returns the holidays and week offs the user may read. A non-nil customerID keeps only that customer's.
// @Description
// @Description  Used by pages: Holidays / week off, Operations Center, Plant Intelligence
// @Tags         Holidays
// @ID           getHolidays
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.Holiday}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/holidays [get]
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
//
// @Summary      Create holiday
// @Description  Adds a holiday or week off. A customer or admin (the logged-in user) adds it to their own company, whatever customer the request names; the superadmin to the customer named.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Holidays / week off
// @Tags         Holidays
// @ID           createHoliday
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateHolidayRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.Holiday}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/holidays [post]
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
//
// @Summary      Update holiday
// @Description  Replaces the holiday (it may change between holiday and week off). A blank status keeps the saved one.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid status.
// @Description
// @Description  Used by pages: Holidays / week off
// @Tags         Holidays
// @ID           updateHoliday
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateHolidayRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.Holiday}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/holidays/{id} [put]
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
//
// @Summary      Delete holiday
// @Description  Soft-deletes the holiday. A customer or admin (the logged-in user) deletes only their company's ("Holiday not found" for others).
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Holidays / week off
// @Tags         Holidays
// @ID           deleteHoliday
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/holidays/{id} [delete]
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
