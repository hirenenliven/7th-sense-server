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

// UserCountryService holds the business rules for the user's countries. Every method
// acts for userID only: a user never sees or changes another user's countries.
type UserCountryService struct {
	userCountryRepository *repository.UserCountryRepository
}

// NewUserCountryService creates a UserCountryService that uses userCountryRepository.
func NewUserCountryService(userCountryRepository *repository.UserCountryRepository) *UserCountryService {
	return &UserCountryService{userCountryRepository: userCountryRepository}
}

// GetAll returns the user's countries.
func (userCountryService *UserCountryService) GetAll(ctx context.Context, userID *int64) ([]model.UserCountry, error) {
	return userCountryService.userCountryRepository.FindAll(ctx, userID)
}

func (userCountryService *UserCountryService) GetByID(ctx context.Context, userID, id int64) (model.UserCountry, error) {
	country, err := userCountryService.userCountryRepository.FindByID(ctx, userID, id)
	return country, toUserCountryResponseError(err)
}

// Create adds a country for userID.
func (userCountryService *UserCountryService) Create(ctx context.Context, userID int64, countryRequest dto.CreateUserCountryRequest) (model.UserCountry, error) {
	created, err := userCountryService.userCountryRepository.Create(ctx, model.UserCountry{UserID: userID, Name: countryRequest.Name})
	return created, toUserCountryResponseError(err)
}

// Delete soft-deletes the country. A country that still has states is kept:
// its states must be deleted first.
func (userCountryService *UserCountryService) Delete(ctx context.Context, userID, id int64) error {
	hasStates, err := userCountryService.userCountryRepository.HasStates(ctx, userID, id)
	if err != nil {
		return err
	}
	if hasStates {
		return response.NewError(http.StatusBadRequest, "This country has states. Delete its states first")
	}
	return toUserCountryResponseError(userCountryService.userCountryRepository.Delete(ctx, userID, id))
}

// toUserCountryResponseError maps repository errors to country responses. Other errors
// pass through and are sent as ServerError.
func toUserCountryResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Country not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "Country already exists")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please enter valid user")
	default:
		return err
	}
}
