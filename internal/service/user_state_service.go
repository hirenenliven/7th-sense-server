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

// UserStateService holds the business rules for the user's states. Every method
// acts for userID only: a user never sees or changes another user's states.
type UserStateService struct {
	userStateRepository *repository.UserStateRepository
	userCountryService  *UserCountryService
}

// NewUserStateService creates a UserStateService that stores states with userStateRepository
// and checks their country with userCountryService.
func NewUserStateService(userStateRepository *repository.UserStateRepository, userCountryService *UserCountryService) *UserStateService {
	return &UserStateService{userStateRepository: userStateRepository, userCountryService: userCountryService}
}

// GetAll returns the user's states. A non-nil userCountryID keeps only that country's states.
func (userStateService *UserStateService) GetAll(ctx context.Context, userID *int64, userCountryID *int64) ([]model.UserState, error) {
	return userStateService.userStateRepository.FindAll(ctx, userID, userCountryID)
}

func (userStateService *UserStateService) GetByID(ctx context.Context, userID, id int64) (model.UserState, error) {
	state, err := userStateService.userStateRepository.FindByID(ctx, userID, id)
	return state, toUserStateResponseError(err)
}

// Create adds a state for userID under one of the user's countries.
func (userStateService *UserStateService) Create(ctx context.Context, userID int64, stateRequest dto.CreateUserStateRequest) (model.UserState, error) {
	if err := userStateService.checkCountry(ctx, userID, stateRequest.UserCountryID); err != nil {
		return model.UserState{}, err
	}
	created, err := userStateService.userStateRepository.Create(ctx, model.UserState{
		UserID:        userID,
		UserCountryID: stateRequest.UserCountryID,
		Name:          stateRequest.Name,
	})
	return created, toUserStateResponseError(err)
}

// Delete soft-deletes the state. A state that still has cities is kept:
// its cities must be deleted first.
func (userStateService *UserStateService) Delete(ctx context.Context, userID, id int64) error {
	hasCities, err := userStateService.userStateRepository.HasCities(ctx, userID, id)
	if err != nil {
		return err
	}
	if hasCities {
		return response.NewError(http.StatusBadRequest, "This state has cities. Delete its cities first")
	}
	return toUserStateResponseError(userStateService.userStateRepository.Delete(ctx, userID, id))
}

// checkCountry returns "Please enter valid country" unless userCountryID is one of the
// user's countries.
func (userStateService *UserStateService) checkCountry(ctx context.Context, userID, userCountryID int64) error {
	_, err := userStateService.userCountryService.GetByID(ctx, userID, userCountryID)
	if err != nil && err.Error() == "Country not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid country")
	}
	return err
}

// toUserStateResponseError maps repository errors to state responses. Other errors
// pass through and are sent as ServerError.
func toUserStateResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "State not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "State already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid country")
	default:
		return err
	}
}
