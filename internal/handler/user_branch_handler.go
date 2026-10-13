package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserBranchHandler handles HTTP requests for the user's branches (/user-branch).
// Every request acts for the current user (middleware.Authenticate), never for
// a user id sent by the client.
type UserBranchHandler struct {
	service *service.UserBranchService
}

// NewUserBranchHandler creates a UserBranchHandler that uses userBranchService.
func NewUserBranchHandler(userBranchService *service.UserBranchService) *UserBranchHandler {
	return &UserBranchHandler{service: userBranchService}
}

// GetAll handles GET /api/v1/user-branch[?userCityId=<id>].
func (userBranchHandler *UserBranchHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	userCityID, err := queryID(request, "userCityId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	branches, err := userBranchHandler.service.GetAll(request.Context(), userID, userCityID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, branches)
}

// GetByID handles GET /api/v1/user-branch/{id}.
func (userBranchHandler *UserBranchHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	branch, err := userBranchHandler.service.GetByID(request.Context(), userID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, branch)
}

// Create handles POST /api/v1/user-branch.
func (userBranchHandler *UserBranchHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var branchRequest dto.CreateUserBranchRequest
	if err := decodeAndValidate(writer, request, &branchRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	branch, err := userBranchHandler.service.Create(request.Context(), userID, branchRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, branch)
}

// Update handles PUT /api/v1/user-branch/{id}.
func (userBranchHandler *UserBranchHandler) Update(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var branchRequest dto.UpdateUserBranchRequest
	if err := decodeAndValidate(writer, request, &branchRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	branch, err := userBranchHandler.service.Update(request.Context(), userID, id, branchRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, branch)
}

// Delete handles DELETE /api/v1/user-branch/{id}.
func (userBranchHandler *UserBranchHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := userBranchHandler.service.Delete(request.Context(), userID, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Branch deleted")
}
