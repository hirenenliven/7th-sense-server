package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// MachineAccuracyService holds the rules of the machine accuracy list (server ADR 0017): one list
// shared by every customer, read by every logged-in user (machine and part forms), written by the
// superadmin only (superAdminOnly routes, Account Center › Machine Accuracy). Names are unique,
// ignoring case. An accuracy used by machines or part processes cannot be deleted.
type MachineAccuracyService struct {
	accuracyRepository *repository.MachineAccuracyRepository
}

// NewMachineAccuracyService creates a MachineAccuracyService that stores accuracies with
// accuracyRepository.
func NewMachineAccuracyService(accuracyRepository *repository.MachineAccuracyRepository) *MachineAccuracyService {
	return &MachineAccuracyService{accuracyRepository: accuracyRepository}
}

// GetAll returns every accuracy, ordered by id.
func (accuracyService *MachineAccuracyService) GetAll(ctx context.Context) ([]model.MachineAccuracy, error) {
	return accuracyService.accuracyRepository.FindAll(ctx)
}

// Create adds an accuracy; a name already in the list answers MachineAccuracyExists.
func (accuracyService *MachineAccuracyService) Create(ctx context.Context, accuracyRequest dto.MachineAccuracyRequest) (model.MachineAccuracy, error) {
	accuracy, err := accuracyService.accuracyRepository.Create(ctx, accuracyRequest.Name)
	return accuracy, toMachineAccuracyResponseError(err)
}

// Update renames accuracy id and returns it.
func (accuracyService *MachineAccuracyService) Update(ctx context.Context, id int64, accuracyRequest dto.MachineAccuracyRequest) (model.MachineAccuracy, error) {
	if err := accuracyService.accuracyRepository.Update(ctx, id, accuracyRequest.Name); err != nil {
		return model.MachineAccuracy{}, toMachineAccuracyResponseError(err)
	}
	accuracy, err := accuracyService.accuracyRepository.FindByID(ctx, id)
	return accuracy, toMachineAccuracyResponseError(err)
}

// Delete soft-deletes accuracy id unless machines or part processes use it (they would lose
// their accuracy).
func (accuracyService *MachineAccuracyService) Delete(ctx context.Context, id int64) error {
	accuracy, err := accuracyService.accuracyRepository.FindByID(ctx, id)
	if err != nil {
		return toMachineAccuracyResponseError(err)
	}
	if accuracy.UsageCount > 0 {
		return response.NewError(http.StatusBadRequest,
			fmt.Sprintf("This accuracy is used by %d machine(s) or part process(es). Change it there first", accuracy.UsageCount))
	}
	return toMachineAccuracyResponseError(accuracyService.accuracyRepository.Delete(ctx, id))
}

// checkExists returns invalid unless id is an accuracy of the list.
func (accuracyService *MachineAccuracyService) checkExists(ctx context.Context, id int64, invalid error) error {
	_, err := accuracyService.accuracyRepository.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return invalid
	}
	return err
}

// toMachineAccuracyResponseError maps repository errors to machine accuracy responses. Other
// errors pass through and are sent as ServerError.
func toMachineAccuracyResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.MachineAccuracyMissing
	case errors.Is(err, repository.ErrDuplicateName):
		return response.MachineAccuracyExists
	default:
		return err
	}
}
