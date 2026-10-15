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
