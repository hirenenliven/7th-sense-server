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

// UserCityService holds the business rules for the user's cities. Every method
// acts for userID only: a user never sees or changes another user's cities.
type UserCityService struct {
	userCityRepository *repository.UserCityRepository
	userStateService   *UserStateService
}

// NewUserCityService creates a UserCityService that stores cities with userCityRepository
// and checks their state with userStateService.
func NewUserCityService(userCityRepository *repository.UserCityRepository, userStateService *UserStateService) *UserCityService {
	return &UserCityService{userCityRepository: userCityRepository, userStateService: userStateService}
}

// GetAll returns the user's cities. A non-nil userStateID keeps only that state's cities.
func (userCityService *UserCityService) GetAll(ctx context.Context, userID int64, userStateID *int64) ([]model.UserCity, error) {
	return userCityService.userCityRepository.FindAll(ctx, userID, userStateID)
}

func (userCityService *UserCityService) GetByID(ctx context.Context, userID, id int64) (model.UserCity, error) {
	city, err := userCityService.userCityRepository.FindByID(ctx, userID, id)
	return city, toUserCityResponseError(err)
}

// Create adds a city for userID under one of the user's states.
func (userCityService *UserCityService) Create(ctx context.Context, userID int64, cityRequest dto.CreateUserCityRequest) (model.UserCity, error) {
	if err := userCityService.checkState(ctx, userID, cityRequest.UserStateID); err != nil {
		return model.UserCity{}, err
	}
	created, err := userCityService.userCityRepository.Create(ctx, model.UserCity{
		UserID:      userID,
		UserStateID: cityRequest.UserStateID,
		Name:        cityRequest.Name,
	})
	return created, toUserCityResponseError(err)
}

// Update renames the city and moves it under another of the user's states.
func (userCityService *UserCityService) Update(ctx context.Context, userID, id int64, cityRequest dto.UpdateUserCityRequest) (model.UserCity, error) {
	if err := userCityService.checkState(ctx, userID, cityRequest.UserStateID); err != nil {
		return model.UserCity{}, err
	}
	updated, err := userCityService.userCityRepository.Update(ctx, model.UserCity{
		ID:          id,
		UserID:      userID,
		UserStateID: cityRequest.UserStateID,
		Name:        cityRequest.Name,
	})
	return updated, toUserCityResponseError(err)
}

// Delete soft-deletes the city. A city that still has branches is kept:
// its branches must be deleted first.
func (userCityService *UserCityService) Delete(ctx context.Context, userID, id int64) error {
	hasBranches, err := userCityService.userCityRepository.HasBranches(ctx, id)
	if err != nil {
		return err
	}
	if hasBranches {
		return response.NewError(http.StatusBadRequest, "This city has branches. Delete its branches first")
	}
	return toUserCityResponseError(userCityService.userCityRepository.Delete(ctx, userID, id))
}

// checkState returns "Please enter valid state" unless userStateID is one of the
// user's states.
func (userCityService *UserCityService) checkState(ctx context.Context, userID, userStateID int64) error {
	_, err := userCityService.userStateService.GetByID(ctx, userID, userStateID)
	if err != nil && err.Error() == "State not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid state")
	}
	return err
}

// toUserCityResponseError maps repository errors to city responses. Other errors
// pass through and are sent as ServerError.
func toUserCityResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "City not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "City already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid state")
	default:
		return err
	}
}
