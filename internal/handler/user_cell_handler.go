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
