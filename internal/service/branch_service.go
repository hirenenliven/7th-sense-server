package service

import (
	"context"
	"errors"
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// BranchService holds the business rules for branches. A branch can belong to one
// of the current user's cities (user_city) and has locations under it.
type BranchService struct {
	branchRepository *repository.BranchRepository
	userCityService  *UserCityService
}

// NewBranchService creates a BranchService that stores branches with branchRepository
// and checks their city with userCityService.
func NewBranchService(branchRepository *repository.BranchRepository, userCityService *UserCityService) *BranchService {
	return &BranchService{branchRepository: branchRepository, userCityService: userCityService}
}

// GetAll returns the branches. A non-nil userCityID keeps only that city's branches.
func (branchService *BranchService) GetAll(ctx context.Context, userCityID *int64) ([]model.Branch, error) {
	return branchService.branchRepository.FindAll(ctx, userCityID)
}

func (branchService *BranchService) GetByID(ctx context.Context, id int64) (model.Branch, error) {
	branch, err := branchService.branchRepository.FindByID(ctx, id)
	return branch, toBranchResponseError(err)
}

// Create adds a branch. A city, when given, must be one of userID's cities.
func (branchService *BranchService) Create(ctx context.Context, userID int64, branchRequest dto.CreateBranchRequest) (model.Branch, error) {
	if err := branchService.checkCity(ctx, userID, branchRequest.UserCityID); err != nil {
		return model.Branch{}, err
	}
	branch, err := branchService.branchRepository.Create(ctx, model.Branch{
		Code:       branchRequest.Code,
		Name:       branchRequest.Name,
		Email:      branchRequest.Email,
		Address:    branchRequest.Address,
		City:       branchRequest.City,
		StateID:    branchRequest.StateID,
		UserCityID: branchRequest.UserCityID,
		CustomerID: branchRequest.CustomerID,
		GSTNumber:  branchRequest.GSTNumber,
		PANNumber:  branchRequest.PANNumber,
		ColorCode:  branchRequest.ColorCode,
		Status:     branchRequest.Status,
	})
	return branch, toBranchResponseError(err)
}

// Update replaces the branch. A city, when given, must be one of userID's cities.
func (branchService *BranchService) Update(ctx context.Context, userID, id int64, branchRequest dto.UpdateBranchRequest) (model.Branch, error) {
	if err := branchService.checkCity(ctx, userID, branchRequest.UserCityID); err != nil {
		return model.Branch{}, err
	}
	branch, err := branchService.branchRepository.Update(ctx, model.Branch{
		ID:         id,
		Code:       branchRequest.Code,
		Name:       branchRequest.Name,
		Email:      branchRequest.Email,
		Address:    branchRequest.Address,
		City:       branchRequest.City,
		StateID:    branchRequest.StateID,
		UserCityID: branchRequest.UserCityID,
		CustomerID: branchRequest.CustomerID,
		GSTNumber:  branchRequest.GSTNumber,
		PANNumber:  branchRequest.PANNumber,
		ColorCode:  branchRequest.ColorCode,
		Status:     branchRequest.Status,
	})
	return branch, toBranchResponseError(err)
}

// Delete soft-deletes the branch. A branch that still has locations is kept:
// its locations must be deleted first.
func (branchService *BranchService) Delete(ctx context.Context, id int64) error {
	hasLocations, err := branchService.branchRepository.HasLocations(ctx, id)
	if err != nil {
		return err
	}
	if hasLocations {
		return response.NewError(http.StatusBadRequest, "This branch has locations. Delete its locations first")
	}
	return toBranchResponseError(branchService.branchRepository.Delete(ctx, id))
}

// checkCity returns "Please enter valid city" unless userCityID is nil or one of
// userID's cities.
func (branchService *BranchService) checkCity(ctx context.Context, userID int64, userCityID *int64) error {
	if userCityID == nil {
		return nil
	}
	_, err := branchService.userCityService.GetByID(ctx, userID, *userCityID)
	if err != nil && err.Error() == "City not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid city")
	}
	return err
}

// toBranchResponseError maps repository errors to branch responses. Other errors
// pass through and are sent as ServerError.
func toBranchResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Branch not found")
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.NewError(http.StatusBadRequest, "Code already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid city")
	default:
		return err
	}
}
