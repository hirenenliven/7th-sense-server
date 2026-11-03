package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserCellHandler handles HTTP requests for the user's cells (/user-cell).
// Every request acts for the current user (middleware.Authenticate); an admin may act
// for another user with ?userId=<id> (see targetUserID).
type UserCellHandler struct {
	service *service.UserCellService
}

// NewUserCellHandler creates a UserCellHandler that uses userCellService.
func NewUserCellHandler(userCellService *service.UserCellService) *UserCellHandler {
	return &UserCellHandler{service: userCellService}
}

// GetAll handles GET /api/v1/user-cell[?userLineId=<id>][&userId=<id>].
//
// @Summary      Get user cells
// @Description  Returns the user's cells. A non-nil userLineID keeps only that line's cells.
// @Description
// @Description  Used by pages: Cells, Supervisors
// @Tags         Organization
// @ID           getUserCells
// @Produce      json
// @Security     cookieAuth
// @Param        userLineId query int false "User line id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Success      200 {object} response.Envelope{data=[]model.UserCell}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-cell [get]
func (userCellHandler *UserCellHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := listUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userLineID, err := queryID(request, "userLineId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	cells, err := userCellHandler.service.GetAll(request.Context(), userID, userLineID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, cells)
}

// GetByID handles GET /api/v1/user-cell/{id}[?userId=<id>].
func (userCellHandler *UserCellHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	cell, err := userCellHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, cell)
}

// Create handles POST /api/v1/user-cell[?userId=<id>].
//
// @Summary      Create user cell
// @Description  Adds a cell for userID under one of the user's lines.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Cells
// @Tags         Organization
// @ID           createUserCell
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.CreateUserCellRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.UserCell}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-cell [post]
func (userCellHandler *UserCellHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := targetUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var cellRequest dto.CreateUserCellRequest
	if err := decodeAndValidate(writer, request, &cellRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	cell, err := userCellHandler.service.Create(request.Context(), userID, cellRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, cell)
}

// Update handles PUT /api/v1/user-cell/{id}[?userId=<id>].
//
// @Summary      Update user cell
// @Description  Renames the cell and moves it under another of the user's lines.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Cells
// @Tags         Organization
// @ID           updateUserCell
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "Customer / company user id (superadmin only)"
// @Param        request body dto.UpdateUserCellRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.UserCell}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/user-cell/{id} [put]
func (userCellHandler *UserCellHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var cellRequest dto.UpdateUserCellRequest
	if err := decodeAndValidate(writer, request, &cellRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	cell, err := userCellHandler.service.Update(request.Context(), userID, id, cellRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, cell)
}

// Delete handles DELETE /api/v1/user-cell/{id}[?userId=<id>].
//
// @Summary      Delete user cell
// @Description  Soft-deletes the cell. A cell that still has supervisors is kept: its supervisors must be deleted first.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Cells
// @Tags         Organization
// @ID           deleteUserCell
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
// @Router       /v1/user-cell/{id} [delete]
func (userCellHandler *UserCellHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := userCellHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Cell deleted")
}
