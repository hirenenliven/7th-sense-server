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

// UserLocationService holds the business rules for the user's locations. Every method
// acts for userID only: a user never sees or changes another user's locations.
type UserLocationService struct {
	userLocationRepository *repository.UserLocationRepository
	branchService          *BranchService
}

// NewUserLocationService creates a UserLocationService that stores locations with userLocationRepository
// and checks their branch with branchService.
func NewUserLocationService(userLocationRepository *repository.UserLocationRepository, branchService *BranchService) *UserLocationService {
	return &UserLocationService{userLocationRepository: userLocationRepository, branchService: branchService}
}

// GetAll returns the user's locations. A non-nil branchID keeps only that branch's locations.
func (userLocationService *UserLocationService) GetAll(ctx context.Context, userID *int64, branchID *int64) ([]model.UserLocation, error) {
	return userLocationService.userLocationRepository.FindAll(ctx, userID, branchID)
}

func (userLocationService *UserLocationService) GetByID(ctx context.Context, userID, id int64) (model.UserLocation, error) {
	location, err := userLocationService.userLocationRepository.FindByID(ctx, userID, id)
	return location, toUserLocationResponseError(err)
}

// Create adds a location for userID under a branch.
func (userLocationService *UserLocationService) Create(ctx context.Context, userID int64, locationRequest dto.CreateUserLocationRequest) (model.UserLocation, error) {
	if err := userLocationService.checkBranch(ctx, locationRequest.BranchID); err != nil {
		return model.UserLocation{}, err
	}
	created, err := userLocationService.userLocationRepository.Create(ctx, model.UserLocation{
		UserID:   userID,
		BranchID: locationRequest.BranchID,
		Name:     locationRequest.Name,
	})
	return created, toUserLocationResponseError(err)
}

// Update renames the location and moves it under another branch.
func (userLocationService *UserLocationService) Update(ctx context.Context, userID, id int64, locationRequest dto.UpdateUserLocationRequest) (model.UserLocation, error) {
	if err := userLocationService.checkBranch(ctx, locationRequest.BranchID); err != nil {
		return model.UserLocation{}, err
	}
	updated, err := userLocationService.userLocationRepository.Update(ctx, model.UserLocation{
		ID:       id,
		UserID:   userID,
		BranchID: locationRequest.BranchID,
		Name:     locationRequest.Name,
	})
	return updated, toUserLocationResponseError(err)
}

// Delete soft-deletes the location. A location that still has shops is kept:
// its shops must be deleted first.
func (userLocationService *UserLocationService) Delete(ctx context.Context, userID, id int64) error {
	hasShops, err := userLocationService.userLocationRepository.HasShops(ctx, userID, id)
	if err != nil {
		return err
	}
	if hasShops {
		return response.NewError(http.StatusBadRequest, "This location has shops. Delete its shops first")
	}
	return toUserLocationResponseError(userLocationService.userLocationRepository.Delete(ctx, userID, id))
}

// checkBranch returns "Please enter valid branch" unless branchID is an existing branch.
func (userLocationService *UserLocationService) checkBranch(ctx context.Context, branchID int64) error {
	_, err := userLocationService.branchService.GetByID(ctx, branchID)
	if err != nil && err.Error() == "Branch not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid branch")
	}
	return err
}

// toUserLocationResponseError maps repository errors to location responses. Other errors
// pass through and are sent as ServerError.
func toUserLocationResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Location not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "Location already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid branch")
	default:
		return err
	}
}
