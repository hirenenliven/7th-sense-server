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

// UserShopService holds the business rules for the user's shops. Every method
// acts for userID only: a user never sees or changes another user's shops.
type UserShopService struct {
	userShopRepository  *repository.UserShopRepository
	userLocationService *UserLocationService
}

// NewUserShopService creates a UserShopService that stores shops with userShopRepository
// and checks their location with userLocationService.
func NewUserShopService(userShopRepository *repository.UserShopRepository, userLocationService *UserLocationService) *UserShopService {
	return &UserShopService{userShopRepository: userShopRepository, userLocationService: userLocationService}
}

// GetAll returns the user's shops. A non-nil userLocationID keeps only that location's shops.
func (userShopService *UserShopService) GetAll(ctx context.Context, userID int64, userLocationID *int64) ([]model.UserShop, error) {
	return userShopService.userShopRepository.FindAll(ctx, userID, userLocationID)
}

func (userShopService *UserShopService) GetByID(ctx context.Context, userID, id int64) (model.UserShop, error) {
	shop, err := userShopService.userShopRepository.FindByID(ctx, userID, id)
	return shop, toUserShopResponseError(err)
}

// Create adds a shop for userID under one of the user's locations.
func (userShopService *UserShopService) Create(ctx context.Context, userID int64, shopRequest dto.CreateUserShopRequest) (model.UserShop, error) {
	if err := userShopService.checkLocation(ctx, userID, shopRequest.UserLocationID); err != nil {
		return model.UserShop{}, err
	}
	created, err := userShopService.userShopRepository.Create(ctx, model.UserShop{
		UserID:         userID,
		UserLocationID: shopRequest.UserLocationID,
		Name:           shopRequest.Name,
	})
	return created, toUserShopResponseError(err)
}

// Update renames the shop and moves it under another of the user's locations.
func (userShopService *UserShopService) Update(ctx context.Context, userID, id int64, shopRequest dto.UpdateUserShopRequest) (model.UserShop, error) {
	if err := userShopService.checkLocation(ctx, userID, shopRequest.UserLocationID); err != nil {
		return model.UserShop{}, err
	}
	updated, err := userShopService.userShopRepository.Update(ctx, model.UserShop{
		ID:             id,
		UserID:         userID,
		UserLocationID: shopRequest.UserLocationID,
		Name:           shopRequest.Name,
	})
	return updated, toUserShopResponseError(err)
}

// Delete soft-deletes the shop. A shop that still has lines is kept:
// its lines must be deleted first.
func (userShopService *UserShopService) Delete(ctx context.Context, userID, id int64) error {
	hasLines, err := userShopService.userShopRepository.HasLines(ctx, userID, id)
	if err != nil {
		return err
	}
	if hasLines {
		return response.NewError(http.StatusBadRequest, "This shop has lines. Delete its lines first")
	}
	return toUserShopResponseError(userShopService.userShopRepository.Delete(ctx, userID, id))
}

// checkLocation returns "Please enter valid location" unless userLocationID is one of the
// user's locations.
func (userShopService *UserShopService) checkLocation(ctx context.Context, userID, userLocationID int64) error {
	_, err := userShopService.userLocationService.GetByID(ctx, userID, userLocationID)
	if err != nil && err.Error() == "Location not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid location")
	}
	return err
}

// toUserShopResponseError maps repository errors to shop responses. Other errors
// pass through and are sent as ServerError.
func toUserShopResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Shop not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "Shop already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid location")
	default:
		return err
	}
}
