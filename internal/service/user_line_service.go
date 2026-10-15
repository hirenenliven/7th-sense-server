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

// UserLineService holds the business rules for the user's lines. Every method
// acts for userID only: a user never sees or changes another user's lines.
type UserLineService struct {
	userLineRepository *repository.UserLineRepository
	userShopService    *UserShopService
}

// NewUserLineService creates a UserLineService that stores lines with userLineRepository
// and checks their shop with userShopService.
func NewUserLineService(userLineRepository *repository.UserLineRepository, userShopService *UserShopService) *UserLineService {
	return &UserLineService{userLineRepository: userLineRepository, userShopService: userShopService}
}

// GetAll returns the user's lines. A non-nil userShopID keeps only that shop's lines.
func (userLineService *UserLineService) GetAll(ctx context.Context, userID *int64, userShopID *int64) ([]model.UserLine, error) {
	return userLineService.userLineRepository.FindAll(ctx, userID, userShopID)
}

func (userLineService *UserLineService) GetByID(ctx context.Context, userID, id int64) (model.UserLine, error) {
	line, err := userLineService.userLineRepository.FindByID(ctx, userID, id)
	return line, toUserLineResponseError(err)
}

// Create adds a line for userID under one of the user's shops.
func (userLineService *UserLineService) Create(ctx context.Context, userID int64, lineRequest dto.CreateUserLineRequest) (model.UserLine, error) {
	if err := userLineService.checkShop(ctx, userID, lineRequest.UserShopID); err != nil {
		return model.UserLine{}, err
	}
	created, err := userLineService.userLineRepository.Create(ctx, model.UserLine{
		UserID:     userID,
		UserShopID: lineRequest.UserShopID,
		Name:       lineRequest.Name,
	})
	return created, toUserLineResponseError(err)
}

// Update renames the line and moves it under another of the user's shops.
func (userLineService *UserLineService) Update(ctx context.Context, userID, id int64, lineRequest dto.UpdateUserLineRequest) (model.UserLine, error) {
	if err := userLineService.checkShop(ctx, userID, lineRequest.UserShopID); err != nil {
		return model.UserLine{}, err
	}
	updated, err := userLineService.userLineRepository.Update(ctx, model.UserLine{
		ID:         id,
		UserID:     userID,
		UserShopID: lineRequest.UserShopID,
		Name:       lineRequest.Name,
	})
	return updated, toUserLineResponseError(err)
}

// Delete soft-deletes the line. A line that still has cells is kept:
// its cells must be deleted first.
func (userLineService *UserLineService) Delete(ctx context.Context, userID, id int64) error {
	hasCells, err := userLineService.userLineRepository.HasCells(ctx, userID, id)
	if err != nil {
		return err
	}
	if hasCells {
		return response.NewError(http.StatusBadRequest, "This line has cells. Delete its cells first")
	}
	return toUserLineResponseError(userLineService.userLineRepository.Delete(ctx, userID, id))
}

// checkShop returns "Please enter valid shop" unless userShopID is one of the
// user's shops.
func (userLineService *UserLineService) checkShop(ctx context.Context, userID, userShopID int64) error {
	_, err := userLineService.userShopService.GetByID(ctx, userID, userShopID)
	if err != nil && err.Error() == "Shop not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid shop")
	}
	return err
}

// toUserLineResponseError maps repository errors to line responses. Other errors
// pass through and are sent as ServerError.
func toUserLineResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Line not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "Line already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid shop")
	default:
		return err
	}
}
