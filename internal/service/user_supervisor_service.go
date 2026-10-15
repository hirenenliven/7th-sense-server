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

// UserSupervisorService holds the business rules for the user's supervisors. Every method
// acts for userID only: a user never sees or changes another user's supervisors.
type UserSupervisorService struct {
	userSupervisorRepository *repository.UserSupervisorRepository
	userCellService          *UserCellService
}

// NewUserSupervisorService creates a UserSupervisorService that stores supervisors with userSupervisorRepository
// and checks their cell with userCellService.
func NewUserSupervisorService(userSupervisorRepository *repository.UserSupervisorRepository, userCellService *UserCellService) *UserSupervisorService {
	return &UserSupervisorService{userSupervisorRepository: userSupervisorRepository, userCellService: userCellService}
}

// GetAll returns the user's supervisors. A non-nil userCellID keeps only that cell's supervisors.
func (userSupervisorService *UserSupervisorService) GetAll(ctx context.Context, userID *int64, userCellID *int64) ([]model.UserSupervisor, error) {
	return userSupervisorService.userSupervisorRepository.FindAll(ctx, userID, userCellID)
}

func (userSupervisorService *UserSupervisorService) GetByID(ctx context.Context, userID, id int64) (model.UserSupervisor, error) {
	supervisor, err := userSupervisorService.userSupervisorRepository.FindByID(ctx, userID, id)
	return supervisor, toUserSupervisorResponseError(err)
}

// Create adds a supervisor for userID under one of the user's cells.
func (userSupervisorService *UserSupervisorService) Create(ctx context.Context, userID int64, supervisorRequest dto.CreateUserSupervisorRequest) (model.UserSupervisor, error) {
	if err := userSupervisorService.checkCell(ctx, userID, supervisorRequest.UserCellID); err != nil {
		return model.UserSupervisor{}, err
	}
	created, err := userSupervisorService.userSupervisorRepository.Create(ctx, model.UserSupervisor{
		UserID:     userID,
		UserCellID: supervisorRequest.UserCellID,
		Name:       supervisorRequest.Name,
	})
	return created, toUserSupervisorResponseError(err)
}

// Update renames the supervisor and moves it under another of the user's cells.
func (userSupervisorService *UserSupervisorService) Update(ctx context.Context, userID, id int64, supervisorRequest dto.UpdateUserSupervisorRequest) (model.UserSupervisor, error) {
	if err := userSupervisorService.checkCell(ctx, userID, supervisorRequest.UserCellID); err != nil {
		return model.UserSupervisor{}, err
	}
	updated, err := userSupervisorService.userSupervisorRepository.Update(ctx, model.UserSupervisor{
		ID:         id,
		UserID:     userID,
		UserCellID: supervisorRequest.UserCellID,
		Name:       supervisorRequest.Name,
	})
	return updated, toUserSupervisorResponseError(err)
}

// Delete soft-deletes the supervisor.
func (userSupervisorService *UserSupervisorService) Delete(ctx context.Context, userID, id int64) error {
	return toUserSupervisorResponseError(userSupervisorService.userSupervisorRepository.Delete(ctx, userID, id))
}

// checkCell returns "Please enter valid cell" unless userCellID is one of the
// user's cells.
func (userSupervisorService *UserSupervisorService) checkCell(ctx context.Context, userID, userCellID int64) error {
	_, err := userSupervisorService.userCellService.GetByID(ctx, userID, userCellID)
	if err != nil && err.Error() == "Cell not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid cell")
	}
	return err
}

// toUserSupervisorResponseError maps repository errors to supervisor responses. Other errors
// pass through and are sent as ServerError.
func toUserSupervisorResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Supervisor not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "Supervisor already exists")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please enter valid cell")
	default:
		return err
	}
}
