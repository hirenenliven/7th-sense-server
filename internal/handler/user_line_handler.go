package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserLineHandler handles HTTP requests for the user's lines (/user-line).
// Every request acts for the current user (middleware.Authenticate); an admin may act
// for another user with ?userId=<id> (see targetUserID).
type UserLineHandler struct {
	service *service.UserLineService
}

// NewUserLineHandler creates a UserLineHandler that uses userLineService.
func NewUserLineHandler(userLineService *service.UserLineService) *UserLineHandler {
	return &UserLineHandler{service: userLineService}
}

// GetAll handles GET /api/v1/user-line[?userShopId=<id>][&userId=<id>].
//
// @Summary      Get user lines
// @Description  Returns the user's lines. A non-nil userShopID keeps only that shop's lines.
// @Description
// @Description  Used by pages: Add / edit machine, Lines, Cells, Add / edit customer or user, Customer / user detail
// @Tags         Organization
// @ID           getUserLines
// @Produce      json
// @Security     cookieAuth
// @Param        userShopId query int false "User shop id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope{data=[]model.UserLine}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-line [get]
func (userLineHandler *UserLineHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := listUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userShopID, err := queryID(request, "userShopId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	lines, err := userLineHandler.service.GetAll(request.Context(), userID, userShopID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, lines)
}

// GetByID handles GET /api/v1/user-line/{id}[?userId=<id>].
func (userLineHandler *UserLineHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	line, err := userLineHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, line)
}

// Create handles POST /api/v1/user-line[?userId=<id>].
//
// @Summary      Create user line
// @Description  Adds a line for userID under one of the user's shops.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Lines
// @Tags         Organization
// @ID           createUserLine
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.CreateUserLineRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.UserLine}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-line [post]
func (userLineHandler *UserLineHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var lineRequest dto.CreateUserLineRequest
	if err := decodeAndValidate(writer, request, &lineRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	line, err := userLineHandler.service.Create(request.Context(), userID, lineRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, line)
}

// Update handles PUT /api/v1/user-line/{id}[?userId=<id>].
//
// @Summary      Update user line
// @Description  Renames the line and moves it under another of the user's shops.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Lines
// @Tags         Organization
// @ID           updateUserLine
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.UpdateUserLineRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.UserLine}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-line/{id} [put]
func (userLineHandler *UserLineHandler) Update(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var lineRequest dto.UpdateUserLineRequest
	if err := decodeAndValidate(writer, request, &lineRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	line, err := userLineHandler.service.Update(request.Context(), userID, id, lineRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, line)
}

// Delete handles DELETE /api/v1/user-line/{id}[?userId=<id>].
//
// @Summary      Delete user line
// @Description  Soft-deletes the line. A line that still has cells is kept: its cells must be deleted first. A line with machines, or that a user's machine access names, is kept too.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Lines
// @Tags         Organization
// @ID           deleteUserLine
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-line/{id} [delete]
func (userLineHandler *UserLineHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := userLineHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Line deleted")
}
