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
//
// @Summary      Get branches
// @Description  Returns the branches. A non-nil userCityID keeps only that city's branches.
// @Description
// @Description  Used by pages: Add / edit customer or user, Branches, Locations
// @Tags         Branches
// @ID           getBranches
// @Produce      json
// @Security     cookieAuth
// @Param        userCityId query int false "User city id"
// @Success      200 {object} response.Envelope{data=[]model.Branch}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/branches [get]
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
//
// @Summary      Get branch by id
// @Description  GetByID handles GET /api/v1/branches/{id}.
// @Description
// @Description  Used by pages: Customer / user detail, Add / edit branch
// @Tags         Branches
// @ID           getBranchById
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope{data=model.Branch}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/branches/{id} [get]
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
//
// @Summary      Create branch
// @Description  Adds a branch. A city, when given, must be one of the branch customer's cities.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Add / edit branch
// @Tags         Branches
// @ID           createBranch
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateBranchRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.Branch}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/branches [post]
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
//
// @Summary      Update branch
// @Description  Replaces the branch. A city, when given, must be one of the branch customer's cities.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Branches, Add / edit branch
// @Tags         Branches
// @ID           updateBranch
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateBranchRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.Branch}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/branches/{id} [put]
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
//
// @Summary      Delete branch
// @Description  Soft-deletes the branch. A branch that still has locations is kept: its locations must be deleted first. A branch a user's machine access names is kept too.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Branches
// @Tags         Branches
// @ID           deleteBranch
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/branches/{id} [delete]
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
