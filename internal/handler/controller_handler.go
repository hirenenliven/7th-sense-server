package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// ControllerHandler handles HTTP requests for controllers.
type ControllerHandler struct {
	service *service.ControllerService
}

// NewControllerHandler creates a ControllerHandler that uses controllerService.
func NewControllerHandler(controllerService *service.ControllerService) *ControllerHandler {
	return &ControllerHandler{service: controllerService}
}

// GetAll handles GET /api/v1/controllers. ?brandId= keeps only that brand's controllers.
func (controllerHandler *ControllerHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	brandID, err := queryID(request, "brandId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	controllers, err := controllerHandler.service.GetAll(request.Context(), brandID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, controllers)
}

// GetByID handles GET /api/v1/controllers/{id}.
func (controllerHandler *ControllerHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	controller, err := controllerHandler.service.GetByID(request.Context(), id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, controller)
}

// Create handles POST /api/v1/controllers.
func (controllerHandler *ControllerHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var controllerRequest dto.CreateControllerRequest
	if err := decodeAndValidate(writer, request, &controllerRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	controller, err := controllerHandler.service.Create(request.Context(), controllerRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, controller)
}

// Update handles PUT /api/v1/controllers/{id}.
func (controllerHandler *ControllerHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var controllerRequest dto.UpdateControllerRequest
	if err := decodeAndValidate(writer, request, &controllerRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	controller, err := controllerHandler.service.Update(request.Context(), id, controllerRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, controller)
}

// Delete handles DELETE /api/v1/controllers/{id}.
func (controllerHandler *ControllerHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := controllerHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Controller deleted")
}
