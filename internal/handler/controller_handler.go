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
//
// @Summary      Get controllers
// @Description  Returns the controllers. A non-nil brandID keeps only that brand's controllers.
// @Description
// @Description  Used by pages: Add / edit machine, Controllers, Machine models
// @Tags         Machine catalog
// @ID           getControllers
// @Produce      json
// @Security     cookieAuth
// @Param        brandId query int false "Brand id"
// @Success      200 {object} response.Envelope{data=[]model.Controller}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/controllers [get]
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
//
// @Summary      Create controller
// @Description  Create handles POST /api/v1/controllers.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Controllers
// @Tags         Machine catalog
// @ID           createController
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateControllerRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.Controller}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/controllers [post]
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
//
// @Summary      Update controller
// @Description  Replaces the controller's fields. Its brand cannot change while machines use it, because each machine stores the brand of its controller.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Controllers
// @Tags         Machine catalog
// @ID           updateController
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateControllerRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.Controller}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/controllers/{id} [put]
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
//
// @Summary      Delete controller
// @Description  Soft-deletes the controller. A controller that still has models or is used by machines is kept: those must be deleted or changed first.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Controllers
// @Tags         Machine catalog
// @ID           deleteController
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/controllers/{id} [delete]
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
