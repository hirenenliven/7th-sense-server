package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// BranchHandler handles HTTP requests for branches.
type BranchHandler struct {
	service *service.BranchService
}

// NewBranchHandler creates a BranchHandler that uses branchService.
func NewBranchHandler(branchService *service.BranchService) *BranchHandler {
	return &BranchHandler{service: branchService}
}

// GetAll handles GET /api/v1/branches[?userCityId=<id>].
func (branchHandler *BranchHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	userCityID, err := queryID(request, "userCityId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	branches, err := branchHandler.service.GetAll(request.Context(), userCityID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, branches)
}

// GetByID handles GET /api/v1/branches/{id}.
func (branchHandler *BranchHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	branch, err := branchHandler.service.GetByID(request.Context(), id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, branch)
}

// Create handles POST /api/v1/branches.
func (branchHandler *BranchHandler) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var branchRequest dto.CreateBranchRequest
	if err := decodeAndValidate(writer, request, &branchRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	branch, err := branchHandler.service.Create(request.Context(), userID, branchRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, branch)
}

// Update handles PUT /api/v1/branches/{id}.
func (branchHandler *BranchHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var branchRequest dto.UpdateBranchRequest
	if err := decodeAndValidate(writer, request, &branchRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	branch, err := branchHandler.service.Update(request.Context(), userID, id, branchRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, branch)
}

// Delete handles DELETE /api/v1/branches/{id}.
func (branchHandler *BranchHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := branchHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Branch deleted")
}
