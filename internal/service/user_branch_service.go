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

// UserBranchService holds the business rules for the user's branches. Every method
// acts for userID only: a user never sees or changes another user's branches.
type UserBranchService struct {
	userBranchRepository *repository.UserBranchRepository
	userCityService      *UserCityService
}

// NewUserBranchService creates a UserBranchService that stores branches with userBranchRepository
// and checks their city with userCityService.
func NewUserBranchService(userBranchRepository *repository.UserBranchRepository, userCityService *UserCityService) *UserBranchService {
	return &UserBranchService{userBranchRepository: userBranchRepository, userCityService: userCityService}
}

// GetAll returns the user's branches. A non-nil userCityID keeps only that city's branches.
func (userBranchService *UserBranchService) GetAll(ctx context.Context, userID int64, userCityID *int64) ([]model.UserBranch, error) {
	return userBranchService.userBranchRepository.FindAll(ctx, userID, userCityID)
}

func (userBranchService *UserBranchService) GetByID(ctx context.Context, userID, id int64) (model.UserBranch, error) {
	branch, err := userBranchService.userBranchRepository.FindByID(ctx, userID, id)
	return branch, toUserBranchResponseError(err)
}

// Create adds a branch for userID under one of the user's cities.
func (userBranchService *UserBranchService) Create(ctx context.Context, userID int64, branchRequest dto.CreateUserBranchRequest) (model.UserBranch, error) {
	if err := userBranchService.checkCity(ctx, userID, branchRequest.UserCityID); err != nil {
		return model.UserBranch{}, err
	}
	created, err := userBranchService.userBranchRepository.Create(ctx, model.UserBranch{
		UserID:     userID,
		UserCityID: branchRequest.UserCityID,
		Name:       branchRequest.Name,
	})
	return created, toUserBranchResponseError(err)
}

// Update renames the branch and moves it under another of the user's cities.
func (userBranchService *UserBranchService) Update(ctx context.Context, userID, id int64, branchRequest dto.UpdateUserBranchRequest) (model.UserBranch, error) {
	if err := userBranchService.checkCity(ctx, userID, branchRequest.UserCityID); err != nil {
		return model.UserBranch{}, err
	}
	updated, err := userBranchService.userBranchRepository.Update(ctx, model.UserBranch{
		ID:         id,
		UserID:     userID,
		UserCityID: branchRequest.UserCityID,
		Name:       branchRequest.Name,
	})
	return updated, toUserBranchResponseError(err)
}

// Delete soft-deletes the branch. A branch that still has locations is kept:
// its locations must be deleted first.
func (userBranchService *UserBranchService) Delete(ctx context.Context, userID, id int64) error {
	hasLocations, err := userBranchService.userBranchRepository.HasLocations(ctx, userID, id)
	if err != nil {
		return err
	}
	if hasLocations {
		return response.NewError(http.StatusBadRequest, "This branch has locations. Delete its locations first")
	}
	return toUserBranchResponseError(userBranchService.userBranchRepository.Delete(ctx, userID, id))
}

// checkCity returns "Please enter valid city" unless userCityID is one of the
// user's cities.
func (userBranchService *UserBranchService) checkCity(ctx context.Context, userID, userCityID int64) error {
	_, err := userBranchService.userCityService.GetByID(ctx, userID, userCityID)
	if err != nil && err.Error() == "City not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid city")
	}
	return err
}

// toUserBranchResponseError maps repository errors to branch responses. Other errors
// pass through and are sent as ServerError.
func toUserBranchResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Branch not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "Branch already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid city")
	default:
		return err
	}
}
