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

// UserCellService holds the business rules for the user's cells. Every method
// acts for userID only: a user never sees or changes another user's cells.
type UserCellService struct {
	userCellRepository *repository.UserCellRepository
	userLineService    *UserLineService
}

// NewUserCellService creates a UserCellService that stores cells with userCellRepository
// and checks their line with userLineService.
func NewUserCellService(userCellRepository *repository.UserCellRepository, userLineService *UserLineService) *UserCellService {
	return &UserCellService{userCellRepository: userCellRepository, userLineService: userLineService}
}

// GetAll returns the user's cells. A non-nil userLineID keeps only that line's cells.
func (userCellService *UserCellService) GetAll(ctx context.Context, userID *int64, userLineID *int64) ([]model.UserCell, error) {
	return userCellService.userCellRepository.FindAll(ctx, userID, userLineID)
}

func (userCellService *UserCellService) GetByID(ctx context.Context, userID, id int64) (model.UserCell, error) {
	cell, err := userCellService.userCellRepository.FindByID(ctx, userID, id)
	return cell, toUserCellResponseError(err)
}

// Create adds a cell for userID under one of the user's lines.
func (userCellService *UserCellService) Create(ctx context.Context, userID int64, cellRequest dto.CreateUserCellRequest) (model.UserCell, error) {
	if err := userCellService.checkLine(ctx, userID, cellRequest.UserLineID); err != nil {
		return model.UserCell{}, err
	}
	created, err := userCellService.userCellRepository.Create(ctx, model.UserCell{
		UserID:     userID,
		UserLineID: cellRequest.UserLineID,
		Name:       cellRequest.Name,
	})
	return created, toUserCellResponseError(err)
}

// Update renames the cell and moves it under another of the user's lines.
func (userCellService *UserCellService) Update(ctx context.Context, userID, id int64, cellRequest dto.UpdateUserCellRequest) (model.UserCell, error) {
	if err := userCellService.checkLine(ctx, userID, cellRequest.UserLineID); err != nil {
		return model.UserCell{}, err
	}
	updated, err := userCellService.userCellRepository.Update(ctx, model.UserCell{
		ID:         id,
		UserID:     userID,
		UserLineID: cellRequest.UserLineID,
		Name:       cellRequest.Name,
	})
	return updated, toUserCellResponseError(err)
}

// Delete soft-deletes the cell. A cell that still has supervisors is kept:
// its supervisors must be deleted first.
func (userCellService *UserCellService) Delete(ctx context.Context, userID, id int64) error {
	hasSupervisors, err := userCellService.userCellRepository.HasSupervisors(ctx, userID, id)
	if err != nil {
		return err
	}
	if hasSupervisors {
		return response.NewError(http.StatusBadRequest, "This cell has supervisors. Delete its supervisors first")
	}
	return toUserCellResponseError(userCellService.userCellRepository.Delete(ctx, userID, id))
}

// checkLine returns "Please enter valid line" unless userLineID is one of the
// user's lines.
func (userCellService *UserCellService) checkLine(ctx context.Context, userID, userLineID int64) error {
	_, err := userCellService.userLineService.GetByID(ctx, userID, userLineID)
	if err != nil && err.Error() == "Line not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid line")
	}
	return err
}

// toUserCellResponseError maps repository errors to cell responses. Other errors
// pass through and are sent as ServerError.
func toUserCellResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Cell not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "Cell already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid line")
	default:
		return err
	}
}
